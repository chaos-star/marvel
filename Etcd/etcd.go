package Etcd

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

type Engine struct {
	*clientv3.Client
}

func Initialize(config map[string]interface{}) (*Engine, error) {
	var (
		points           []string
		username         string
		password         string
		timeout          int64 = 5
		keepAlive        int64
		keepAliveTimeout int64
	)

	if len(config) <= 0 {
		return nil, errors.New("the ETCD configuration is missing")
	}

	if val, ok := config["endpoints"].([]interface{}); ok {
		if len(val) > 0 {
			for _, item := range val {
				if sVal, is := item.(string); is {
					points = append(points, sVal)
				}
			}
		}
	}

	if val, ok := config["username"].(string); ok {
		username = val
	}

	if val, ok := config["password"].(string); ok {
		password = val
	}

	if val, ok := config["timeout"].(int64); ok {
		timeout = val
	}

	if val, ok := config["keep_alive"].(int64); ok {
		keepAlive = val
	}

	if val, ok := config["keep_alive_timeout"].(int64); ok {
		keepAliveTimeout = val
	}

	if len(points) <= 0 {
		return nil, errors.New("etcd requires at least one endpoint")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: points,
		Username:  username,
		Password:  password,
		//多久未连接上则超时
		DialTimeout: time.Duration(timeout) * time.Second,
		//客户端ping服务器以查看TCP是否存活的时间
		DialKeepAliveTime: time.Duration(keepAlive) * time.Second,
		//客户端等待keep-alive探测响应的时间。如果在此时间内没有收到响应，则关闭连接。
		DialKeepAliveTimeout: time.Duration(keepAliveTimeout) * time.Second,
	})

	return &Engine{cli}, err
}


func (x *Engine) NewGrpcServiceDiscovery(schema string,service string) resolver.Builder  {
	return &etcdGrpcDiscoveryService{client: x.Client,schema: schema,service: service}
}

// RegisterService 注册服务
func (x *Engine) RegisterService(key,value string, expire int64, on bool)(*etcdRegisterService, error)  {
	srv := &etcdRegisterService{
		client: x.Client,
		ctx: context.Background(),
	}
	err := srv.createLease(expire)
	if err != nil{
		return nil,err
	}
	err = srv.BindLease(key,value)
	if err != nil{
		return nil,err
	}
	err = srv.KeepAlive(on)
	if err != nil{
		return nil,err
	}
	return srv, nil
}
