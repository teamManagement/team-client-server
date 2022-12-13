package remoteserver

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/goextension"
	"github.com/go-base-lib/logs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/teamManagement/common/conn"
	"net"
	"strconv"
	"strings"
	"sync"
	"team-client-server/tools"
	"team-client-server/vos"
	"time"
)

type JwtTokenDataWrapper struct {
	Expire string `json:"expire,omitempty"`
	Data   []byte `json:"data,omitempty"`
}

type TcpTransferCmdCode byte

const (
	// TcpTransferCmdCodeBlockingConnection 阻断服务器链接
	TcpTransferCmdCodeBlockingConnection TcpTransferCmdCode = iota
	// TcpTransferCmdCodeRestoreServerConnErr 重新恢复用户连接失败, 通知页面退回登录
	TcpTransferCmdCodeRestoreServerConnErr
	// TcpTransferCmdCodeRestoreServerConnOK 远程服务恢复成功
	TcpTransferCmdCodeRestoreServerConnOK
	// TcpTransferCmdCodeOtherUserStatusChange 其他用户状态变更
	TcpTransferCmdCodeOtherUserStatusChange
	// TcpTransferCmCodeOtherLogin 用户在他处登录
	TcpTransferCmCodeOtherLogin
)

// TcpTransferInfo tcp转移信息, 用于tcp消息转换为ws消息
type TcpTransferInfo struct {
	// CmdCode 命令码
	CmdCode TcpTransferCmdCode `json:"cmdCode"`
	// Data 数据
	Data goextension.Bytes `json:"data,omitempty"`
	// ErrMsg 错误信息
	ErrMsg string `json:"err,omitempty"`
	// DataType 数据类型, 0: 对象, 其他: 字符串
	DataType byte `json:"dataType"`
}

var tcpTransferChan chan *TcpTransferInfo

func StartTcpTransfer() {
	tcpTransferChan = make(chan *TcpTransferInfo, 16)
}

func StopTcpTransfer() {
	close(tcpTransferChan)
	tcpTransferChan = nil
}

// GetTcpTransfer 获取tcp信息转移数据
func GetTcpTransfer() <-chan *TcpTransferInfo {
	return tcpTransferChan
}

func sendTcpTransfer(info *TcpTransferInfo) {
	defer func() { recover() }()
	if tcpTransferChan != nil {
		tcpTransferChan <- info
	}
}

var (
	nowUserInfo *vos.UserInfo
	rawToken    [][]byte
	loginOk     = false
	connCloseCh chan struct{}

	lock = sync.Mutex{}

	chanLock = sync.Mutex{}
)

const (
	LocalWebServerAddress = "https://apps.byzk.cn:443"
	LocalWSServerAddress  = "wss://apps.byzk.cn:443"

	ServerAddress = "apps.byzk.cn:80"

	autoLoginSettingKey = "USER_CREDENTIALS"
)

func ClearAutoLoginInfo() {
	autoLoginSetting := &vos.Setting{
		Name: autoLoginSettingKey,
	}
	sqlite.Db().Model(autoLoginSetting).Delete(autoLoginSetting)
	Logout()
}

func AutoLogin() (res bool) {
	autoLoginSetting := &vos.Setting{
		Name: autoLoginSettingKey,
	}
	if err := sqlite.Db().Model(autoLoginSetting).First(&autoLoginSetting).Error; err != nil {
		return
	}

	val := autoLoginSetting.Value
	valSplit := strings.Split(string(val), ".")
	if len(valSplit) != 2 {
		return
	}

	username, err := base64.StdEncoding.DecodeString(valSplit[0])
	if err != nil {
		return
	}

	password, err := base64.StdEncoding.DecodeString(valSplit[1])
	if err != nil {
		return
	}

	if err = Login(base64.StdEncoding.EncodeToString([]byte(username)), string(password)); err != nil {
		return
	}

	return true
}

