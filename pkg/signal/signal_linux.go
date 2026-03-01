//go:build !windows
// +build !windows

package signal

import (
	"errors"
	"os"
	osSignal "os/signal"

	"golang.org/x/sys/unix"
)

func osRouterSignalNotify(c chan<- os.Signal) {
	osSignal.Notify(c, os.Interrupt, os.Kill, unix.SIGUSR1, unix.SIGUSR2)
}

func osGetRestartSignal(system string) (os.Signal, error) {
	switch system {
	case "router":
		return unix.SIGUSR1, nil
	case "main":
		return unix.SIGUSR2, nil
	}
	return nil, errors.New("wrong system signal input")
}
