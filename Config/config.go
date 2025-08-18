package Config

import (
	"bytes"
	"sync"

	"github.com/spf13/viper"
)

var (
	once    sync.Once
	xLoader Loading
)

type Config struct {
	*viper.Viper
	loader Loading
}

func SetLoader(loader Loading) {
	xLoader = loader
}
func Initialize() (err error, config *Config) {
	once.Do(func() {
		if xLoader == nil {
			err, config = NewConfig(&DefaultLoading{})
		} else {
			err, config = NewConfig(xLoader)
		}
	})
	if err != nil {
		return
	}
	return
}

func NewConfig(load Loading) (err error, config *Config) {
	var (
		viper = viper.New()
		body  []byte
	)
	viper.SetConfigType("toml")
	err, body = load.Load()
	if err != nil {
		return
	}
	err = viper.ReadConfig(bytes.NewReader(body))
	if err != nil {
		return
	}
	config = &Config{viper, load}
	return
}

func (cfg *Config) SetConfigFile(filepath string, filename string, filetype string) error {
	cfg.AddConfigPath(filepath)
	cfg.SetConfigName(filename)
	cfg.SetConfigType(filetype)
	err := cfg.ReadInConfig()
	return err
}
func (cfg *Config) Refresh() (err error) {
	var body []byte
	err, body = cfg.loader.Load()
	err = cfg.Viper.ReadConfig(bytes.NewReader(body))
	if err != nil {
		return
	}
	return
}
