package main

import (
	"context"
	"github.com/rainu/wow-quest-reader/internal/companion/config"
	"github.com/rainu/wow-quest-reader/internal/companion/processor"
	"github.com/rainu/wow-quest-reader/internal/companion/store"
	"github.com/rainu/wow-quest-reader/internal/companion/system"
	"github.com/rainu/wow-quest-reader/internal/speech/sound/aws"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var ApplicationVersion = "dev"
var ApplicationCodeRev = "revision"

func main() {
	cfg := config.Read()

	println("+------------------------------------+")
	println("| Rainu's WoW Quest Reader Companion | ")
	println("+------------------------------------+")
	println(" Version: " + ApplicationVersion + "(" + ApplicationCodeRev + ")")
	println()

	logrus.SetLevel(cfg.LogLevel)

	speechStore, err := store.NewSpeech(cfg.Sound.Directory)
	if err != nil {
		panic(err)
	}

	speechPool := aws.NewPool(cfg.Sound.AmazonWebService.Region, cfg.Sound.AmazonWebService.Key, cfg.Sound.AmazonWebService.Secret)
	p, err := processor.New(speechPool, system.NewSpeaker(), speechStore, processor.KeyConfiguration{
		HotKeyReading:    cfg.Key.Read,
		AddonKeyPressing: cfg.Key.Addon,
	})
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
