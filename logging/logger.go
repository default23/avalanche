package logging

import "github.com/sirupsen/logrus"

func New() *logrus.Entry {
	f := new(logrus.TextFormatter)
	f.FullTimestamp = true
	f.TimestampFormat = "02.01.2006 15:04:05"

	logger := logrus.New()
	logger.SetFormatter(f)

	return logrus.NewEntry(logger)
}
