package userchat

import (
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

func InitUserChatWebService(engine *gin.RouterGroup) {
	engine.Group("chat").
		POST("msg/put", ginmiddleware.WrapperResponseHandle(chatMsgPut)).
		POST("msg/query", ginmiddleware.WrapperResponseHandle(chatMsgQuery))
}

var (
	// chatMsgPut 聊天消息推送
	chatMsgPut ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var (
			chatMsgInfo *vos.UserChatMsg

			err error
		)
		if err = ctx.ShouldBindJSON(&chatMsgInfo); err != nil {
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

		if chatMsgInfo.ClientUniqueId == "" {
			return errors.New("客户端唯一索引不能为空")
		}

		var count int64
		if err := sqlite.Db().Model(&vos.UserChatMsg{}).Where("client_unique_id=?", chatMsgInfo.ClientUniqueId).Count(&count).Error; err != nil {
			return fmt.Errorf("查询消息暂存数量失败: %s", err.Error())
		}

		if count > 0 {
			return errors.New("消息已存在")
		}

		userChatMsg := &vos.UserChatMsg{
			TargetId:       chatMsgInfo.TargetId,
			ChatType:       chatMsgInfo.ChatType,
			MsgType:        chatMsgInfo.MsgType,
			ClientUniqueId: chatMsgInfo.ClientUniqueId,
			Content:        chatMsgInfo.Content,
		}

		if userChatMsg, err = remoteserver.UserChatPut(userChatMsg); err != nil {
			return err
		}

		sqlite.Db().Model(&userChatMsg).Create(userChatMsg)

		return userChatMsg
	}

	// chatMsgQuery 消息列表查询
	chatMsgQuery ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}
		var param *QueryParam
		if err := ctx.ShouldBindJSON(&param); err != nil {
			return fmt.Errorf("解析请求参数失败: %s", err.Error())
		}

		if param.TargetId == "" {
			return errors.New("接收对象ID不能为空")
		}

		userChatListWhere := sqlite.Db().Model(&vos.UserChatMsg{}).Where("(target_id=? and source_id=?) or (target_id=? and source_id=?)", param.TargetId, currentUser.Id, currentUser.Id, param.TargetId)
		if param.ClientTimeId != "" {
			userChatListWhere = userChatListWhere.Where("client_unique_id > ?", param.ClientTimeId)
		}

		var userChatMsgList []*vos.UserChatMsg
		if err := userChatListWhere.Order("client_unique_id").Find(&userChatMsgList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询用户消息列表失败: %s", err.Error())
		}

		return userChatMsgList
	}
)
