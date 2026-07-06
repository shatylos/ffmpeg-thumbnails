package app

import (
	"github.com/shatylos/ffmpeg-screenshots/tools/apperrors"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Outputdir string         `yaml:"output_dir"`
	Frequency time.Duration  `yaml:"frequency"`
	Timeout   time.Duration  `yaml:"timeout"`
	Forks     int            `yaml:"forks"`
	Streams   []StreamConfig `yaml:"streams"`
}

type StreamConfig struct {
	Src    string `yaml:"src"`
	Output string `yaml:"output"`
}

func GetConfig() (config Config, err error) {
	fileName := "config.yml"
	var yamlFile []byte
	yamlFile, err = os.ReadFile(fileName)
	if err != nil {
		err = apperrors.Wrap(err, "error reading config file: %s", fileName)
		return
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		err = apperrors.Wrap(err, "error unmarshal config file")
		return
	}

	if config.Forks < 1 {
		err = apperrors.New("forks must be greater than zero")
		return
	}

	seenOutputs := make(map[string]bool, len(config.Streams))
	for _, stream := range config.Streams {
		if seenOutputs[stream.Output] {
			err = apperrors.New("duplicate stream output: %s", stream.Output)
			return
		}
		seenOutputs[stream.Output] = true
	}

	config.Frequency = config.Frequency * time.Second
	config.Timeout = config.Timeout * time.Second

	return
}
