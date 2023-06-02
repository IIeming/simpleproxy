package tools

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"regexp"
	"simpleproxy/iiptables"
	"strconv"

	"github.com/coreos/go-iptables/iptables"
)

func IsFormatMatch(inputs ...string) {
	// 定义正则表达式
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\.]+:[0-9]+$`)

	// 判断输入字符串是否符合正则表达式
	for _, str := range inputs {
		if !pattern.MatchString(str) {
			log.Fatalln("srcaddr or destaddr format does not match")
		}
	}
}

func KillProcess(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	pid, err := strconv.Atoi(string(data))
	process, err := os.FindProcess(pid)
	if err = process.Kill(); err != nil {
		log.Printf("Failed to parse pid: %v\n", err)
	} else {
		log.Println("The HTTP service has shutdown success")
	}
	if err = os.Remove(path); err != nil {
		log.Fatal(err)
	}
}

func GetHexStr(str string) string {
	// 计算字符串的 md5 哈希值
	hash := md5.Sum([]byte(str))

	// 将哈希值转换为十六进制字符串
	hexStr := fmt.Sprintf("%x", hash)

	return hexStr
	// fmt.Printf("string: %s\nhash: %s\n", str, hexStr)
}

func AddRules(ipt *iptables.IPTables, protocol, destaddr *string, srcUrl []string) {
	// 创建iptables规则字符串
	if *protocol == "https" {
		ruleSpec := []string{"-p", "tcp", "--dport", srcUrl[1], "-d", srcUrl[0], "-j", "DNAT", "--to-destination", "127.0.0.1:8080"}
		iiptables.Add(ipt, ruleSpec)
	} else {
		ruleSpec := []string{"-p", "tcp", "--dport", srcUrl[1], "-d", srcUrl[0], "-j", "DNAT", "--to-destination", *destaddr}
		iiptables.Add(ipt, ruleSpec)
	}
}
