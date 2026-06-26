package plog

import "github.com/sirupsen/logrus"

func SetJsonFormatter() {
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "content",
		},
		DisableTimestamp: false,
	})
}