func Login(username, password string) (err error) {
	if lock.TryLock() {
		defer lock.Unlock()
	}
	Logout()

	decodePasswd, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return fmt.Errorf("密码格式解析失败: %s", err.Error())
	}

	localInter, ok := tools.TelnetHostRangeNetInterfaces(ServerAddress)
	if !ok {
		return fmt.Errorf("未识别到可用的网卡信息")
	}

	localInterIpStr := localInter.String()

	localAddr, err := net.ResolveTCPAddr("tcp", localInterIpStr+":0")
	if err != nil {
		return errors.New("获取本机网卡IP失败")
	}

	usernameOriginBytes, err := base64.StdEncoding.DecodeString(username)
	if err != nil {
		return errors.New("解析用户名格式失败")
	}
	username = base64.StdEncoding.EncodeToString([]byte(localInterIpStr + "@" + string(usernameOriginBytes)))

	dialer := &net.Dialer{
		LocalAddr: localAddr,
	}

	dial, err := tls.DialWithDialer(dialer, "tcp", ServerAddress, tools.GenerateTLSConfig())
	if err != nil {
		return fmt.Errorf("连接远程服务失败: %s", err.Error())
	}

	connWrapper := conn.NewWrapper(dial)
	defer func() {
		if err != nil {
			connWrapper.WriteErrMessage(err.Error())
			_ = dial.Close()
		}
	}()

	if err = connWrapper.
		WriteByte(1).
		WriteFormatBytesData([]byte(username)).
		WriteFormatBytesData([]byte(password)).Error(); err != nil {
		return fmt.Errorf("向服务器端发送数据指令失败")
	}

	jwtTokenBytes, err := connWrapper.ReadeFormatBytesData()
	if err != nil {
		return err
	}

	refreshTokenBytes, err := connWrapper.ReadeFormatBytesData()
	if err != nil {
		return err
	}

	jwtSplit := bytes.Split(jwtTokenBytes, []byte("."))
	if len(jwtSplit) != 3 {
		return fmt.Errorf("用户凭证格式不正确")
	}

	rawToken = [][]byte{jwtTokenBytes, refreshTokenBytes}

	jwtTokenStr := string(jwtTokenBytes)

	//jwtDataBytes, err := base64.StdEncoding.DecodeString(string(jwtSplit[1]))
	jwtDataBytes, err := jwt.DecodeSegment(string(jwtSplit[1]))
	if err != nil {
		return fmt.Errorf("解析token内的数据失败")
	}

	var dataWrapper *JwtTokenDataWrapper
	if err = json.Unmarshal(jwtDataBytes, &dataWrapper); err != nil {
		return fmt.Errorf("解析token包装数据结构失败")
	}

	if err = json.Unmarshal(dataWrapper.Data, &nowUserInfo); err != nil {
		return fmt.Errorf("解析用户信息失败： %s", err.Error())
	}

	if nowUserInfo.TokenExpire, err = strconv.ParseInt(dataWrapper.Expire, 10, 64); err != nil {
		return fmt.Errorf("解析token有效期失败: %s", err.Error())
	}

	nowUserInfo.Token = jwtTokenStr

	cachePasswd := nowUserInfo.Id + "_teamwork_cache_" + string(decodePasswd)
	cachePasswdHash, err := coderutils.Hash(sha1.New(), []byte(cachePasswd))
	if err != nil {
		return fmt.Errorf("用户缓存帐号计算失败: %s", err.Error())
	}
	nowUserInfo.CachePassword = cachePasswdHash.ToHexStr()

	if err = connWrapper.WriteByte(6).Error(); err != nil {
		return fmt.Errorf("消息签收失败")
	}

	if confirm, err := connWrapper.ReadByte(); err != nil || confirm != 6 {
		return fmt.Errorf("登录消息签收失败")
	}

	sqlite.Db().Save(&vos.Setting{
		Name:  autoLoginSettingKey,
		Value: vos.EncryptValue(base64.StdEncoding.EncodeToString([]byte(nowUserInfo.Username)) + "." + base64.StdEncoding.EncodeToString([]byte(password))),
	})

	loginOk = true

	_ = sqlite.Db().Table("app-" + nowUserInfo.Id + "-0").AutoMigrate(&vos.Setting{})

	_ = FlushAllCache()

	if err = startChatWs(); err != nil {
		return err
	}

	connCloseCh = make(chan struct{}, 1)
	go userConnHandler(connWrapper, dial)

	return nil
}

