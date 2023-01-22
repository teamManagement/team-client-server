package db

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/goextension"
	"team-client-server/config"
	"team-client-server/queue"
	"team-client-server/vos"
	"time"
)

const encryptValCipherKeyLen = 512

type EncryptValue string

func (e EncryptValue) Value() (driver.Value, error) {
	key := coderutils.Sm4RandomKey()

	cipherData, err := coderutils.Sm4Encrypt(key, []byte(e))
	if err != nil {
		return nil, err
	}

	encryptKey, err := coderutils.RsaEncrypt(key, config.ClientPublicKey)
	if err != nil {
		return nil, err
	}

	return goextension.Bytes(bytes.Join([][]byte{
		encryptKey, cipherData,
	}, nil)).ToBase64Str(), nil

}

func (e *EncryptValue) Scan(src any) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", src))
	}

	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.New("Failed to decode base64: " + str)
	}

	if len(bytes) < encryptValCipherKeyLen {
		return errors.New("failed data length")
	}

	encryptSm4Key := bytes[:encryptValCipherKeyLen]

	sm4Key, err := coderutils.RsaDecrypt(encryptSm4Key, config.ClientPrivateKey)
	if err != nil {
		return err
	}

	originData, err := coderutils.Sm4Decrypt(sm4Key, bytes[encryptValCipherKeyLen:])
	if err != nil {
		return err
	}

	*e = EncryptValue(originData.ToString())
	return nil
}

// Setting 加密配置
type Setting struct {
	Name  string       `json:"name,omitempty" gorm:"primary_key"`
	Value EncryptValue `json:"value,omitempty"`
}

type ApplicationType uint

const (
	// ApplicationTypeRemoteWeb 远程web
	ApplicationTypeRemoteWeb ApplicationType = iota
	// ApplicationTypeLocalWeb 本地web
	ApplicationTypeLocalWeb
)

type ApplicationStatus uint

const (
	// ApplicationStatusTakeDown 下架
	ApplicationStatusTakeDown ApplicationStatus = iota
	// ApplicationStatusTakeAudit 审核中
	ApplicationStatusTakeAudit
	// ApplicationStatusTakeReject 审核拒绝
	ApplicationStatusTakeReject
	// ApplicationStatusRectification 整改中
	ApplicationStatusRectification
	// ApplicationStatusNormal 正常
	ApplicationStatusNormal
)

type IconType uint

const (
	// IconTypeUrl 图片URL, 使用img渲染
	IconTypeUrl IconType = iota
	// IconTypeIconfont iconFont图标名称, 使用className渲染
	IconTypeIconfont
	// IconTypeInsideTd 使用腾讯的ui组件库中的icon
	IconTypeInsideTd
	// IconTypeInsideAnt 使用antd组件库中的icon
	IconTypeInsideAnt
)

// Application 应用信息
type Application struct {
	// Id 应用ID
	Id string `json:"id,omitempty" gorm:"primaryKey"`
	// Name 应用名称
	Name string `json:"name,omitempty"`
	// CategoryId 应用类别ID
	CategoryId string `json:"categoryId,omitempty"`
	// AuthorId 作者ID或者为主导者ID
	AuthorId string `json:"authorId,omitempty"`
	// Inside 是否为平台内部应用
	Inside bool `json:"inside,omitempty"`
	// Type 应用类型
	Type ApplicationType `json:"type"`
	// RemoteSiteUrl 远程站点地址, 当 Type 为 ApplicationTypeRemoteWeb 时必填
	RemoteSiteUrl string `json:"remoteSiteUrl,omitempty" gorm:"type:longtext"`
	// Url 访问url
	Url string `json:"url,omitempty" gorm:"type:longtext"`
	// BackgroundFileDir 背景文件路径
	BackgroundFileDir string `json:"backgroundFileDir,omitempty"`
	// BackgroundFileHash 背景文件的hash, 格式为json格式的字符串, {"文件名": "hash hex"}
	BackgroundFileHash string `json:"-"`
	// LocalFileHash 本地文件HASH
	LocalFileHash []byte `json:"localFileHash,omitempty"`
	// Icon 图标地址
	Icon string `json:"icon,omitempty" gorm:"type:longtext"`
	// IconType 图标类型
	IconType IconType `json:"iconType" gorm:"type:longtext"`
	// Desc 描述, 支持markdown
	Desc string `json:"desc,omitempty" gorm:"type:longtext"`
	// ShortDesc 短描述最多100个字
	ShortDesc string `json:"shortDesc,omitempty" gorm:"type:varchar(100)"`
	// Slideshow 轮播图, 最多支持9张图片
	Slideshow string `json:"slideshow,omitempty" gorm:"type:longtext"`
	// Version 版本
	Version string `json:"version,omitempty"`
	// ReleaseList 发布过的版本列表
	ReleaseList []*vos.ReleaseDigest `json:"releaseList,omitempty" gorm:"-"`
	// Status 应用状态
	Status ApplicationStatus `json:"status"`
	// Recommend 是否推荐
	Recommend bool `json:"recommend,omitempty"`
	// HideInStore 在应用商店内隐藏
	HideInStore bool `json:"hideInStore,omitempty"`
	// Debugging 是否正在调试中
	Debugging bool `json:"debugging,omitempty"`
	// UserId 用户Id
	UserId string `json:"userId,omitempty" gorm:"primaryKey"`
	//HaveRemoteDb 是否具有远程db
	HaveRemoteDb bool `json:"-"`
}

