package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

// Simple setup over charmbacelet/log
func Setup(verbose bool, logRelFile string) {
	wrs := []io.Writer{os.Stderr}
	loggerOpts := log.Options{
		TimeFormat: time.Kitchen,
		Level:      log.InfoLevel,
	}

	if verbose {
		loggerOpts.Level = log.DebugLevel
		loggerOpts.ReportCaller = true
	}

	if logRelFile != "" {
		if err := os.MkdirAll(filepath.Dir(logRelFile), 0755); err != nil {
			log.Fatal("unable to create log file", "error", err)
		}
		f, err := os.OpenFile(logRelFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("unable to open log file", "error", err)
		}
		wrs = append(wrs, f)
	}
	ml := log.NewWithOptions(io.MultiWriter(wrs...), loggerOpts)
	if len(wrs) == 1 {
		ml.SetColorProfile(termenv.TrueColor)
	} else {
		ml.SetColorProfile(termenv.Ascii)
	}
	log.SetDefault(ml)
}
