package utility

import (
	"github.com/sirupsen/logrus"
	"github.com/yasser-sobhy/sparrow/env"
)

// InitLogger intializers logrus logger
func init() {
	logrus.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	if env.IsDevelopment() == true {
		// Only log the warning severity or above.
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}
