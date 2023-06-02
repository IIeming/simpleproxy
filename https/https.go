package https

import (
	"log"
	"os/exec"
	"simpleproxy/certificate"
	"simpleproxy/client"
	"simpleproxy/server"
	"simpleproxy/setting"
	"strings"
)

func Init(srcaddr, destaddr, responsebody, responsecode *string) {
	// 读取配置文件
	vp := setting.Init()

	// 创建证书文件
	certificate.Init(*srcaddr)

	// 定义查询命令
	cmd := exec.Command("bash", "-c", "(netstat -tnlp | grep 8080) 2> /dev/null")

	// 执行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		log.Println("failed to perform the command:", err)
	}

	// 判断进程是否已存在
	if len(output) == 0 {
		log.Println("simpleproxy process is not running")
		// 启动服务
		server.Init(srcaddr, destaddr, responsebody, responsecode, vp)
	} else if strings.Contains(string(output), "simpleproxy") {
		log.Println("simpleproxy process is running")
		// 添加配置
		client.Init(srcaddr, destaddr, responsebody, responsecode, vp)
	} else {
		log.Fatalln("filed listen tcp 8080, address already in use")
	}
}
