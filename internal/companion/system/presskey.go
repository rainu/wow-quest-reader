package system

import (
	"github.com/micmonay/keybd_event"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

type keyPresser struct {
	kb *keybd_event.KeyBonding
}

type KeyPressing struct {
	Ctrl  bool
	Shift bool
	Keys  []int
}

func NewKeyPresser() (*keyPresser, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return nil, err
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	return &keyPresser{
		kb: &kb,
	}, nil
}

func (k *keyPresser) PressCopy() {
	k.PressRaw(true, false, keybd_event.VK_C)
}

func (k *keyPresser) Press(kp KeyPressing) {
	k.PressRaw(kp.Ctrl, kp.Shift, kp.Keys...)
}

func (k *keyPresser) PressRaw(ctrl, shift bool, keys ...int) {
	k.kb.SetKeys(keys...)
	k.kb.HasCTRL(ctrl)
	k.kb.HasSHIFT(shift)

	log := logrus.WithField("ctrl", ctrl).WithField("shift", shift).WithField("keys", keys)

	err := k.kb.Launching()
	if err != nil {
		log.WithError(err).Error("Unable to launch key!")
		return
	}

	err = k.kb.Press()
	if err != nil {
		log.WithError(err).Error("Unable to press key!")
		return
	}

	//wait a little
	time.Sleep(50 * time.Millisecond)

	err = k.kb.Release()
	if err != nil {
		log.WithError(err).Error("Unable to release key!")
		return
	}

	log.Info("Pressed key.")
}
