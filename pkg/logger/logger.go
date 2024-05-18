package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func InitLogger() {
	// Set the output to stdout
	log.Out = os.Stdout

	// Set the log level
	log.SetLevel(logrus.InfoLevel)

	// Set the log format to a simple text format
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

func GetLogger() *logrus.Logger {
	return log
}
