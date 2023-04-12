package config

import (
	"github.com/micmonay/keybd_event"
	"github.com/rainu/wow-quest-reader/internal/companion/system"
	"github.com/sirupsen/logrus"
	"golang.design/x/hotkey"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type Config struct {
	Debug    bool         `yaml:"debug"`
	LogLevel logrus.Level `yaml:"logLevel"`

	Sound struct {
		Directory        string `yaml:"dir"`
		AmazonWebService struct {
			Region     string `yaml:"region"`
			Key        string `yaml:"key"`
			Secret     string `yaml:"secret"`
			SpeechRate string `yaml:"speechRate"`
		} `yaml:"aws"`
	} `yaml:"sound"`

	Key key `yaml:"key"`
}

type key struct {
	Read  system.Hotkey      `yaml:"read"`
	Addon system.KeyPressing `yaml:"addon"`
}

var appDir = determineApplicationDirectory()

func determineApplicationDirectory() string {
	app, err := os.Executable()
	if err == nil {
		return path.Dir(app)
	}

	worDir, err := os.Getwd()
	if err == nil {
		return worDir
	}

	return "./"
}

func Read() Config {
	// default values
	result := Config{
		Debug:    true,
		LogLevel: logrus.InfoLevel,
		Key: key{
			Read:  system.Hotkey{Modifier: []hotkey.Modifier{}, Key: 0x22}, //PDOWN
			Addon: system.KeyPressing{Ctrl: true, Keys: []int{keybd_event.VK_F12}},
		},
	}

	// read potential yaml config file
	readConfigs(&result,
		path.Join(appDir, "config.yml"),
		path.Join(appDir, "config.yaml"),
		path.Join("./", "config.yml"),
		path.Join("./", "config.yaml"),
	)

	if result.Sound.Directory == "" {
		result.Sound.Directory = path.Join(appDir, "sounds")
	}

	if result.Debug {
		result.LogLevel = logrus.DebugLevel
	}

	result.Validate()
	return result
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func readConfigs(cfg *Config, path ...string) {
	for _, curPath := range path {
		if exists(curPath) {
			readConfig(curPath, cfg)
		}
	}
}

func readConfig(path string, cfg *Config) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		logrus.WithError(err).WithField("path", path).Fatalf("Unable to read configuration file!")
	}
}
