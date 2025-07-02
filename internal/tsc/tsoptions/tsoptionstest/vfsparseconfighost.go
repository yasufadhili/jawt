package tsoptionstest

import (
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
	"github.com/yasufadhili/jawt/internal/tsc/vfs"
	"github.com/yasufadhili/jawt/internal/tsc/vfs/vfstest"
)

func fixRoot(path string) string {
	rootLength := tspath.GetRootLength(path)
	if rootLength == 0 {
		return path
	}
	if len(path) == rootLength {
		return "."
	}
	return path[rootLength:]
}

type VfsParseConfigHost struct {
	Vfs              vfs.FS
	CurrentDirectory string
}

var _ tsoptions.ParseConfigHost = (*VfsParseConfigHost)(nil)

func (h *VfsParseConfigHost) FS() vfs.FS {
	return h.Vfs
}

func (h *VfsParseConfigHost) GetCurrentDirectory() string {
	return h.CurrentDirectory
}

func NewVFSParseConfigHost(files map[string]string, currentDirectory string, useCaseSensitiveFileNames bool) *VfsParseConfigHost {
	return &VfsParseConfigHost{
		Vfs:              vfstest.FromMap(files, useCaseSensitiveFileNames),
		CurrentDirectory: currentDirectory,
	}
}
