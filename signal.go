package easyshutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var shutdownSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

var initOnce = make(chan struct{})

// signals subscribes to shutdownSignals and
// returns a closed stop channel when any of those signals are received.
func signals() <-chan struct{} {
	close(initOnce)

	stop := make(chan struct{})
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, shutdownSignals...)

	go func() {
		<-signalChan
		close(stop)
		<-signalChan
		os.Exit(1) // Two exit signals in a row? Exit directly.
	}()

	return stop
}
