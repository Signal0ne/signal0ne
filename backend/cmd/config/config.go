package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Server struct {
	Mode           string `mapstructure:"MODE"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	ServerDomain   string `mapstructure:"SERVER_DOMAIN"`
	ServerIsSecure bool   `mapstructure:"SERVER_IS_SECURE"`
}

type Config struct {
	Debug           bool   `mapstructure:"DEBUG"`
	IPCSocket       string `mapstructure:"IPC_SOCKET"`
	MongoUri        string `mapstructure:"MONGO_URI"`
	Server          Server `mapstructure:",squash"`
	SignalOneSecret string `mapstructure:"SIGNAL_ONE_SECRET"`
	SkipAuth        bool   `mapstructure:"SKIP_AUTH"`
}

var (
	config *Config
	once   sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		viper.SetConfigName(".default")
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			return
		}

		config = &Config{}

		err = viper.Unmarshal(config)
		if err != nil {
			return
		}
	})

	return config
}
