package main

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/locale"
	questProcessor "github.com/rainu/wow-quest-client/internal/quest/processor"
	"github.com/rainu/wow-quest-client/internal/quest/store/sqlite"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	store, err := sqlite.New("/tmp/database.sql")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	p := questProcessor.New(store, locale.English, locale.German)
	pDone := make(chan bool, 1)
	pWg := sync.WaitGroup{}

	appCtx, ctxCancelFn := context.WithCancel(context.Background())

	// Catch signals to enable graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	pWg.Add(1)
	go func() {
		defer pWg.Done()

		p.Run(appCtx, 500)
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
