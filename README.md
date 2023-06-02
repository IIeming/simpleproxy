# simpleproxy

一个简单的http和https代理工具



### 概述

在日常开发过程中，经常遇到需要代理的问题。例如想要访问企微的api，但是想模拟下接口不通的情况下如何处理，这个时候一般都是去修改本地hosts，或者配个nginx代理。但是都比较麻烦，simpleproxy可以简单解决这个问题，它支持http和https代理，https模式下支持自定义状态码和返回值。具体使用示例如下。

支持`ubuntu 20`和`centos 7`操作系统，需要`root`权限执行

### 使用说明

**构建可执行文件**

```bash
# 构建可执行文件 (simpleproxy名称尽量不要该，代码里有相关判断。可根据情况在代码里修改后在修改)
go mod tidy
go build -o simpleproxy main.go
```

启动程序需要一些必要的配置文件，整体如下

```bash
├── certs                # https 证书存放目录
│   ├── localhost.crt
│   └── localhost.key
├── logs
│   └── http-server.log  # http服务日志
├── run                  # 存放pid文件
├── simpleproxy          # 程序
└── simpleproxy.yaml     # https 代理配置文件（无须手动改）
```

**命令行参数**

```bash
usage: simpleproxy [<flags>] <command> [<args> ...]

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  add --srcaddr=SRCADDR [<flags>]
    Refresh the local rules.

  list [<flags>]
    List all rules.

  refresh [<flags>]
    Refresh the local rules.
```

**添加一条http规则**

```bash
# 将发往192.168.40.7:80的http请求代理到127.0.0.1:80
./simpleproxy add -s 192.168.40.7:80 -d 127.0.0.1:80
```

**添加一条https规则**

```bash
# 将发往www.baidu.com的https请求代理到默认的127.0.0.1:8080,也可以加 -d 自定义目标地址
./simpleproxy add -p https -s www.baidu.com:443 --responsebody="hello baidu"
```

结果如下

```bash
curl https://www.baidu.com
hello baidu
```

**列出已添加规则**

```bash
# 展示当前添加的规则
./simpleproxy list
Id                                  	Protocol  	Srcaddr                  	To	Destaddr
0e4071b1c684b221d3e224f020a7da25    	http      	192.168.40.7:80          	to	127.0.0.1:80
434f10d8497c906345812301abcc58ff    	https     	www.baidu.com:443        	to	127.0.0.1:8080
# 展示当前iptables规则
./simpleproxy list -i
-P OUTPUT ACCEPT
-A OUTPUT -d 192.168.40.7/32 -p tcp -m tcp --dport 80 -j DNAT --to-destination 127.0.0.1:80
-A OUTPUT -d 127.0.0.1/32 -p tcp -m tcp --dport 443 -j DNAT --to-destination 127.0.0.1:8080
```

**删除指定规则**

```bash
# -n 选项后填写list展示出来的规则id
./simpleproxy refresh -n 0e4071b1c684b221d3e224f020a7da25
# 再次查看当前规则
./simpleproxy list
Id                                  	Protocol  	Srcaddr                  	To	Destaddr
434f10d8497c906345812301abcc58ff    	https     	www.baidu.com:443        	to	127.0.0.1:8080
```

**删除所有规则**

```bash
# 会清空以有的配置和iptables的nat表output链下的所有规则（注意:会无差别清空所有，如果有自定义的规则还请先备份，谨慎操作)
./simpleproxy refresh -a
```

至此基本使用已经了解，如果是由中有什么问题可以提交issues。



**在此特别感谢知乎的[A7kaou](https://www.zhihu.com/people/ding-xin-85-36)大佬给我的灵感**

### 参考链接

https://zhuanlan.zhihu.com/p/507883101