package tools

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ReadFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func CompileRequest(data string, ssl bool) (targetUrl string, method string, body string, req http.Request, httpversion string) {
	scheme := "http"
	datas := strings.Split(strings.Replace(data, "\r", "", -1), "\n\n")
	headers := datas[0]
	body = datas[1]
	lines := strings.Split(headers, "\n")
	line := strings.Split(lines[0], " ")
	method = line[0]
	path := line[1]
	httpversion = line[2]
	targetDomain := ""
	for i := 1; i < len(lines); i++ {
		head := strings.Split(lines[i], ": ")
		if head[0] == "Host" {
			targetDomain = head[0]
		}
		req.Header.Add(head[0], head[1])
	}
	if ssl {
		scheme = "https"
	}
	targetUrl = scheme + "://" + targetDomain + path
	return targetUrl, method, body, req, httpversion
}
