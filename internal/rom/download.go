package rom

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/logging"
)

const (
	defaultChunkSize   = int64(16 << 20)
	defaultConcurrency = 16
	probeSize          = int64(1 << 20)
)

type downloadState struct {
	Version   int    `json:"version"`
	Size      int64  `json:"size"`
	ChunkSize int64  `json:"chunk_size"`
	MD5       string `json:"md5"`
	Done      []bool `json:"done"`
}

type mirrorProbe struct {
	URL      string
	Size     int64
	Duration time.Duration
}

type DownloadOptions struct {
	Mirrors     []string
	Connections int
}

func Download(ctx context.Context, client *http.Client, release Release, outputDir string, output ...io.Writer) (string, error) {
	return DownloadWithOptions(ctx, client, release, outputDir, DownloadOptions{}, output...)
}

func DownloadWithOptions(ctx context.Context, client *http.Client, release Release, outputDir string, options DownloadOptions, output ...io.Writer) (string, error) {
	if err := ValidateOfficialURL(release.URL); err != nil {
		return "", err
	}
	for _, mirror := range options.Mirrors {
		if err := ValidateMirrorBaseURL(mirror); err != nil {
			return "", i18n.WrapError("error.mirror_invalid", err, mirror)
		}
	}
	connections := options.Connections
	if connections == 0 {
		connections = defaultConcurrency
	}
	if connections < 1 || connections > 32 {
		return "", i18n.NewError("error.connections_invalid")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	u, err := url.Parse(release.URL)
	if err != nil {
		return "", err
	}
	name := filepath.Base(u.Path)
	if filepath.Ext(name) != ".zip" {
		return "", i18n.NewError("error.rom_url_not_zip")
	}
	destination := filepath.Join(outputDir, name)
	out := io.Discard
	if len(output) > 0 && output[0] != nil {
		out = output[0]
	}
	if existing, ok, err := reuseLocalROM(ctx, outputDir, destination, release, out); err != nil {
		return "", err
	} else if ok {
		return existing, nil
	}

	mirrors, err := probeMirrors(ctx, client, release.URL, options.Mirrors, out)
	if err != nil {
		return "", err
	}
	return downloadRanges(ctx, client, release, destination, mirrors, connections, out)
}

// reuseLocalROM looks for an already downloaded ROM of the same release before
// touching the network mirrors. The expected destination name is preferred, but
// any ZIP carrying the release version is accepted so a previously fetched OTA
// package is reused even if its filename differs. A candidate is only reused
// when its MD5 matches the online index, so the catalog stays authoritative.
func reuseLocalROM(ctx context.Context, outputDir, destination string, release Release, out io.Writer) (string, bool, error) {
	logger := logging.WithWriter(out, logging.Language(ctx))
	candidates, err := localROMCandidates(outputDir, destination, release.Version)
	if err != nil {
		return "", false, err
	}
	for _, candidate := range candidates {
		actual, err := fileMD5(candidate)
		if err != nil {
			return "", false, err
		}
		if actual == release.MD5 {
			logger.Info("rom.local_reused", "path", candidate)
			return candidate, true, nil
		}
		logger.Warn("rom.local_md5_mismatch", "path", candidate)
		if candidate == destination {
			if err := os.Remove(candidate); err != nil {
				return "", false, err
			}
		}
	}
	return "", false, nil
}

func localROMCandidates(outputDir, destination, version string) ([]string, error) {
	seen := map[string]bool{}
	var candidates []string
	if info, err := os.Stat(destination); err == nil && !info.IsDir() {
		candidates = append(candidates, destination)
		seen[destination] = true
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	version = strings.ToLower(strings.TrimSpace(version))
	if version == "" {
		return candidates, nil
	}
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return candidates, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".zip") {
			continue
		}
		path := filepath.Join(outputDir, entry.Name())
		if seen[path] || !strings.Contains(strings.ToLower(entry.Name()), version) {
			continue
		}
		seen[path] = true
		candidates = append(candidates, path)
	}
	return candidates, nil
}

func candidateMirrors(raw string, configured []string) []string {
	u, err := url.Parse(raw)
	if err != nil {
		return []string{raw}
	}
	seen := map[string]bool{}
	result := make([]string, 0, len(configured)+4)
	for _, base := range configured {
		candidate, err := configuredMirrorURL(base, u)
		if err != nil || seen[candidate] {
			continue
		}
		seen[candidate] = true
		result = append(result, candidate)
	}
	hosts := []string{u.Host, "bkt-sgp-miui-ota-update-alisgp.oss-ap-southeast-1.aliyuncs.com", "bigota.d.miui.com", "hugeota.d.miui.com"}
	for _, host := range hosts {
		copyURL := *u
		copyURL.Host = host
		candidate := copyURL.String()
		if seen[candidate] || ValidateOfficialURL(candidate) != nil {
			continue
		}
		seen[candidate] = true
		result = append(result, candidate)
	}
	return result
}

