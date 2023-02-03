package queue

import (
	"errors"
	"github.com/rabbitmq/amqp091-go"
)

func ListenQueue(queueName string, args amqp091.Table, handler MsgHandler, notice chan<- error) {
	if consumerConnection == nil {
		notice <- errors.New("队列连接未被初始化")
	}

	consumerConnection.listenQueue(queueName, args, handler, notice)
}

func Ack(tag uint64) error {
	if consumerConnection == nil {
		return errors.New("队列连接未被初始化")
	}
	return consumerConnection.ch.Ack(tag, false)
}

func Nack(tag uint64, requeue bool) error {
	if consumerConnection == nil {
		return errors.New("队列连接未被初始化")
	}
	return consumerConnection.ch.Nack(tag, false, requeue)
}
