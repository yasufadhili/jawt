package tsc

import (
	"errors"
	"flag"
	"fmt"
	"github.com/yasufadhili/jawt/internal/tsc/api"
	"github.com/yasufadhili/jawt/internal/tsc/bundled"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"io"
	"os"
)

func runAPI(args []string) int {
	flag := flag.NewFlagSet("api", flag.ContinueOnError)
	cwd := flag.String("cwd", core.Must(os.Getwd()), "current working directory")
	if err := flag.Parse(args); err != nil {
		return 2
	}

	defaultLibraryPath := bundled.LibPath()

	s := api.NewServer(&api.ServerOptions{
		In:                 os.Stdin,
		Out:                os.Stdout,
		Err:                os.Stderr,
		Cwd:                *cwd,
		NewLine:            "\n",
		DefaultLibraryPath: defaultLibraryPath,
	})

	if err := s.Run(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Println(err)
		return 1
	}
	return 0
}
