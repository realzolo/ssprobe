package notify

import (
	"encoding/json"
	"net/http"
	"ssprobe-server/model"
	"strconv"
	"strings"
)

func SentToTelegramByHttp(tg *model.Telegram, node *model.Node, actionType int64) {
	if !tg.Enable || tg.UseEmbed {
		return
	}
	bytes, _ := json.Marshal(node)
	_, err := http.Get("http://localhost/?" +
		"id=" + tg.UserId +
		"&content=" + string(bytes) +
		"&language=" + strings.ToUpper(tg.Language) +
		"&action=" + strconv.FormatInt(actionType, 10))
	logger.LogWithError(err, "")
}
