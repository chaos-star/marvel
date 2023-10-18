package Cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	*redis.Client
}

func Initialize(config map[string]interface{}) *Cache {
	var (
		address string
		passwd  string
		db      int64
	)
	if addr, ok := config["address"].(string); ok {
		address = addr
	} else {
		panic("redis address exception")
	}

	if pwd, ok := config["password"].(string); ok {
		passwd = pwd
	}

	if cdb, ok := config["db"].(int64); ok {
		db = cdb
	}
	return &Cache{
		redis.NewClient(&redis.Options{
			Addr:     address,
			Password: passwd,
			DB:       int(db),
		}),
	}
}
