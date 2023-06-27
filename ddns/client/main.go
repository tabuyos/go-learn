package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kardianos/service"
	"gopkg.in/ini.v1"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func loadConfig() (domains []string, remdoms []string) {
	cfg, err := ini.Load("report.ini")
	if err != nil {
		log.Println(err)
		return
	}

	domainSection := cfg.Section("domain")
	remdomSection := cfg.Section("remdom")
	for _, key := range domainSection.KeyStrings() {
		domains = append(domains, domainSection.Key(key).String())
	}
	for _, key := range remdomSection.KeyStrings() {
		remdoms = append(remdoms, remdomSection.Key(key).String())
	}
	return
}

func fetchIp() (ip string, err error) {
	ip = ""
	defer func() {
		if len(ip) == 0 {
			err = errors.New("未查找到 ip 信息")
		}
	}()

	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		return
	}

	address, err := net.LookupHost(hostname)
	if err != nil {
		log.Println(err)
		return
	}

	for _, addr := range address {
		ip = addr
		if len(ip) > 0 {
			break
		}
	}
	return
}

type SystemService struct{}

func (ss *SystemService) Start(s service.Service) error {
	log.Println("即将启动.......")
	go ss.run()
	return nil
}

func (ss *SystemService) run() {
	ip, err := fetchIp()
	if err != nil {
		log.Println("获取 ip 失败")
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	currentIp := ip
	doRequest(currentIp)

	for range ticker.C {
		ip, err := fetchIp()
		if err != nil {
			log.Println("获取 ip 失败")
		}
		if strings.Compare(currentIp, ip) == 0 {
			log.Println("ip 未发生变化")
			continue
		}

		currentIp = ip
		go doRequest(currentIp)
	}
}

func doRequest(ip string) {
	body := make(map[string]interface{})
	domains, remdoms := loadConfig()

	body["ip"] = ip
	if len(domains) != 0 {
		body["domain"] = domains
	}
	if len(remdoms) != 0 {
		body["remdom"] = remdoms
	}
	marshal, _ := json.Marshal(body)

	log.Println(string(marshal))

	_, _ = http.Post("http://localhost:8080/report", "application/json", bytes.NewBuffer(marshal))
}

func (ss *SystemService) Stop(s service.Service) error {
	log.Println("即将停止.......")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "domain-ip-report-service",
		DisplayName: "Domain-IP Report Service",
		Description: "Domain-IP Report Service",
	}

	ss := &SystemService{}
	s, err := service.New(ss, svcConfig)
	if err != nil {
		log.Printf("服务新建失败, 错误: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		serviceAction := os.Args[1]
		switch serviceAction {
		case "install":
			err := s.Install()
			if err != nil {
				log.Println("安装服务失败: ", err.Error())
			} else {
				log.Println("安装服务成功")
			}
			return
		case "uninstall":
			err := s.Uninstall()
			if err != nil {
				log.Println("卸载服务失败: ", err.Error())
			} else {
				log.Println("卸载服务成功")
			}
			return
		case "start":
			err := s.Start()
			if err != nil {
				log.Println("运行服务失败: ", err.Error())
			} else {
				log.Println("运行服务成功")
			}
			return
		case "stop":
			err := s.Stop()
			if err != nil {
				log.Println("停止服务失败: ", err.Error())
			} else {
				log.Println("停止服务成功")
			}
			return
		}
	}

	err = s.Run()
	if err != nil {
		log.Println(err)
	}
}
