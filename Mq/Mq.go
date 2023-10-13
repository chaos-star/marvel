package Mq

import (
	"errors"
	mate "github.com/chaos-star/queue-mate"
	"time"
)

type Mq struct {
	MqRabbits map[string]*mate.Rabbit
	MqBase    *mate.MQBase
}

func (m *Mq) Instance(name string) *mate.Rabbit {
	if rabbit, ok := m.MqRabbits[name]; ok {
		return rabbit
	}
	return nil
}

type mqConfig struct {
	Alias       string `json:"alias"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"maxIdle"`
	MaxLifeTime int    `json:"maxLifeTime"`
	TimeOut     int    `json:"timeout"`
	VHost       string `json:"vhost""`
}

func Initialize(mqConfigs interface{}, logger mate.Logger) (*Mq, error) {
	var (
		mcs []map[string]interface{}
	)
	if v, ok := mqConfigs.(map[string]interface{}); ok {
		mcs = append(mcs, v)
	}

	if v, ok := mqConfigs.([]interface{}); ok {
		if len(v) > 0 {
			for _, iv := range v {
				if ivm, is := iv.(map[string]interface{}); is {
					mcs = append(mcs, ivm)
				}
			}
		}
	}

	var mqInst Mq
	mqInst.MqBase = new(mate.MQBase).With(logger)
	for _, mc := range mcs {
		mConf, err := parseMqConfig(mc)
		if err != nil {
			panic(err)
		}
		mqInst.MqRabbits[mConf.Alias] = mate.NewRabbit(
			mConf.Address,
			mConf.Port,
			mConf.Username,
			mConf.Password,
			mConf.VHost,
			mConf.MaxIdle,
			time.Duration(mConf.MaxLifeTime),
			time.Duration(mConf.TimeOut),
			logger,
		)

	}
	return &mqInst, nil
}

func parseMqConfig(conf map[string]interface{}) (*mqConfig, error) {

	var mc mqConfig

	if alias, ok := conf["alias"]; ok {
		if val, is := alias.(string); is && val != "" {
			mc.Alias = val
		} else {
			return nil, errors.New("the data type of alias is incorrect or empty")
		}
	} else {
		return nil, errors.New("alias is not exists")
	}

	if address, ok := conf["address"]; ok {
		if val, is := address.(string); is && val != "" {
			mc.Address = val
		} else {
			return nil, errors.New("the data type of address is incorrect or empty")
		}
	} else {
		return nil, errors.New("address is not exists")
	}

	if username, ok := conf["username"]; ok {
		if val, is := username.(string); is && val != "" {
			mc.Username = val
		} else {
			return nil, errors.New("the data type of username is incorrect or empty")
		}
	} else {
		return nil, errors.New("username is not exists")
	}

	if passwd, ok := conf["password"]; ok {
		if val, is := passwd.(string); is && val != "" {
			mc.Password = val
		} else {
			return nil, errors.New("the data type of password is incorrect or empty")
		}
	} else {
		return nil, errors.New("password is not exists")
	}

	if vhost, ok := conf["vhost"]; ok {
		if val, is := vhost.(string); is && val != "" {
			mc.VHost = val
		} else {
			return nil, errors.New("the data type of vhost is incorrect or empty")
		}
	} else {
		return nil, errors.New("vhost is not exists")
	}

	if port, ok := conf["port"]; ok {
		if val, is := port.(int64); is && val > 0 {
			mc.Port = int(val)
		} else {
			return nil, errors.New("the data type of port is incorrect or empty")
		}
	} else {
		return nil, errors.New("port is not exists")
	}

	if timeOut, ok := conf["timeout"]; ok {
		if val, is := timeOut.(int); is && val > 0 {
			mc.TimeOut = val
		}
	}

	if maxIdle, ok := conf["maxIdle"]; ok {
		if val, is := maxIdle.(int); is && val > 0 {
			mc.MaxIdle = val
		}
	}

	if maxLifeTime, ok := conf["maxLifeTime"]; ok {
		if val, is := maxLifeTime.(int); is && val > 0 {
			mc.MaxLifeTime = val
		}
	}

	return &mc, nil
}
