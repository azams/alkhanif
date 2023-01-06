package tools

import (
	b64 "encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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

func isFileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func getFileExtension(filename string) string {
	splitted := strings.Split(filename, ".")
	return splitted[len(splitted)-1]
}

func base64encode(source string) string {
	return b64.StdEncoding.EncodeToString([]byte(source))
}

func base64decode(source string) string {
	dec, err := b64.StdEncoding.DecodeString(source)
	if err != nil {
		return ""
	}
	return string(dec)
}

func grabString(source string, regex string) [][]string {
	r, err := regexp.Compile(regex)
	if err != nil {
		return [][]string{}
	}
	return r.FindAllStringSubmatch(source, -1)
}
