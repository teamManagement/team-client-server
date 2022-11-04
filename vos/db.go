package vos

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/goextension"
	"team-client-server/tools"
)

const encryptValCipherKeyLen = 512

type EncryptValue string

func (e EncryptValue) Value() (driver.Value, error) {
	key := coderutils.Sm4RandomKey()

	cipherData, err := coderutils.Sm4Encrypt(key, []byte(e))
	if err != nil {
		return nil, err
	}

	encryptKey, err := coderutils.RsaEncrypt(key, tools.ClientPublicKey)
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

	sm4Key, err := coderutils.RsaDecrypt(encryptSm4Key, tools.ClientPrivateKey)
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
	// ApplicationNormal 正常
	ApplicationNormal
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
	Id string `json:"id,omitempty" gorm:"primary_key"`
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
	ReleaseList []*ReleaseDigest `json:"releaseList,omitempty" gorm:"-"`
	// Status 应用状态
	Status ApplicationStatus `json:"status"`
	// Recommend 是否推荐
	Recommend bool `json:"recommend,omitempty"`
	// HideInStore 在应用商店内隐藏
	HideInStore bool `json:"hideInStore,omitempty"`
	// Debugging 是否正在调试中
	Debugging bool `json:"debugging,omitempty"`
	// UserId 用户Id
	UserId string `json:"userId,omitempty"`
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
