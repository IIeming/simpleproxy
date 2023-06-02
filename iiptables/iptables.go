package iiptables

import (
	"fmt"
	"log"

	"github.com/coreos/go-iptables/iptables"
)

var (
	chain = "OUTPUT"
	table = "nat"
)

// 初始化防火墙
func Init() *iptables.IPTables {
	// 创建一个 iptables 实例
	ipt, err := iptables.New()
	if err != nil {
		log.Println("Error creating iptables instance:", err)
	}
	return ipt
}

// 添加一条规则
func Add(ipt *iptables.IPTables, ruleSpec []string) {
	// ruleSpec := []string{"-p", "tcp", "--dport", "443", "-d", "www.baidu.com", "-j", "DNAT", "--to-destination", "127.0.0.1:8080"}
	err := ipt.AppendUnique(table, chain, ruleSpec...)
	if err != nil {
		log.Fatalln("Error adding iptables rule:", err)
	}
	log.Println("iptables rule added successfully")
}

// 列出所有规则
func List(ipt *iptables.IPTables) {
	rules, err := ipt.List(table, chain)
	if err != nil {
		log.Fatalln("Error listing iptables rules:", err)
	}

	// log.Println("iptables rules:", rules)
	for _, rule := range rules {
		fmt.Println(rule)
	}
}

// 删除规则
func Delete(ipt *iptables.IPTables, ruleSpec []string) {
	if err := ipt.Delete(table, chain, ruleSpec...); err != nil {
		log.Fatalln("Error deleting iptables rule:", err)
	}
	log.Println("iptables rule deleted successfully")
}

// 清空nat表中OUTPUT链的所有规则
func DeleteAll(ipt *iptables.IPTables) {
	if err := ipt.ClearChain("nat", "OUTPUT"); err != nil {
		log.Fatalln(err)
	}
}
