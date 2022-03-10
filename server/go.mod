module ssprobe-server

go 1.17

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/gobuffalo/packr/v2 v2.8.3
	github.com/gorilla/websocket v1.5.0
	gopkg.in/telebot.v3 v3.0.0
	gopkg.in/yaml.v2 v2.4.0
	ssprobe-common v0.0.0
)

replace ssprobe-common => ../common

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/gobuffalo/logger v1.0.6 // indirect
	github.com/gobuffalo/packd v1.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/markbates/errx v1.1.0 // indirect
	github.com/markbates/oncer v1.0.0 // indirect
	github.com/markbates/safe v1.0.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)
