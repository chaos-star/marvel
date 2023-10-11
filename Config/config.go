package Config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config *Config
)

type Config struct {
	config *viper.Viper
}

func Initialize() (error, *Config) {
	var err error
	once.Do(func() {
		err, config = New("./Conf", "app", "toml")
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
		return err, nil
	}
	return err, &Config{viper}
}

func (cfg *Config) SetConfigFile(filepath string, filename string, filetype string) error {
	cfg.config.AddConfigPath(filepath)
	cfg.config.SetConfigName(filename)
	cfg.config.SetConfigType(filetype)
	err := cfg.config.ReadInConfig()
	return err
}

func (cfg *Config) GetConfig(key string) interface{} {
	return cfg.config.Get(key)
}

func (cfg *Config) GetString(key string) string {
	return cfg.config.GetString(key)
}

func (cfg *Config) GetBool(key string) bool {
	return cfg.config.GetBool(key)
}

func (cfg *Config) GetInt(key string) int {
	return cfg.config.GetInt(key)
}

func (cfg *Config) GetInt32(key string) int32 {
	return cfg.config.GetInt32(key)
}

func (cfg *Config) GetInt64(key string) int64 {
	return cfg.config.GetInt64(key)
}

func (cfg *Config) GetUint(key string) uint {
	return cfg.config.GetUint(key)
}

func (cfg *Config) GetUint32(key string) uint32 {
	return cfg.config.GetUint32(key)
}

func (cfg *Config) GetUint64(key string) uint64 {
	return cfg.config.GetUint64(key)
}

func (cfg *Config) GetFloat64(key string) float64 {
	return cfg.config.GetFloat64(key)
}

func (cfg *Config) GetTime(key string) time.Time {
	return cfg.config.GetTime(key)
}

func (cfg *Config) GetDuration(key string) time.Duration {
	return cfg.config.GetDuration(key)
}

func (cfg *Config) GetStringSlice(key string) []string {
	return cfg.config.GetStringSlice(key)
}

func (cfg *Config) GetStringMap(key string) map[string]interface{} {
	return cfg.config.GetStringMap(key)
}

func (cfg *Config) GetStringMapString(key string) map[string]string {
	return cfg.config.GetStringMapString(key)
}

func (cfg *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cfg.config.GetStringMapStringSlice(key)
}

func (cfg *Config) IsSet(key string) bool {
	return cfg.config.IsSet(key)
}
