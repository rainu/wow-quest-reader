package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micmonay/keybd_event"
	processorModel "github.com/rainu/wow-quest-reader/internal/companion/model"
	"github.com/rainu/wow-quest-reader/internal/companion/system"
	"github.com/rainu/wow-quest-reader/internal/locale"
	"github.com/rainu/wow-quest-reader/internal/model"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

type event byte

const (
	eventNone = event(iota)
	eventRead = event(iota)
)

type processor struct {
	cba ClipboardAccess

	keyPresser KeyPresser

	hkForReading system.Hotkey

	kpAddon system.KeyPressing

	lastEvent event

	mp3Player      Mp3Player
	speechPool     SpeechPool
	speechCtx      context.Context
	speechCancel   context.CancelFunc
	speechDoneChan chan bool
	speechStore    SoundStore
}

func New(speechPool SpeechPool, mp3Player Mp3Player, store SoundStore, keyConfig KeyConfiguration) (*processor, error) {
	cba, err := system.NewClipboardAccess()
	if err != nil {
		return nil, fmt.Errorf("unable to create clipboard access: %w", err)
	}

	result := &processor{
		cba: cba,

		hkForReading: keyConfig.HotKeyReading,
		kpAddon:      keyConfig.AddonKeyPressing,

		mp3Player:      mp3Player,
		speechPool:     speechPool,
		speechStore:    store,
		speechCtx:      nil,
		speechCancel:   nil,
		speechDoneChan: nil,
	}

	result.keyPresser, err = system.NewKeyPresser()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *processor) Run(ctx context.Context) {
	hkEventChan := make(chan system.Hotkey)
	cbContentChan := make(chan string)

	go p.runHotkeyListener(ctx, hkEventChan)

	go p.runClipboardListener(ctx, cbContentChan)
	go p.cba.Watch(ctx, cbContentChan)

	system.ListenForKeys(ctx, hkEventChan, p.hkForReading)
}

func (p *processor) runHotkeyListener(ctx context.Context, hkEventChan chan system.Hotkey) {
	for {
		select {
		case <-ctx.Done():
			return
		case hk := <-hkEventChan:
			if isEq(hk, p.hkForReading) {
				p.lastEvent = eventRead
			} else {
				logrus.WithField("hotkey", fmt.Sprintf("%#v", hk)).Error("Unknown incoming hotkey!")
				continue
			}

			if p.speechCtx != nil {
				// user wants to interrupt the current speech
				logrus.Info("Interrupt current playing speech.")
				p.stopPlay(ctx)
				continue
			}

			// press the key which open the addon-frame inside WoW
			p.keyPresser.Press(p.kpAddon)

			// the opening frame will already select the whole content inside this frame
			// just wait a bit for opening...
			select {
			case <-ctx.Done():
				return
			case <-time.NewTimer(50 * time.Millisecond).C:
			}

			// press <Ctrl+C> to copy the selected content into clipboard
			p.keyPresser.PressCopy()

			// the clipboard listener will than receive the content and do the corresponding action!

			// close the addon-frame
			p.keyPresser.PressRaw(false, false, keybd_event.VK_0)
		}
	}
}

func isEq(hk1, hk2 system.Hotkey) bool {
	if hk1.Key != hk2.Key {
		return false
	}
	if fmt.Sprintf("%v", hk1.Modifier) != fmt.Sprintf("%v", hk2.Modifier) {
		return false
	}

	return true
}

