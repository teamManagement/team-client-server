package remoteserver

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"sync"
	"team-client-server/vos"
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
}

type ChatMsgWrapper struct {
	// Cmd 命令码
	Cmd ChatMsgCmd `json:"cmd,omitempty"`
	// ErrMsg 错误消息
	ErrMsg string `json:"errMsg,omitempty"`
	// ChatData 单个聊天消息
	ChatData *vos.UserChatMsg `json:"chatData,omitempty"`
	// ChatListData 聊天消息列表
	ChatListData []*vos.UserChatMsg `json:"chatListData,omitempty"`
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

	dialer := &websocket.Dialer{}
	wsChatConn, wsChatHttpResponse, err = dialer.Dial(LocalWSServerAddress+"/ws/chat", http.Header{
		"_t":         []string{Token()},
		"User-Agent": []string{"teamManageLocal"},
	})

	if err != nil {
		return fmt.Errorf("连接用户即时通讯通道失败: %s", err.Error())
	}

	go chatWsLoop()
	return nil
}

func chatWsLoop() {
	defer stopWsChat()

	for {

		var chatMsgWrapper *ChatMsgWrapper
		if err := wsChatConn.ReadJSON(&chatMsgWrapper); err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		if chatMsgWrapper.Cmd == ChatMsgCmdResponse {
			wsChatResponseChan <- chatMsgWrapper
		}

		continue
	}
}

func UserChatPut(chatData *vos.UserChatMsg) (*vos.UserChatMsg, error) {
	if wsChatConn == nil {
		return nil, errors.New("用户未登录")
	}

	currentUser, err := NowUser()
	if err != nil {
		return nil, err
	}

	chatData.SourceId = currentUser.Id

	if err = wsChatConn.WriteJSON(chatData); err != nil {
		return nil, fmt.Errorf("消息发送失败: %s", err.Error())
	}

	res, isOpen := <-wsChatResponseChan
	if !isOpen {
		return nil, errors.New("用户登录凭证异常")
	}

	if res.Error {
		return nil, errors.New(res.ErrMsg)
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
