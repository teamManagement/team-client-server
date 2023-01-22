package queue

import (
	"errors"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/logs"
	"github.com/panjf2000/ants/v2"
	"github.com/tjfoc/gmsm/sm3"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var (
	consumerConnection *Connection
	lock               sync.Mutex
)

func StartListenerQueue(userId, userPassword, address, vHost string) error {
	lock.Lock()
	defer lock.Unlock()

	userId = "u_" + userId

	if targetPassword, err := coderutils.Hash(sm3.New(), []byte(userPassword)); err != nil {
		logs.Errorf("转换用户的队列密钥失败: %s", err.Error())
		return errors.New("转换队列密钥失败")
	} else {
		userPassword = targetPassword.ToBase64Str()
	}

	consumerConnection = &Connection{}

	waitConnSuccess := make(chan error)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()

	_ = ants.Submit(func() {
		consumerConnection.initRabbitMQConsumer(address, userId, userPassword, vHost, false, waitConnSuccess)
	})

	select {
	case err := <-waitConnSuccess:
		return err
	case <-ctx.Done():
		return errors.New("队列连接超时, 请检查网络")
	}
}

func StopListenerQueue() {
	lock.Lock()
	defer lock.Unlock()

	if consumerConnection == nil {
		return
	}
	consumerConnection.stop()
}
