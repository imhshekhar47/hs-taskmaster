package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger = logrus.New()

	LogTextFormatter = &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02T15:04:05",
		FullTimestamp:   true,
	}
)

func init() {
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil || level == logrus.PanicLevel {
		level = logrus.TraceLevel
	}

	Logger.SetLevel(level)

	Logger.SetFormatter(LogTextFormatter)
}

func GetLogger(origin string) *logrus.Entry {
	return Logger.WithField("origin", origin)
}