func userConnHandler(connWrapper *conn.Wrapper, dial net.Conn) {
	var (
		cmdCh       = make(chan byte, 1)
		operationCh = make(chan byte, 1)
		dataCh      = make(chan goextension.Bytes, 1)
		errCh       = make(chan error, 1)
	)

	logs.Debugln("监听用户远程TCP消息通道")

	isClose := false
	defer func() {
		stopWsChat()
		operationCh <- 1
		close(cmdCh)
		close(operationCh)
		close(dataCh)
		close(errCh)
		_ = dial.Close()
		nowUserInfo = nil
		rawToken = nil

		if !isClose {
			logs.Debugln("发送tcp关闭换消息到tcp与ws交互通道内...")
			tcpTransferChan <- &TcpTransferInfo{
				CmdCode: TcpTransferCmdCodeBlockingConnection,
			}
			logs.Debugln("成功退出用户远程TCP服务消息通道")
			return
		}
		close(connCloseCh)
		logs.Debugln("成功关闭用户远程TCP服务通道")
	}()

	go connReadHandler(connWrapper, cmdCh, operationCh, dataCh, errCh)

	readDataFn := func() (goextension.Bytes, error) {
		timeout := time.After(30 * time.Second)
		select {
		case <-timeout:
			return nil, errors.New("数据读取超时")
		case d := <-dataCh:
			return d, <-errCh
		}
	}

	writeDataFn := func(data []byte) error {
		operationCh <- 2
		dataCh <- data
		return <-errCh
	}

	for {
		t := time.After(8 * time.Minute)
		select {
		case <-t:
			chanLock.Lock()
			logs.Debugln("token即将过期, 主动尝试用户token延期业务...")
			if err := tokenDelay(readDataFn, writeDataFn); err != nil {
				chanLock.Unlock()
				logs.Debugf("token延期失败, 将要关闭用户远程TCP服务通道, 本次错误信息: %s", err.Error())
				return
			}
			chanLock.Unlock()
			logs.Debugln("token延期成功")
		case <-connCloseCh:
			logs.Debugln("主动关闭用户远程TCP服务通道")
			isClose = true
			return
		case code := <-cmdCh:
			data := <-dataCh
			err := <-errCh
			logs.Debugf("调用代码")
			if err = serverCmdHandler(err, code, data, readDataFn, writeDataFn); err != nil {
				return
			}
		case err := <-errCh:
			if err == nil {
				continue
			}
			logs.Debugf("TCP通道监听返回错误: %s， 将要关闭TCP通道", err.Error())
			return
		}
	}
}

func tokenDelay(readDataFn func() (goextension.Bytes, error), writeDataFn func(data []byte) error) error {
	lock.Lock()
	defer lock.Unlock()

	if len(rawToken) != 2 {
		return errors.New("缺失的Token信息")
	}

	if err := writeDataFn(bytes.Join([][]byte{
		{0, 3},
		rawToken[0],
		{'.'},
		rawToken[1],
	}, nil)); err != nil {
		return err
	}

	jwtTokenBytes, err := readDataFn()
	if err != nil {
		return err
	}

	refreshTokenBytes, err := readDataFn()
	if err != nil {
		return err
	}

	jwtSplit := bytes.Split(jwtTokenBytes, []byte("."))
	if len(jwtSplit) != 3 {
		return fmt.Errorf("用户凭证格式不正确")
	}

	jwtDataBytes, err := jwt.DecodeSegment(string(jwtSplit[1]))
	if err != nil {
		return fmt.Errorf("解析token内的数据失败")
	}

	var dataWrapper *JwtTokenDataWrapper
	if err = json.Unmarshal(jwtDataBytes, &dataWrapper); err != nil {
		return fmt.Errorf("解析token包装数据结构失败")
	}

	if err = json.Unmarshal(dataWrapper.Data, &nowUserInfo); err != nil {
		return fmt.Errorf("解析用户信息失败： %s", err.Error())
	}

	if nowUserInfo.TokenExpire, err = strconv.ParseInt(dataWrapper.Expire, 10, 64); err != nil {
		return fmt.Errorf("解析token有效期失败: %s", err.Error())
	}

	nowUserInfo.Token = string(jwtTokenBytes)

	rawToken = [][]byte{jwtTokenBytes, refreshTokenBytes}

	return nil
}

