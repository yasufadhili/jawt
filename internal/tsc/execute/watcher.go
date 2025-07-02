package execute

import (
	"reflect"
	"time"

	"github.com/yasufadhili/jawt/internal/tsc/collections"
	"github.com/yasufadhili/jawt/internal/tsc/compiler"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

type watcher struct {
	sys              System
	configFileName   string
	options          *tsoptions.ParsedCommandLine
	reportDiagnostic diagnosticReporter

	host           compiler.CompilerHost
	program        *compiler.Program
	prevModified   map[string]time.Time
	configModified bool
}

func createWatcher(sys System, configParseResult *tsoptions.ParsedCommandLine, reportDiagnostic diagnosticReporter) *watcher {
	w := &watcher{
		sys:              sys,
		options:          configParseResult,
		reportDiagnostic: reportDiagnostic,
		// reportWatchStatus: createWatchStatusReporter(sys, configParseResult.CompilerOptions().Pretty),
	}
	if configParseResult.ConfigFile != nil {
		w.configFileName = configParseResult.ConfigFile.SourceFile.FileName()
	}
	return w
}

func (w *watcher) compileAndEmit() {
	// !!! output/error reporting is currently the same as non-watch mode
	// diagnostics, emitResult, exitStatus :=
	emitFilesAndReportErrors(w.sys, w.program, w.reportDiagnostic)
}

func (w *watcher) hasErrorsInTsConfig() bool {
	// only need to check and reparse tsconfig options/update host if we are watching a config file
	if w.configFileName != "" {
		extendedConfigCache := collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]{}
		// !!! need to check that this merges compileroptions correctly. This differs from non-watch, since we allow overriding of previous options
		configParseResult, errors := tsoptions.GetParsedCommandLineOfConfigFile(w.configFileName, &core.CompilerOptions{}, w.sys, &extendedConfigCache)
		if len(errors) > 0 {
			for _, e := range errors {
				w.reportDiagnostic(e)
			}
			return true
		}
		// CompilerOptions contain fields which should not be compared; clone to get a copy without those set.
		if !reflect.DeepEqual(w.options.CompilerOptions().Clone(), configParseResult.CompilerOptions().Clone()) {
			// fmt.Fprint(w.sys.Writer(), "build triggered due to config change", w.sys.NewLine())
			w.configModified = true
		}
		w.options = configParseResult
		w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath(), &extendedConfigCache)
	}
	return false
}

func (w *watcher) hasBeenModified(program *compiler.Program) bool {
	// checks watcher's snapshot against program file modified times
	currState := map[string]time.Time{}
	filesModified := w.configModified
	for _, sourceFile := range program.SourceFiles() {
		fileName := sourceFile.FileName()
		s := w.sys.FS().Stat(fileName)
		if s == nil {
			// do nothing; if file is in program.SourceFiles() but is not found when calling Stat, file has been very recently deleted.
			// deleted files are handled outside of this loop
			continue
		}
		currState[fileName] = s.ModTime()
		if !filesModified {
			if currState[fileName] != w.prevModified[fileName] {
				// fmt.Fprint(w.sys.Writer(), "build triggered from ", fileName, ": ", w.prevModified[fileName], " -> ", currState[fileName], w.sys.NewLine())
				filesModified = true
			}
			// catch cases where no files are modified, but some were deleted
			delete(w.prevModified, fileName)
		}
	}
	if !filesModified && len(w.prevModified) > 0 {
		// fmt.Fprint(w.sys.Writer(), "build triggered due to deleted file", w.sys.NewLine())
		filesModified = true
	}
	w.prevModified = currState

	// reset state for next cycle
	w.configModified = false
	return filesModified
}
