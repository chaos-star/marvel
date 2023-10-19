package Cache

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis map[string]*redis.Client
}

func Initialize(configs interface{}) (*Cache, error) {
	var (
		mcs []map[string]interface{}
	)
	if v, ok := configs.(map[string]interface{}); ok {
		mcs = append(mcs, v)
	}
	if v, ok := configs.([]interface{}); ok {
		if len(v) > 0 {
			for _, iv := range v {
				if ivm, is := iv.(map[string]interface{}); is {
					mcs = append(mcs, ivm)
				}
			}
		}
	}
	var cacheInst = &Cache{}
	cacheInst.redis = make(map[string]*redis.Client)
	for _, mc := range mcs {
		redisInst, alias, err := cacheInst.newRedis(mc)
		if err != nil {
			return nil, err
		}
		cacheInst.redis[alias] = redisInst
	}
	return cacheInst, nil
}

func (c *Cache) Instance(name string) *redis.Client {
	if inst, ok := c.redis[name]; ok {
		return inst
	}
	return nil
}

func (c *Cache) newRedis(config map[string]interface{}) (*redis.Client, string, error) {
	var (
		alias    string
		address  string
		username string
		passwd   string
		db       int64
	)
	if nick, ok := config["alias"].(string); ok {
		alias = nick
	} else {
		return nil, "", errors.New("redis alias exception")
	}

	if addr, ok := config["address"].(string); ok {
		address = addr
	} else {
		return nil, "", errors.New("redis address exception")
	}

	if name, ok := config["username"].(string); ok {
		username = name
	}

	if pwd, ok := config["password"].(string); ok {
		passwd = pwd
	}

	if cdb, ok := config["db"].(int64); ok {
		db = cdb
	}
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: passwd,
		DB:       int(db),
	}), alias, nil
}
