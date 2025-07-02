package module

import (
	"sync"

	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/packagejson"
)

type ModeAwareCache[T any] map[ModeAwareCacheKey]T

type caches struct {
	packageJsonInfoCache *packagejson.InfoCache

	// Cached representation for `core.CompilerOptions.paths`.
	// Doesn't handle other path patterns like in `typesVersions`.
	parsedPatternsForPathsOnce sync.Once
	parsedPatternsForPaths     *ParsedPatterns
}

func newCaches(
	currentDirectory string,
	useCaseSensitiveFileNames bool,
	options *core.CompilerOptions,
) caches {
	return caches{
		packageJsonInfoCache: packagejson.NewInfoCache(currentDirectory, useCaseSensitiveFileNames),
	}
}
