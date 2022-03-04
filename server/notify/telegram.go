package notify

import (
	"fmt"
	telegram "gopkg.in/telebot.v3"
	"log"
	"math/rand"
	"ssprobe-common/util"
	"ssprobe-server/consts"
	"ssprobe-server/model"
	"strings"
	"time"
)

var logger = util.Logger{}
var (
	bot      *telegram.Bot
	receiver *telegram.User
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
		log.Printf("Telegram bot initialization failed! %v\n", err)
		return
	}
	var (
		token   = randString(32)
		hasInit = false
	)
	logger.LogWithFormat("Your Telegram Bot token is %v", token)
	_bot.Handle(telegram.OnText, func(c telegram.Context) error {
		if hasInit {
			return c.Delete()
		}
		if strings.Compare(token, c.Text()) != 0 {
			return c.Delete()
		}
		hasInit = true
		bot = _bot
		receiver = c.Sender()
		if tg.Language == consts.CHINESE {
			return c.Send("绑定成功,你将会收到来自此机器人的通知!")
		}
		return c.Send("Bind successfully, you will receive notification from this robot!")
	})

	_bot.Start()
}

func SendToTelegram(tg *model.Telegram, node *model.Node, actionType int64) {
	if bot == nil || !tg.Enable {
		return
	}
	var message string
	language := strings.ToUpper(tg.Language)
	switch language {
	case consts.ENGLISH:
		switch actionType {
		case consts.NEW:
			message = fmt.Sprintf("Meow ~, [%s - %s](%s) is online and running normally.", node.Name, node.Location, node.Host)
			break
		case consts.RENEW:
			message = fmt.Sprintf("Meow ~, your node [%s - %s](%s) returns to normal.", node.Name, node.Location, node.Host)
			break
		case consts.DOWN:
			message = fmt.Sprintf("Meow ~, your node[%s - %s](%s) failed and went offline.", node.Name, node.Location, node.Host)
			break
		}
	case consts.CHINESE:
		switch actionType {
		case consts.NEW:
			message = fmt.Sprintf("喵喵喵~, 您的机器[%s - %s](%s)已上线,状态正常!", node.Name, node.Location, node.Host)
			break
		case consts.RENEW:
			message = fmt.Sprintf("喵喵喵~ [%s - %s](%s)节点恢复了喵~", node.Name, node.Location, node.Host)
			break
		case consts.DOWN:
			message = fmt.Sprintf("喵喵喵~ [%s - %s](%s)节点掉线了喵~", node.Name, node.Location, node.Host)
			break
		}
	}
	_, err := bot.Send(receiver, message)
	logger.LogWithError(err, "")
}

func randString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
