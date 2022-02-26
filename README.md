## 👻Introduce

[server-minotor](https://github.com/realzolo/server-monitor) is a server monitor that provides a web page that displays your server data in real time, such as CPU usage, memory usage, network speed, etc.  

## 🎉Installation & Usage

### Server

Download the corresponding `server program` on the [release page](https://github.com/realzolo/server-monitor/releases). Create a `config.yaml` file with the following contents. 

```yaml
# config.yaml
token: 123456   # Used to authenticate identity.
port:	
  server: 3384   # Server port
  web-api: 9000  # Http request port
```

Place the two files in the same directory, and execute the program.

```bash
chmod 755 server-monitor-linux-server
./server-monitor-linux-server
```



### Client

Download the corresponding `client program` on the [release page](https://github.com/realzolo/server-monitor/releases). Use the following command to execute the program.

```bash
chmod 755 server-monitor-linux-client
./server-monitor-linux-client --name=CLIENT_NAME --server=SERVER_ADDRESS --port=SERVER_PORT --token=YOUR_TOKEN
```

For example, `./server-monitor-linux-client --name=ClientA --server=110.42.133.216 --port=3384 --token=123456`

### Web

Deploy files in a web [directory](https://github.com/realzolo/server-monitor/tree/master/web) to your HTTP server or static web hosting. You can change your configuration information in [config.json](https://github.com/realzolo/server-monitor/blob/master/web/config.json). After the deployment, the monitoring page is displayed.

