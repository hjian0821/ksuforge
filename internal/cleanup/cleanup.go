// Package cleanup inventories and removes only recognized artifacts inside a
// pinned application output directory. Callers select scan IDs, never paths.
package cleanup

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hjian0821/ksuforge/internal/rom"
)

type File struct {
	Path         string `json:"path"`
	Bytes        int64  `json:"bytes"`
	LogicalBytes int64  `json:"logicalBytes"`
	Modified     string `json:"modified"`
}
type Item struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Device      string `json:"device"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Files       []File `json:"files"`
	Bytes       int64  `json:"bytes"`
	Recommended bool   `json:"recommended"`
}
type Report struct {
	Token      string `json:"token"`
	Root       string `json:"root"`
	Items      []Item `json:"items"`
	TotalBytes int64  `json:"totalBytes"`
}
type Failure struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}
type Result struct {
	DeletedFiles int       `json:"deletedFiles"`
	DeletedBytes int64     `json:"deletedBytes"`
	Failed       []Failure `json:"failed"`
	Cancelled    bool      `json:"cancelled"`
}
type savedFile struct {
	path  string
	info  fs.FileInfo
	bytes int64
}
type Plan struct {
	rootInfo fs.FileInfo
	Report   Report
	root     *os.Root
	files    map[string][]savedFile
}

func (p *Plan) Close() {
	if p != nil && p.root != nil {
		p.root.Close()
		p.root = nil
	}
}

func Scan(output string) (*Plan, error) {
	root, err := filepath.Abs(output)
	if err != nil {
		return nil, err
	}
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	p := &Plan{Report: Report{Root: root, Token: hex.EncodeToString(nonce[:]), Items: []Item{}}, files: map[string][]savedFile{}}
	p.root, err = os.OpenRoot(root)
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return nil, err
	}
	p.rootInfo, err = p.root.Stat(".")
	if err != nil {
		p.Close()
		return nil, err
	}
	ok := false
	defer func() {
		if !ok {
			p.Close()
		}
	}()
	entries := map[string]fs.FileInfo{}
	for _, category := range []string{"rom", "stock", "patched"} {
		info, err := p.root.Lstat(category)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		err = fs.WalkDir(p.root.FS(), category, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			parts := strings.Split(filepath.ToSlash(path), "/")
			if d.IsDir() {
				if len(parts) > 3 {
					return fs.SkipDir
				}
				return nil
			}
			if len(parts) != 4 || d.Type()&os.ModeSymlink != 0 {
				return nil
			}
			info, err := p.root.Lstat(path)
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				entries[path] = info
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	// Hash only the small stock images, never multi-gigabyte ROM packages.
	verified := map[string][]string{}
	checked := map[string]bool{}
	for path, info := range entries {
		parts := strings.Split(filepath.ToSlash(path), "/")
		category, device, version, name := parts[0], parts[1], parts[2], parts[3]
		kind := ""
		switch category {
		case "rom":
			switch {
			case strings.HasSuffix(name, ".zip"):
				kind = "rom"
			case strings.HasSuffix(name, ".zip.part"):
				kind = "partial"
			case strings.HasSuffix(name, ".zip.part.json"):
				if _, exists := entries[strings.TrimSuffix(path, ".json")]; !exists {
					kind = "partial"
				}
			}
		case "stock":
			if name == "boot.img" || name == "init_boot.img" {
				kind = "stock"
			}
		case "patched":
			if name == "kernelsu-patched-boot.img" || name == "kernelsu-patched-init_boot.img" {
				kind = "patched"
			}
		}
		if kind == "" {
			continue
		}
		digest := sha256.Sum256([]byte(filepath.ToSlash(path)))
		id := hex.EncodeToString(digest[:16])
		item := Item{ID: id, Kind: kind, Device: device, Version: version, Name: name, Files: []File{}}
		paths := []string{path}
		if _, exists := entries[path+".json"]; exists && (kind == "stock" || kind == "patched" || (kind == "partial" && strings.HasSuffix(path, ".part"))) {
			paths = append(paths, path+".json")
		}
		p.files[id] = nil
		for _, file := range paths {
			st := info
			if file != path {
				st = entries[file]
			}
			size := allocatedBytes(st)
			item.Files = append(item.Files, File{Path: filepath.ToSlash(file), Bytes: size, LogicalBytes: st.Size(), Modified: st.ModTime().Format(time.RFC3339)})
			item.Bytes += size
			p.files[id] = append(p.files[id], savedFile{path: file, info: st, bytes: size})
		}
		if kind == "rom" {
			key := device + "/" + version
			if !checked[key] {
				checked[key] = true
				// Reject symbolic-link directories before validating a cached backup.
				stockRel := filepath.Join("stock", device, version)
				if p.noLinks(stockRel) == nil {
					for _, partition := range []string{"boot", "init_boot"} {
						if _, release, valid := rom.FindVerifiedStock(filepath.Join(root, stockRel), partition, version); valid {
							verified[key] = append(verified[key], release.MD5)
						}
					}
				}
			}
			for _, md5 := range verified[key] {
				if strings.Contains(strings.ToLower(name), md5[:10]) {
					item.Recommended = true
				}
			}
		}
		p.Report.TotalBytes += item.Bytes
		p.Report.Items = append(p.Report.Items, item)
	}
	sort.Slice(p.Report.Items, func(i, j int) bool {
		a, b := p.Report.Items[i], p.Report.Items[j]
		return a.Device+"/"+a.Version+"/"+a.Kind+"/"+a.Name < b.Device+"/"+b.Version+"/"+b.Kind+"/"+b.Name
	})
	ok = true
	return p, nil
}

func (p *Plan) Selection(ids []string) ([]Item, error) {
	result := []Item{}
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if _, ok := p.files[id]; !ok {
			return nil, fmt.Errorf("unknown or expired cleanup item")
		}
		for _, item := range p.Report.Items {
			if item.ID == id {
				result = append(result, item)
				break
			}
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no cleanup items selected")
	}
	return result, nil
}
func (p *Plan) noLinks(path string) error {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i := range parts {
		info, err := p.root.Lstat(filepath.Join(parts[:i+1]...))
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symbolic link changed since scan")
		}
	}
	return nil
}
func (p *Plan) Validate(ids []string) error {
	if p.root == nil {
		return fmt.Errorf("cleanup scan is no longer available")
	}
	currentRoot, err := os.Stat(p.Report.Root)
	if err != nil {
		return err
	}
	if !os.SameFile(p.rootInfo, currentRoot) {
		return fmt.Errorf("output directory changed since scan")
	}
	items, err := p.Selection(ids)
	if err != nil {
		return err
	}
	for _, item := range items {
		for _, f := range p.files[item.ID] {
			if err := p.noLinks(f.path); err != nil {
				return err
			}
			current, err := p.root.Lstat(f.path)
			if err != nil {
				return err
			}
			if !current.Mode().IsRegular() || !os.SameFile(f.info, current) || f.info.Size() != current.Size() || !f.info.ModTime().Equal(current.ModTime()) {
				return fmt.Errorf("file changed since scan: %s", f.path)
			}
		}
	}
	return nil
}
func (p *Plan) Clean(ids []string) (Result, error) {
	result := Result{Failed: []Failure{}}
	if err := p.Validate(ids); err != nil {
		return result, err
	}
	items, _ := p.Selection(ids)
	for _, item := range items {
		for _, f := range p.files[item.ID] {
			if err := p.noLinks(f.path); err != nil {
				result.Failed = append(result.Failed, Failure{Path: filepath.ToSlash(f.path), Error: err.Error()})
				continue
			}
			if err := p.root.Remove(f.path); err != nil {
				result.Failed = append(result.Failed, Failure{Path: filepath.ToSlash(f.path), Error: err.Error()})
				continue
			}
			result.DeletedFiles++
			result.DeletedBytes += f.bytes
		}
	}
	return result, nil
}
