package main

import (
	"fmt"
	"log"
	"os"
	"simpleproxy/certificate"
	"simpleproxy/database"
	"simpleproxy/https"
	"simpleproxy/iiptables"
	"simpleproxy/tools"
	"strings"

	"github.com/alecthomas/kingpin"
)

// 定义命令行参数
var (
	add          = kingpin.Command("add", "Refresh the local rules.")
	protocol     = add.Flag("protocol", "Agent HTTPS or HTTP protocol.").Short('p').Default("http").String()
	srcaddr      = add.Flag("srcaddr", "Need the agent's domain name.").Short('s').Required().String()
	destaddr     = add.Flag("destaddr", "Need the agent's destination address.").Short('d').Default("127.0.0.1:8080").String()
	responsebody = add.Flag("responsebody", "Custom content accordingly (Is only effective when \"--destaddr\" is empty).").Default("hello world").String()
	responsecode = add.Flag("responsecode", "Custom status code accordingly  (Is only effective when \"--destaddr\" is empty).").Default("200").String()
	list         = kingpin.Command("list", "List all rules.")
	listIpt      = list.Flag("ipt", "List the iptables rules.").Short('i').Bool()
	refresh      = kingpin.Command("refresh", "Refresh the local rules.")
	refreshAll   = refresh.Flag("all", "Delete all the rules").Short('a').Bool()
	refreshNums  = refresh.Flag("num", "Delete the specified rules").Short('n').String()
)

func main() {
	// 定义程序版本
	kingpin.Version("simpleproxy app version is 1.1")

	// 初始化数据库
	db := database.Init()
	defer db.Close()

	// 初始化防火墙
	ipt := iiptables.Init()

	// 解析命令行参数
	switch kingpin.Parse() {
	case "add":
		// 校验输入格式
		tools.IsFormatMatch(*srcaddr, *destaddr)
		// 执行sql添加语句
		sumStr := *protocol + *srcaddr + *destaddr + *responsebody + *responsecode
		// fmt.Println("test_sumstr", sumStr)
		hexStr := tools.GetHexStr(sumStr)
		isExist := *database.QuertDB(db, srcaddr)
		srcUrl := strings.Split(*srcaddr, ":")

		if *protocol == "https" {
			// 初始化https服务
			https.Init(&srcUrl[0], destaddr, responsebody, responsecode)
		}

		if len(isExist) == 0 {
			// 插入数据库数据
			database.InsertDB(db, protocol, srcaddr, destaddr, responsebody, responsecode, &hexStr)
			// 创建iptables规则字符串
			tools.AddRules(ipt, protocol, destaddr, srcUrl)
		} else if isExist[0].DestAddr == *destaddr {
			log.Println("Duplicate data, do not modify")
			// 更新数据库数据
			database.UpdataDB(db, protocol, destaddr, responsebody, responsecode, srcaddr, &hexStr)
		} else {
			tools.AddRules(ipt, protocol, destaddr, srcUrl)
			oldSurl := strings.Split(isExist[0].SrcAddr, ":")
			// 删除重复的iptables规则
			if *protocol == "https" {
				oldRuleSpec := []string{"-p", "tcp", "--dport", oldSurl[1], "-d", oldSurl[0], "-j", "DNAT", "--to-destination", "127.0.0.1:8080"}
				iiptables.Delete(ipt, oldRuleSpec)
			} else {
				oldRuleSpec := []string{"-p", "tcp", "--dport", oldSurl[1], "-d", oldSurl[0], "-j", "DNAT", "--to-destination", isExist[0].DestAddr}
				iiptables.Delete(ipt, oldRuleSpec)
			}

			// 更新数据库数据
			database.UpdataDB(db, protocol, destaddr, responsebody, responsecode, srcaddr, &hexStr)
		}

	case "refresh":
		if *refreshAll {
			// 删除所有iptables规则
			var input string
			fmt.Printf("Will delete all rules on the OUTPUT chain under the NAT table. Are you sure you want to continue? (y/n): ")
			fmt.Scanln(&input)
			if input == "y" || input == "yes" {
				// 删除iptables规则
				iiptables.DeleteAll(ipt)
				// 删除本地db文件
				if err := os.Remove("./proxys.db"); err != nil {
					log.Fatal(err)
				}
				// 删除创建的证书文件
				certificate.DelCert()
				// 关闭http进程
				tools.KillProcess("./run/http.pid")
				log.Println("The rule list has been cleared.")
			} else {
				log.Fatalln("Aborting...")
			}

		} else if len(*refreshNums) != 0 {
			// 删除指定iptables规则
			delData := *database.QuertNumDB(db, *refreshNums)
			delUrl := strings.Split(delData[0].SrcAddr, ":")
			delRuleSpec := []string{"-p", "tcp", "--dport", delUrl[1], "-d", delUrl[0], "-j", "DNAT", "--to-destination", delData[0].DestAddr}
			iiptables.Delete(ipt, delRuleSpec)

			// 删除数据库数据
			database.DeleteDbHexStr(db, *refreshNums)
			log.Println("The rule deleted successfully")
		} else {
			fmt.Println(`Warning please input "-a" or the "-n" option`)
		}
	case "list":
		if *listIpt {
			// 列出当前iptables规则nat表的output链的所有规则
			iiptables.List(ipt)
		} else {
			// 列出已有的规则
			// proxy := database.QuertDB(db, nil)
			fmt.Printf("%-36s\t%-10s\t%-25s\tTo\t%s\n", "Id", "Protocol", "Srcaddr", "Destaddr")
			for _, v := range *database.QuertDB(db, nil) {
				fmt.Printf("%-36s\t%-10s\t%-25s\tto\t%s\n", v.HexStr, v.Protocol, v.SrcAddr, v.DestAddr)
			}
		}

	}
	// fmt.Println("simpleproxy_end")
}
