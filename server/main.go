package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"server-monitor-server/model"
	"server-monitor-server/util"
	"strconv"
	"sync"
)

var data = sync.Map{}
var (
	token      string
	serverPort int
	webApiPort int
)

func init() {
	var c util.Conf
	conf := c.GetConf()
	token = conf.Token
	serverPort = conf.Port.Server
	webApiPort = conf.Port.WebApi
}
func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(serverPort))
	if err != nil {
		log.Fatalf("Service initialization failed! %v\n", err)
	}
	log.Printf("Server initialized successfully, listening on port %d...\n", serverPort)
	go openHttpServe()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go createConn(conn)
	}
}

// openHttpServe Enable the HTTP service so that data can be obtained through HTTP requests.
func openHttpServe() {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("content-type", "application/json")

		var tempArray []*interface{}
		data.Range(func(key, value interface{}) bool {
			tempArray = append(tempArray, &value)
			return true
		})
		bytes, _ := json.Marshal(tempArray)
		w.Write(bytes)
	})
	log.Printf("The HTTP service is enabled and you can get data via [ip:%d/json].\n", webApiPort)
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(webApiPort), nil)
}

func createConn(conn net.Conn) {
	defer conn.Close()
	clientIp := conn.RemoteAddr().String()
	if !auth(conn) {
		bytes, _ := json.Marshal(-1)
		conn.Write(bytes)
		return
	}
	bytes, _ := json.Marshal(0)
	_, err := conn.Write(bytes)
	if err != nil {
		return
	}
	// Listen for the client to read data.
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		// If the client goes offline, change its state to false
		if err != nil {
			log.Printf("%v logout!\n", clientIp)
			if v, ok := data.Load(clientIp); ok {
				osModel := v.(*model.OSModel)
				osModel.State = false
			}
			return
		}
		var osModel model.OSModel
		err = json.Unmarshal(buf[:n], &osModel)
		if err != nil {
			log.Printf("Unmarshal: %v", err)
			continue
		}
		data.Store(clientIp, &osModel)
	}
}

// auth Authentication client.
func auth(conn net.Conn) bool {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false
	}
	_token := string(buf[:n])
	return token == _token
}
