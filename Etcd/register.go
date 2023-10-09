package Etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdRegisterService struct {
	client  *clientv3.Client
	leaseId clientv3.LeaseID
	ctx     context.Context
}

// 创建租约
func (e *etcdRegisterService) createLease(expire int64) error {
	lgr, err := e.client.Grant(e.ctx, expire)
	if err != nil {
		return err
	}
	e.leaseId = lgr.ID
	return nil
}

// BindLease 键绑定租约
func (e *etcdRegisterService) BindLease(key string, value string) error {
	_, err := e.client.Put(e.ctx, key, value, clientv3.WithLease(e.leaseId))
	if err != nil {
		return err
	}
	return nil
}

// KeepAlive 租约续约
func (e *etcdRegisterService) KeepAlive(isWatch bool) error {
	leaseChan, err := e.client.KeepAlive(e.ctx, e.leaseId)
	if err != nil {
		return err
	}
	if isWatch {
		go e.watch(leaseChan)
	}
	return nil
}

// 监控续租
func (e *etcdRegisterService) watch(leaseChan <-chan *clientv3.LeaseKeepAliveResponse) {
	for _ = range leaseChan {

		//fmt.Printf("续约成功; content:%+v\n",k)
	}
	fmt.Println("续租失败")
}

func (e *etcdRegisterService) Close() error {
	_, err := e.client.Revoke(e.ctx, e.leaseId)
	return err
}
