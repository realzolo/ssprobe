## 👻简介

[SSProbe](https://github.com/realzolo/ssprobe) 是一个服务器监控程序，它提供了一个可视化的界面，实时显示你的服务器状态，如CPU占用率、内存使用情况和网络速度等。

![](https://image.onezol.com/img/ssprobe.jpg)

## 🎉下载和使用

### 服务端

在[release页面](https://github.com/realzolo/ssprobe/releases)找到对应系统的**服务端版本**并下载到你的服务器上，你可以在`config.yaml`文件中配置你的端口和Token令牌。

```yaml
# config.yaml
token: 123456   # 用于验证你的客户端
port:	
  server: 3384   # 服务器端口
  web-api: 9000  # Http请求端口
```

确保两个文件在同一个目录中，然后执行程序。  

```bash
chmod a+x server
./server

# 或者后台运行
nohup ./server &
```



### 客户端

在[release页面](https://github.com/realzolo/ssprobe/releases)找到对应系统的**客户端版本**并下载到你的服务器上，并使用如下命令启动你的客户端程序：

```bash
chmod a+x ./client
./client --name=客户端名称 --server=服务器地址 --port=服务器端口 --token=你的Token令牌

# 或者后台运行
nohup ./client --name=客户端名称 --server=服务器地址 --port=服务器端口 --token=你的Token令牌 &
```

例如: `./client --name=ClientA --server=110.42.133.216 --port=3384 --token=123456`

### Web端

将[web目录](https://github.com/realzolo/ssprobe/tree/master/web)中的文件部署到你的HTTP服务器或其他静态网站托管平台。你可以在[config.json](https://github.com/realzolo/ssprobe/blob/master/web/config.json)中更改你的配置信息。 部署完成后，就可以进入监控页面了。  

```json
{
    "API": "ws://服务器地址:服务器端口/json",   
    "SITE_TITLE":"这是网站标题" 
}
```



