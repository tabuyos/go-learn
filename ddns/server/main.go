package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func report(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	token := request.Header.Get("DeepseaQt-Auth")
	log.Printf("请求方法: %s\n", method)
	log.Printf("访问令牌: %s\n", token)
	var jsonBody map[string]interface{}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("读取请求体错误, %v\n", err)
		return
	}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.Printf("解析请求体错误, %v\n", err)
		return
	}
	log.Printf("请求体: %v\n", jsonBody)
	ip := jsonBody["ip"]
	domain := jsonBody["domain"]
	remdom := jsonBody["remdom"]
	var dnsEntries []string
	var domains []string
	switch t := domain.(type) {
	case []interface{}:
		for _, dm := range t {
			log.Printf("待更新: %s\n", dm)
			domains = append(domains, fmt.Sprint(dm))
			dnsEntries = append(dnsEntries, fmt.Sprintf("%s %s", ip, dm))
		}
	}

	switch t := remdom.(type) {
	case []interface{}:
		for _, dm := range t {
			log.Printf("待移除: %s\n", dm)
			domains = append(domains, fmt.Sprint(dm))
		}
	}

	outputFile, outputError := os.OpenFile("report.hosts", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if outputError != nil {
		log.Println("打开或创建文件失败")
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(outputFile)

	inputFile, inputError := os.OpenFile("coredns.hosts", os.O_CREATE, 0666)
	if inputError != nil {
		log.Println("无法打开文件")
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(inputFile)

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	outputWriter := bufio.NewWriter(outputFile)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			continue
		}
		if !isContains(line, domains) {
			_, _ = outputWriter.WriteString(fmt.Sprintln(line))
		}
	}

	log.Printf("接收的域名数量: %d\n", len(dnsEntries))

	for _, entry := range dnsEntries {
		_, _ = outputWriter.WriteString(fmt.Sprintln(entry))
	}

	err = outputWriter.Flush()
	if err != nil {
		return
	}

	_, _ = copyFile("report.hosts", "coredns.hosts")

}

func isContains(line string, domains []string) (rs bool) {
	rs = false
	for _, domain := range domains {
		rs = strings.Contains(line, domain)
		if rs {
			break
		}
	}
	return
}

func copyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	return io.Copy(dst, src)
}

func main() {
	http.HandleFunc("/report", report)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
