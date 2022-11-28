package userchat

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

func InitUserChatWebService(engine *gin.RouterGroup) {
	engine.Group("chat").
		POST("msg/put", ginmiddleware.WrapperResponseHandle(chatMsgPut))
}

type ChatMsgInfo struct {
	ChatType vos.ChatType    `json:"type,omitempty"`
	MsgType  vos.ChatMsgType `json:"msgType,omitempty"`
	TargetId string          `json:"targetId,omitempty"`
	Content  string          `json:"content,omitempty"`
}

var (
	// chatMsgPut 聊天消息推送
	chatMsgPut ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var (
			chatMsgInfo *ChatMsgInfo

			err error
		)
		if err = ctx.ShouldBindJSON(chatMsgInfo); err != nil {
			return fmt.Errorf("解析消息内容失败: %s", err.Error())
		}

		if chatMsgInfo.Content == "" {
			return errors.New("消息内容不能为空")
		}

		if chatMsgInfo.TargetId == "" {
			return errors.New("接收者ID不能为空")
		}

		if chatMsgInfo.MsgType < vos.ChatMsgTypeText || chatMsgInfo.MsgType > vos.ChatMsgTypeImg {
			return errors.New("不支持的消息内容")
		}

		if chatMsgInfo.ChatType <= vos.ChatUnknown || chatMsgInfo.ChatType > vos.ChatTypeApp {
			return errors.New("不支持的消息类型")
		}

		userChatMsg := &vos.UserChatMsg{
			TargetId: chatMsgInfo.TargetId,
			ChatType: chatMsgInfo.ChatType,
			MsgType:  chatMsgInfo.MsgType,
		}

		if chatMsgInfo.MsgType == vos.ChatMsgTypeText {
			userChatMsg.ContentText = chatMsgInfo.Content
		} else {
			userChatMsg.ContentFileId = chatMsgInfo.Content
		}

		if userChatMsg, err = remoteserver.UserChatPut(userChatMsg); err != nil {
			return err
		}

		return nil
	}
)
