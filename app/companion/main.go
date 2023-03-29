package main

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/companion/processor"
	"github.com/rainu/wow-quest-client/internal/companion/store"
	"github.com/rainu/wow-quest-client/internal/companion/system"
	"github.com/rainu/wow-quest-client/internal/speech/sound/aws"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	speechStore, err := store.NewSpeech("./wow")
	if err != nil {
		panic(err)
	}

	speechPool := aws.NewPool("<region>", "<key>", "<secret>")
	p, err := processor.New(speechPool, system.Speaker(), speechStore)
	if err != nil {
		panic(err)
	}
	pDone := make(chan bool, 1)
	pWg := sync.WaitGroup{}

	appCtx, ctxCancelFn := context.WithCancel(context.Background())

	// Catch signals to enable graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	pWg.Add(1)
	go func() {
		defer pWg.Done()

		p.Run(appCtx)
		pDone <- true
	}()

	// wait until application is cancelled or finished
	select {
	case <-sigs:
		//use interrupt
		ctxCancelFn()
		logrus.Info("Interrupt signal received. Cancel application...")
	case <-pDone:
		logrus.Info("Application is done...")
	}

	pWg.Wait()
}
