//go:build windows
// +build windows

package signal

import (
	"errors"
	"os"
	osSignal "os/signal"
)

func osRouterSignalNotify(c chan<- os.Signal) {
	osSignal.Notify(c, os.Interrupt, os.Kill)
}

func osGetRestartSignal(backendSystem string) (os.Signal, error) {
	return nil, errors.New("not support restart signal on windows")
}
