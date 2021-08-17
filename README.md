# 介绍
MySQL-SOCKS5-client是一款使用socks5协议的数据库代理工具

传统的mysql客户端没有socks5的代理，而是只有ssh的代理。使用此款产品是在本地启动此软件提供mysql连接socks5的中转服务，这样mysql客户端就无需知道背后的服务也能访问数据库。

![运行原理图.png](https://i.loli.net/2021/08/17/k5RVi82zNXZpcD6.png)

### 支持的数据库
以下数据库是经过使用得出的，其他类型的数据库也可以尝试看看~
+ MySQL5.6
+ MySQL5.7
+ ClickHouse 21.4.4.30

# 使用

### 编译
> 这儿是windows编译出Linux的可执行文件，经测试适用于centos，其他Linux版本可自行尝试
1. 进入项目文件
2. 运行 `build.cmd`

可执行文件的同目录需要配置文件 `config.json`，同一个配置可配置多个需要转发的MySQL源
```
[
  {
    "name": "xxx源",                        # 目标源名称
    "socks5_server_ip": "172.18.12.70",     # socks5服务器的ip
    "socks5_server_port": "10080",          # socks5服务器的端口
    "mysql_ip": "10.104.20.42",             # 远程MySQL的地址（远程socks5服务器能连的ip）
    "mysql_port": "3306",                   # 程MySQL的端口
    "listen_addr": "0.0.0.0",               # 本地监听的地址
    "listen_port": "6666"                   # 本地监听的端口
  },
  ...
]
```
### 运行
1. 进入可执行文件的目录，确保配置文件也在这个目录！！！
2. 执行 ./build 即可

### 一些说明
+ 使用命令行的话，不会后台执行，需要使用screen/nohub/supervisor等方式运行
+ 日志只直接输出的，可以重定向到本地等方式

