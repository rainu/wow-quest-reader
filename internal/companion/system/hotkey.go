package system

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
	"sync"
)

type Hotkey struct {
	Key      hotkey.Key
	Modifier []hotkey.Modifier
}

func ListenForKeys(ctx context.Context, evtChan chan Hotkey, hotkeys ...Hotkey) {
	defer close(evtChan)

	mainthread.Init(func() {
		wg := &sync.WaitGroup{}

		for _, hkm := range hotkeys {
			hk := hotkey.New(hkm.Modifier, hkm.Key)
			if err := hk.Register(); err != nil {
				logrus.WithField("key", hk.String()).WithError(err).Error("Unable to register key event listener!")
				return
			}

			wg.Add(1)
			go func(hkm Hotkey) {
				defer wg.Done()

				listenForKey(ctx, evtChan, hkm, hk)
			}(hkm)
		}

		wg.Wait()
	})
}

func listenForKey(ctx context.Context, evtChan chan Hotkey, keyModel Hotkey, hk *hotkey.Hotkey) {
	log := logrus.WithField("key", hk.String())

	log.Info("Start listening for key events.")
	for {
		select {
		case <-ctx.Done():
			log.Info("Stop listening for key events.")
			return
		case <-hk.Keydown():
			log.Debug("Received key down event.")
			evtChan <- keyModel
		}
	}
}
