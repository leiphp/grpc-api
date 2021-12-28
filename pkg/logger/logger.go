package logger

import (
	"github.com/sirupsen/logrus"
)

var me *logrus.Entry

func init() {
	logEntry := logrus.New()
	//logEntry.Level = logrus.WarnLevel
	logEntry.Level = logrus.InfoLevel
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "20060102T15:04:05.999"
	Formatter.FullTimestamp = true
	logEntry.SetFormatter(Formatter)
	me = logrus.NewEntry(logEntry).WithField("service", "grpc-gateway")
}

// Pre is used to get a prepared logrus.Entry
func Pre() *logrus.Entry {
	return me
}
