package remoteserver

const (
	// remoteModelCodeLogin 登录
	remoteModelCodeLogin byte = iota + 1
	// remoteModeCodeLoginWithToken 通过token进行登录
	remoteModeCodeLoginWithToken
	// remoteModelCodeRefresh 刷新token
	remoteModelCodeRefresh
	// remoteModelCodeLogout 退出登录
	remoteModelCodeLogout
	// remoteModelCodeOtherLoginNotify 他处登录
	remoteModelCodeOtherLoginNotify
	// remoteModelCodeLoginOk 登录成功
	remoteModelCodeLoginOk
	// remoteModelCodeOtherUserStatusChange 其他用户的状态变更, 返回数据, id+"__"+('online' | 'offline')
	remoteModelCodeOtherUserStatusChange
	// remoteModelCodeUpdateCache 更新缓存
	remoteModelCodeUpdateCache
)
