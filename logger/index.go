package logger

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Fields = logrus.Fields

func init() {
	Log.SetFormatter(&logrus.JSONFormatter{})
}