func configuredMirrorURL(base string, original *url.URL) (string, error) {
	baseURL, err := url.Parse(strings.TrimSpace(base))
	if err != nil || ValidateMirrorBaseURL(baseURL.String()) != nil {
		return "", i18n.NewError("error.mirror_https")
	}
	baseURL.Path = strings.TrimRight(baseURL.Path, "/") + "/" + strings.TrimLeft(original.Path, "/")
	baseURL.RawPath = ""
	baseURL.RawQuery = original.RawQuery
	baseURL.Fragment = ""
	return baseURL.String(), nil
}

func ValidateMirrorBaseURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return i18n.NewError("error.mirror_base_url")
	}
	return nil
}

func probeMirrors(ctx context.Context, client *http.Client, raw string, configured []string, out io.Writer) ([]mirrorProbe, error) {
	logger := logging.WithWriter(out, logging.Language(ctx))
	type result struct {
		probe mirrorProbe
		err   error
	}
	candidates := candidateMirrors(raw, configured)
	results := make(chan result, len(candidates))
	for _, candidate := range candidates {
		go func(candidate string) {
			probeCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()
			start := time.Now()
			req, err := rangeRequest(probeCtx, candidate, 0, probeSize-1)
			if err != nil {
				results <- result{err: err}
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				results <- result{err: err}
				return
			}
			defer resp.Body.Close()
			total, err := rangeResponseSize(resp)
			if err == nil {
				toRead := probeSize
				if total < toRead {
					toRead = total
				}
				_, err = io.CopyN(io.Discard, resp.Body, toRead)
			}
			results <- result{probe: mirrorProbe{URL: candidate, Size: total, Duration: time.Since(start)}, err: err}
		}(candidate)
	}
	var probes []mirrorProbe
	for range candidates {
		result := <-results
		if result.err == nil && result.probe.Size > 0 {
			probes = append(probes, result.probe)
		} else if result.err != nil {
			slog.Debug("rom.mirror_probe_failed", "error", result.err)
		}
	}
	if len(probes) == 0 {
		return nil, i18n.NewError("error.mirrors_unavailable")
	}
	sort.Slice(probes, func(i, j int) bool { return probes[i].Duration < probes[j].Duration })
	expected := probes[0].Size
	filtered := probes[:0]
	for _, probe := range probes {
		if probe.Size == expected {
			filtered = append(filtered, probe)
			logger.Info("rom.mirror_benchmark", "host", hostOf(probe.URL), "mib_per_second", float64(probeSize)/(1<<20)/probe.Duration.Seconds())
		}
	}
	if len(filtered) > 0 {
		slog.Info("rom.mirror_benchmark_completed", "selected", hostOf(filtered[0].URL), "available", len(filtered), "size", filtered[0].Size)
	}
	return filtered, nil
}

func downloadRanges(ctx context.Context, client *http.Client, release Release, destination string, mirrors []mirrorProbe, connections int, out io.Writer) (string, error) {
	logger := logging.WithWriter(out, logging.Language(ctx))
	size := mirrors[0].Size
	partPath := destination + ".part"
	statePath := partPath + ".json"
	chunks := int((size + defaultChunkSize - 1) / defaultChunkSize)
	state := downloadState{Version: 1, Size: size, ChunkSize: defaultChunkSize, MD5: release.MD5, Done: make([]bool, chunks)}
	partInfo, partErr := os.Stat(partPath)
	if saved, err := loadDownloadState(statePath); err == nil && partErr == nil && partInfo.Size() == size && saved.Size == size && saved.ChunkSize == defaultChunkSize && saved.MD5 == release.MD5 && len(saved.Done) == chunks {
		state = saved
		logger.Info("rom.download_resumed")
	}
	file, err := os.OpenFile(partPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := file.Truncate(size); err != nil {
		return "", err
	}
	if err := saveDownloadState(statePath, state); err != nil {
		return "", err
	}

	var completed atomic.Int64
	for i, done := range state.Done {
		if done {
			completed.Add(chunkLength(int64(i), size, defaultChunkSize))
		}
	}
	jobs := make(chan int)
	errCh := make(chan error, connections)
	var stateMu sync.Mutex
	var workers sync.WaitGroup
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	for worker := 0; worker < connections; worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for index := range jobs {
				if err := downloadChunk(workerCtx, client, file, index, size, mirrors); err != nil {
					select {
					case errCh <- err:
					default:
					}
					cancel()
					return
				}
				stateMu.Lock()
				state.Done[index] = true
				err := saveDownloadState(statePath, state)
				stateMu.Unlock()
				if err != nil {
					select {
					case errCh <- err:
					default:
					}
					cancel()
					return
				}
				completed.Add(chunkLength(int64(index), size, defaultChunkSize))
			}
		}()
	}
	progressDone := make(chan struct{})
	go reportProgress(workerCtx, &completed, size, out, progressDone)
