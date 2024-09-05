package router

import (
	"embed"
	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/gin-gonic/gin"
	nh "net/http"
)

var (
	//go:embed frontend/build
	content embed.FS
)

func init() {
	go serveBackend()
	go serveFrontend()
}

func serveBackend() {
	engine := http.NewEngine(global.Config.HttpEngine.ServeURL)

	engine.AddEndPoints(OpenAiCompatibleRouterGroup...)
	engine.AddEndPoints(ManagementRouterGroup...)

	engine.ServeAsync(global.Config.HttpEngine.ServeAddr, make(chan struct{}, 1))
}

func serveFrontend() {
	engine := gin.New()
	engine.StaticFS("/static", nh.FS(content))

	engine.GET("/", func(c *gin.Context) {
		data, err := content.ReadFile("frontend/build/index.html")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	engine.Run("0.0.0.0:3000")
}
