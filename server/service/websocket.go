package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"ssprobe-server/util"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// StartWebsocketService enable the websocket service to transmit data to the web in real time.
func StartWebsocketService(conf *util.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
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
	logger.LogWithFormat("The websocket is enabled and you can get data via \"ws://127.0.0.1:%d\"", conf.Server.WebsocketPort)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Server.WebsocketPort), nil)
	logger.ErrorWithExit(err, "")
}
