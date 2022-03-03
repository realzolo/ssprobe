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
	//Email    `yaml:"email"`
}
type Telegram struct {
	Enable   bool   `yaml:"enable"`
	BotToken string `yaml:"botToken"`
}

/*
type Email struct {
	Enable   bool `yaml:"enable"`
	Sender   Mail `yaml:"sender"`
	Receiver Mail `yaml:"receiver"`
}

type Mail struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
*/
