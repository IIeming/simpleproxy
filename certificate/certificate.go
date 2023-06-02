package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func GenerateSelfSignedCertKey(domain string) {
	// 创建一个RSA私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln("Failed to generate private key:", err)
	}

	// 创建一个自签名的证书
	template := x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().Unix()),
		Subject:               pkix.Name{CommonName: domain},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// 使用私钥签名证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalln("Failed to create certificate:", err)
	}

	// 将证书和私钥保存到文件
	certOut, err := os.Create("./certs/" + domain + ".crt")
	if err != nil {
		log.Fatalf("Failed to create %s.crt: %s\n", domain, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Printf("Written %s.crt\n", domain)

	keyOut, err := os.OpenFile("./certs/"+domain+".key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open %s.key for writing: %s\n", domain, err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	keyOut.Close()
	log.Printf("Written %s.key\n", domain)
}

func isUbuntu() bool {
	_, err := exec.LookPath("apt")
	return err == nil
}

func isCentos() bool {
	_, err := exec.LookPath("yum")
	return err == nil
}

func execUbuntuUpdateCaCert() {
	// 执行update-ca-certificates命令
	err := exec.Command("/usr/sbin/update-ca-certificates").Run()
	if err != nil {
		log.Println("update-ca-certificates command execution failed", err)
	} else {
		log.Println("update-ca-certificates command execution success")
	}
}

func execCentosUpdateCaCert() {
	// 执行update-ca-certificates命令
	err := exec.Command("/usr/bin/update-ca-trust").Run()
	if err != nil {
		log.Println("update-ca-trust command execution failed", err)
	} else {
		log.Println("update-ca-trust command execution success")
	}
}

func moveUbuntuCertificate(domain string) {
	// 定义源文件与目标文件路径
	srcCertFile := "./certs/" + domain
	destCertFile := "/usr/local/share/ca-certificates/simpleproxy_" + domain
	if _, err := os.Stat(destCertFile); os.IsNotExist(err) {
		// 打开源文件
		src, err := os.Open(srcCertFile)
		if err != nil {
			log.Fatalln("failed open source cert files:", err)
		}
		defer src.Close()
		// 打开目标文件
		dest, err := os.Create(destCertFile)
		if err != nil {
			log.Fatalln("failed open destination cert files:", err)
		}
		defer dest.Close()
		// 复制文件内容
		_, err = io.Copy(dest, src)
		if err != nil {
			log.Fatalln("failed copying cert file:", err)
		}
		// 执行update-ca-certificates命令
		execUbuntuUpdateCaCert()
	} else {
		log.Println("The certificate file already exists:", destCertFile)
	}
}

func moveCentosCertificate(domain string) {
	// 定义源文件与目标文件路径
	srcCertFile := "./certs/" + domain
	destCertFile := "/etc/pki/ca-trust/source/anchors/simpleproxy_" + domain
	if _, err := os.Stat(destCertFile); os.IsNotExist(err) {
		// 打开源文件
		src, err := os.Open(srcCertFile)
		if err != nil {
			log.Fatalln("failed open source cert files:", err)
		}
		defer src.Close()
		// 打开目标文件
		dest, err := os.Create(destCertFile)
		if err != nil {
			log.Fatalln("failed open destination cert files:", err)
		}
		defer dest.Close()
		// 复制文件内容
		_, err = io.Copy(dest, src)
		if err != nil {
			log.Fatalln("failed copying cert file:", err)
		}
		// 执行update-ca-certificates命令
		execCentosUpdateCaCert()
	} else {
		log.Println("The certificate file already exists:", destCertFile)
	}
}

func Init(domain string) {
	certFileName := domain + ".crt"
	if _, err := os.Stat("./certs/" + certFileName); os.IsNotExist(err) {
		GenerateSelfSignedCertKey(domain)
	} else {
		log.Println("warning the certificate file already exists")
	}
	if isUbuntu() {
		moveUbuntuCertificate(certFileName)
	} else if isCentos() {
		moveCentosCertificate(certFileName)
	} else {
		moveUbuntuCertificate(certFileName)
	}

}

func DelUbuntuCert() {
	// 定义文件路径
	pattern := "/usr/local/share/ca-certificates/simpleproxy_*"

	// 查找匹配的文件
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("glob error: %v\n", err)
	}

	// 删除匹配的文件
	for _, file := range files {
		// fmt.Println("test_file", file)
		if err := os.Remove(file); err != nil {
			log.Fatalf("remove file %s error: %v\n", file, err)
		} else {
			log.Printf("remove file %s success\n", file)
		}
	}

	// 执行命令更新证书命令
	execUbuntuUpdateCaCert()
}

func DelCentosCert() {
	// 定义文件路径
	pattern := "/etc/pki/ca-trust/source/anchors/simpleproxy_*"

	// 查找匹配的文件
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("glob error: %v\n", err)
	}

	// 删除匹配的文件
	for _, file := range files {
		// fmt.Println("test_file", file)
		if err := os.Remove(file); err != nil {
			log.Fatalf("remove file %s error: %v\n", file, err)
		} else {
			log.Printf("remove file %s success\n", file)
		}
	}

	// 执行命令更新证书命令
	execCentosUpdateCaCert()
}

func DelCert() {
	if isUbuntu() {
		DelUbuntuCert()
	} else if isCentos() {
		DelCentosCert()
	} else {
		DelUbuntuCert()
	}
}
