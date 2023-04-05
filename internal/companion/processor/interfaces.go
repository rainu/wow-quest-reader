package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/companion/system"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/speech/sound"
	"io"
)

type SpeechPool interface {
	SpeechGeneratorFor(locale locale.Locale) sound.SpeechGenerator
}

type ClipboardAccess interface {
	Clear(ctx context.Context)
	Watch(ctx context.Context, contentChan chan string)
}

type KeyPresser interface {
	PressCopy()
	Press(kp system.KeyPressing)
	PressRaw(ctrl, shift bool, keys ...int)
}

type KeyConfiguration struct {
	HotKeyReading system.Hotkey

	AddonKeyPressing system.KeyPressing
}

type Mp3Player interface {
	Play(ctx context.Context, mp3Stream io.ReadCloser) error
}

type SoundStore interface {
	GetDescription(questId string, locale locale.Locale) io.ReadCloser
	GetProgress(questId string, locale locale.Locale) io.ReadCloser
	GetCompletion(questId string, locale locale.Locale) io.ReadCloser
	GetGossip(gId string, locale locale.Locale) io.ReadCloser

	GetFileLocationForDescription(questId string, locale locale.Locale) string
	GetFileLocationForProgress(questId string, locale locale.Locale) string
	GetFileLocationForCompletion(questId string, locale locale.Locale) string
	GetFileLocationForGossip(gId string, locale locale.Locale) string
}
