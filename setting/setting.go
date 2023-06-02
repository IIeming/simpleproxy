package setting

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() *viper.Viper {
	v := viper.NewWithOptions(viper.KeyDelimiter("/"))
	// v.SetConfigFile("./" + conf)
	v.SetConfigName("simpleproxy")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("failed error config file: %s\n", err)
	}

	// fmt.Println("test11", v.GetString("www.baidu.com/certfile"))

	// 监控配置文件变化
	v.WatchConfig()
	// fmt.Println("test12 走到这里额")
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("Configuration change......")
		// fmt.Println("test11", v.GetStringMap("www.baidu.com")["destaddr"])
	})
	return v
}
