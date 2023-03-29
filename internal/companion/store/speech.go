package store

import (
	"fmt"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
)

const (
	kindDescription = "DESC"
	kindCompletion  = "COMP"
	kindProgress    = "PROG"
)

type speech struct {
	baseDir string
}

func NewSpeech(baseDir string) (*speech, error) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err = os.MkdirAll(baseDir, 0750)
		if err != nil {
			return nil, err
		}
	}

	return &speech{
		baseDir: baseDir,
	}, nil
}

func (s *speech) GetDescription(questId string, locale locale.Locale) io.ReadCloser {
	return s.get(questId, locale, kindDescription)
}

func (s *speech) GetProgress(questId string, locale locale.Locale) io.ReadCloser {
	return s.get(questId, locale, kindProgress)
}

func (s *speech) GetCompletion(questId string, locale locale.Locale) io.ReadCloser {
	return s.get(questId, locale, kindCompletion)
}

func (s *speech) get(questId string, locale locale.Locale, kind string) io.ReadCloser {
	soundPath := s.path(questId, locale, kind)
	if _, err := os.Stat(soundPath); os.IsNotExist(err) {
		return nil
	}

	stream, err := os.Open(soundPath)
	if err != nil {
		logrus.WithField("path", soundPath).WithError(err).Error("Unable to open sound file for reading!")
		return nil
	}

	return stream
}

func (s *speech) GetFileLocationForDescription(questId string, locale locale.Locale) string {
	return s.getFileLocationFor(questId, locale, kindDescription)
}

func (s *speech) GetFileLocationForProgress(questId string, locale locale.Locale) string {
	return s.getFileLocationFor(questId, locale, kindProgress)
}

func (s *speech) GetFileLocationForCompletion(questId string, locale locale.Locale) string {
	return s.getFileLocationFor(questId, locale, kindCompletion)
}

func (s *speech) getFileLocationFor(questId string, locale locale.Locale, kind string) string {
	p := s.path(questId, locale, kind)
	os.MkdirAll(path.Dir(p), 0750)

	return p
}

func (s *speech) path(questId string, locale locale.Locale, kind string) string {
	return path.Join(s.baseDir, string(locale), fmt.Sprintf("%s_%s.mp3", questId, kind))
}
