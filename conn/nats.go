package conn

import (
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	Conn *nats.Conn
}

// 创建etcd的客户端
func NewNatsClient(natsUrl, pwd string) (NC *NatsClient, err error) {
	NC = &NatsClient{}
	NC.Conn, err = nats.Connect(natsUrl)
	return
}

func (n *NatsClient) Close() {
	_ = n.Conn.Close
}
