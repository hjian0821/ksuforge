package image

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/hjian0821/ksuforge/internal/i18n"
)

type Info struct {
	Path   string
	Size   int64
	SHA256 string
}

func Inspect(path string) (Info, error) {
	f, err := os.Open(path)
	if err != nil {
		return Info{}, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return Info{}, err
	}
	if !st.Mode().IsRegular() || st.Size() < 4096 {
		return Info{}, i18n.NewError("error.image_invalid")
	}
	magic := make([]byte, 8)
	if _, err := io.ReadFull(f, magic); err != nil {
		return Info{}, err
	}
	// Android boot image v0-v4 and vendor_boot use these magic values.
	if string(magic) != "ANDROID!" && string(magic) != "VNDRBOOT" {
		return Info{}, i18n.NewError("error.image_magic")
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return Info{}, err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return Info{}, err
	}
	return Info{Path: path, Size: st.Size(), SHA256: hex.EncodeToString(h.Sum(nil))}, nil
}
