package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/goextension"
	"github.com/go-base-lib/logs"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"strings"
	"sync"
	"team-client-server/remoteserver"
	"team-client-server/tools"
	"time"
)

type SocketMessageType uint

const (
	SocketMessageTypePush SocketMessageType = iota
	SocketMessageTypeCallback
)

type SocketMessageContent struct {
	Type    SocketMessageType `json:"type"`
	Data    any               `json:"data,omitempty"`
	Id      string            `json:"id"`
	wrapper *SocketWrapper
}

func (s *SocketMessageContent) Callback(data any) error {
	return s.wrapper.CallbackMessage(s.Id, data)
}

func (s *SocketMessageContent) CallbackAndReceive(data any) (*SocketMessageContent, error) {
	return s.wrapper.CallbackAndReceiveMessage(s.Id, data)
}

type SocketHandlerFn func(content *SocketMessageContent, wrapper *SocketWrapper) error

type socketWriteDataInfo struct {
	JsonData any
	Err      chan error
}

type SocketWrapper struct {
	l    sync.Mutex
	conn *websocket.Conn

	handlerMap      map[string]SocketHandlerFn
	callbackChanMap map[string]chan *SocketMessageContent
	writeMsgChan    chan *socketWriteDataInfo

	// closeRemoteServerRestoreCh 关闭远程服务恢复chan
	closeRemoteServerRestoreCh chan struct{}
}

func (s *SocketWrapper) Close() {
	s.conn.Close()
}

func (s *SocketWrapper) CallbackMessage(id string, sendData any) error {
	s.l.Lock()
	defer s.l.Unlock()

	if s.conn == nil {
		return io.EOF
	}

	return s.conn.WriteJSON(&SocketMessageContent{
		Type: SocketMessageTypeCallback,
		Id:   id,
		Data: sendData,
	})
}

func (s *SocketWrapper) CallbackAndReceiveMessage(id string, sendData any) (*SocketMessageContent, error) {
	s.l.Lock()
	defer s.l.Unlock()

	if s.conn == nil || s.callbackChanMap == nil {
		return nil, io.EOF
	}

	ch := make(chan *SocketMessageContent, 1)
	s.callbackChanMap[id] = ch
	if err := s.conn.WriteJSON(&SocketMessageContent{
		Type: SocketMessageTypeCallback,
		Id:   id,
		Data: sendData,
	}); err != nil {
		return nil, err
	}
	data := <-ch
	if data == nil {
		return nil, io.EOF
	}

	return data, nil
}

func (s *SocketWrapper) registryHandler(id string, fn SocketHandlerFn) {
	s.handlerMap[id] = fn
}

func (s *SocketWrapper) checkRemoteServerLoop() {
	s.l.Lock()
	defer s.l.Unlock()

	if s.closeRemoteServerRestoreCh != nil {
		return
	}

	s.closeRemoteServerRestoreCh = make(chan struct{}, 1)

	go func() {
		defer func() { recover() }()
		for {
			restoreTime := time.After(5 * time.Second)
			select {
			case <-s.closeRemoteServerRestoreCh:
				close(s.closeRemoteServerRestoreCh)
				return
			case <-restoreTime:
				if !tools.TelnetHost(remoteserver.ServerAddress) {
					continue
				}

				cmdCode := remoteserver.TcpTransferCmdCodeRestoreServerConnErr
				if remoteserver.AutoLogin() {
					cmdCode = remoteserver.TcpTransferCmdCodeRestoreServerConnOK
				}

				tcpTransferMarshal, _ := json.Marshal(&remoteserver.TcpTransferInfo{
					CmdCode: cmdCode,
				})

				writeData := &socketWriteDataInfo{
					JsonData: &SocketMessageContent{
						Type: SocketMessageTypePush,
						Data: base64.StdEncoding.EncodeToString(tcpTransferMarshal),
					},
					Err: make(chan error),
				}

				s.writeMsgChan <- writeData
				<-writeData.Err
				close(s.closeRemoteServerRestoreCh)
				s.closeRemoteServerRestoreCh = nil

				return
			}
		}
	}()
}

