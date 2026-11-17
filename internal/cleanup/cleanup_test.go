package cleanup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const version = "OS4.0.4.0.XOCCNXM"
const md5 = "0123456789abcdef0123456789abcdef"

func put(t *testing.T, root, path string, data []byte) string {
	t.Helper()
	file := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return file
}
func stock(t *testing.T, root string) {
	t.Helper()
	data := make([]byte, 4096)
	copy(data, "ANDROID!")
	sum := sha256.Sum256(data)
	meta, _ := json.Marshal(map[string]string{"version": version, "md5": md5, "sha256": hex.EncodeToString(sum[:])})
	put(t, root, "stock/dada/"+version+"/init_boot.img", data)
	put(t, root, "stock/dada/"+version+"/init_boot.img.json", meta)
}
func item(t *testing.T, p *Plan, kind string) Item {
	t.Helper()
	for _, i := range p.Report.Items {
		if i.Kind == kind {
			return i
		}
	}
	t.Fatalf("missing %s", kind)
	return Item{}
}
func exists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestScanAndDefaultCleanupPreserveImagesAndExternalFiles(t *testing.T) {
	root := t.TempDir()
	stock(t, root)
	zip := put(t, root, "rom/dada/"+version+"/dada-ota_full-"+version+"-0123456789.zip", []byte("downloaded ROM"))
	part := put(t, root, "rom/dada/"+version+"/next.zip.part", []byte("partial"))
	put(t, root, "rom/dada/"+version+"/next.zip.part.json", []byte("{}"))
	patched := put(t, root, "patched/dada/"+version+"/kernelsu-patched-init_boot.img", []byte("patched"))
	put(t, root, "patched/dada/"+version+"/kernelsu-patched-init_boot.img.json", []byte("{}"))
	other := put(t, root, "rom/dada/"+version+"/notes.txt", []byte("keep"))
	unowned := put(t, root, "rom/dada/"+version+"/dada-ota_full-"+version+"-0123456789.zip.json", []byte("unrecognized"))
	external := put(t, root, "my-rom.zip", []byte("external"))
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if len(p.Report.Items) != 4 {
		t.Fatalf("items=%+v", p.Report.Items)
	}
	rom := item(t, p, "rom")
	if !rom.Recommended {
		t.Fatal("verified ROM should be recommended")
	}
	for _, i := range p.Report.Items {
		if i.Kind != "rom" && i.Recommended {
			t.Fatalf("unexpected recommendation: %+v", i)
		}
	}
	if len(item(t, p, "partial").Files) != 2 || len(item(t, p, "stock").Files) != 2 {
		t.Fatal("sidecars must be paired")
	}
	result, err := p.Clean([]string{rom.ID, rom.ID})
	if err != nil {
		t.Fatal(err)
	}
	if result.DeletedFiles != 1 || len(result.Failed) != 0 {
		t.Fatalf("result=%+v", result)
	}
	if _, err := os.Stat(zip); !os.IsNotExist(err) {
		t.Fatal("ROM not removed")
	}
	for _, path := range []string{part, patched, other, unowned, external, filepath.Join(root, "stock/dada", version, "init_boot.img")} {
		exists(t, path)
	}
}

func TestDefaultDoesNotSelectROMWithoutVerifiedBackup(t *testing.T) {
	root := t.TempDir()
	stock(t, root)
	put(t, root, "rom/dada/"+version+"/ota-0123456789.zip", []byte("ROM"))
	f, err := os.OpenFile(filepath.Join(root, "stock/dada", version, "init_boot.img"), os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.WriteAt([]byte{1}, 100); err != nil {
		t.Fatal(err)
	}
	f.Close()
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if item(t, p, "rom").Recommended {
		t.Fatal("tampered stock backup must not authorize recommendation")
	}
}

func TestExplicitImageSelectionDeletesItsSidecar(t *testing.T) {
	root := t.TempDir()
	stock(t, root)
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	i := item(t, p, "stock")
	result, err := p.Clean([]string{i.ID})
	if err != nil {
		t.Fatal(err)
	}
	if result.DeletedFiles != 2 {
		t.Fatalf("result=%+v", result)
	}
	for _, f := range i.Files {
		if _, err := os.Stat(filepath.Join(root, f.Path)); !os.IsNotExist(err) {
			t.Fatalf("not removed: %s", f.Path)
		}
	}
}

func TestChangedOrUnknownSelectionDeletesNothing(t *testing.T) {
	root := t.TempDir()
	first := put(t, root, "rom/dada/"+version+"/a.zip", []byte("a"))
	second := put(t, root, "rom/dada/"+version+"/b.zip", []byte("b"))
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if _, err := p.Clean([]string{"../../outside"}); err == nil {
		t.Fatal("unknown ID accepted")
	}
	if err := os.WriteFile(second, []byte("changed"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Clean([]string{p.Report.Items[0].ID, p.Report.Items[1].ID}); err == nil {
		t.Fatal("changed file accepted")
	}
	exists(t, first)
	exists(t, second)
}

func TestSymlinksAreNeverScannedOrFollowedOnDelete(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	target := put(t, outside, "foreign.zip", []byte("private"))
	if err := os.Symlink(outside, filepath.Join(root, "rom")); err != nil {
		t.Skip(err)
	}
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Report.Items) != 0 {
		t.Fatal("followed linked category")
	}
	p.Close()
	os.Remove(filepath.Join(root, "rom"))
	path := put(t, root, "rom/dada/"+version+"/owned.zip", []byte("owned"))
	p, err = Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	id := item(t, p, "rom").ID
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Clean([]string{id}); err == nil {
		t.Fatal("replacement link accepted")
	}
	exists(t, target)
}

func TestReplacedDirectoryInvalidatesPinnedPlan(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "output")
	path := put(t, root, "rom/dada/"+version+"/owned.zip", []byte("owned"))
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	moved := filepath.Join(parent, "moved")
	if err := os.Rename(root, moved); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Clean([]string{item(t, p, "rom").ID}); err == nil {
		t.Fatal("directory replacement accepted")
	}
	exists(t, filepath.Join(moved, stringsRelative(t, root, path)))
}
func stringsRelative(t *testing.T, root, path string) string {
	t.Helper()
	r, e := filepath.Rel(root, path)
	if e != nil {
		t.Fatal(e)
	}
	return r
}

func TestMissingOutputDirectoryIsEmptyAndNotCreated(t *testing.T) {
	root := filepath.Join(t.TempDir(), "absent")
	p, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if len(p.Report.Items) != 0 {
		t.Fatal("nonempty")
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatal("scan created directory")
	}
}
