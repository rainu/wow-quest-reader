package aws

import (
	"github.com/rainu/wow-quest-reader/internal/locale"
	common "github.com/rainu/wow-quest-reader/internal/speech/sound"
	"github.com/sirupsen/logrus"
)

type speechPool struct {
	awsRegion string
	awsKey    string
	awsSecret string

	speechRate string

	pool map[locale.Locale]common.SpeechGenerator
}

func NewPool(region, key, secret, speechRate string) *speechPool {
	return &speechPool{
		awsRegion: region,
		awsKey:    key,
		awsSecret: secret,

		speechRate: speechRate,

		pool: map[locale.Locale]common.SpeechGenerator{},
	}
}

func (s *speechPool) SpeechGeneratorFor(l locale.Locale) common.SpeechGenerator {
	client := s.pool[l]
	if client == nil {
		var err error
		s.pool[l], err = New(s.awsRegion, s.awsKey, s.awsSecret, s.speechRate, l)
		if err != nil {
			logrus.WithError(err).WithField("locale", l).Error("Unable to initialise new speech generator!")
			return nil
		}
	}

	return s.pool[l]
}
