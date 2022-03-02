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
	_util "ssprobe-server/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	logger     util.Logger
	clientKeys []string
	data       = sync.Map{}
	conf       *_util.Conf
	upgrader   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func init() {
	var c _util.Conf
	_conf, err := c.GetConf()
	logger.ErrorWithExit(err, "Configuration file parsing failed!")
	conf = _conf
}

func main() {
	go openWebServe()
	go openWebsocketServe()
	openSocketServe()
}

func openWebServe() {
	// Disable console logging.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	router.LoadHTMLFiles("static/index.html")
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	err := router.Run(":" + strconv.Itoa(conf.Port.Web))
	logger.LogWithError(err, "")
}

// openSocketServe Enable the socket service to receive data from clients.
func openSocketServe() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Port.Server))
	logger.ErrorWithExit(err, "Service initialization failed!")
	logger.LogWithFormat("Server initialized successfully, listening on port %d...", conf.Port.Server)
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
	logger.LogWithFormat("The websocket service is enabled and you can get data via \"ws://[ip:%d/json]\".", conf.Port.Websocket)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port.Websocket), nil)
	logger.ErrorWithExit(err, "")
}

func createConn(conn net.Conn) {
	defer conn.Close()
	if err := authClient(&conn); err != nil {
		return
	}

	clientIp := parseIpAddress(conn.RemoteAddr().String())
	logger.LogWithFormat("[%s] is connected!", clientIp)
	// Listen to the client and read the data.
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		// If the client goes offline, change its state to false.
		if err != nil {
			logger.LogWithFormat("[%s]: connection break!", clientIp)
			if v, ok := data.Load(clientIp); ok {
				osModel := v.(*model.OSModel)
				osModel.State = false
			}
			return
		}
		var osModel model.OSModel
		if err = json.Unmarshal(buf[:n], &osModel); err != nil {
			continue
		}

		if _, ok := data.Load(clientIp); !ok {
			clientKeys = append(clientKeys, clientIp)
		}
		data.Store(clientIp, &osModel)
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

// Get the IP address, for example: 127.0.0.1:8000 -> 127.0.0.1
func parseIpAddress(address string) string {
	return address[:strings.LastIndex(address, ":")]
}
