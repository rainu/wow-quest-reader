package aws

import (
	"fmt"
	"regexp"
	"strings"
)

var loudSpeechParts = regexp.MustCompile(`\b[A-Z][A-Z ]*[A-Z]\b`)

func transformText(text, rate string) string {
	result := text
	result = strings.ReplaceAll(result, "<", "")
	result = strings.ReplaceAll(result, ">", "")

	parts := loudSpeechParts.FindStringSubmatch(result)
	for _, part := range parts {
		result = strings.ReplaceAll(result, part, fmt.Sprintf(`<prosody volume="loud">%s</prosody>`, part))
	}

	if rate == "" {
		result = fmt.Sprintf(`<speak>%s</speak>`, result)
	} else {
		result = fmt.Sprintf(`<speak><prosody rate="%s">%s</prosody></speak>`, rate, result)
	}
	return result
}
