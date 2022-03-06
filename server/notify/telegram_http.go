package notify

import (
	"encoding/json"
	"net/http"
	"ssprobe-server/model"
	"strconv"
	"strings"
)

const TgBotApi = "http://tgbot.onezol.com"

func SentToTelegramByHttp(tg *model.Telegram, node *model.Node, actionType int64) {
	if !tg.Enable || tg.UseEmbed {
		return
	}
	bytes, _ := json.Marshal(node)
	_, err := http.Get(TgBotApi + "?" +
		"id=" + strconv.FormatInt(tg.UserId, 10) +
		"&content=" + string(bytes) +
		"&language=" + strings.ToUpper(tg.Language) +
		"&action=" + strconv.FormatInt(actionType, 10))
	logger.LogWithError(err, "")
}
