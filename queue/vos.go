package queue

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageType uint16

const (
	MessageTypeStart MessageType = iota
	// MessageTypeChatMsg 聊天消息
	MessageTypeChatMsg
	// MessageTypeApplicationMsg 应用消息
	MessageTypeApplicationMsg
	MessageTypeEnd
)

type MessageInfo[T any] struct {
	Id string `json:"id,omitempty"`
	// Type 消息类型
	Type MessageType `json:"type,omitempty"`
	// Content 消息内容
	Content T `json:"content,omitempty"`
	// TargetId 目标ID
	TargetId string `json:"targetId,omitempty"`
	// SenderId 发送者ID
	SenderId string `json:"senderId,omitempty"`
}

type MsgType int

const (
	// MsgTypeChatPutConfirm 消息推送确认
	MsgTypeChatPutConfirm MsgType = iota + 1
	// MsgTypeChatSendOut 聊天消息下发
	MsgTypeChatSendOut
)

type MsgInfo[M any] struct {
	// Type 类型
	Type MsgType `json:"type,omitempty"`
	// Content 内容
	Content []byte `json:"content,omitempty"`
	// Meta 附加属性
	Meta M `json:"meta,omitempty"`
	// SendTime 发送时间
	SendTime time.Time
}

func (m *MsgInfo[M]) BindContent(i any) error {
	if err := json.Unmarshal(m.Content, &i); err != nil {
		return fmt.Errorf("转换队列内容失败: %s", err.Error())
	}

	return nil
}