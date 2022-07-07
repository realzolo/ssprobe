package monitor

import (
	"encoding/json"
	"net"
	"os"
	"ssprobe-common/model"
	"ssprobe-common/util"
	"strings"
)

type RequestParam struct {
	Name   string
	Server string
	Port   string
	Token  string
}

type AuthResult struct {
	Ok   bool
	Conn *net.Conn
	Name string
}

var logger util.Logger

// parseParam parse user parameters.
func parseParam() *RequestParam {
	rp := &RequestParam{
		Name:   "onezol.com",
		Server: "127.0.0.1",
		Port:   "3384",
		Token:  "123456",
	}
	for _, arg := range os.Args {
		if strings.Contains(arg, "--name=") {
			rp.Name = arg[7:]
		} else if strings.Contains(arg, "--server=") {
			rp.Server = arg[9:]
		} else if strings.Contains(arg, "--port=") {
			rp.Port = arg[7:]
		} else if strings.Contains(arg, "--token=") {
			rp.Token = arg[8:]
		}
	}
	if len(rp.Server) == 0 || len(rp.Token) == 0 {
		logger.LogWithExit("The argument you provided does not match [--server,--token].")
	}
	return rp
}

// RequestAuth authenticate the client and return the connection.
func RequestAuth() *AuthResult {
	rp := parseParam()

	conn, err := net.Dial("tcp", rp.Server+":"+rp.Port)
	if err != nil {
		logger.LogWithError(err, "Failed to connect to server!")
		return &AuthResult{Ok: false}
	}
	// Authentication.
	bytes, _ := json.Marshal(rp.Token)
	_, err = conn.Write(bytes)
	if err != nil {
		logger.OnlyLog("Failed to send authentication request.")
		return &AuthResult{Ok: false}
	}
	var buf = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		logger.LogWithError(err, "Authentication failed.")
		return &AuthResult{Ok: false}
	}
	var res model.AuthResponse
	_ = json.Unmarshal(buf[:n], &res)
	// The token is incorrect.
	if res.Code == -1 {
		logger.LogWithExit("Authentication failed, incorrect token.")
	}
	logger.OnlyLog("Server connection successful!")
	return &AuthResult{
		Ok:   true,
		Conn: &conn,
		Name: rp.Name,
	}
}
