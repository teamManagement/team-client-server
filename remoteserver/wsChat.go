package remoteserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"net/http"
	"sync"
	"team-client-server/config"
	"team-client-server/db"
)

type ChatMsgCmd uint

const (
	// ChatMsgCmdPut 推送消息
	ChatMsgCmdPut ChatMsgCmd = iota + 1
	// ChatMsgCmdResponse 消息的响应
	ChatMsgCmdResponse
	// ChatMsgCmdGetList 获取消息列表
	ChatMsgCmdGetList
)

type ChatMsgQueryParam struct {
	// TimeId 查询该Id之后的所有聊天结果
	TimeId string `json:"timeId,omitempty"`
}

type ChatMsgWrapper struct {
	// Cmd 命令码
	Cmd ChatMsgCmd `json:"cmd,omitempty"`
	// ErrMsg 错误消息
	ErrMsg string `json:"errMsg,omitempty"`
	// ChatData 单个聊天消息
	ChatData *db.UserChatMsg `json:"chatData,omitempty"`
	// ChatListData 聊天消息列表
	ChatListData []*db.UserChatMsg `json:"chatListData,omitempty"`
	// QueryParam 消息查询参数
	QueryParam ChatMsgQueryParam `json:"queryParam,omitempty"`
	// Error 是否有错
	Error bool `json:"error,omitempty"`
	// ErrCode 错误代码
	ErrCode string `json:"errCode,omitempty"`
}

var (
	wsChatConn         *websocket.Conn
	wsChatHttpResponse *http.Response
	wsChatLock         sync.Mutex
	wsChatResponseChan chan *ChatMsgWrapper
)

func startChatWs() (err error) {
	wsChatLock.Lock()
	defer wsChatLock.Unlock()
	stopWsChat()

	defer func() {
		if err != nil {
			stopWsChat()
		}
	}()

	wsChatResponseChan = make(chan *ChatMsgWrapper)

	dialer := &websocket.Dialer{
		TLSClientConfig: config.HttpsTlsConfig,
	}
	wsChatConn, wsChatHttpResponse, err = dialer.Dial(config.LocalWSServerAddress+"/ws/chat", http.Header{
		"_t":         []string{Token()},
		"_a":         []string{LoginIp()},
		"User-Agent": []string{"teamManageLocal"},
	})

	if err != nil {
		return fmt.Errorf("连接用户即时通讯通道失败: %s", err.Error())
	}

	go chatWsLoop()

	return UserChatFlushLocalByRemoteServer()

}

func chatWsLoop() {
	defer stopWsChat()
	defer func() { recover() }()

	for {

		var chatMsgWrapper *ChatMsgWrapper
		if err := wsChatConn.ReadJSON(&chatMsgWrapper); err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		switch chatMsgWrapper.Cmd {
		case ChatMsgCmdResponse:
			wsChatResponseChan <- chatMsgWrapper
		case ChatMsgCmdPut:
			chatData := chatMsgWrapper.ChatData
			if chatData.ClientUniqueId == "" {
				continue
			}

			var count int64
			if err := sqlite.Db().Model(&db.UserChatMsg{}).Where("client_unique_id=?", chatData.ClientUniqueId).Count(&count).Error; err != nil {
				continue
			}

			if count > 0 {
				if err := sqlite.Db().Model(&db.UserChatMsg{}).Where("client_unique_id=?", chatData.ClientUniqueId).UpdateColumns(&chatData).Error; err != nil {
					continue
				}
			} else {
				if err := sqlite.Db().Model(&db.UserChatMsg{}).Create(&chatData).Error; err != nil {
					continue
				}
			}

			marshal, _ := json.Marshal(chatData)

			sendTcpTransfer(&TcpTransferInfo{
				CmdCode: TcpTransferCmdCodeChatMsgChange,
				Data:    marshal,
			})
		}

		continue
	}
}

func UserChatFlushLocalByRemoteServer() (err error) {
	var (
		chatMsgList   []*db.UserChatMsg
		chatEndTimeId string
	)

	chatMsgGlobalTimeIdSetting := &db.Setting{
		Name: "chat_msg_global_time_id",
	}
	if err = sqlite.Db().Model(&db.Setting{}).Where(&chatMsgGlobalTimeIdSetting).Find(&chatMsgGlobalTimeIdSetting).Error; err != nil || chatMsgGlobalTimeIdSetting.Value == "" {
		chatMsgList, chatEndTimeId, err = UserChatQueryAll()
	} else {
		chatMsgList, chatEndTimeId, err = UserChatQueryTimeIdAfter(string(chatMsgGlobalTimeIdSetting.Value))
	}

	if err != nil {
		return err
	}

	return sqlite.Db().Transaction(func(tx *gorm.DB) error {
		if chatEndTimeId != "" {
			chatMsgGlobalTimeIdSetting.Value = db.EncryptValue(chatEndTimeId)
			if err = tx.Save(&chatMsgGlobalTimeIdSetting).Error; err != nil {
				return fmt.Errorf("更新消息标识失败: %s", err.Error())
			}
		}

		for i := range chatMsgList {
			chatMsg := chatMsgList[i]
			if chatMsg.ClientUniqueId == "" {
				continue
			}

			if err = tx.Save(&chatMsg).Error; err != nil {
				return fmt.Errorf("更新本地消息失败: %s", err.Error())
			}
		}

		return nil
	})
}

func userChatWriteWrapperDataAndResponse(wrapper *ChatMsgWrapper) (*ChatMsgWrapper, error) {
	if wsChatConn == nil {
		return nil, errors.New("用户未登录")
	}

	if err := wsChatConn.WriteJSON(wrapper); err != nil {
		return nil, fmt.Errorf("消息查询失败: %s", err.Error())
	}

	res, isOpen := <-wsChatResponseChan
	if !isOpen {
		return nil, errors.New("用户登录凭证异常")
	}

	return res, nil
}

func UserChatQueryAll() ([]*db.UserChatMsg, string, error) {
	res, err := userChatWriteWrapperDataAndResponse(&ChatMsgWrapper{
		Cmd: ChatMsgCmdGetList,
	})
	if err != nil {
		return nil, "", err
	}

	return res.ChatListData, res.QueryParam.TimeId, err
}

func UserChatQueryTimeIdAfter(timeId string) ([]*db.UserChatMsg, string, error) {

	res, err := userChatWriteWrapperDataAndResponse(&ChatMsgWrapper{
		Cmd: ChatMsgCmdGetList,
		QueryParam: ChatMsgQueryParam{
			TimeId: timeId,
		},
	})

	if err != nil {
		return nil, "", err
	}

	return res.ChatListData, res.QueryParam.TimeId, nil
}

func UserChatPut(chatData *db.UserChatMsg) (*db.UserChatMsg, error) {

	res, err := userChatWriteWrapperDataAndResponse(&ChatMsgWrapper{
		Cmd:      ChatMsgCmdPut,
		ChatData: chatData,
	})

	if err != nil {
		return nil, err
	}

	return res.ChatData, nil
}

func stopWsChat() {
	if wsChatLock.TryLock() {
		defer wsChatLock.Unlock()
	}
	if wsChatHttpResponse != nil {
		_ = wsChatHttpResponse.Body.Close()
	}

	if wsChatConn != nil {
		_ = wsChatConn.Close()
	}
	wsChatCloseChan()
	wsChatHttpResponse = nil
	wsChatConn = nil
}

func wsChatCloseChan() {
	defer func() { recover() }()
	if wsChatResponseChan == nil {
		return
	}
	close(wsChatResponseChan)
}
