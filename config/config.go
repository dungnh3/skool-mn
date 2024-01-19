package config

import (
	"bytes"
	_ "embed"
	"github.com/dungnh3/skool-mn/pkg/db"
	l "github.com/dungnh3/skool-mn/pkg/log"
	"github.com/spf13/viper"
	"strings"
)

var logger = l.New()

//go:embed default.yaml
var defaultConfig []byte

type (
	Config struct {
		Base  `mapstructure:",squash"`
		MySQL *db.MySQL `json:"mysql" yaml:"mysql" mapstructure:"mysql"`
	}

	Env    string
	Listen struct {
		Host string `json:"host" mapstructure:"host" yaml:"host"`
		Port int    `json:"port" mapstructure:"port" yaml:"port"`
	}

	ServerConfig struct {
		HTTP Listen `json:"http" mapstructure:"http" yaml:"http"`
	}

	Base struct {
		Server ServerConfig `yaml:"server" mapstructure:"server"`
		Env    Env          `yaml:"env" mapstructure:"env"`
	}
)

func Load() *Config {
	cfg := &Config{}
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		logger.Fatal("Failed to read viper config", l.Error(err))
		panic(err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		logger.Fatal("Failed to unmarshal config", l.Error(err))
		panic(err)
	}
	return cfg
}
