package main

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/wymli/relaxadmin/common/config"
	"github.com/wymli/relaxadmin/common/utils"
	"github.com/wymli/relaxadmin/dal/rds"
	"github.com/wymli/relaxadmin/handler"
	"k8s.io/klog/v2"
)

//go:embed ui/*
var uiFS embed.FS

//go:embed static/*
var staticFS embed.FS

var (
	indexHtml []byte
)

type HtmlConfig struct {
	Env string
}

func main() {
	config.Init()

	gconfig := config.GetConfig()

	klog.Infof("global config init: %+v", utils.Json(gconfig))

	if gconfig.Server.ConnectDB {
		rds.Init(gconfig.DB)
	}

	indexHtmlConfig := &HtmlConfig{
		Env: gconfig.Server.Env,
	}

	// 初始化html渲染
	InitIndexHtml(indexHtmlConfig)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/fe") })
	r.GET("/index", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/fe") })
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	r.GET("/fe", handleIndexHtml(gconfig.Server.ReloadUI, indexHtmlConfig))

	r.StaticFS("/static", handleEmbedFS(staticFS, "static"))

	_ = r.Group("/api/v1")

	proxyGroup := r.Group("/proxy")
	p := handler.GetDefaultReverseProxy()
	proxyGroup.Any("/*wildcard", p.Proxy)

	r.Run(gconfig.Server.Host + ":" + gconfig.Server.Port)
}

func handleEmbedFS(embedfs embed.FS, subpath string) http.FileSystem {
	sub, err := fs.Sub(staticFS, subpath)
	if err != nil {
		panic(err)
	}

	return http.FS(sub)
}

func InitIndexHtml(c *HtmlConfig) {
	tmpl := template.Must(template.New("").ParseFS(uiFS, "ui/*/*.html", "ui/*.html"))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.ExecuteTemplate(buf, "index.html", c); err != nil {
		panic(err)
	}

	indexHtml = buf.Bytes()
}

func handleIndexHtml(reload bool, hc *HtmlConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		if reload {
			tmpl := template.Must(template.Must(template.New("").ParseGlob("ui/*/*.html")).ParseGlob("ui/*.html"))
			tmpl.ExecuteTemplate(c.Writer, "index.html", hc)
			klog.Infof("ui reloaded..., config=%s", utils.Json(hc))
		} else {
			c.Writer.Write(indexHtml)
		}

		c.Header("ContentType", "text/html; charset=utf-8")
	}
}
