package system

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Db struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type S3 struct {
	key    string `yaml:"key"`
	Secret string `yaml:"secret"`
	Bucket string `yaml:"bucket"`
	Path   string `yaml:"path"`
	Region string `yaml:"region"`
}

type Config struct {
	Views []string `yaml:"views"`
	S3    S3       `yaml:"s3"`
	Db    Db       `yaml:"db"`
}

func ReadConfigFile() (Config, error) {
	//read  yaml config file
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
