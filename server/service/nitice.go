package service

import (
	"ssprobe-common/model"
	m "ssprobe-server/model"
	"ssprobe-server/notify"
)

func NoticeDispatcher(notifier *m.Notifier, osModel *model.OSModel, actionType int64) {
	node := &m.Node{
		Name:     osModel.Name,
		Location: osModel.Location,
		Host:     osModel.Host,
	}
	notify.SendToTelegram(&notifier.Telegram, node, actionType)
	notify.SentToTelegramByHttp(&notifier.Telegram, node, actionType)
	// ...
}