func (s *SocketWrapper) loop() {
	defer s.conn.Close()
	defer func() {
		s.l.Lock()
		defer s.l.Unlock()

		for k := range s.callbackChanMap {
			s.callbackChanMap[k] <- nil
		}

		s.callbackChanMap = nil
		s.handlerMap = nil
		remoteserver.Logout()
		close(s.writeMsgChan)

		if s.closeRemoteServerRestoreCh != nil {
			s.closeRemoteServerRestoreCh <- struct{}{}
			<-s.closeRemoteServerRestoreCh
			s.closeRemoteServerRestoreCh = nil
		}

	}()

	go func() {
		remoteserver.StartTcpTransfer()
		transfer := remoteserver.GetTcpTransfer()
		for {
			select {
			case serverMsg := <-transfer:
				logs.Debugf("接收到TCP服务向客户端转发的数据内容, 将要想客户端进行推送...")

				if serverMsg.CmdCode == remoteserver.TcpTransferCmdCodeBlockingConnection {
					logs.Debugf("检测到TCP服务已断开, 开启远程TCP服务重连")
					s.checkRemoteServerLoop()
				}

				marshal, _ := json.Marshal(serverMsg)
				_ = s.conn.WriteJSON(&SocketMessageContent{
					Type: SocketMessageTypePush,
					Data: base64.StdEncoding.EncodeToString(marshal),
				})
			case writeData, isOpen := <-s.writeMsgChan:
				if !isOpen {
					remoteserver.StopTcpTransfer()
					return
				}

				writeData.Err <- s.conn.WriteJSON(writeData.JsonData)
			}

		}
	}()

	for {
		var messageContent *SocketMessageContent
		if err := s.conn.ReadJSON(&messageContent); err != nil {
			return
		}

		messageContent.wrapper = s

		switch messageContent.Type {
		case SocketMessageTypePush:
			id := messageContent.Id
			index := strings.Index(id, ":")
			if index == -1 {
				return
			}
			cmd := id[:index]
			messageContent.Id = id[index+1:]

			fn, ok := s.handlerMap[cmd]
			if !ok {
				return
			}
			go func() {
				if err := fn(messageContent, s); err != nil {
					_ = s.conn.Close()
				}
			}()
		case SocketMessageTypeCallback:
			ch, ok := s.callbackChanMap[messageContent.Id]
			if ok {
				delete(s.callbackChanMap, messageContent.Id)
				ch <- messageContent
			}
		}
	}
}

func createSocketWrapper(conn *websocket.Conn) *SocketWrapper {
	return &SocketWrapper{conn: conn, handlerMap: make(map[string]SocketHandlerFn), callbackChanMap: make(map[string]chan *SocketMessageContent), writeMsgChan: make(chan *socketWriteDataInfo, 1)}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func initWs(engine *gin.Engine) {
	engine.Any("/ws", func(context *gin.Context) {
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		wrapper := createSocketWrapper(conn)
		wrapper.registryHandler("conn", socketHandlerConn)
		wrapper.registryHandler("checkIsLogin", socketHandlerCheckIsLogin)
		wrapper.registryHandler("login", socketHandlerLogin)
		wrapper.registryHandler("autoLogin", socketHandlerAutoLogin)
		wrapper.loop()
	})
}

var (
	socketHandlerAutoLogin SocketHandlerFn = func(content *SocketMessageContent, wrapper *SocketWrapper) (err error) {
		return content.Callback(remoteserver.AutoLogin())
	}

	socketHandlerLogin SocketHandlerFn = func(content *SocketMessageContent, wrapper *SocketWrapper) (err error) {
		res := make(map[string]interface{})
		defer func() {
			if err != nil {
				res["error"] = true
				res["message"] = err.Error()
			}
			err = content.Callback(res)
		}()
		usernameAndPassword, ok := content.Data.(string)
		if !ok {
			return fmt.Errorf("错误的数据包格式")
		}

		split := strings.Split(usernameAndPassword, ".")
		if len(split) != 2 {
			return fmt.Errorf("错误的用户数据组包格式")
		}

		return remoteserver.Login(split[0], split[1])
	}

	socketHandlerCheckIsLogin SocketHandlerFn = func(content *SocketMessageContent, wrapper *SocketWrapper) error {
		return content.Callback(remoteserver.LoginOk())
	}
	socketHandlerConn SocketHandlerFn = func(content *SocketMessageContent, wrapper *SocketWrapper) error {
		clientEncRandom, ok := content.Data.(string)
		if !ok {
			return fmt.Errorf("类型转换失败")
		}

		if len(clientEncRandom) < 1024 {
			return fmt.Errorf("数据长过短")
		}

		encAesKey, err := goextension.Bytes(clientEncRandom[:1024]).DecodeHex()
		if err != nil {
			return err
		}

		aesKey, err := coderutils.RsaDecrypt(encAesKey, tools.ClientServerPrivateKey)
		if err != nil {
			return err
		}

		aesKey, err = aesKey.DecodeHex()
		if err != nil {
			return err
		}

		cipherText, err := goextension.Bytes(clientEncRandom[1024:]).DecodeHex()
		clientRandom, err := coderutils.AesCBCDecrypt(cipherText, aesKey)
		if err != nil {
			return err
		}

		serverRandomStr := coderutils.GetRandomString(16)

		clientMsg, err := content.CallbackAndReceive(fmt.Sprintf("%s%s", clientRandom, serverRandomStr))
		if err != nil {
			return err
		}

		if clientMsg.Data != serverRandomStr {
			return fmt.Errorf("握手失败")
		}

		return clientMsg.Callback("ok")
	}
)
