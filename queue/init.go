package queue

import (
	"github.com/nsqio/go-nsq"
)

var consumer *nsq.Consumer

func StartListenerQueue(userId string, loginIp string, nsqdLookupUrl string, handlerFunc nsq.HandlerFunc) (err error) {
	cfg := nsq.NewConfig()
	consumer, err = nsq.NewConsumer("u"+userId, loginIp, cfg)
	if err != nil {
		return err
	}

	consumer.AddHandler(handlerFunc)

	if err = consumer.ConnectToNSQLookupds([]string{nsqdLookupUrl}); err != nil {
		return err
	}

	return nil
}

func StopListenerQueue() {
	if consumer == nil {
		return
	}
	defer func() {
		recover()
		consumer = nil
	}()
	consumer.Stop()
}
