package service

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"ssprobe-server/config"
	"ssprobe-server/util"
)

func StartWebService(conf *util.Conf) {
	if !conf.Web.Enable {
		return
	}
	// Disable console logging.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	router.Use(config.Cors())
	router.LoadHTMLFiles("static/index.html")
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/ws", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"site_title":     conf.Web.Title,
			"websocket_port": conf.Server.WebsocketPort,
			"github":         conf.Web.Github,
			"telegram":       conf.Web.Telegram,
		})
	})
	err := router.Run(":10240")
	logger.LogWithError(err, "")
}
