package CacheCluster

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type Cluster struct {
	redisCluster map[string]*redis.ClusterClient
}

func Initialize(configs interface{}) (*Cluster, error) {
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
	var clusterInst = &Cluster{}
	clusterInst.redisCluster = make(map[string]*redis.ClusterClient)
	for _, mc := range mcs {
		redisInst, alias, err := clusterInst.newCluster(mc)
		if err != nil {
			return nil, err
		}
		clusterInst.redisCluster[alias] = redisInst
	}
	return clusterInst, nil
}

func (c *Cluster) Instance(name string) *redis.ClusterClient {
	if inst, ok := c.redisCluster[name]; ok {
		return inst
	}
	return nil
}

func (c *Cluster) newCluster(config map[string]interface{}) (*redis.ClusterClient, string, error) {
	var (
		alias    string
		address  []string
		username string
		passwd   string
	)
	if nick, ok := config["alias"].(string); ok {
		alias = nick
	} else {
		return nil, "", errors.New("redis alias exception")
	}

	if adders, ok := config["address"].([]interface{}); ok {
		for _, addr := range adders {
			address = append(address, addr.(string))
		}
	} else {
		return nil, "", errors.New("redis address exception")
	}

	if name, ok := config["username"].(string); ok {
		username = name
	}

	if pwd, ok := config["password"].(string); ok {
		passwd = pwd
	}

	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    address,
		Username: username,
		Password: passwd,
	})

	_, err := cluster.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}

	return cluster, alias, nil
}
