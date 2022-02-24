## 👻Introduce

[server-minotor](/) is a server monitor that provides a web page that displays your server data in real time, such as CPU usage, memory usage, network speed, etc.  

## 🎉Installation & Usage

### Server

Open 'server/main.go' to modify 'token', compile the server directory to generate an executable, then run it.  

```bash
go build -o server
./server
```

### Client

Copy the Python script file from the client directory to your client (the server you want to monitor), run it, and fill in the information as required. 

```py
python ./client.py
-> Server address: xxx.xxx.xxx.xxx   # Server IP Address
-> Client name: my_first_client
-> Token: my_token  # This parameter must be consistent with the Token configured on the server.
```

### Web

Deploy files in a Web directory to your HTTP server or static web hosting. After the deployment is complete, the monitoring page is displayed.  