package remoteserver

import (
	"encoding/base64"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/go-base-lib/coderutils"
	"github.com/nsqio/go-nsq"
	"gorm.io/gorm"
	"sync"
	"team-client-server/vos"
	"time"
)

var (
	teamworkSm4Key = []byte("teamwork queue message transfer fixed key!!!")[:16]
	queueLock      = sync.Mutex{}
	queueMsgMap    = make(map[string]*nsq.Message)
	requeueDelay   = 5 * time.Second
)

var queueHandler nsq.HandlerFunc = func(message *nsq.Message) error {
	var queueChannelMsgInfo *vos.QueueChannelMsgInfo

	msgId := string(message.ID[:])
	queueLock.Lock()
	defer queueLock.Unlock()
	queueChannelMsgInfoModal := sqlite.Db().Model(&vos.QueueChannelMsgInfo{})

	if err := queueChannelMsgInfoModal.Where("id=? and queue_type=?", msgId, vos.QueueTypeReceive).First(&queueChannelMsgInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
		message.Requeue(requeueDelay)
		return nil
	}

	if queueChannelMsgInfo.Id != "" {
		if !queueChannelMsgInfo.Ack {
			message.Requeue(requeueDelay)
		}
		return nil
	}

	transferMsgBody, err := coderutils.Sm4Decrypt(teamworkSm4Key, message.Body)
	if err != nil {
		return nil
	}

	if transferMsgBody, err = base64.StdEncoding.DecodeString(string(transferMsgBody)); err != nil {
		return nil
	}

	pos := len(transferMsgBody) - 16
	sm4RandomKey := transferMsgBody[pos:]
	transferMsgBody = transferMsgBody[:pos]

	if transferMsgBody, err = coderutils.Sm4Decrypt(sm4RandomKey, transferMsgBody); err != nil {
		message.Finish()
		return nil
	}

	queueMsgMap[msgId] = message

	sendTcpTransfer(&TcpTransferInfo{
		CmdCode: TcpTransferCmdCodeQueue,
		Data:    transferMsgBody,
	})

	return nil

}
