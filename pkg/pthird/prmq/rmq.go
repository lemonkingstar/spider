package prmq

import (
	"fmt"
	"sync"
	"time"

	"github.com/lemonkingstar/spider/pkg/plog"
	"github.com/lemonkingstar/spider/pkg/psafe"
	"github.com/lemonkingstar/spider/pkg/putil"
	"github.com/streadway/amqp"
)

var (
	logger = plog.WithField("[PACKET]", "prmq")
)

type Client interface {
	Publish(msg string) error
	StartConsume(f func([]byte) error)
	PeriodicConsume(hf func([]byte) error, cf func() bool, periodic int)
	Get(f func([]byte) error)
	Close() error
}

func NewDefault(addr, exType, exName, rtKey, quName string) (Client, error) {
	return (&Option{
		ExType: exType, ExName: exName,
		RtKey: rtKey, QuName: quName,
		Uri:      addr,
		AutoBind: true, Durable: true,
	}).Build()
}

type client struct {
	sync.Mutex
	conn    *amqp.Connection
	channel *amqp.Channel
	opt     *Option

	reconnecting bool
}

func (c *client) connect() error {
	mqConn, err := amqp.Dial(c.opt.Uri)
	if err != nil {
		logger.Errorf("Dial rmq error: %v", err)
		return err
	}
	mqChan, err := mqConn.Channel()
	if err != nil {
		logger.Errorf("Get rmq channel error: %v", err)
		return err
	}
	c.conn = mqConn
	c.channel = mqChan
	if c.opt.AutoBind {
		return c.bind()
	}
	return nil
}

func (c *client) bind() error {
	err := c.channel.ExchangeDeclare(c.opt.ExName, c.opt.ExType, c.opt.Durable, c.opt.AutoDelete, false,
		false, nil)
	if err != nil {
		return err
	}
	_, err = c.channel.QueueDeclare(c.opt.QuName, c.opt.Durable, c.opt.AutoDelete, false, false, nil)
	if err != nil {
		return err
	}
	// bind queue &exchange with key
	err = c.channel.QueueBind(c.opt.QuName, c.opt.RtKey, c.opt.ExName, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Close() error {
	err := c.channel.Close()
	if err != nil {
		logger.Errorf("Close rmq channel error: %v", err)
		return err
	}
	err = c.conn.Close()
	if err != nil {
		logger.Errorf("Close rmq connection error: %v", err)
		return err
	}
	return nil
}

func (c *client) Publish(msg string) error {
	err := c.channel.Publish(c.opt.ExName, c.opt.RtKey, false, false, amqp.Publishing{
		ContentType: "text/plain", Body: []byte(msg),
	})
	if err != nil {
		logger.Errorf("Publish error: %v", err)
		psafe.Go(c.reconnect)
	}
	return nil
}

func (c *client) Consume(f func([]byte) error) (error, bool) {
	c.channel.Qos(1, 0, true)
	msgChan, err := c.channel.Consume(c.opt.QuName, "", false, false,
		false, false, nil)
	if err != nil {
		logger.Errorf("Consume error: %v", err)
		return err, true
	}
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				logger.Warn("Consume done")
				return nil, false
			}
			if err := f(msg.Body); err != nil {
				logger.Errorf("Consume handle error: %v", err)
				if c.opt.Requeue {
					if err := msg.Nack(false, true); err != nil {
						logger.Errorf("Nack current message error: %v", err)
					}
					continue
				}
			}
			if err := msg.Ack(false); err != nil {
				logger.Errorf("Ack current message error: %v", err)
			}
		}
	}
	return nil, false
}

func (c *client) StartConsume(f func([]byte) error) {
	psafe.Go(func() {
		for {
			_, reset := c.Consume(f)
			time.Sleep(9 * time.Second)
			if reset {
				c.reconnect()
			}
		}
	})
}

func (c *client) reconnect() {
	if c.reconnecting == true {
		return
	}
	c.Lock()
	logger.Warn("Start for reconnect...")
	c.reconnecting = true
	defer func() {
		c.reconnecting = false
		c.Unlock()
	}()
	c.Close()
	for {
		err := c.connect()
		if err == nil {
			break
		} else {
			logger.Warn("Wait for reconnect...")
		}
		time.Sleep(9 * time.Second)
	}
}

// Get 同步从队列中获取单个消息
// 考虑到性能问题，生产环境建议使用 Consume
func (c *client) Get(f func([]byte) error) {
	msg, ok, err := c.channel.Get(c.opt.QuName, false)
	if err != nil {
		logger.Errorf("Get message error: %v", err)
		return
	}
	if !ok {
		logger.Warn("No message waiting")
		return
	}
	if err = f(msg.Body); err != nil {
		logger.Errorf("Get handle error: %v", err)
		if c.opt.Requeue {
			if err = msg.Nack(false, true); err != nil {
				logger.Errorf("Nack current message error: %v", err)
			}
			return
		}
	}
	if err = msg.Ack(false); err != nil {
		logger.Errorf("Ack current message error: %v", err)
	}
}

// PeriodicConsume 间隔非持续消费
// hf: 消费处理函数
// cf: 当前节点消费检测函数
// periodic: 检测间隔
func (c *client) PeriodicConsume(hf func([]byte) error, cf func() bool, periodic int) {
	psafe.Go(func() {
		consume := func(hf func([]byte) error, cf func() bool, periodic int) (error, bool) {
			consumer := fmt.Sprintf("prmq-%s", putil.UUID())
			ticker := time.NewTicker(time.Duration(periodic) * time.Second)
			defer ticker.Stop()

			logger.Infof("Start Channel consume: %s", consumer)
			c.channel.Qos(1, 0, true)
			msgChan, err := c.channel.Consume(c.opt.QuName, consumer, false, false,
				false, false, nil)
			if err != nil {
				logger.Errorf("Consume error: %v", err)
				return err, true
			}
			for {
				select {
				case msg, ok := <-msgChan:
					if !ok {
						logger.Warn("Consume done")
						return nil, false
					}
					if err := hf(msg.Body); err != nil {
						logger.Errorf("Consume handle error: %v", err)
						if c.opt.Requeue {
							if err := msg.Nack(false, true); err != nil {
								logger.Errorf("Nack current message error: %v", err)
							}
							continue
						}
					}
					if err := msg.Ack(false); err != nil {
						logger.Errorf("Ack current message error: %v", err)
					}
				case <-ticker.C:
					if !cf() {
						logger.Infof("Start Cancel channel consume: %s", consumer)
						if err := c.channel.Cancel(consumer, false); err != nil {
							logger.Errorf("Cancel channel error: %v", err)
						}
						return nil, false
					}
				}
			}
			return nil, false
		}

		for {
			if cf() {
				_, reset := consume(hf, cf, periodic)
				time.Sleep(9 * time.Second)
				if reset {
					c.reconnect()
				}
			}
			time.Sleep(9 * time.Second)
		}
	})
}
