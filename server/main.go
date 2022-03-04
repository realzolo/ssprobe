package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net"
	"net/http"
	"ssprobe-common/model"
	"ssprobe-common/util"
	"ssprobe-server/config"
	"ssprobe-server/consts"
	m "ssprobe-server/model"
	"ssprobe-server/notify"
	u "ssprobe-server/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	logger     util.Logger
	clientKeys []string
	data       = sync.Map{}
	conf       *u.Conf
	upgrader   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func init() {
	var c u.Conf
	_conf, err := c.LoadConfig()
	logger.ErrorWithExit(err, "Configuration file parsing failed!")
	conf = &u.Conf{
		Server: m.Server{
			Token:         c.SetOrDefault(_conf.Server.Token, "123456").(string),
			Port:          c.SetOrDefault(_conf.Server.Port, 3384).(int),
			WebsocketPort: c.SetOrDefault(_conf.Server.WebsocketPort, 9000).(int),
		},
		Web: m.Web{
			Enable: _conf.Web.Enable,
			Title:  c.SetOrDefault(_conf.Web.Title, "SSProbe").(string),
		},
		Notifier: m.Notifier{
			Telegram: m.Telegram{
				Enable:   c.Telegram.Enable,
				UseEmbed: c.Telegram.UseEmbed,
				Language: c.SetOrDefault(c.Telegram.Language, consts.ENGLISH).(string),
				BotToken: c.Telegram.BotToken,
				UserId:   c.Telegram.UserId,
			},
		},
	}
}

func main() {
	go openWebServe()
	go openWebsocketServe()
	go notify.InitTelegramBot(conf.Telegram)
	openSocketServe()
}

func openWebServe() {
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
			"title":          conf.Web.Title,
			"websocket_port": conf.Server.WebsocketPort,
		})
	})
	err := router.Run(":10240")
	logger.LogWithError(err, "")
}

// openSocketServe Enable the socket service to receive data from clients.
func openSocketServe() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Server.Port))
	logger.ErrorWithExit(err, "Service initialization failed!")
	logger.LogWithFormat("Server initialized successfully, listening on port %d...", conf.Server.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go createConn(conn)
	}
}

// openWebsocketServe Enable the websocket service to transmit data to the Web in real time.
func openWebsocketServe() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		var tempArray []*interface{}
		for {
			for _, clientKey := range clientKeys {
				v, _ := data.Load(clientKey)
				tempArray = append(tempArray, &v)
			}
			bytes, _ := json.Marshal(tempArray)
			tempArray = nil
			_ = conn.WriteMessage(websocket.TextMessage, bytes)
			time.Sleep(time.Second * 2)
		}
	})
	logger.LogWithFormat("The websocket service is enabled and you can get data via \"ws://ip:%d\".", conf.Server.WebsocketPort)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Server.WebsocketPort), nil)
	logger.ErrorWithExit(err, "")
}

func createConn(conn net.Conn) {
	defer conn.Close()
	if err := authClient(&conn); err != nil {
		return
	}

	clientIp := u.ParseAddress(conn.RemoteAddr().String())
	logger.LogWithFormat("%s is connected!", clientIp)
	// Listen to the client and read the data.
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		// If the client goes offline, change its state to false.
		if err != nil {
			logger.LogWithFormat("%s connection break!", clientIp)
			if v, ok := data.Load(clientIp); ok {
				osModel := v.(*model.OSModel)
				osModel.State = false
				noticeDispatcher(osModel, consts.DOWN)
			}
			return
		}
		var osModel *model.OSModel
		if err = json.Unmarshal(buf[:n], &osModel); err != nil {
			continue
		}

		if value, ok := data.Load(clientIp); ok && !value.(*model.OSModel).State {
			go noticeDispatcher(osModel, consts.RENEW)
		} else if !ok {
			go noticeDispatcher(osModel, consts.NEW)
			clientKeys = append(clientKeys, clientIp)
		}
		data.Store(clientIp, osModel)
	}
}

// auth Authentication client.
func authClient(conn *net.Conn) error {
	buf := make([]byte, 1024)
	n, err := (*conn).Read(buf)
	if err != nil {
		return err
	}
	token := string(buf[:n])
	token = strings.Trim(token, "\"")

	resModel := model.AuthResponse{Code: -1}
	// Authentication failed!
	if strings.Compare(conf.Token, token) != 0 {
		bytes, _ := json.Marshal(resModel)
		_, _ = (*conn).Write(bytes)
		return err
	}
	// Authentication successful!
	resModel.Code = 0
	bytes, _ := json.Marshal(resModel)
	_, _ = (*conn).Write(bytes)
	return nil
}

func noticeDispatcher(osModel *model.OSModel, actionType int64) {
	node := &m.Node{
		Name:     osModel.Name,
		Location: osModel.Location,
		Host:     osModel.Host,
	}
	notify.SendToTelegram(&conf.Telegram, node, actionType)
	notify.SentToTelegramByHttp(&conf.Telegram, node, actionType)
	// ...
}
