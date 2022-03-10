package service

import (
	"embed"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"net/http"
	"ssprobe-server/config"
	"ssprobe-server/util"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWebService(conf *util.Conf, f embed.FS) {
	// Disable console logging.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	router.Use(config.Cors())
	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "static/index.html")))
	router.StaticFS("/public", http.FS(f))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"site_title": conf.Web.Title,
			"github":     conf.Web.Github,
			"telegram":   conf.Web.Telegram,
		})
	})
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		var tempArray []*interface{}
		for {
			for _, clientKey := range ClientKeys {
				v, _ := LMDB.Load(clientKey)
				tempArray = append(tempArray, &v)
			}
			bytes, _ := json.Marshal(tempArray)
			tempArray = nil
			_ = conn.WriteMessage(websocket.TextMessage, bytes)
			time.Sleep(time.Second * 2)
		}
	})

	err := router.Run(":10240")
	logger.LogWithError(err, "")
}
