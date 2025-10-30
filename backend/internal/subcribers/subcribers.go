package subcribers

import (
	"command-dispatcher/internal/subcribers/device"

	"github.com/sirupsen/logrus"
)

func Init() {
	logrus.Info("hello")
	device.Register()
}
