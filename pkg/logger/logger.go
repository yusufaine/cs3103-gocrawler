package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

// Simple setup over charmbacelet/log
func Setup(verbose bool) {
	loggerOpts := log.Options{
		TimeFormat: time.Kitchen,
		Level:      log.InfoLevel,
	}

	if verbose {
		loggerOpts.Level = log.DebugLevel
		loggerOpts.ReportCaller = true
	}

	log.SetDefault(log.NewWithOptions(os.Stderr, loggerOpts))
}
