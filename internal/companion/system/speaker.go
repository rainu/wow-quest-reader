package system

import (
	"context"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"sync"
	"time"
)

type player struct {
	initialised bool
	mutex       sync.Mutex
}

func Speaker() *player {
	return &player{}
}

func (p *player) Play(ctx context.Context, mp3Stream io.ReadCloser) error {
	streamer, format, err := mp3.Decode(mp3Stream)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = p.initialise(format)
	if err != nil {
		return err
	}

	done := make(chan bool, 1)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	defer speaker.Clear()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *player) initialise(format beep.Format) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.initialised {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(100*time.Microsecond))
		if err != nil {
			return err
		}
		p.initialised = true
	}

	return nil
}
