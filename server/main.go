package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"server-monitor/model"
	"strings"
	"sync"
)

const TOKEN = "123456"

var data = sync.Map{}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3384")
	if err != nil {
		log.Fatalf("服务初始化失败! %v\n", err)
	}
	log.Println("服务器初始化成功,正在监听3384端口...")
	go httpServe()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept err. %v\n", err)
			continue
		}
		go createConn(conn)
	}
}

func httpServe() {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json

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
	//host := parseIpAddress(conn.RemoteAddr().String())
	host := conn.RemoteAddr().String()
	defer conn.Close()
	// 认证客户机
	if !auth(conn) {
		bytes, _ := json.Marshal(model.AuthModel{Code: -1})
		conn.Write(bytes)
		return
	}
	bytes, _ := json.Marshal(
		model.AuthModel{
			Code:      0,
			Host:      host,
			IpVersion: getIpVersion(host),
		})
	_, err := conn.Write(bytes)
	if err != nil {
		return
	}
	// 持续读取数据
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf) // n: 读到的字节数
		// 若客户机下线,则将其状态更改为false
		if err != nil {
			log.Printf("%v 已退出!\n", host)
			if v, ok := data.Load(host); ok {
				osModel := v.(*model.OSModel)
				osModel.State = false
			}
			return
		}
		var osInfo model.OSModel
		err = json.Unmarshal(buf[:n], &osInfo)
		if err != nil {
			log.Fatalf("数据解析失败: %#v", err)
		}
		data.Store(host, &osInfo)
	}
}

// auth 根据用户传过来的Token验证用户身份
func auth(conn net.Conn) bool {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // n: 读到的字节数
	if err != nil {
		return false
	}
	token := string(buf[:n])
	if TOKEN == token {
		log.Printf("---用户 %v 连接成功!\n", conn.RemoteAddr())
		return true
	}
	return false
}

// 获取IP版本
func getIpVersion(address string) string {
	if strings.Count(address, ":") < 1 {
		return "IPv4"
	} else { // strings.Count(address, ":") >= 1
		return "IPv6"
	}
}

// 获取ip地址,例如: 127.0.0.1:8000 -> 127.0.0.1
func parseIpAddress(address string) string {
	return address[:strings.LastIndex(address, ":")]
}
