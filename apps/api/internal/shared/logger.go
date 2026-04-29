package shared

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(appEnv string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	if appEnv == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
		return logger
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	return logger
}
