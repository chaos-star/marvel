package Config

import (
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	once   sync.Once
	config *Config
)

type Config struct {
	*viper.Viper
}

const EnvConfPath = "ARTEMIS_CONFIG_PATH"

func Initialize() (error, *Config) {
	var err error
	once.Do(func() {
		path := os.Getenv(EnvConfPath)
		if path == "" {
			path = "./Conf"
		}
		err, config = New(path, "app", "toml")
	})
	if err != nil {
		return err, nil
	}
	return err, config
}

func New(filepath string, filename string, filetype string) (error, *Config) {
	viper := viper.New()
	viper.AddConfigPath(filepath)
	viper.SetConfigName(filename)
	viper.SetConfigType(filetype)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return err, &Config{viper}
}

func (cfg *Config) SetConfigFile(filepath string, filename string, filetype string) error {
	cfg.AddConfigPath(filepath)
	cfg.SetConfigName(filename)
	cfg.SetConfigType(filetype)
	err := cfg.ReadInConfig()
	return err
}
