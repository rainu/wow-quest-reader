package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

var percentageRegexp = regexp.MustCompile(`^([0-9]{1,3})%$`)

func (c *Config) Validate() {
	var errors []error
	if c.Sound.Directory == "" {
		errors = append(errors, fmt.Errorf("missing sound directory"))
	}
	if c.Sound.AmazonWebService.Region == "" {
		errors = append(errors, fmt.Errorf("aws region is missing"))
	}
	if c.Sound.AmazonWebService.Key == "" {
		errors = append(errors, fmt.Errorf("aws key is missing"))
	}
	if c.Sound.AmazonWebService.Secret == "" {
		errors = append(errors, fmt.Errorf("aws secret is missing"))
	}
	if c.Sound.AmazonWebService.SpeechRate != "" {
		r := strings.ToLower(c.Sound.AmazonWebService.SpeechRate)
		r = strings.ReplaceAll(r, " ", "")

		if percentageRegexp.MatchString(r) {
			sVal := percentageRegexp.FindStringSubmatch(r)[1]
			iVal, err := strconv.ParseInt(sVal, 10, 8)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid percentage value for speech rate: %w", err))
			} else if iVal < 20 || iVal > 200 {
				errors = append(errors, fmt.Errorf("invalid percentage value for speech rate: must be between 20 and 200"))
			}
		} else if r != "x-slow" && r != "slow" && r != "medium" && r != "fast" && r != "x-fast" {
			errors = append(errors, fmt.Errorf("invalid speech rate"))
		}
	}

	if len(errors) > 0 {
		log := logrus.NewEntry(logrus.StandardLogger())
		for i, err := range errors {
			log = log.WithField(fmt.Sprintf("error#%d", i), err)
		}

		log.Fatal("Configuration is invalid!")
	}
}
