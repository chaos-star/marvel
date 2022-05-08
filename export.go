package marvel

import (
	"github.com/khaos-star/marvel/Config"
	etcd "github.com/khaos-star/marvel/Etcd"
	"github.com/khaos-star/marvel/Log"
	"github.com/khaos-star/marvel/Mysql/Gorm"
)

var (
	Conf  *Config.Config
	Logger Log.ILogger
	Mysql *Gorm.Engine
	Etcd  *etcd.Engine
)
