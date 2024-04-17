package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AKYC-chat/akyc-chatting/runner"
)

func handler(signal os.Signal) {
	switch signal {
	case syscall.SIGTERM, syscall.SIGINT:
		runner.SessionStorage.CloseCurrSessions()
		defer os.Exit(0)
	}
}

func main() {
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	exitchnl := make(chan int)

	go func() {
		for {
			s := <-sigchnl
			handler(s)
		}
	}()

	runner.Run()

	exitcode := <-exitchnl
	os.Exit(exitcode)
}
