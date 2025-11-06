package log

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func Init() {
	// Make sure log folder exists
	path := os.ExpandEnv("$HOME/logs/go.log")
	os.MkdirAll(filepath.Dir(path), 0755)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		PadLevelText:     true,
	})

	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(7*24*time.Hour),
	)
	if err != nil {
		log.Fatalf("Failed to create log writer: %v", err)
	}

	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  writer,
			log.ErrorLevel: writer,
			log.WarnLevel:  writer,
			log.DebugLevel: writer,
		},
		&log.JSONFormatter{},
	))

	// Print logs to stdout too
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
