package subcribers

import (
	"command-dispatcher/internal/subcribers/device"
)

func Init() {
	device.Register()
}
