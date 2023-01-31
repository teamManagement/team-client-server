package remoteserver

import (
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/nsqio/go-nsq"
	"github.com/rabbitmq/amqp091-go"
	"sync"
	"team-client-server/db"
	"team-client-server/queue"
	"time"
)

var (
	teamworkSm4Key = []byte("teamwork queue message transfer fixed key!!!")[:16]
	queueLock      = sync.Mutex{}
	queueMsgMap    = make(map[string]*nsq.Message)
	requeueDelay   = 5 * time.Second
)

var queueHandler queue.MsgHandler = func(data *queue.MsgInfoWrapper, delivery amqp091.Delivery) {
	queueLock.Lock()
	defer queueLock.Unlock()

	msgInfo := data.Info

	if msgInfo.Type == queue.MsgTypeChatPutConfirm || msgInfo.Type == queue.MsgTypeChatSendOut {
		var userChatMsg *db.UserChatMsg
		if err := msgInfo.BindContent(&userChatMsg); err != nil {
			_ = queue.Nack(delivery.DeliveryTag, false)
			return
		}

		if userChatMsg.ClientUniqueId == "" || userChatMsg.Id == "" ||
			userChatMsg.TargetId == "" || userChatMsg.SourceId == "" ||
			userChatMsg.ChatType <= db.ChatUnknown || userChatMsg.ChatType > db.ChatTypeApp ||
			userChatMsg.MsgType < db.ChatMsgTypeText || userChatMsg.MsgType > db.ChatMsgTypeImg {
			_ = queue.Nack(delivery.DeliveryTag, false)
			return
		}

		userChatMsg.Status = "ok"
		if err := sqlite.Db().Model(&db.UserChatMsg{}).Where("client_unique_id=?", userChatMsg.ClientUniqueId).Save(userChatMsg).Error; err != nil {
			_ = queue.Nack(delivery.DeliveryTag, false)
			return
		}

	} else {
		_ = queue.Nack(delivery.DeliveryTag, false)
		return
	}

	_ = queue.Ack(delivery.DeliveryTag)

	sendTcpTransfer(&TcpTransferInfo{
		CmdCode: TcpTransferCmdCodeQueue,
		Data:    data.RawData,
	})
	//var queueChannelMsgInfo *db.QueueChannelMsgInfo
	//
	//msgId := string(message.ID[:])
	//queueLock.Lock()
	//defer queueLock.Unlock()
	//queueChannelMsgInfoModal := sqlite.Db().Model(&db.QueueChannelMsgInfo{})
	//
	//if err := queueChannelMsgInfoModal.Where("id=? and queue_type=?", msgId, db.QueueTypeReceive).First(&queueChannelMsgInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
	//	message.Requeue(requeueDelay)
	//	return nil
	//}
	//
	//if queueChannelMsgInfo.Id != "" {
	//	if !queueChannelMsgInfo.Ack {
	//		message.Requeue(requeueDelay)
	//	}
	//	return nil
	//}
	//
	//transferMsgBody, err := coderutils.Sm4Decrypt(teamworkSm4Key, message.Body)
	//if err != nil {
	//	return nil
	//}
	//
	//if transferMsgBody, err = base64.StdEncoding.DecodeString(string(transferMsgBody)); err != nil {
	//	return nil
	//}
	//
	//pos := len(transferMsgBody) - 16
	//sm4RandomKey := transferMsgBody[pos:]
	//transferMsgBody = transferMsgBody[:pos]
	//
	//if transferMsgBody, err = coderutils.Sm4Decrypt(sm4RandomKey, transferMsgBody); err != nil {
	//	message.Finish()
	//	return nil
	//}
	//
	//queueMsgMap[msgId] = message
	//
	//sendTcpTransfer(&TcpTransferInfo{
	//	CmdCode: TcpTransferCmdCodeQueue,ju89
	//	Data:    transferMsgBody,
	//})
	//
	//return nil

}
