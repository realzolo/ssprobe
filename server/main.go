package main

import (
	"ssprobe-common/util"
	"ssprobe-server/consts"
	"ssprobe-server/model"
	"ssprobe-server/notify"
	"ssprobe-server/service"
	u "ssprobe-server/util"
)

var (
	logger util.Logger
	conf   *u.Conf
)

func init() {
	var c u.Conf
	err := c.LoadConfig()
	logger.ErrorWithExit(err, "Configuration file parsing failed.")
	conf = &u.Conf{
		Server: model.Server{
			Token:         c.SetOrDefault(c.Server.Token, consts.ServerToken).(string),
			Port:          c.SetOrDefault(c.Server.Port, consts.ServerPort).(int),
			WebsocketPort: c.SetOrDefault(c.Server.WebsocketPort, consts.WebsocketPort).(int),
		},
		Web: model.Web{
			Enable: c.Web.Enable,
			Title:  c.SetOrDefault(c.Web.Title, consts.SiteTitle).(string),
		},
		Notifier: model.Notifier{
			Telegram: model.Telegram{
				Enable:   c.Telegram.Enable,
				UseEmbed: c.Telegram.UseEmbed,
				Language: c.SetOrDefault(c.Telegram.Language, consts.English).(string),
				BotToken: c.Telegram.BotToken,
				UserId:   c.Telegram.UserId,
			},
		},
	}
}

func main() {
	go service.StartWebService(conf)
	go service.StartWebsocketService(conf)
	go notify.InitTelegramBot(conf.Telegram)
	service.StartSocketService(conf)
}
