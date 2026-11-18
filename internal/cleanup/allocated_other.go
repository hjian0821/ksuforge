//go:build !darwin && !linux

package cleanup

import "io/fs"

func allocatedBytes(info fs.FileInfo) int64 { return info.Size() }
