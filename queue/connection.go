package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/goextension"
	"github.com/go-base-lib/logs"
	"github.com/panjf2000/ants/v2"
	"github.com/rabbitmq/amqp091-go"
	"sync"
	"team-client-server/config"
	"time"
)

// Connection Rabbit连接
type Connection struct {
	lock sync.Mutex
	//连接
	Conn *amqp091.Connection
	//通道
	ch *amqp091.Channel
	//连接异常结束
	ConnNotifyClose chan *amqp091.Error
	//通道异常接收
	ChNotifyClose chan *amqp091.Error
	url           string
	config        *amqp091.Config
	//用于关闭进程
	CloseProcess chan bool
	consumeMap   map[string]amqp091.Queue
	isClosed     bool
	//消费者信息
	//RabbitConsumerList []RabbitConsumerInfo
	//生产者信息
	//RabbitProducerMap map[string]string
	//自定义消费者处理函数
	//ConsumeHandle func(<-chan amqp091.Delivery)
}

func (c *Connection) initRabbitMQConsumer(address, username, password, vHost string, isClose bool, waitSuccessChan chan<- error) {
	if isClose {
		c.CloseProcess <- true
	}

	c.config = &amqp091.Config{
		SASL: []amqp091.Authentication{
			&amqp091.AMQPlainAuth{
				Username: username,
				Password: password,
			},
		},
		Vhost:      vHost,
		Properties: amqp091.NewConnectionProperties(),
	}

	c.url = address
	conn, err := amqp091.DialConfig(c.url, *c.config)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	if err != nil {
		logs.Errorf("rabbit连接异常:%s", err.Error())
		//logs.Error("休息5S,开始重连rabbitMq消费者")
		//time.Sleep(5 * time.Second)
		//ants.Submit(func() { c.initRabbitMQConsumer(address, username, password, vHost, false, waitSuccessChan) })
		if waitSuccessChan != nil {
			waitSuccessChan <- fmt.Errorf("队列服务连接失败: %w", err)
		}
		return
	}
	c.Conn = conn
	logs.Info("与rabbitmq建立连接成功")

	ch, err := conn.Channel()
	if err != nil {
		logs.Errorf("rabbitMQ打开通道异常:%s", err.Error())
		if waitSuccessChan != nil {
			waitSuccessChan <- fmt.Errorf("打开队列通道失败: %w", err)
		}
		return
	}
	defer ch.Close()
	c.ch = ch

	if waitSuccessChan != nil {
		waitSuccessChan <- nil
	}

	c.isClosed = false

	c.CloseProcess = make(chan bool, 1)
	c.consumerReConnect(address, username, password, vHost)
	logs.Info("结束消费者旧主进程")
}

// consumerReConnect 消费者重连
func (c *Connection) consumerReConnect(address, username, password, vHost string) {
closeTag:
	for {
		c.ConnNotifyClose = c.Conn.NotifyClose(make(chan *amqp091.Error))
		c.ChNotifyClose = c.ch.NotifyClose(make(chan *amqp091.Error))
		var err *amqp091.Error
		select {
		case err, _ = <-c.ConnNotifyClose:
		case err, _ = <-c.ChNotifyClose:
			if err != nil {
				logs.Errorf("rabbit消费者连接异常:%s", err.Error())
			}
			// 判断连接是否关闭
			if !c.Conn.IsClosed() {
				if err := c.Conn.Close(); err != nil {
					logs.Errorf("rabbit连接关闭异常:%s", err.Error())
				}
			}
			_, isConnChannelOpen := <-c.ConnNotifyClose
			if isConnChannelOpen {
				close(c.ConnNotifyClose)
			}
			ants.Submit(func() {
				c.initRabbitMQConsumer(address, username, password, vHost, false, nil)
			})
			break closeTag
		case <-c.CloseProcess:
			c.isClosed = true
			_ = c.ch.Close()
			_ = c.Conn.Close()
			c.ch = nil
			c.Conn = nil
			break closeTag
		}
	}
	logs.Info("结束消费者旧进程")
}

type MsgInfoMeta map[string]string

type MsgInfoWrapper struct {
	Info    *MsgInfo[MsgInfoMeta]
	RawData goextension.Bytes
}

type MsgHandler func(data *MsgInfoWrapper, delivery amqp091.Delivery)

func (c *Connection) listenQueue(queueName string, args amqp091.Table, handler MsgHandler, notice chan<- error) {
	if c.ch == nil {
		notice <- errors.New("未与队列服务成功创建监听")
	}

	consume, err := c.ch.Consume(queueName, "", false, false, false, false, args)
	if err != nil {
		notice <- fmt.Errorf("创建队列[%s]失败: %w", queueName, err)
	}

	notice <- nil
	for {
		delivery, isOpen := <-consume
		if !isOpen {
			if c.isClosed {
				logs.Infof("队列[%s]退出监听", queueName)
				return
			}
		Retry:
			logs.Infof("队列[%s]异常断开, 5s之后进行重新监听")
			time.Sleep(5 * time.Second)
			noticeChan := make(chan error)
			_ = ants.Submit(func() {
				c.listenQueue(queueName, args, handler, noticeChan)
			})

			if err = <-noticeChan; err != nil {
				goto Retry
			}
			return
		}

		if len(delivery.Body) == 0 {
			_ = Nack(delivery.DeliveryTag, false)
			continue
		}

		body, err := goextension.Bytes(delivery.Body).DecodeBase64()
		if err != nil {
			_ = Nack(delivery.DeliveryTag, false)
			continue
		}

		rawData, err := coderutils.Sm4Decrypt(config.TeamworkSm4Key, body)
		if err != nil || len(rawData) == 0 {
			_ = Nack(delivery.DeliveryTag, false)
			continue
		}

		msgInfo := &MsgInfo[MsgInfoMeta]{}
		if err = json.Unmarshal(rawData, &msgInfo); err != nil {
			_ = Nack(delivery.DeliveryTag, false)
			continue
		}

		_ = ants.Submit(func() {
			handler(&MsgInfoWrapper{
				Info:    msgInfo,
				RawData: rawData,
			}, delivery)
		})
	}
}

func (c *Connection) stop() {
	if c.isClosed || c.CloseProcess == nil {
		return
	}
	c.CloseProcess <- true
}
