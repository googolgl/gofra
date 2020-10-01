package mod

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	cfg = Config{}
)

// Config struct
type Config struct {
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
		Enable      bool   `yaml:"enable"`
		Application string `yaml:"application"`
		URL         string `yaml:"url"`
		WS          string `yaml:"ws"`
		UserName    string `yaml:"username"`
		Password    string `yaml:"password"`
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
		Database string `yaml:"database"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`
	FilePath string        `yaml:"filepath"`
	LogLevel *logrus.Level `yaml:"loglevel"`
	Log      *logrus.Entry //*logrus.Logger

}

//ConfigNew - construct config
func ConfigNew() *Config {
	configPath := os.Getenv("GOFRA_CFG")
	if len(configPath) == 0 {
		configPath = "config.yml"
	}

	// Create config structure
	cfg = Config{
		Log: logrus.WithFields(logrus.Fields{"mod": "config", "func": "ConfigNew"}),
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

	// set loglevel
	cfg.Log.Logger.SetLevel(*cfg.LogLevel)

	return &cfg
}