type ProxyHttpServerInfo struct {
	Name               string `json:"name,omitempty" gorm:"primary_key"`
	Host               string `json:"host,omitempty"`
	Schema             string `json:"schema,omitempty"`
	AllowReturnType    string `json:"allowReturnType,omitempty" gorm:"type:text"`
	NotAllowReturnType string `json:"notAllowReturnType,omitempty" gorm:"type:text"`
}

type ProxyHttpResponseCache struct {
	// RequestHash hash
	RequestHash string `gorm:"primary_key"`
	// ContentPath 内容路径
	ContentPath string `gorm:"type:text"`
	// ContentHash 内容hash采用sha256
	ContentHash        []byte
	ResponseHeader     string
	ResponseStatusCode int
}

// ChatGroupInfo 聊天群信息
type ChatGroupInfo struct {
	// Id 聊天群ID
	Id string `json:"id,omitempty" gorm:"primaryKey"`
	// Name 名称
	Name string `json:"name,omitempty"`
	// Icon 图标url
	Icon string `json:"icon,omitempty"`
	// Desc
	Desc string `json:"desc,omitempty" gorm:"type:longtext"`
	// CreateUserId 创建者ID
	CreateUserId string `json:"createUserId,omitempty"`
	// MainManagerUserId 群主ID
	MainManagerUserId string `json:"mainManagerUserId,omitempty"`
	// CreateAt 创建时间
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdateAt 更新时间
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// ChatType 聊天类型
type ChatType uint8

const (
	ChatUnknown ChatType = iota
	// ChatTypeUser 用户<->用户
	ChatTypeUser
	// ChatTypeGroup 用户<->群组
	ChatTypeGroup
	// ChatTypeApp 用户<->App
	ChatTypeApp
)

// ChatMsgType 聊天消息类新
type ChatMsgType uint8

const (
	// ChatMsgTypeText 文本聊天消息
	ChatMsgTypeText ChatMsgType = iota + 1
	// ChatMsgTypeFile 文件消息
	ChatMsgTypeFile
	// ChatMsgTypeImg 图片消息
	ChatMsgTypeImg
)

// UserChatMsg 用户聊天信息
type UserChatMsg struct {
	Id string `json:"id,omitempty"`
	// TargetId 目标ID
	TargetId string `json:"targetId,omitempty"`
	// TargetInfo 目标信息
	TargetInfo any `json:"targetInfo,omitempty" gorm:"-"`
	// SourceId 源Id
	SourceId string `json:"sourceId,omitempty"`
	// SourceInfo 源信息
	SourceInfo any `json:"sourceInfo,omitempty" gorm:"-"`
	// Content 内容
	Content string `json:"content,omitempty" gorm:"type:longtext"`
	// FileIcon 当消息内容为文件类型时, 存放图标
	FileIcon string `json:"fileIcon,omitempty"`
	// ChatType 聊天类别
	ChatType ChatType `json:"chatType,omitempty"`
	// MsgType 消息类别
	MsgType ChatMsgType `json:"msgType,omitempty"`
	// CreateAt 创建时间
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdateAt 更新时间
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// ClientUniqueId 客户端唯一ID
	ClientUniqueId string `json:"clientUniqueId,omitempty" gorm:"primaryKey"`
	// TimeStamp 时间戳
	TimeStamp string `json:"timeStamp,omitempty"`
	// Status 当前状态
	Status string `json:"status,omitempty"`
	// ErrMsg 错误信息
	ErrMsg string `json:"errMsg,omitempty"`
}

type QueueType uint

const (
	QueueTypeSend QueueType = iota
	QueueTypeReceive
)

// QueueChannelMsgInfo 通道消息签收信息
type QueueChannelMsgInfo struct {
	Id string `json:"id,omitempty" gorm:"primaryKey"`
	// QueueMsgId 队列消息
	QueueMsgId  string            `json:"-"`
	Content     string            `json:"content,omitempty"`
	Ack         bool              `json:"ack,omitempty"`
	AckTime     int64             `json:"ackTime,omitempty"`
	SendTime    int64             `json:"sendTime,omitempty"`
	ReceiveTime int64             `json:"receiveTime,omitempty"`
	QueueType   QueueType         `json:"-" gorm:"primaryKey"`
	Type        queue.MessageType `json:"type,omitempty"`
	AppId       string            `json:"appId,omitempty"`
}
