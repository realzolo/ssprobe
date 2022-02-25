package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"server-monitor/model"
	"sync"
)

const TOKEN = "123456"

var data = sync.Map{}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:3384")
	if err != nil {
		log.Fatalf("Service initialization failed! %v\n", err)
	}
	log.Println("Server initialized successfully, listening on port 3384...")
	go httpServe()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go createConn(conn)
	}
}

func httpServe() {
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
	http.ListenAndServe("0.0.0.0:9000", nil)
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
			log.Printf("Data parsing failed: %#v", err.Error())
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
	token := string(buf[:n])
	return TOKEN == token
}
