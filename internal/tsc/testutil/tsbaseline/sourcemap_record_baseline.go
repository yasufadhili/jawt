package tsbaseline

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/baseline"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/harnessutil"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

func DoSourcemapRecordBaseline(
	t *testing.T,
	baselinePath string,
	header string,
	options *core.CompilerOptions,
	result *harnessutil.CompilationResult,
	harnessSettings *harnessutil.HarnessOptions,
	opts baseline.Options,
) {
	actual := baseline.NoContent
	if options.SourceMap.IsTrue() || options.InlineSourceMap.IsTrue() || options.DeclarationMap.IsTrue() {
		record := removeTestPathPrefixes(result.GetSourceMapRecord(), false /*retainTrailingDirectorySeparator*/)
		if !(options.NoEmitOnError.IsTrue() && len(result.Diagnostics) > 0) && len(record) > 0 {
			actual = record
		}
	}

	if tspath.FileExtensionIsOneOf(baselinePath, []string{tspath.ExtensionTs, tspath.ExtensionTsx}) {
		baselinePath = tspath.ChangeExtension(baselinePath, ".sourcemap.txt")
	}

	baseline.Run(t, baselinePath, actual, opts)
}
