package model

type Server struct {
	Token         string `yaml:"token"`
	Port          int    `yaml:"port"`
	WebsocketPort int    `yaml:"websocketPort"`
}
type Web struct {
	Enable bool   `yaml:"enable"`
	Title  string `yaml:"title"`
}
type Notifier struct {
	Telegram `yaml:"telegram"`
}
type Telegram struct {
	Enable   bool   `yaml:"enable"`
	UseEmbed bool   `yaml:"useEmbed"`
	Language string `yaml:"language"`
	BotToken string `yaml:"botToken"`
	UserId   string `yaml:"userId"`
}