sendJobs:
	for index, done := range state.Done {
		if done {
			continue
		}
		select {
		case jobs <- index:
		case <-workerCtx.Done():
			break sendJobs
		}
	}
	close(jobs)
	workers.Wait()
	cancel()
	<-progressDone
	select {
	case err := <-errCh:
		return "", err
	default:
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := file.Sync(); err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	logger.Info("rom.download_verifying")
	actual, err := fileMD5(partPath)
	if err != nil {
		return "", err
	}
	if actual != release.MD5 {
		return "", i18n.NewError("error.rom_md5_parts", actual, release.MD5)
	}
	if err := os.Rename(partPath, destination); err != nil {
		return "", err
	}
	_ = os.Remove(statePath)
	return destination, nil
}

func downloadChunk(ctx context.Context, client *http.Client, file *os.File, index int, size int64, mirrors []mirrorProbe) error {
	start := int64(index) * defaultChunkSize
	end := start + defaultChunkSize - 1
	if end >= size {
		end = size - 1
	}
	want := end - start + 1
	var lastErr error
	for attempt := 0; attempt < len(mirrors)*2; attempt++ {
		mirror := mirrors[(index+attempt)%len(mirrors)]
		req, err := rangeRequest(ctx, mirror.URL, start, end)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err == nil {
			if err = validateRangeResponse(resp, start, end, size); err == nil {
				writer := io.NewOffsetWriter(file, start)
				var n int64
				n, err = io.CopyN(writer, resp.Body, want)
				if err == nil && n != want {
					err = io.ErrUnexpectedEOF
				}
			}
			resp.Body.Close()
		}
		if err == nil {
			return nil
		}
		lastErr = err
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 300 * time.Millisecond):
		}
	}
	return i18n.WrapError("error.chunk_download", lastErr, index)
}

func rangeRequest(ctx context.Context, raw string, start, end int64) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	req.Header.Set("Referer", "https://www.mi.com/")
	req.Header.Set("User-Agent", "KSUForge/1.0")
	return req, nil
}

func rangeResponseSize(resp *http.Response) (int64, error) {
	if resp.StatusCode != http.StatusPartialContent {
		return 0, fmt.Errorf("HTTP %s", resp.Status)
	}
	value := resp.Header.Get("Content-Range")
	slash := strings.LastIndexByte(value, '/')
	if slash < 0 {
		return 0, i18n.NewError("error.content_range_missing")
	}
	size, err := strconv.ParseInt(value[slash+1:], 10, 64)
	if err != nil || size <= 0 {
		return 0, i18n.NewError("error.content_range_missing")
	}
	return size, nil
}

func validateRangeResponse(resp *http.Response, start, end, size int64) error {
	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("HTTP %s", resp.Status)
	}
	var gotStart, gotEnd, gotSize int64
	if _, err := fmt.Sscanf(resp.Header.Get("Content-Range"), "bytes %d-%d/%d", &gotStart, &gotEnd, &gotSize); err != nil {
		return i18n.NewError("error.content_range_missing")
	}
	if gotStart != start || gotEnd != end || gotSize != size {
		return i18n.NewError("error.content_range_mismatch", resp.Header.Get("Content-Range"))
	}
	return nil
}

func loadDownloadState(path string) (downloadState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return downloadState{}, err
	}
	var state downloadState
	err = json.Unmarshal(data, &state)
	return state, err
}

func saveDownloadState(path string, state downloadState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func reportProgress(ctx context.Context, completed *atomic.Int64, total int64, out io.Writer, done chan<- struct{}) {
	defer close(done)
	logger := logging.WithWriter(out, logging.Language(ctx))
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	started := time.Now()
	base := completed.Load()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := completed.Load()
			speed := float64(current-base) / time.Since(started).Seconds()
			logger.Info("rom.download_progress",
				"downloaded_mib", float64(current)/(1<<20),
				"total_mib", float64(total)/(1<<20),
				"percent", float64(current)*100/float64(total),
				"mib_per_second", speed/(1<<20),
			)
		}
	}
}

func chunkLength(index, size, chunkSize int64) int64 {
	start := index * chunkSize
	end := start + chunkSize
	if end > size {
		end = size
	}
	return end - start
}

func hostOf(raw string) string {
	u, _ := url.Parse(raw)
	return u.Hostname()
}

func fileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// VerifyFileMD5 verifies an existing local ROM against a catalog checksum.
func VerifyFileMD5(path, expected string) error {
	expected = strings.ToLower(strings.TrimSpace(expected))
	if len(expected) != md5.Size*2 {
		return i18n.NewError("error.md5_invalid")
	}
	if _, err := hex.DecodeString(expected); err != nil {
		return i18n.WrapError("error.md5_invalid", err)
	}
	actual, err := fileMD5(path)
	if err != nil {
		return err
	}
	if actual != expected {
		return i18n.NewError("error.rom_md5", actual, expected)
	}
	return nil
}
