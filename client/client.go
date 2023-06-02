package client

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init(srcaddr, destaddr, responsebody, responsecode *string, v *viper.Viper) {
	// 创建配置文件
	v.Set(*srcaddr, map[string]interface{}{
		"protocol":     "https",
		"srcaddr":      *srcaddr + ":443",
		"certfile":     *srcaddr + ".crt",
		"certkey":      *srcaddr + ".key",
		"destaddr":     *destaddr,
		"responsebody": *responsebody,
		"responsecode": *responsecode,
	})

	// 合并修改后的配置
	v.MergeInConfig()
	if err := v.WriteConfig(); err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}

	// 	fmt.Println("test00", v.GetStringMap("www.baidu.com"))
	// 	fmt.Println("test01", v.GetStringMap("www"))
}
