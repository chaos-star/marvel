package Etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

type etcdGrpcDiscoveryService struct {
	client      *clientv3.Client
	conn        resolver.ClientConn
	serviceList sync.Map
	schema      string
	service     string
}

// Build 实现 resolver.Build接口
func (ds *etcdGrpcDiscoveryService) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ds.conn = cc
	prefix := fmt.Sprintf("/%s.%s/", target.URL.Scheme, target.URL.Host)
	fmt.Println(prefix)
	gr, err := ds.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range gr.Kvs {
		ds.store(kv.Key, kv.Value)
	}
	ds.updateStatue()
	go ds.watch(prefix)
	return ds, nil
}

func (ds *etcdGrpcDiscoveryService) watch(prefix string) {
	wCh := ds.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for w := range wCh {
		for _, event := range w.Events {
			switch event.Type {
			case 0:
				ds.store(event.Kv.Key, event.Kv.Value)
				ds.updateStatue()
			case 1:
				ds.delete(event.Kv.Key)
				ds.updateStatue()
			}
		}
	}
}

func (ds *etcdGrpcDiscoveryService) store(k, v []byte) {
	ds.serviceList.Store(string(k), string(v))
}

func (ds *etcdGrpcDiscoveryService) delete(k []byte) {
	ds.serviceList.Delete(string(k))
}

func (ds *etcdGrpcDiscoveryService) updateStatue() {
	var ipList resolver.State
	ds.serviceList.Range(func(key, value interface{}) bool {
		ip, ok := value.(string)
		if !ok {
			return false
		}
		ipList.Addresses = append(ipList.Addresses, resolver.Address{Addr: ip})
		return true
	})
	ds.conn.UpdateState(ipList)
}

func (ds *etcdGrpcDiscoveryService) Scheme() string {
	return ds.schema
}

func (ds *etcdGrpcDiscoveryService) Service() string {
	return ds.service
}

func (ds *etcdGrpcDiscoveryService) ResolveNow(options resolver.ResolveNowOptions) {
}

func (ds *etcdGrpcDiscoveryService) Close() {
	fmt.Println("discovery close over")
}
