package userchat

type QueryParam struct {
	// TargetId 目标对象ID
	TargetId string `json:"targetId,omitempty"`
	// ClientTimeId 客户端时间ID结束时间
	ClientTimeId string `json:"clientTimeId,omitempty"`
}
