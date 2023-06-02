package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CertFile     string `yaml:"certfile"`
	CertKey      string `yaml:"certkey"`
	DestAddr     string `yaml:"destaddr"`
	Protocol     string `yaml:"protocol"`
	ResponseBody string `yaml:"responsebody"`
	ResponseCode string `yaml:"responsecode"`
	SrcAddr      string `yaml:"srcaddr"`
}

func createReverseProxy(targetURL string) *httputil.ReverseProxy {
	// 解析目标URL
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	// 创建反向代理
	return httputil.NewSingleHostReverseProxy(target)
}

func Init(srcaddr, destaddr, responsebody, responsecode *string, v *viper.Viper) {
	// 初始化配置文件
	config := map[string]Config{
		*srcaddr: struct {
			CertFile     string `yaml:"certfile"`
			CertKey      string `yaml:"certkey"`
			DestAddr     string `yaml:"destaddr"`
			Protocol     string `yaml:"protocol"`
			ResponseBody string `yaml:"responsebody"`
			ResponseCode string `yaml:"responsecode"`
			SrcAddr      string `yaml:"srcaddr"`
		}{
			CertFile:     *srcaddr + ".crt",
			CertKey:      *srcaddr + ".key",
			DestAddr:     *destaddr,
			Protocol:     "https",
			ResponseBody: *responsebody,
			ResponseCode: *responsecode,
			SrcAddr:      *srcaddr + ":443",
		},
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = os.WriteFile("simpleproxy.yaml", data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 解析默认证书文件
	defaultCert, err := tls.LoadX509KeyPair("./certs/localhost.crt", "./certs/localhost.key")
	if err != nil {
		fmt.Println("cert parse failed: ", err)
	}

	// 配置TLS
	certConfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			hostMap := v.GetStringMapString(info.ServerName)
			if len(hostMap) > 0 {
				// 从证书文件中解析公钥
				// fmt.Println("test001:" + "./certs/" + string(hostMap["certfile"]))
				certfile, err := tls.LoadX509KeyPair("./certs/"+string(hostMap["certfile"]), "./certs/"+hostMap["certkey"])
				if err != nil {
					fmt.Println("cert parse failed: ", err)
				}
				// fmt.Println("test-有这个值", value.Certfile, value.Certkey)
				return &certfile, err
			} else {
				log.Println("failed to certificate file does not exist")
				// fmt.Println("test06", &defaultCert)
				return &defaultCert, nil
			}
		},
	}

	// 添加/路径代理处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname := strings.Split(r.Host, ":")
		hostsMap := v.GetStringMapString(hostname[0])
		// fmt.Println("test333", hostsMap)
		if len(hostsMap) > 1 {
			// 创建反向代理
			if hostsMap["destaddr"] != "127.0.0.1:8080" {
				// fmt.Println("test00000000", hostsMap["destaddr"])
				proxy := createReverseProxy("http://" + hostsMap["destaddr"])
				proxy.ServeHTTP(w, r)
			} else {
				code, _ := strconv.Atoi(hostsMap["responsecode"])
				w.WriteHeader(code)
				fmt.Fprintln(w, hostsMap["responsebody"])
			}
		} else {
			w.WriteHeader(200)
			fmt.Fprintln(w, "hello world")
		}
	})

	// 定义 daemon 配置
	daemonConfig := &daemon.Context{
		PidFileName: "./run/http.pid",
		PidFilePerm: 0644,
		LogFileName: "./logs/http-server.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	// 启动 daemon
	child, err := daemonConfig.Reborn()
	if err != nil {
		log.Fatalf("daemon error: %v\n", err)
	}
	if child != nil {
		return
	}
	defer daemonConfig.Release()

	// 自定义服务器
	// fmt.Println("test_server")
	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: certConfig,
	}

	// 启动服务器
	log.Printf("Starting server on %s\n", server.Addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
