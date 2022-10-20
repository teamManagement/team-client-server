package server

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/logs"
	"io/fs"
	"net/http"
	"os"
)

//go:embed icons
var icons embed.FS

func InitIcons(engine *gin.Engine) {

	f, err := fs.Sub(icons, "icons")
	if err != nil {
		logs.Panicf("读取图标文件失败: %s", err.Error())
		os.Exit(11)
	}

	engine.StaticFS("/icons", http.FS(f))
}
