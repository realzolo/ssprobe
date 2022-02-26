package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"server-monitor-server/model"
	"server-monitor-server/util"
	"strconv"
	"strings"
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
	go openHttpServe()
	openSocketServe()
}

// openSocketServe Enable the socket service to receive data from clients.
func openSocketServe() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(serverPort))
	if err != nil {
		log.Fatalf("Service initialization failed! %v", err)
	}
	log.Printf("Server initialized successfully, listening on port %d...\n", serverPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
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

	err := authClient(&conn)
	if err != nil {
		return
	}

	clientIp := parseIpAddress(conn.RemoteAddr().String())
	log.Println(clientIp, "Connection successful!")
	// Listen for the client to read data.
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		// If the client goes offline, change its state to false
		if err != nil {
			log.Printf("%s: connection break!\n", clientIp)
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
func authClient(conn *net.Conn) error {
	// Read and verify the Token.
	buf := make([]byte, 1024)
	n, err := (*conn).Read(buf)
	if err != nil {
		return err
	}
	_token := string(buf[:n])
	_token = strings.TrimLeft(_token, "\"")
	_token = strings.TrimRight(_token, "\"")
	resModel := struct {
		Code int `json:"code"`
	}{Code: -1}
	// Authentication failed!
	if !strings.EqualFold(token, _token) {
		bytes, _ := json.Marshal(resModel)
		(*conn).Write(bytes)
		return err
	}
	// Authentication successful!
	resModel.Code = 0
	bytes, _ := json.Marshal(resModel)
	(*conn).Write(bytes)
	return nil
}

// Get the IP address, for example: 127.0.0.1:8000 -> 127.0.0.1
func parseIpAddress(address string) string {
	return address[:strings.LastIndex(address, ":")]
}
