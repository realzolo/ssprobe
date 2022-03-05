package service

import (
	"encoding/json"
	"net"
	"ssprobe-common/model"
	u "ssprobe-common/util"
	"ssprobe-server/consts"
	"ssprobe-server/util"
	"strconv"
	"strings"
	"sync"
)

var (
	ClientKeys []string
	LMDB       = sync.Map{}
	logger     = u.Logger{}
)

// StartSocketService enable the socket service to receive data from clients.
func StartSocketService(conf *util.Conf) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Server.Port))
	logger.ErrorWithExit(err, "Socket initialization failed.")
	logger.LogWithFormat("The socket was initialized successfully, listening on port %d...", conf.Server.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go createConn(&conn, conf)
	}
}

func createConn(conn *net.Conn, conf *util.Conf) {
	defer (*conn).Close()
	if err := authClient(conn, conf.Server.Token); err != nil {
		return
	}

	clientIp := util.ParseAddress((*conn).RemoteAddr().String())
	logger.LogWithFormat("%s is connected!", clientIp)

	// Listen to the client and read the data.
	for {
		buf := make([]byte, 1024)
		n, err := (*conn).Read(buf)
		// If the client goes offline, change its state to false.
		if err != nil {
			logger.LogWithFormat("%s disconnects.", clientIp)
			if v, ok := LMDB.Load(clientIp); ok {
				osModel := v.(*model.OSModel)
				osModel.State = false
				NoticeDispatcher(&conf.Notifier, osModel, consts.Offline)
			}
			return
		}
		var osModel *model.OSModel
		if err = json.Unmarshal(buf[:n], &osModel); err != nil {
			continue
		}

		if value, ok := LMDB.Load(clientIp); ok && !value.(*model.OSModel).State {
			go NoticeDispatcher(&conf.Notifier, osModel, consts.Recover)
		} else if !ok {
			go NoticeDispatcher(&conf.Notifier, osModel, consts.Online)
			ClientKeys = append(ClientKeys, clientIp)
		}
		LMDB.Store(clientIp, osModel)
	}
}

// authClient to authenticate client.
func authClient(conn *net.Conn, token string) error {
	buf := make([]byte, 1024)
	n, err := (*conn).Read(buf)
	if err != nil {
		return err
	}
	_token := strings.Trim(string(buf[:n]), "\"")

	resModel := model.AuthResponse{Code: -1}
	// Authentication failed!
	if token != _token {
		bytes, _ := json.Marshal(resModel)
		_, err = (*conn).Write(bytes)
		return err
	}
	// Authentication successful!
	resModel.Code = 0
	bytes, _ := json.Marshal(resModel)
	_, err = (*conn).Write(bytes)
	return err
}
