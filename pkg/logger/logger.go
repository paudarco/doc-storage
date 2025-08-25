package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LogrusWriter struct {
	Logger *logrus.Logger
}

func (w *LogrusWriter) Write(p []byte) (n int, err error) {
	w.Logger.Error(string(p))
	return len(p), nil
}

func InitLogger(env string) *logrus.Logger {
	log := logrus.New()

	if env == "dev" {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetFormatter(&logrus.JSONFormatter{})

		log.SetLevel(logrus.InfoLevel)
	}

	log.SetOutput(os.Stdout)

	return log
}
