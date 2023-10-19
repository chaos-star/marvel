package CacheCluster

import "github.com/redis/go-redis/v9"

type Cluster struct {
	*redis.ClusterClient
}

func Initialize(config map[string]interface{}) *Cluster {
	var (
		address  []string
		username string
		passwd   string
	)
	if addrs, ok := config["address"].([]interface{}); ok {
		for _, addr := range addrs {
			address = append(address, addr.(string))
		}
	} else {
		panic("redis address exception")
	}

	if name, ok := config["username"].(string); ok {
		username = name
	}

	if pwd, ok := config["password"].(string); ok {
		passwd = pwd
	}

	return &Cluster{
		redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    address,
			Username: username,
			Password: passwd,
		}),
	}
}
