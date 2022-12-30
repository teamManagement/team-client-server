package queue

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
