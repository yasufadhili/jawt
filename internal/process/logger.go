package process

import (
	"bufio"
	"github.com/yasufadhili/jawt/internal/core"
	"io"
)

// ProcessLogger pipes the output of a command to a logger with a prefix.
func ProcessLogger(reader io.Reader, logger core.Logger, prefix string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logger.Info(scanner.Text(), core.StringField("process", prefix))
	}
}
