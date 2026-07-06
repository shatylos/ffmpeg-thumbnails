package app

import (
	"github.com/shatylos/ffmpeg-thumbnails/tools/apperrors"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

const (
	StorageDisk   = "disk"
	StorageMemory = "memory"
)

type Config struct {
	Outputdir     string         `yaml:"output_dir"`
	Frequency     time.Duration  `yaml:"frequency"`
	Timeout       time.Duration  `yaml:"timeout"`
	Forks         int            `yaml:"forks"`
	Storage       string         `yaml:"storage"`
	ServerEnabled bool           `yaml:"server_enabled"`
	ServerAddr    string         `yaml:"server_addr"`
	Streams       []StreamConfig `yaml:"streams"`
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

	if config.Storage == "" {
		config.Storage = StorageDisk
	}
	if config.Storage != StorageDisk && config.Storage != StorageMemory {
		err = apperrors.New("storage must be %q or %q, got %q", StorageDisk, StorageMemory, config.Storage)
		return
	}

	if config.ServerEnabled && config.ServerAddr == "" {
		config.ServerAddr = ":8080"
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
