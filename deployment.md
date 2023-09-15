## 项目部署所需文件
```
dbip                    # 后端程序
dbip-full-2023-08.mmdb  # 后端程序读取的ip数据库
log.sh                  # 启动脚本依赖的日志服务
run.sh                  # 启动/监控脚本
```
其中数据库的下载链接为
https://download.db-ip.com/key/d5ee0192c292d866ad6418ed17f626ff498a1b90.mmdb

shell文件在代码库中获取 https://e.gitee.com/youquinc/repos/youquinc/ip-api/tree/master/extends/go/ipcc/dbip

将上述代码库下载到本地并进入主目录，后端程序的编译命令`GOOS=linux go build -o dbip ./cmd/web`

将上述文件上传到服务器的某个目录下，放在一起，并且配置以下文件的可执行权限
```shell
$ chmod +x dbip
$ chmod +x run.sh
```
## 编写系统服务文件
`vi /etc/systemd/system/dbip.service`，增加下列内容（其中/xxxx要改为上面上传的目录）：
```
[Unit]
Description=dbip service

[Service]
Type=simple
WorkingDirectory=/xxxx
ExecStart=/xxxx/run.sh
ExecStop=/bin/kill -s TERM $MAINPID

[Install]
WantedBy=multi-user.target
```
如果是第一次编写该系统服务文件，可以直接跳到[下一步](#查看系统服务文件是否被识别)；如果是在原有文件基础上进行了修改，还需要重新加载一下该文件`systemctl daemon-reload`
## 查看系统服务文件是否被识别
```shell
$ systemctl list-unit-files | grep dbip
dbip.service                               disabled        disabled
```
看到上面的结果说明系统服务文件已经被识别并加载了
## 配置系统服务
### 启动服务
```shell
$ systemctl start dbip.service
```
### 查看服务是否成功启动
```shell
$ systemctl status dbip.service -l
● dbip.service - dbip service
     Loaded: loaded (/etc/systemd/system/dbip.service; disabled; preset: disabled)
     Active: active (running) since Thu 2023-08-03 13:41:27 CST; 15s ago
   Main PID: 63260 (run.sh)
      Tasks: 7 (limit: 2266)
     Memory: 3.0M
        CPU: 20ms
     CGroup: /system.slice/dbip.service
             ├─63260 /bin/bash /home/kali/Desktop/ipcc/run.sh
             ├─63276 ./dbip -addr :29952 -mmdb ./dbip-full-2023-08.mmdb
             └─63277 sleep 60

Aug 03 13:41:27 kali systemd[1]: Started dbip.service - dbip service.
Aug 03 13:41:27 kali run.sh[63260]:  [WARN] 2023-08-03 13:41:27 service abnormal
```
### 停止服务（一般不需要）
```shell
$ systemctl stop dbip.service
```
### 将服务设置成开机自启动
```shell
$ systemctl enable dbip.service
Created symlink /etc/systemd/system/multi-user.target.wants/dbip.service → /etc/systemd/system/dbip.service.
```
### 查看服务是否设置成开机自启动
```shell
$ systemctl list-unit-files | grep dbip
dbip.service                               enabled         disabled
```
### 取消服务开机自启动（一般不需要）
```shell
$ systemctl disable dbip.service
```
### 查看服务的控制台日志
```shell
$ journalctl -flu dbip.service
```