package channel

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"team-client-server/queue"
	"team-client-server/remoteserver"
	"team-client-server/sfnake"
	"team-client-server/vos"
	"time"
)

func InitChannelWebServices(engine *gin.RouterGroup) {
	channelGroup := engine.Group("/channel")

	{
		channelGroup.Group("broadcast").
			POST("appMsg", ginmiddleware.WrapperResponseHandle(channelBroadcastAppMsg))
	}
}

type PushResult struct {
	UserId  string `json:"userId,omitempty"`
	Success bool   `json:"success,omitempty"`
	ErrMsg  string `json:"errMsg,omitempty"`
}

type BroadcastAppMsgInfo struct {
	// TargetUserId 目标用户ID
	TargetUserId []string `json:"targetUserId,omitempty"`
	// AppId 应用ID
	AppId string `json:"appId,omitempty"`
	// Content 内容
	Content any `json:"content,omitempty"`
}

var (
	// channelBroadcastAppMsg 广播应用消息
	channelBroadcastAppMsg ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		var param *BroadcastAppMsgInfo
		if err = ctx.ShouldBindJSON(&param); err != nil {
			return fmt.Errorf("解析请求数据失败: %w", err)
		}

		targetUserLen := len(param.TargetUserId)
		if targetUserLen == 0 {
			return errors.New("目标用户不能为空")
		}

		if param.AppId == "" {
			return errors.New("应用ID不能为空")
		}

		var count int64
		if err = sqlite.Db().Model(&vos.Application{}).Where("id=? and user_id=? and status='4'", param.AppId, currentUser.Id).Count(&count).Error; err != nil {
			return fmt.Errorf("查询用户与应用关系失败: %w", err)
		}

		if count != 1 {
			return errors.New("用户无使用次应用的权限")
		}

		msgData := &queue.MessageInfo[any]{
			SenderId: currentUser.Id,
			TargetId: param.AppId,
			Type:     queue.MessageTypeApplicationMsg,
			Content:  param.Content,
		}

		contentMarshal, err := json.Marshal(param.Content)
		if err != nil {
			return fmt.Errorf("序列化消息内容失败: %s", err.Error())
		}

		channelMsgAckInfo := &vos.QueueChannelMsgInfo{
			Content:  base64.StdEncoding.EncodeToString(contentMarshal),
			Ack:      false,
			SendTime: time.Now().UnixNano(),
			Type:     queue.MessageTypeApplicationMsg,
			AppId:    param.AppId,
		}

		result := make([]*PushResult, targetUserLen, targetUserLen)
		for i := range param.TargetUserId {
			targetUserId := param.TargetUserId[i]
			idStr, err := sfnake.GetIdStr()
			if err != nil {
				return fmt.Errorf("生成数据ID失败: %w", err)
			}
			msgData.Id = currentUser.Id + "_" + targetUserId + "_" + idStr

			if err = sqlite.Db().Transaction(func(tx *gorm.DB) error {
				channelMsgAckInfo.Id = msgData.Id
				if err = tx.Create(&channelMsgAckInfo).Error; err != nil {
					return fmt.Errorf("保存消息签收信息失败: %w", err)
				}

				if err = remoteserver.RequestWebServiceWithData("/channel/push/"+targetUserId, msgData, nil); err != nil {
					return err
				}

				result[i] = &PushResult{
					Success: true,
					UserId:  targetUserId,
				}
				return nil
			}); err != nil {
				result[i] = &PushResult{
					Success: false,
					ErrMsg:  err.Error(),
					UserId:  targetUserId,
				}
			}

		}

		return result
	}
)
