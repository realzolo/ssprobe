<div align="center">
  <h1><a href="https://status.onezol.com">SSProbe</a></h1>
  <p><em>A server status monitoring program, powered by Golang.</em></p>
  <a href="https://github.com/realzolo/ssprobe/blob/master/README_CN.md"><img src="https://img.shields.io/badge/ç®ä½ä¸­æ-000000?style=flat&logo=Academia&logoColor=%23FFFFFF" alt="Language" /><a/>
  <img src="https://img.shields.io/badge/Golang-black?style=flat&logo=Go&logoColor=white" alt="Golang" />
  <img src="https://img.shields.io/badge/React.js-black?style=flat&logo=React&logoColor=white" alt="React.js" />
  <img src="https://img.shields.io/github/last-commit/realzolo/ssprobe?&label=Last%20commit&color=CF2B5B&labelColor=black&logo=github" alt="commit"/>
  <img src="https://img.shields.io/github/stars/realzolo/ssprobe?color=%2300979D&label=Stars&labelColor=black&logo=Apache%20Spark&logoColor=%23FFFFFF" alt="stars"/>
<br/><br/>
</div>



## ðç®ä»

[SSProbe](https://github.com/realzolo/ssprobe) æ¯ä¸æ¬¾æå¡å¨çæ§ç¨åº(ä¹å°±æ¯æè°ç"æ¢é")ï¼å®æä¾äºä¸ä¸ªå¯è§åççé¢ï¼å®æ¶ä¸ºä½ æ¾ç¤ºæå¡å¨ç¶æï¼å¦CPUå ç¨çãåå­ä½¿ç¨æåµåç½ç»éåº¦ç­ã

* ä½CPUååå­å ç¨ãð
* ç¨åºé¨ç½²æä¸ºå®¹æãð
* æ¼äº®çæ°æ®å±ç¤ºé¡µé¢ãð§
* æå¡ç«¯æçº¿æèéå¯æ¶ï¼å®¢æ·ç«¯ä¼èªå¨åè¿(é»è®¤æ¶é´60s)ãâ±ï¸

### Demo

Live demo at [Zolo's SSProbe](https://status.onezol.com/).

![](https://image.onezol.com/img/ssprobe-cn.png)

## ðä¸è½½åä½¿ç¨

å¨ [releaseé¡µé¢](https://github.com/realzolo/ssprobe/releases) æ¾å°å¯¹åºç³»ç»çæ¬çåç¼©å,è§£åä¹åå°`ssprobe-server`å`config.yaml`è¿ä¸¤ä¸ªæä»¶ä¸ä¼ å°ä½ çæå¡å¨ä¸,å°`ssprobe-client`ä¸ä¼ å°ä½ éè¦çæ§çæºå¨ä¸ã

### 1. éç½®è¯´æ

`config.yaml`ä¸­åå«äºä½ çä¸äºéç½®,éç½®è¯´æå¦ä¸:

```yaml
server:
  token: 123456   # æå¡å¨ä»¤ç,ç¨äºå®¢æ·ç«¯(è¢«çæ§çæºå¨)è¿æ¥æå¡å¨æ¶éªè¯èº«ä»½
# port: 3384      # æå¡å¨çå¬çç«¯å£(é»è®¤ä¸º3384)

web:
  title: Zolo's Probe  # çæ§é¡µé¢çç½ç«æ é¢
  github: ""  	# ä½ çgithubå°å(éå¡«)
  telegram: ""  # ä½ çtelegramå°å(éå¡«)

notifier:
  telegram:
    enable: true        # å¯ç¨Telegram Botåééç¥
    useEmbed: false     # æ¯å¦ä½¿ç¨èªå·±åå»ºçBot
    language: chinese   # Botéç¥çè¯­è¨
    botToken: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  # Bot token. å½useEmbedä¸ºtrueæ¶,æ­¤é¡¹ææã
    userId: 1953745499  # Telegram ç¨æ·ID
```

> â ï¸ **`useEmbed`** æ¯å¦ä½¿ç¨èªå·±åå»ºçBotï¼
>
> è¥ä½ å°æ­¤æå¡ç«¯ç¨åºé¨ç½²å¨ä¸­å½å¤§éçæå¡å¨ä¸,ç±äºç½ç»åå åæ æ³ä½¿ç¨Botåééç¥ãå½useEmbedè®¾ç½®ä¸ºfalseæ¶,ä½ å¯ä»¥æå®ä¸ä¸ª**userId**,ä½¿ç¨æå·²ç»æ­å»ºå¥½çBotã

#### ð¤·ââï¸ å¦ä½åå»ºBotå¹¶è·åBotToken?

Telegramæ·»å  `@BotFather`, åä»åé`/newbot `,æ ¹æ®æç¤ºå®æåå»ºãåå»ºå®æ¯ååä»åé `/mybots` éæ©ä½ çBot,ç¶åè·åTokenã

#### ð¤·ââï¸ å¦ä½è·åUserID?

Telegramæ·»å  `@ssprobe_bot`, åä»åé`/me`å³å¯è·å¾ä½ çUserID. (`@ssprobe_bot`å°±æ¯æä»¬å·²ç»åå»ºå¥½çæºå¨äººï¼å¦æä½ å¯ç¨äºTelegramæºå¨äººå¹¶ä¸è®¾ç½®äº`useEmbed: false`, ååç»­çéç¥åå°±ä¼éè¿æ­¤æºå¨äººåæ¨åéã)

### 2. å¯å¨ç¨åº

(1) ð¥ï¸ å¯å¨æå¡ç«¯ç¨åº

ç¡®ä¿`config.yaml`å`ssprobe-server`å¤äºåä¸ç®å½ä¸,ä½¿ç¨å¦ä¸å½ä»¤å¯å¨ä½ çæå¡ç«¯ç¨åº:

```bash
chmod a+x ./ssprobe-server && ./ssprobe-server
```

æ­¤æ¶,æå¼ `http://ip:10240` å°±å¯ä»¥çå°çæ§é¡µé¢äºã 




(2) ð» å¯å¨å®¢æ·ç«¯ç¨åº

```bash
chmod a+x ./ssprobe-client
./ssprobe-client --name=å®¢æ·ç«¯åç§° --server=æå¡å¨å°å --token=ä½ çTokenä»¤ç
```

> â ï¸ å¦æä½ ä¿®æ¹äºæå¡ççå¬ç«¯å£,åè¿éè¦é¢å¤æå®ä¸ä¸ª`--port`åæ°(é»è®¤æå¡å¨ç«¯å£ä¸º: 3384)

ä¾å¦: `./ssprobe-client --name=ClientA --server=110.42.133.216 --port=3384 --token=123456`

å½ä»¤æ§è¡å®æ¯å,å°±å¯ä»¥å¨çæ§é¡µé¢çå°è¿å°æºå¨çæ°æ®äºã

### 3. éç½®Https(å¯é)

ä»¥ `nginx` ä¸ºä¾ï¼å¨ `nginx.conf` ä¸­å å¥å¦ä¸åå®¹:

```nginx
server {
    listen Â  Â  Â  443 ssl;
    server_name Â test.onezol.com;

    ssl_certificate Â  Â  Â /home/zolo/ssprobe/test.onezol.com_chain.crt;  # SSLè¯ä¹¦å°å
    ssl_certificate_key Â /home/zolo/ssprobe/test.onezol.com_key.key;

    ssl_session_cache Â  Â shared:SSL:1m;
    ssl_session_timeout Â 5m;

    ssl_ciphers Â HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers Â on;

    location / {
        proxy_pass http://127.0.0.1:10240/;        # 10240æ¯æå¡ç«¯ç¨åºssprobe-serverçç«¯å£
    }

    location /wss/ {      # æ­¤å¤ /wss/ ä¸å¯ä¿®æ¹                         
        proxy_pass http://127.0.0.1:10240/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header X-real-ip $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
```



### 4. ååç«¯åç¦»é¨ç½²(å¯é)

å¦æä½ éè¦å°ååç«¯åç¦»é¨ç½²,ä½ éè¦å° [`server/static/`](https://github.com/realzolo/ssprobe/tree/master/server/static) ç®å½ä¸çæä»¶ä¸è½½ä¸æ¥,å¹¶ä¸å¨`index.htmlæå¨ç®å½`åå»ºä¸ä¸ª`config.json`æä»¶åå®¹å¦ä¸:

```json
{
    "site_title":"ç½ç«æ é¢",
    "websocket_url": "ws://æå¡å¨å°å:10240",
    "github": "ä½ çGitHubå°å(éå¡«)",
    "telegram": "ä½ çTelegramå°å(éå¡«)"
}
```



