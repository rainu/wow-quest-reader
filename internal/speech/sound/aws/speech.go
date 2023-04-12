package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/smithy-go/logging"
	"github.com/rainu/wow-quest-reader/internal/locale"
	"github.com/rainu/wow-quest-reader/internal/model"
	common "github.com/rainu/wow-quest-reader/internal/speech/sound"
	"github.com/sirupsen/logrus"
	"io"
)

type awsClient struct {
	polly        *polly.Client
	speechRate   string
	languageCode types.LanguageCode
}

func New(region, key, secret, speechRate string, l locale.Locale) (common.SpeechGenerator, error) {
	cfg := aws.Config{
		Credentials: aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     key,
				SecretAccessKey: secret,
			}, nil
		}),
		Region: region,
		Logger: logging.LoggerFunc(func(classification logging.Classification, format string, v ...interface{}) {
			log := logrus.WithField("origin", "aws")
			switch classification {
			case logging.Debug:
				log.Debugf(format, v...)
			case logging.Warn:
				log.Warnf(format, v...)
			default:
				log.Infof(format, v...)
			}
		}),
	}

	result := &awsClient{
		polly:      polly.NewFromConfig(cfg),
		speechRate: speechRate,
	}

	switch l {
	case locale.English:
		result.languageCode = types.LanguageCodeEnUs
	case locale.German:
		result.languageCode = types.LanguageCodeDeDe
	default:
		return nil, fmt.Errorf("unsupported locale: %s", l)
	}

	return result, nil
}

func (a *awsClient) SpeechAsNpc(ctx context.Context, text string, npc model.NonPlayerCharacter) (io.ReadCloser, error) {
	voiceId := types.VoiceId("Vicki")
	if npc.Male {
		voiceId = "Daniel"
	}

	return a.speech(ctx, text, voiceId)
}

func (a *awsClient) SpeechAsObject(ctx context.Context, text string, object model.Object) (io.ReadCloser, error) {
	return a.speech(ctx, text, "Vicki")
}

func (a *awsClient) SpeechAsItem(ctx context.Context, text string, item model.Item) (io.ReadCloser, error) {
	return a.speech(ctx, text, "Vicki")
}

func (a *awsClient) speech(ctx context.Context, text string, voiceId types.VoiceId) (io.ReadCloser, error) {
	logrus.WithField("voice_id", voiceId).
		WithField("language_code", a.languageCode).
		Info("Generate speech via AWS polly.")

	o, err := a.polly.SynthesizeSpeech(ctx, &polly.SynthesizeSpeechInput{
		OutputFormat: types.OutputFormatMp3,
		Text:         aws.String(transformText(text, a.speechRate)),
		VoiceId:      voiceId,
		Engine:       types.EngineNeural,
		LanguageCode: a.languageCode,
		TextType:     types.TextTypeSsml,
	})
	if err != nil {
		return nil, err
	}

	return o.AudioStream, nil
}
