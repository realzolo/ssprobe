package notify

import (
	"fmt"
	telegram "gopkg.in/telebot.v3"
	"ssprobe-common/util"
	"ssprobe-server/consts"
	"ssprobe-server/model"
	"strings"
	"time"
)

var (
	logger = util.Logger{}
	bot    *telegram.Bot
)

func InitTelegramBot(tg model.Telegram) {
	if !tg.Enable || !tg.UseEmbed {
		return
	}
	pref := telegram.Settings{
		Token:  tg.BotToken,
		Poller: &telegram.LongPoller{Timeout: 10 * time.Second},
	}
	_bot, err := telegram.NewBot(pref)
	if err != nil {
		logger.OnlyLog("Telegram Bot initialization failed, the token does not exist.")
		return
	}
	logger.OnlyLog("Telegram Bot started successfully!")
	_bot.Handle(telegram.OnText, func(c telegram.Context) error {
		if c.Text() != "/start" {
			return c.Delete()
		}
		return nil
	})
	_bot.Handle("/me", func(c telegram.Context) error {
		msg := fmt.Sprintf("Hello %s\nUser ID: %d\nUsername: %s\n", c.Sender().FirstName, c.Sender().ID, c.Sender().Username)
		return c.Send(msg)
	})
	bot = _bot
	_bot.Start()
}

func SendToTelegram(tg *model.Telegram, node *model.Node, actionType int64) {
	if bot == nil || !tg.Enable {
		return
	}
	var message string
	language := strings.ToUpper(tg.Language)
	switch language {
	case consts.English:
		switch actionType {
		case consts.Online:
			message = fmt.Sprintf("Meow ~, The node [%s - %s](%s) is online!", node.Name, node.Location, node.Host)
			break
		case consts.Recover:
			message = fmt.Sprintf("Meow ~, The node [%s - %s](%s) is recovered.", node.Name, node.Location, node.Host)
			break
		case consts.Offline:
			message = fmt.Sprintf("Meow ~, The node [%s - %s](%s) is offline.", node.Name, node.Location, node.Host)
			break
		}
	case consts.Chinese:
		switch actionType {
		case consts.Online:
			message = fmt.Sprintf("喵喵喵~, 您的机器 [%s - %s](%s) 已上线~", node.Name, node.Location, node.Host)
			break
		case consts.Recover:
			message = fmt.Sprintf("喵喵喵~, 节点 [%s - %s](%s) 恢复了~", node.Name, node.Location, node.Host)
			break
		case consts.Offline:
			message = fmt.Sprintf("喵喵喵~, 节点 [%s - %s](%s) 掉线了~", node.Name, node.Location, node.Host)
			break
		}
	}
	_, err := bot.Send(&telegram.User{ID: tg.UserId}, message)
	logger.LogWithError(err, "")
}
