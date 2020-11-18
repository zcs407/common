package conn

import (
	"time"

	"github.com/zcs407/common/zlog"

	jsoniter "github.com/json-iterator/go"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	Conn *nats.Conn
}

const (
	reconWaitTime                      = 2 * time.Second
	maxReconnects                      = 5
	errCheckFuncInit                   = "NewNatsClient"
	errCheckFuncPublish                = "Publish"
	errCheckFuncSubscribe              = "Subscribe"
	errCheckFuncQueueSubscribe         = "QueueSubscribe"
	errcheckFuncChanSubscribe          = "ChanSubscribe"
	errcheckFuncChanChanQueueSubscribe = "ChanQueueSubscribe"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 创建etcd的客户端
func NewNatsClient(natsUrl, pwd string) (NC *NatsClient, err error) {
	NC = &NatsClient{}
	NC.Conn, err = nats.Connect(
		natsUrl,
		nats.DontRandomize(),
		nats.MaxReconnects(maxReconnects),
		nats.ReconnectWait(reconWaitTime),
		nats.ClosedHandler(func(nc *nats.Conn) {
			// TODO log
		}),
		nats.ErrorHandler(func(conn *nats.Conn, subscription *nats.Subscription, e error) {
			// TODO log
		}))
	return
}

func (n *NatsClient) Close() {
	_ = n.Conn.Close
}

// 发布消息到nats
func (n *NatsClient) Publish(subj string, msg interface{}) {
	data, err := json.Marshal(&msg)
	if err != nil {
		n.checkErr("Publish", subj, "", err)
		return
	}
	err = n.Conn.Publish(subj, data)
	n.checkErr("Publish", subj, "", err)
}

// 普通订阅,可多次消费,单节点时使用,多用于goroutine 所以错误日志内部处理
func (n *NatsClient) Subscribe(subj string, msgHandle func(data []byte)) {
	_, err := n.Conn.Subscribe(subj, func(msg *nats.Msg) {
		msgHandle(msg.Data)
	})
	n.checkErr(errCheckFuncSubscribe, subj, "", err)
}

// 队列形式的订阅,只消费一次,多服务时使用,多用于goroutine 所以错误日志内部处理
func (n *NatsClient) QueueSubscribe(subj, queue string, msgHandle func(data []byte)) {
	_, err := n.Conn.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		msgHandle(msg.Data)
	})
	n.checkErr(errCheckFuncQueueSubscribe, subj, queue, err)
}

// 通道订阅模式,用于排序处理消息,多次订阅
func (n *NatsClient) ChanSubscribe(subj string, msgHandle func(data []byte)) {
	ch := make(chan *nats.Msg, 1024)
	_, err := n.Conn.ChanSubscribe(subj, ch)
	n.checkErr(errcheckFuncChanSubscribe, subj, "", err)
	for k := range ch {
		msgHandle(k.Data)
	}
}

// 通道订阅模式,用于排序处理消息,单次订阅
func (n *NatsClient) ChanQueueSubscribe(subj, queue string, msgHandle func(data []byte)) {
	ch := make(chan *nats.Msg, 1024)
	_, err := n.Conn.ChanQueueSubscribe(subj, queue, ch)
	n.checkErr(errcheckFuncChanChanQueueSubscribe, subj, queue, err)
	for k := range ch {
		msgHandle(k.Data)
	}

}

// 检查错误
func (n *NatsClient) checkErr(funcName, subj, queue string, err error) {
	if err != nil {
		switch funcName {
		case errCheckFuncPublish:
			zlog.ErrWithStr(err).Str("subj", subj).Msg("nats 发布失败")
		case errCheckFuncSubscribe, errCheckFuncQueueSubscribe, errcheckFuncChanSubscribe,
			errcheckFuncChanChanQueueSubscribe:
			zlog.ErrWithStr(err).Str("subj", subj).Str("queue", queue).Msg("nats 订阅失败")
		default:
		}
	}
}
