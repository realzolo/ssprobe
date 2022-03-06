<div align="center">
  <h1><a href="https://status.onezol.com">SSProbe</a></h1>
  <p><em>A server status monitoring program, powered by Golang.</em></p>
  <a href="https://github.com/realzolo/ssprobe/blob/master/README_CN.md"><img src="https://img.shields.io/badge/简体中文-000000?style=flat&logo=Academia&logoColor=%23FFFFFF" alt="Language" /><a/>
  <img src="https://img.shields.io/badge/Golang-black?style=flat&logo=Go&logoColor=white" alt="Golang" />
  <img src="https://img.shields.io/badge/React.js-black?style=flat&logo=React&logoColor=white" alt="React.js" />
  <img src="https://img.shields.io/github/last-commit/realzolo/ssprobe?&label=Last%20commit&color=CF2B5B&labelColor=black&logo=github" alt="commit"/>
  <img src="https://img.shields.io/github/stars/realzolo/ssprobe?color=%2300979D&label=Starts&labelColor=black&logo=Apache%20Spark&logoColor=%23FFFFFF" alt="stars"/>
<br/><br/>
</div>


## 🏅Introduce

[SSProbe](https://github.com/realzolo/ssprobe) is a server status monitor, it provides a visual interface, real-time for you to display server status, such as CPU usage, memory usage, network speed and so on.  

* Low CPU and memory usage.📉
* The client is very easy to deploy.🚀
* Super fast ⚡ and responsive 💦
* Takes less than 10 minutes to setup ⏱️
* After the server restarts, the client automatically reconnects to the server.⚒️

### Demo

Live demo at [Zolo's SSProbe](https://status.onezol.com/).

![](https://image.onezol.com/img/ssprobe.png)

## 🎉Installation & Usage

On the [release](https://github.com/realzolo/ssprobe/releases) page, find the zip file corresponding to the system version, unzip it and upload the files `ssprobe-server` and `config.yaml` to your server, upload `ssprobe-client` to the machine you need to monitor.

### 1. Configuration

The `config.yaml` contains some of your configuration, which is described below.

```yaml
server:
  token: 123456   # Used to verify the identity of the client (the monitored machine) when connecting to the server.
  port: 3384      # The port the server is listening on.
  websocketPort: 9000    # If you do not have a separate front-end and back-end deployment, please leave this as default.

web:
  enable: true    # Enable web services. If you need separate front and back-end deployment, you can set this to false.
  title: Zolo's Server Monitor  # The page's site title.

notifier:
  telegram:
    enable: true        # Enable Telegram Bot to send notifications.
    useEmbed: false     # Whether to use the Bot created by yourself.
    language: chinese   # Language of Bot notifications
    botToken: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  # Bot token.(This is valid when useEmbed is true)
    userId: 1953745499  # Telegram's user ID
```

> ⚠️ **`useEmbed`**: Do you use a Bot you created?
> 
> If you deploy this server-side application on a server in mainland China, you will not be able to use Bot to send notifications due to network reasons.
> When `useEmbed` is set to `false`, you can specify a UserID, using the Bot we've already created.

#### 🤷‍♂️ How to create a Bot and get a BotToken?

Use Telegram to search for `@BotFather`, send him `/newbot`,follow the steps to complete the creation, then send him `/mybots` to select your bot and get the token.

#### 🤷‍♀️ How to get the UserID?

Use Telegram to search for `@ssprobe_bot`, Send `/me` to him to get your UserID. (`@ssprobe_bot` is the bot we have created. If you have enabled the Telegram bot and set `useEmbed: false`, then subsequent notifications will be sent to you via this bot.)

### 2. Start-up program

(1) 🖥️ Server-side program

Make sure that `config.yaml` and `ssprobe-server` are in the same directory, and start your program with the following command:

```bash
chmod a+x ./ssprobe-server && ./ssprobe-server
```

At this moment，Open `http://ip:10240` and you will see the monitoring page if you set `web.enable: true`

(2) 💻 Client program

```bash
chmod a+x ./ssprobe-client
./client --name=ClientName --server=ServerAddress --token=YourServerToken
```

> ⚠️ If you modify the server's listening port, you need to specify an additional `--port' parameter.(the default server port is: 3384)

Such as `./ssprobe-client --name=ClientA --server=110.42.133.216 --port=3384 --token=123456`

Once the command is executed, you can see the data on the monitoring page.

### 3. Separate front and back-end deployment (optional)

If you need to deploy the front-end and back-end separately, you need to download the files in the [`server/static/`](https://github.com/realzolo/ssprobe/tree/master/server/static) directory and create a `config.json` file in the `index.html directory` with the following contents:

```json
{
  "SITE_TITLE":"site title",
  "WEBSOCKET_URL": "ws://ServerAddress:9000(WebsocketPort)"
}
```
