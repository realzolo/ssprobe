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
			Enable:   c.Web.Enable,
			Title:    c.SetOrDefault(c.Web.Title, consts.SiteTitle).(string),
			Github:   c.SetOrDefault(c.Web.Github, consts.Github).(string),
			Telegram: c.SetOrDefault(c.Web.Telegram, consts.Telegram).(string),
		},
		Notifier: model.Notifier{
			Telegram: model.Telegram{
				Enable:   c.Notifier.Telegram.Enable,
				UseEmbed: c.Notifier.Telegram.UseEmbed,
				Language: c.SetOrDefault(c.Notifier.Telegram.Language, consts.English).(string),
				BotToken: c.Notifier.Telegram.BotToken,
				UserId:   c.Notifier.Telegram.UserId,
			},
		},
	}
}

func main() {
	go service.StartWebService(conf)
	go service.StartWebsocketService(conf)
	go notify.InitTelegramBot(conf.Notifier.Telegram)
	service.StartSocketService(conf)
}
