package signal

import (
	"os"
)

func RouterSignalNotify(c chan<- os.Signal) {
	osRouterSignalNotify(c)
}

func GetRestartSignal(backendSystem string) (os.Signal, error) {
	return osGetRestartSignal(backendSystem)
}
