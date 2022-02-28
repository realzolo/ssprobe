package monitor

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

var (
	name   string = "onezol.com"
	server string = "127.0.0.1"
	port   string = "3384"
	token  string = "123456"
)

type AuthResult struct {
	Ok   bool
	Conn *net.Conn
	Name string
}

// collectConfig Collect client configuration information.
func collectConfig() {
	args := os.Args
	for _, arg := range args {
		if strings.Contains(arg, "--name=") {
			name = arg[7:]
		} else if strings.Contains(arg, "--server=") {
			server = arg[9:]
		} else if strings.Contains(arg, "--port=") {
			port = arg[7:]
		} else if strings.Contains(arg, "--token=") {
			token = arg[8:]
		}
	}
	if len(name) == 0 || len(server) == 0 || len(port) == 0 || len(token) == 0 {
		log.Fatalf("The argument you provided does not match [--name,--server,--port,--token]. \n")
	}
}

// RequestAuth Authenticate the client and return the connection.
func RequestAuth() *AuthResult {
	collectConfig()

	conn, err := net.Dial("tcp", server+":"+port)
	if err != nil {
		log.Println("Failed to connect to server!", err)
		return &AuthResult{Ok: false}
	}
	// Authentication.
	bytes, _ := json.Marshal(token)
	_, err = conn.Write(bytes)
	if err != nil {
		log.Println("Authentication failed!", err)
		return &AuthResult{Ok: false}
	}
	var buf = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Authentication failed!", err)
		return &AuthResult{Ok: false}
	}
	var resModel = struct {
		Code int `json:"code"`
	}{}
	json.Unmarshal(buf[:n], &resModel)
	// The token is incorrect.
	if resModel.Code == -1 {
		log.Fatal("Client authentication failed, token is incorrect!\n")
	}
	log.Println("Server connection successful!")
	return &AuthResult{
		Ok:   true,
		Conn: &conn,
		Name: name,
	}
}
