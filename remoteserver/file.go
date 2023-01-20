package remoteserver

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/eventials/go-tus"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/coderutils"
	"github.com/patrickmn/go-cache"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"io"
	"os"
	"team-client-server/config"
	"team-client-server/sfnake"
	"team-client-server/tools"
	"time"
)

type FileResponseType uint

const (
	FileResponseTypeFile = iota + 1
	FileResponseTypeProgress
)

type FileResponseData struct {
	// Type 类型
	Type FileResponseType `json:"type,omitempty"`
	// Progress 进度
	Progress float32 `json:"progress,omitempty"`
	// FileId 文件ID
	FileId string `json:"fileId,omitempty"`
}

type UploadProgressNotifyInfo struct {
	Err          error
	ProgressChan tus.Upload
}

var (
	uploadAddress string

	tusConfig = tus.DefaultConfig()

	uploadProgressCache = cache.New(30*time.Second, 1*time.Minute)
)

func InitLocalService(engine *gin.Engine) {
	uploadAddress = config.LocalWebServerAddress + "/fileStorage"
	tusConfig.HttpClient = tools.DefaultHttpClient

	// remoteServer/file/userChat
	engine.Group("/r/f/uc").
		POST("upload", ginmiddleware.WrapperResponseHandle(httpServiceUpload))
}

var (
	// httpServiceUpload 业务上传接口
	httpServiceUpload ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		filePath := ctx.Query("_f")
		if filePath == "" {
			return errors.New("缺失文件路径")
		}

		if stat, err := os.Stat(filePath); err != nil || stat.IsDir() {
			return fmt.Errorf("文件[%s]不存在或是一个目录", filePath)
		}

		h, err := coderutils.HashByFilePath(sha1.New(), filePath)
		if err != nil {
			return fmt.Errorf("计算文件[%s]HASH失败: %s", filePath, err.Error())
		}

		var fileId string
		if err = RequestWebServiceWithResponse(uploadAddress+"/check/exist/"+h.ToHexStr(), &fileId); err != nil {
			return err
		}

		if fileId != "" {
			return FileResponseData{
				Type:   FileResponseTypeFile,
				FileId: fileId,
			}
		}

		if fileId, err = sfnake.GetIdStr(); err != nil {
			return fmt.Errorf("生成文件一次性ID失败: %s", err.Error())
		}

		t := Token()
		if t == "" {
			return errors.New("未登录, 请先登录")
		}

		uploadProgressNotifyChan := make(chan UploadProgressNotifyInfo, 8)
		uploadProgressCache.Set(fileId, uploadProgressNotifyChan, -1)
		go func() {
			defer func() {
				uploadProgressCache.SetDefault(fileId, uploadProgressNotifyChan)
				close(uploadProgressNotifyChan)
			}()

			f, err := os.OpenFile(filePath, os.O_RDONLY, 0655)
			if err != nil {
				uploadProgressNotifyChan <- UploadProgressNotifyInfo{Err: fmt.Errorf("打开文件[%s]失败: %s", err.Error())}
				return
			}
			defer f.Close()

			tusClient, err := tus.NewClient(uploadAddress+"/upload", tusConfig)
			if err != nil {
				uploadProgressNotifyChan <- UploadProgressNotifyInfo{Err: fmt.Errorf("创建文件上传客户端失败: %s", err.Error())}
				return
			}

			tusClient.Header.Set("_t", t)
			tusClient.Header.Set("_a", LoginIp())
			tusClient.Header.Set("User-Agent", "teamManageLocal")

			upload, err := tus.NewUploadFromFile(f)
			if err != nil {
				uploadProgressNotifyChan <- UploadProgressNotifyInfo{Err: fmt.Errorf("创建上传文件信息体失败: %s", err.Error())}
				return
			}

			uploader, err := tusClient.CreateUpload(upload)
			if err != nil {
				uploadProgressNotifyChan <- UploadProgressNotifyInfo{Err: fmt.Errorf("创建文件上传器失败: %s", err.Error())}
				return
			}

			uploadProgressChan := make(chan tus.Upload, 8)
			uploader.NotifyUploadProgress(uploadProgressChan)

			go func() {
				for uploadProgress := range uploadProgressChan {
					uploadProgressNotifyChan <- UploadProgressNotifyInfo{
						ProgressChan: uploadProgress,
					}
				}
			}()

			if err = uploader.Upload(); err != nil {
				uploadProgressNotifyChan <- UploadProgressNotifyInfo{Err: err}
			} else {
				uploadProgressNotifyChan <- UploadProgressNotifyInfo{Err: io.EOF}
			}
		}()

		return FileResponseData{
			Type:   FileResponseTypeProgress,
			FileId: fileId,
		}
	}
)
