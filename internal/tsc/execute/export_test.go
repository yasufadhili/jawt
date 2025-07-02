package execute

import (
	"github.com/yasufadhili/jawt/internal/tsc/compiler"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
)

func CommandLineTest(sys System, cb cbType, commandLineArgs []string) (*tsoptions.ParsedCommandLine, ExitStatus) {
	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	e, _ := executeCommandLineWorker(sys, cb, parsedCommandLine)
	return parsedCommandLine, e
}

func CommandLineTestWatch(sys System, cb cbType, commandLineArgs []string) (*tsoptions.ParsedCommandLine, *watcher) {
	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	_, w := executeCommandLineWorker(sys, cb, parsedCommandLine)
	return parsedCommandLine, w
}

func StartForTest(w *watcher) {
	// this function should perform any initializations before w.doCycle() in `start(watcher)`
	w.initialize()
}

func RunWatchCycle(w *watcher) {
	// this function should perform the same stuff as w.doCycle() without printing time-related output
	if w.hasErrorsInTsConfig() {
		// these are unrecoverable errors--report them and do not build
		return
	}
	// todo: updateProgram()
	w.program = compiler.NewProgram(compiler.ProgramOptions{
		Config: w.options,
		Host:   w.host,
	})
	if w.hasBeenModified(w.program) {
		w.compileAndEmit()
	}
}
