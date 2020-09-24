package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Config struct
type config struct {
	Server struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		MaxHeaderBytes int    `yaml:"maxHeaderBytes"`
		Timeout        struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	ARI struct {
		Enable bool `yaml:"enable"`
	} `yaml:"ari"`
	AMI struct {
		Enable   bool   `yaml:"enable"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"ami"`
	DB struct {
		DrvName  string `yaml:"drvname"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`
	Log *logrus.Logger
}

func configNew() *config {
	configPath := os.Getenv("GOFRACFG")
	if len(configPath) == 0 {
		configPath = "config.yml"
	}

	// Create config structure
	cfg := &config{
		Log: logrus.New(),
	}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		cfg.Log.Fatal(err)
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&cfg); err != nil {
		cfg.Log.Fatal(err)
	}

	return cfg
}
