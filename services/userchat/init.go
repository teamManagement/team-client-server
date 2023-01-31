package userchat

import (
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"strconv"
	"team-client-server/db"
	"team-client-server/remoteserver"
	"team-client-server/tools"
)

var (
	emptyArray []struct{}
)

func InitUserChatWebService(engine *gin.RouterGroup) {
	engine.Group("chat").
		POST("msg/put", ginmiddleware.WrapperResponseHandle(chatMsgPut)).
		POST("msg/query", ginmiddleware.WrapperResponseHandle(chatMsgQuery)).
		POST("msg/query/end/:targetId", ginmiddleware.WrapperResponseHandle(chatMsgQueryEnd)).
		POST("msg/loading/get/idList/:targetId", ginmiddleware.WrapperResponseHandle(chatMsgLoadingGetIdList)).
		POST("msg/query/news/:targetId/:id", ginmiddleware.WrapperResponseHandle(chatMsgQueryNews))
}

var (
	// chatMsgLoadingGetIdList 等待发送完成消息的过滤
	chatMsgLoadingGetIdList ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentUser, err := remoteserver.NowUser()
		if err != nil {
			return nil
		}

		targetId := ctx.Param("targetId")

		var idList []string
		if err := ctx.ShouldBindJSON(&idList); err != nil {
			return err
		}

		var userChatMsgList []db.UserChatMsg
		if err := sqlite.Db().Model(&db.UserChatMsg{}).Select("client_unique_id").Where("client_unique_id in (?) and target_id=? and source_id=? and status='loading'", idList, targetId, currentUser.Id).Find(&userChatMsgList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询正在发送的消息失败: %s", err.Error())
		}

		userChatMsgListLen := len(userChatMsgList)
		if userChatMsgListLen == 0 {
			return userChatMsgList
		}

		idList = make([]string, userChatMsgListLen, userChatMsgListLen)
		for i := range userChatMsgList {
			idList[i] = userChatMsgList[i].ClientUniqueId
		}

		return nil
	}

	// chatMsgPut 聊天消息推送
	chatMsgPut ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var (
			chatMsgInfo *db.UserChatMsg

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

		if chatMsgInfo.MsgType < db.ChatMsgTypeText || chatMsgInfo.MsgType > db.ChatMsgTypeImg {
			return errors.New("不支持的消息内容")
		}

		if chatMsgInfo.ChatType <= db.ChatUnknown || chatMsgInfo.ChatType > db.ChatTypeApp {
			return errors.New("不支持的消息类型")
		}

		if chatMsgInfo.ClientUniqueId == "" {
			return errors.New("客户端唯一索引不能为空")
		}

		//var count int64
		//if err := sqlite.Db().Model(&db.UserChatMsg{}).Where("client_unique_id=?", chatMsgInfo.ClientUniqueId).Count(&count).Error; err != nil {
		//	return fmt.Errorf("查询消息暂存数量失败: %s", err.Error())
		//}
		//
		//if count > 0 {
		//	return errors.New("消息已存在")
		//}

		userChatMsg := &db.UserChatMsg{
			TargetId:       chatMsgInfo.TargetId,
			ChatType:       chatMsgInfo.ChatType,
			MsgType:        chatMsgInfo.MsgType,
			ClientUniqueId: chatMsgInfo.ClientUniqueId,
			Content:        chatMsgInfo.Content,
			Status:         "loading",
		}

		//return sqlite.Db().Transaction(func(tx *gorm.DB) error {
		//	if err = tx.Create(&userChatMsg).Error; err != nil {
		//		return fmt.Errorf("保存聊天信息失败: %s", err.Error())
		//	}
		//
		//	return remoteserver.RequestWebServiceWithData("/user/chat/put", userChatMsg, nil)
		//})
		return remoteserver.RequestWebServiceWithData("/user/chat/put", userChatMsg, nil)

		//return userChatMsg
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

		userChatListWhere := sqlite.Db().Model(&db.UserChatMsg{}).Where("(target_id=? and source_id=?) or (target_id=? and source_id=?)", param.TargetId, currentUser.Id, currentUser.Id, param.TargetId)
		if param.ClientTimeId != "" {
			userChatListWhere = userChatListWhere.Where("client_unique_id > ?", param.ClientTimeId)
		}

		var userChatMsgList []*db.UserChatMsg
		if err := userChatListWhere.Order("client_unique_id").Find(&userChatMsgList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询用户消息列表失败: %s", err.Error())
		}

		return userChatMsgList
	}

	// chatMsgQueryEnd 消息列表查询最后的几条消息
	chatMsgQueryEnd ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		targetId := ctx.Param("targetId")
		queryNumStr := ctx.Query("num")

		endTime := ctx.Query("end")

		//reverse := ctx.Query("reverse")
		//if reverse != "" {
		//	reverse = ""
		//} else {
		//	reverse = "desc"
		//}

		userChatListWhere := sqlite.Db().Model(&db.UserChatMsg{}).Where("((target_id=? and source_id=?) or (source_id=? and target_id=?) or (target_id=? and chat_type=?)) and status='ok'", targetId, currentUser.Id, targetId, currentUser.Id, targetId, db.ChatTypeGroup)
		if endTime != "" {
			//var t time.Time
			//if err = json.Unmarshal([]byte("\""+endTime+"\""), &t); err != nil {
			//	return fmt.Errorf("时间格式错误: %w", err)
			//}

			userChatListWhere = userChatListWhere.Where("time_stamp < ?", endTime)
		}

		userChatListWhere = userChatListWhere.Order("time_stamp desc")

		queryNum, err := strconv.Atoi(queryNumStr)
		if err == nil {
			userChatListWhere = userChatListWhere.Limit(queryNum)
		}

		var userChatMsgList []*db.UserChatMsg
		if err = userChatListWhere.Find(&userChatMsgList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询用户消息列表失败: %s", err.Error())
		}

		return tools.SliceReverse[*db.UserChatMsg](userChatMsgList)
	}

	// chatMsgQueryNews 查询最新的消息
	chatMsgQueryNews ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		targetId := ctx.Param("targetId")
		endId := ctx.Param("id")

		var currentUserChatMsg *db.UserChatMsg
		if err := sqlite.Db().Model(&db.UserChatMsg{}).Select("id, time_stamp, target_id, source_id").Where("id=?", endId, targetId).First(&currentUserChatMsg).Error; err != nil {
			return emptyArray
		}

		if currentUserChatMsg.TargetId != targetId && currentUserChatMsg.SourceId != targetId {
			return emptyArray
		}

		var userChatMsgList []*db.UserChatMsg
		if err := sqlite.Db().Model(&db.UserChatMsg{}).Where("((target_id=? and source_id=?) or (source_id=? and target_id=?) or (target_id=? and chat_type=?)) and status='ok' and time_stamp >= ?",
			targetId, currentUser.Id, targetId, currentUser.Id, targetId, db.ChatTypeGroup, currentUserChatMsg.TimeStamp).Order("time_stamp").Find(&userChatMsgList).Error; err != nil {
			return emptyArray
		}

		index := -1
		for i := range userChatMsgList {
			userChatMsg := userChatMsgList[i]
			if userChatMsg.Id == endId {
				index = i + 1
				break
			}
		}

		if index != -1 {
			userChatMsgList = userChatMsgList[index:]
		}
		return userChatMsgList
	}
)
