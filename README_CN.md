<div align="center">
  <h1><a href="https://status.onezol.com">SSProbe</a></h1>
  <p><em>A server status monitoring program, powered by Golang.</em></p>
  <a href="https://github.com/realzolo/ssprobe/blob/master/README_CN.md"><img src="https://img.shields.io/badge/简体中文-000000?style=flat&logo=Academia&logoColor=%23FFFFFF" alt="Language" /><a/>
  <img src="https://img.shields.io/badge/Golang-black?style=flat&logo=Go&logoColor=white" alt="Golang" />
  <img src="https://img.shields.io/badge/React.js-black?style=flat&logo=React&logoColor=white" alt="React.js" />
  <img src="https://img.shields.io/github/last-commit/realzolo/ssprobe?&label=Last%20commit&color=CF2B5B&labelColor=black&logo=github" alt="commit"/>
  <img src="https://img.shields.io/github/stars/realzolo/ssprobe?color=%2300979D&label=Stars&labelColor=black&logo=Apache%20Spark&logoColor=%23FFFFFF" alt="stars"/>
<br/><br/>
</div>



## 🏅简介

[SSProbe](https://github.com/realzolo/ssprobe) 是一款服务器监控程序(也就是所谓的"探针")，它提供了一个可视化的界面，实时为你显示服务器状态，如CPU占用率、内存使用情况和网络速度等。

* 低CPU和内存占用。📉
* 程序部署极为容易。🚀
* 漂亮的数据展示页面。🧙
* 服务端掉线或者重启时，客户端会自动回连(默认时间60s)。⏱️

### Demo

Live demo at [Zolo's SSProbe](https://status.onezol.com/).

![](https://image.onezol.com/img/ssprobe-cn.png)

## 🎉下载和使用

在 [release页面](https://github.com/realzolo/ssprobe/releases) 找到对应系统版本的压缩包,解压之后将`ssprobe-server`和`config.yaml`这两个文件上传到你的服务器上,将`ssprobe-client`上传到你需要监控的机器上。

### 1. 配置说明

`config.yaml`中包含了你的一些配置,配置说明如下:

```yaml
server:
  token: 123456   # 服务器令牌,用于客户端(被监控的机器)连接服务器时验证身份
# port: 3384      # 服务器监听的端口(默认为3384)

web:
  title: Zolo's Probe  # 监控页面的网站标题
  github: ""  	# 你的github地址(选填)
  telegram: ""  # 你的telegram地址(选填)

notifier:
  telegram:
    enable: true        # 启用Telegram Bot发送通知
    useEmbed: false     # 是否使用自己创建的Bot
    language: chinese   # Bot通知的语言
    botToken: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  # Bot token. 当useEmbed为true时,此项有效。
    userId: 1953745499  # Telegram 用户ID
```

> ⚠️ **`useEmbed`** 是否使用自己创建的Bot？
>
> 若你将此服务端程序部署在中国大陆的服务器上,由于网络原因则无法使用Bot发送通知。当useEmbed设置为false时,你可以指定一个**userId**,使用我已经搭建好的Bot。

#### 🤷‍♂️ 如何创建Bot并获取BotToken?

Telegram添加 `@BotFather`, 向他发送`/newbot `,根据提示完成创建。创建完毕后向他发送 `/mybots` 选择你的Bot,然后获取Token。

#### 🤷‍♀️ 如何获取UserID?

Telegram添加 `@ssprobe_bot`, 向他发送`/me`即可获得你的UserID. (`@ssprobe_bot`就是我们已经创建好的机器人，如果你启用了Telegram机器人并且设置了`useEmbed: false`, 则后续的通知则就会通过此机器人向您发送。)

### 2. 启动程序

(1) 🖥️ 启动服务端程序

确保`config.yaml`和`ssprobe-server`处于同一目录下,使用如下命令启动你的服务端程序:

```bash
chmod a+x ./ssprobe-server && ./ssprobe-server
```

此时,打开 `http://ip:10240` 就可以看到监控页面了。 




(2) 💻 启动客户端程序

```bash
chmod a+x ./ssprobe-client
./ssprobe-client --name=客户端名称 --server=服务器地址 --token=你的Token令牌
```

> ⚠️ 如果你修改了服务的监听端口,则还需要额外指定一个`--port`参数(默认服务器端口为: 3384)

例如: `./ssprobe-client --name=ClientA --server=110.42.133.216 --port=3384 --token=123456`

命令执行完毕后,就可以在监控页面看到这台机器的数据了。

### 3. 配置Https(可选)

以 `nginx` 为例，在 `nginx.conf` 中加入如下内容:

```nginx
server {
    listen       443 ssl;
    server_name  test.onezol.com;

    ssl_certificate      /home/zolo/ssprobe/test.onezol.com_chain.crt;  # SSL证书地址
    ssl_certificate_key  /home/zolo/ssprobe/test.onezol.com_key.key;

    ssl_session_cache    shared:SSL:1m;
    ssl_session_timeout  5m;

    ssl_ciphers  HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers  on;

    location / {
        proxy_pass http://127.0.0.1:10240/;        # 10240是服务端程序ssprobe-server的端口
    }

    location /wss/ {      # 此处 /wss/ 不可修改                         
        proxy_pass http://127.0.0.1:10240/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header X-real-ip $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
```



### 4. 前后端分离部署(可选)

如果你需要将前后端分离部署,你需要将 [`server/static/`](https://github.com/realzolo/ssprobe/tree/master/server/static) 目录下的文件下载下来,并且在`index.html所在目录`创建一个`config.json`文件内容如下:

```json
{
    "site_title":"网站标题",
    "websocket_url": "ws://服务器地址:10240",
    "github": "你的GitHub地址(选填)",
    "telegram": "你的Telegram地址(选填)"
}
```



