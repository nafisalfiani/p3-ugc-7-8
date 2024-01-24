package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func InitLogger(cfg *Value) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.ReportCaller = true
	logger.Formatter = &logrus.JSONFormatter{
		PrettyPrint:      true,
		CallerPrettyfier: caller(),
		TimestampFormat:  time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFile: "caller",
		},
	}

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, err
	}
	logger.Level = level

	return logger, nil
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()

		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}
