package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	Conn *clientv3.Client
}

// 创建etcd的客户端
func NewEtcdClient(etcdNodes []string, user, pwd string) (EC *EtcdClient, err error) {
	EC = &EtcdClient{}
	EC.Conn, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdNodes,
		DialTimeout: 5 * time.Second,
		Username:    user,
		Password:    pwd,
	})
	return
}

// 解析etcd内的value
func (e *EtcdClient) Parse(key string, v interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	result, err := e.Conn.Get(ctx, key)
	if err != nil {
		return err
	}
	if result == nil || len(result.Kvs) == 0 {
		return fmt.Errorf("no more config value,the key is %s, value is %v", key, result)
	}
	return json.Unmarshal(result.Kvs[0].Value, v)
}

// 关闭etcd的连接
func (e *EtcdClient) Close() {
	_ = e.Conn.Close()
}