func serverCmdHandler(err error, cmdCode byte, data goextension.Bytes, readDataFn func() (goextension.Bytes, error), writeDataFn func(data []byte) error) error {
	switch cmdCode {
	case remoteModelCodeUpdateCache:
		_ = FlushAllCache()
	case remoteModelCodeOtherUserStatusChange:
		sendTcpTransfer(&TcpTransferInfo{
			CmdCode: TcpTransferCmdCodeOtherUserStatusChange,
			Data:    data,
		})
	case remoteModelCodeOtherLoginNotify:
		sendTcpTransfer(&TcpTransferInfo{
			CmdCode: TcpTransferCmCodeOtherLogin,
			Data:    data,
		})
	}
	return nil
}

func connReadHandler(connWrapper *conn.Wrapper, cmdCh chan byte, operationCh chan byte, dataCh chan goextension.Bytes, errCh chan error) {
	connDataCh := make(chan goextension.Bytes, 1)
	connErrCh := make(chan error, 1)
	defer func() {
		close(connDataCh)
		close(connErrCh)
	}()

	defer func() { recover() }()

	go func() {
		defer func() { recover() }()
		for {
			data, err := connWrapper.ReadeFormatBytesData()
			if err != nil {
				connErrCh <- err
				connDataCh <- nil
				return
			}

			connErrCh <- nil
			connDataCh <- data
		}
	}()
	for {

		select {
		case d := <-connDataCh:
			if err := <-connErrCh; err != nil {
				errCh <- err
				continue
			}
			switch d[0] {
			case 0: // 发送新的命令
				cmdCh <- d[1]
				errCh <- nil
				dataCh <- d[2:]
				errCh <- nil
			case 1: // 响应数据
				dataCh <- d[1:]
				errCh <- nil
			}
		case c := <-operationCh:
			switch c {
			case 1: // 退出
				return
			case 2: // 发送数据
				d := <-dataCh
				errCh <- connWrapper.WriteFormatBytesData(d).Error()

			}
		}
	}
}

func Logout() {
	if lock.TryLock() {
		defer lock.Unlock()
	}
	stopWsChat()
	chanLock.Lock()
	defer chanLock.Unlock()
	loginOk = false

	if connCloseCh == nil {
		return
	}

	logs.Debugln("向TCP服务通道发送关闭指令....")
	connCloseCh <- struct{}{}
	<-connCloseCh
	logs.Debugln("TCP服务通道成功返回关闭成功")

	connCloseCh = nil
}

func LoginOk() bool {
	lock.Lock()
	defer lock.Unlock()
	return loginOk
}

func Token() string {
	user, err := NowUser()
	if err != nil {
		return ""
	}

	return user.Token
}

func LoginIp() string {
	user, err := NowUser()
	if err != nil {
		return ""
	}
	return user.LoginIp
}

func NowUser() (*vos.UserInfo, error) {
	if lock.TryLock() {
		defer lock.Unlock()
	}

	if nowUserInfo == nil {
		return nil, errors.New("用户未登录")
	}

	if time.Now().Unix() > nowUserInfo.TokenExpire {
		if !AutoLogin() {
			return nil, errors.New("用户登录信息失效")
		}
	}

	return nowUserInfo, nil
}
