package compiler

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/collections"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/parser"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
	"github.com/yasufadhili/jawt/internal/tsc/vfs"
	"github.com/yasufadhili/jawt/internal/tsc/vfs/cachedvfs"
)

type CompilerHost interface {
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string
	NewLine() string
	Trace(msg string)
	GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile
	GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine
}

type FileInfo struct {
	Name string
	Size int64
}

var _ CompilerHost = (*compilerHost)(nil)

type compilerHost struct {
	options             *core.CompilerOptions
	currentDirectory    string
	fs                  vfs.FS
	defaultLibraryPath  string
	extendedConfigCache *collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]
}

func NewCachedFSCompilerHost(
	options *core.CompilerOptions,
	currentDirectory string,
	fs vfs.FS,
	defaultLibraryPath string,
	extendedConfigCache *collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry],
) CompilerHost {
	return NewCompilerHost(options, currentDirectory, cachedvfs.From(fs), defaultLibraryPath, extendedConfigCache)
}

func NewCompilerHost(
	options *core.CompilerOptions,
	currentDirectory string,
	fs vfs.FS,
	defaultLibraryPath string,
	extendedConfigCache *collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry],
) CompilerHost {
	return &compilerHost{
		options:             options,
		currentDirectory:    currentDirectory,
		fs:                  fs,
		defaultLibraryPath:  defaultLibraryPath,
		extendedConfigCache: extendedConfigCache,
	}
}

func (h *compilerHost) FS() vfs.FS {
	return h.fs
}

func (h *compilerHost) DefaultLibraryPath() string {
	return h.defaultLibraryPath
}

func (h *compilerHost) SetOptions(options *core.CompilerOptions) {
	h.options = options
}

func (h *compilerHost) GetCurrentDirectory() string {
	return h.currentDirectory
}

func (h *compilerHost) NewLine() string {
	if h.options == nil {
		return "\n"
	}
	return h.options.NewLine.GetNewLineCharacter()
}

func (h *compilerHost) Trace(msg string) {
	//!!! TODO: implement
}

func (h *compilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	text, ok := h.FS().ReadFile(opts.FileName)
	if !ok {
		return nil
	}
	return parser.ParseSourceFile(opts, text, core.GetScriptKindFromFileName(opts.FileName))
}

func (h *compilerHost) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	commandLine, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, h, nil)
	return commandLine
}