func (p *processor) runClipboardListener(ctx context.Context, cbContentChan chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case content := <-cbContentChan:
			content = strings.TrimSpace(content)
			content = strings.ReplaceAll(content, "\n", "")
			content = strings.ReplaceAll(content, "\r", "")
			var info processorModel.Info

			logrus.WithField("content", content).Debug("New clipboard content detected!")

			if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
				if err := json.Unmarshal([]byte(content), &info); err != nil {
					continue
				}
			} else {
				continue
			}
			if !info.IsValid() {
				logrus.Debug("Quest is invalid. Skip this one.")
				continue
			}

			// interesting content available in clipboard!
			log := logrus.WithField("quest_id", info.Quest.Id)
			log.Info("Quest information in clipboard detected.")

			switch p.lastEvent {
			case eventRead:
				if info.Gossip != "" {
					p.handleRead(ctx, info, p.speechStore.GetGossip, p.speechStore.GetFileLocationForGossip, gossipId, func(i processorModel.Info) string {
						return i.Gossip
					})
				} else if info.Quest.Description != "" {
					p.handleRead(ctx, info, p.speechStore.GetDescription, p.speechStore.GetFileLocationForDescription, questId, func(i processorModel.Info) string {
						return i.Quest.Description
					})
				} else if info.Quest.Progress != "" {
					p.handleRead(ctx, info, p.speechStore.GetProgress, p.speechStore.GetFileLocationForProgress, questId, func(i processorModel.Info) string {
						return i.Quest.Progress
					})
				} else if info.Quest.Completion != "" {
					p.handleRead(ctx, info, p.speechStore.GetCompletion, p.speechStore.GetFileLocationForCompletion, questId, func(i processorModel.Info) string {
						return i.Quest.Completion
					})
				}
			default:
				log.Warn("Ignore clipboard content because no event fired before.")
			}

			// we handled the current event ... wait for the next one
			p.lastEvent = eventNone
		}
	}
}

func questId(i processorModel.Info) string {
	return fmt.Sprintf("%d", i.Quest.Id)
}

func gossipId(i processorModel.Info) string {
	return i.GossipId()
}

func (p *processor) handleRead(ctx context.Context, info processorModel.Info,
	getFromStore func(string, locale.Locale) io.ReadCloser,
	getStoreFileLocation func(string, locale.Locale) string,
	chooseId func(processorModel.Info) string,
	chooseText func(processorModel.Info) string,
) {
	id := chooseId(info)
	log := logrus.WithField("id", id).WithField("locale", info.Locale())

	p.stopPlay(ctx)

	// generate new speech context
	speechCtx, speechCancel := context.WithCancel(ctx)
	var mp3Stream io.ReadCloser

	localSpeech := getFromStore(id, info.Locale())
	if localSpeech != nil {
		//use the local mp3 file instead of generating them again
		mp3Stream = localSpeech
		log.Info("Speech from local file.")
	} else {
		speechGenerator := p.speechPool.SpeechGeneratorFor(info.Locale())
		if speechGenerator == nil {
			log.Error("No speech generator found for given locale.")
			return
		}

		//ATTENTION: here we use the "global" context because if the current speech is interrupting (because another
		// speech is incoming), the generation should NOT be interrupted. Because the speech should be saved to disk.
		// Only the playing should be interrupted.
		speechStream, err := speechGenerator.SpeechAsNpc(ctx, chooseText(info), model.NonPlayerCharacter{
			Male: info.Npc.Sex == processorModel.SexMale,
		})
		if err != nil {
			log.WithError(err).Error("Unable to generate speech!")
			return
		}

		// fork the speech stream to disk, so we can later use the file from disk instead of generating them again
		mp3Stream, err = system.NewTeeReader(speechStream, getStoreFileLocation(id, info.Locale()))
		if err != nil {
			log.WithError(err).Warn("Unable to initialise tee reader! The speech will not saved to disk!")
			mp3Stream = speechStream
		}
	}

	go p.play(log, speechCtx, speechCancel, mp3Stream)
}

func (p *processor) stopPlay(ctx context.Context) {
	if p.speechCancel == nil {
		return
	}

	// interrupt current speech
	p.speechCancel()

	// wait until speech is done (interrupted)
	select {
	case <-ctx.Done():
		return
	case <-p.speechDoneChan:
	}
}

func (p *processor) play(log *logrus.Entry, ctx context.Context, cancel context.CancelFunc, stream io.ReadCloser) {
	defer stream.Close()

	p.speechCtx = ctx
	p.speechCancel = cancel
	p.speechDoneChan = make(chan bool, 1)

	log.Info("Speech start.")
	err := p.mp3Player.Play(ctx, stream)
	if err != nil {
		log.WithError(err).Warn("Error while playing speech.")
	}
	log.Info("Speech done.")

	cancel()

	p.speechCtx = nil
	p.speechCancel = nil
	p.speechDoneChan <- true
	close(p.speechDoneChan)
}
