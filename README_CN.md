<div align="center">
  <h3><a href="https://status.onezol.com">SSProbe</a></h3>
  <p><em>一个由Golang编写的服务器状态监控程序..</em></p>
  <a href="https://github.com/realzolo/ssprobe/blob/master/README_CN.md"><img src="https://img.shields.io/badge/简体中文-000000?style=flat&logo=Academia&logoColor=%23FFFFFF" alt="Language" /><a/>
  <img src="https://img.shields.io/badge/Golang-black?style=flat&logo=Go&logoColor=white" alt="Golang" />
  <img src="https://img.shields.io/badge/React.js-black?style=flat&logo=React&logoColor=white" alt="React.js" />
  <img src="https://img.shields.io/github/last-commit/realzolo/ssprobe?&label=Last%20commit&color=CF2B5B&labelColor=black&logo=github" alt="commit"/>
  <img src="https://img.shields.io/github/stars/realzolo/ssprobe?color=%2300979D&label=Starts&labelColor=black&logo=Apache%20Spark&logoColor=%23FFFFFF" alt="stars"/>
</div>
<br/>
  
## 👻简介

[SSProbe](https://github.com/realzolo/ssprobe) 是一个服务器监控程序，它提供了一个可视化的界面，实时为你显示服务器状态，如CPU占用率、内存使用情况和网络速度等。

![](https://image.onezol.com/img/ssprobe.png)

## 🎉下载和使用
在 [release页面](https://github.com/realzolo/ssprobe/releases) 找到对应系统版本的压缩包,解压之后将`ssprobe-server`和`config.yaml`
这两个文件上传到你的服务器上,将`ssprobe-client`上传到你需要监控的机器上。
### 1. 配置
`config.yaml`中包含了你的一些配置,配置说明如下:
```yaml
server:
  token: 123456   # 服务器令牌,用于客户端(被监控的机器)连接服务器时验证身份
  port: 3384      # 服务器监听的端口
  websocketPort: 9000    # 如果没有前后端分离部署,请保持此项默认

web:
  enable: true    # 启用web服务。如果你需要前后端分离部署,你可以将此项设置为false
  title: Zolo's Server Monitor  # 监控页面的网站标题

notifier:
  telegram:
    enable: true        # 启用Telegram Bot发送通知
    useEmbed: false     # 使用本程序内部的TelegramBot接口,即使用自己创建的Bot。
    language: chinese   # Bot通知的语言
    botToken: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  # Bot token. 当useEmbed为true时,此项有效。
    userId: 1953745499  # Telegram 用户ID
```
> **`useEmbed`** 是否使用本程序内部的TelegramBot接口,即使用自己创建的Bot。
> 
> 若你将此服务端程序部署在中国大陆的服务器上,由于网络原因则无法使用Bot发送通知。当useEmbed设置为false时,你可以指定一个**userId**,使用我已经搭建好的Bot。
#### 🤷‍♂️如何创建Bot并获取BotToken?
Telegram搜索 `BotFather`, 向他发送`/newbot `,根据提示创建Bot。创建完毕后向他发送 `/mybots` 选择你的Bot,然后获取Token。

#### 🤷‍♀️如何获取UserID?
Telegram搜索 `SSProbe Bot`, 向他发送`/me`即可获得你的UserID.

### 2. 启动程序
(1) 服务端程序

确保`config.yaml`和`ssprobe-server`处于同一目录下,使用如下命令启动你的服务端程序:
```bash
chmod a+x ./ssprobe-server && ./ssprobe-server
```
此时,打开 `http://ip:10240` 就可以看到监控页面了。 
若你启用了TelegramBot,并且设置了 `useEmbed: true`, 则会在控制台日志中看到如下内容:

![](https://image.onezol.com/img/ssprobe-console-cn.png)
  
将方框中的Token值发送给你的Telegram Bot, 验证成功之后就会为你推送通知了。

![](https://image.onezol.com/img/bot-bind-cn.png)


(2) 客户端程序
```bash
chmod a+x ./ssprobe-client
./ssprobe-client --name=客户端名称 --server=服务器地址 --token=你的Token令牌
```
> 如果你修改了服务的监听端口,则还需要额外指定一个`--port`参数(默认服务器端口为: 3384)

例如: `./ssprobe-client --name=ClientA --server=110.42.133.216 --port=3384 --token=123456`

命令执行完毕后,就可以在监控页面看到监控数据了。

### 3. 前后端分离部署(可选)
如果你需要将前后端分离部署,你需要将 [`server/static/`](https://github.com/realzolo/ssprobe/tree/master/server/static) 目录下的文件下载下来,并且在`index.html所在目录`创建一个`config.json`文件内容如下:
```json
{
  "SITE_TITLE":"网站标题",
  "WEBSOCKET_URL": "ws://服务器地址:9000(Websocket端口)"
}
```



