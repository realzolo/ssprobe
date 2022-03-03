package notify

import (
	telegram "gopkg.in/telebot.v3"
	"log"
	"math/rand"
	"ssprobe-common/util"
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
	if !tg.Enable {
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
		return c.Send("Bind successfully, you will receive notification from this robot!")
	})

	_bot.Start()
}

func SendToTelegram(message string) {
	if bot == nil {
		return
	}
	_, err := bot.Send(receiver, message)
	if err != nil {
		log.Printf("消息发送失败！%v\n", err)
	}
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
