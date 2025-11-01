package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init(level string) error {
	Logger = logrus.New()

	Logger.SetOutput(os.Stdout)

	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.SetLevel(logrus.InfoLevel)
		Logger.Warnf("Invalid log level '%s', useing 'info' level", level)
	} else {
		Logger.SetLevel(logLevel)
	}
	Logger.SetReportCaller(true)

	Logger.Info("Logger initialized succesfully")

	return nil
}

func GetLogger() *logrus.Logger {
	if Logger == nil {
		Logger = logrus.New()
		Logger.SetLevel(logrus.InfoLevel)
	}

	return Logger
}
