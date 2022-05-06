package marvel

import (
	"github.com/khaos-star/marvel/Config"
	"github.com/khaos-star/marvel/Log"
	"github.com/khaos-star/marvel/Mysql/Gorm"
)

var (
	Conf  *Config.Config
	Logger Log.ILogger
	Mysql *Gorm.Engine
)
