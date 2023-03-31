package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

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

	if len(errors) > 0 {
		log := logrus.NewEntry(logrus.StandardLogger())
		for i, err := range errors {
			log = log.WithField(fmt.Sprintf("error#%d", i), err)
		}

		log.Fatal("Configuration is invalid!")
	}
}
