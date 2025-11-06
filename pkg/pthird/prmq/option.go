package prmq

import (
	"github.com/streadway/amqp"
)

type Option struct {
	// ExType 交换机类型
	ExType  string
	// ExName 交换机名称
	ExName  string
	// RtKey 路由规则
	RtKey   string
	// QuName 队列名称
	QuName  string
	// 是否持久化 - 重启mq服务不会丢失消息/默认false
	Durable bool
	// 是否自动删除 - 没有连接时自动删除队列/默认false
	AutoDelete bool
	// 处理错误是否重新排队
	Requeue	bool
	// 是否自定申明绑定队列
	AutoBind   bool

	Uri 	 string
	Host     string
	Port     int
	Username string
	Password string
	Vhost    string
}

func (o *Option) uri() string {
	// default vhost: amqp://admin:123456@10.10.6.65:5672/
	// named vhost: amqp://admin:123456@10.10.6.65:5672/app
	u := &amqp.URI{
		Username: o.Username,
		Password: o.Password,
		Host:     o.Host,
		Port:     o.Port,
		Scheme:   "amqp",
		Vhost:    o.Vhost,
	}
	return u.String()
}

func (o *Option) Build() (Client, error) {
	if o.Uri == "" {
		o.Uri = o.uri()
	}
	cli := &client{
		opt: o,
	}
	if err := cli.connect(); err != nil {
		cli.Close()
		return nil, err
	}
	return cli, nil
}
