package website

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed appStore/build
var appStoreWebSite embed.FS

// InitAppWebSite 初始化内部应用站点
func InitAppWebSite(engine *gin.Engine) {
	//initOrganizationManaged(engine)
	//initRosterApp(engine)
	initInsideAppWebSite(engine, appStoreWebSite, "appStore/build", "appStore")
}

// initInsideAppWebSite 初始化内部应用站点
func initInsideAppWebSite(engine *gin.Engine, embedFs embed.FS, subDir string, appRouteName string) {

	tmpl := template.Must(template.New("").ParseFS(embedFs, subDir+"/*.html"))

	f, err := fs.Sub(embedFs, subDir)
	if err != nil {
		panic(err)
	}
	engine.StaticFS(appRouteName, http.FS(f)).
		GET(appRouteName+"View", func(ctx *gin.Context) {
			ctx.Header("Content-Type", "text/html;charset=utf-8")
			ctx.Status(200)
			_ = tmpl.ExecuteTemplate(ctx.Writer, "index.html", nil)
		})
}
