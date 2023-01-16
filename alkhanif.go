package alkhanif

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

func IsFileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func GetFileExtension(filename string) string {
	splitted := strings.Split(filename, ".")
	return splitted[len(splitted)-1]
}

func Base64encode(source string) string {
	return b64.StdEncoding.EncodeToString([]byte(source))
}

func Base64decode(source string) string {
	dec, err := b64.StdEncoding.DecodeString(source)
	if err != nil {
		return ""
	}
	return string(dec)
}

func GrabAllString(source string, regex string) [][]string {
	r, err := regexp.Compile(regex)
	if err != nil {
		return [][]string{}
	}
	return r.FindAllStringSubmatch(source, -1)
}

func GrabString(source string, regex string) []string {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}
	}
	return r.FindStringSubmatch(source)
}

func GetFilePermission(filename string) string {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", fileInfo.Mode())
}

func GetFileSize(filename string) string {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return ""
	}
	return strconv.FormatInt(fileinfo.Size(), 10)
}

func Hasher(origin string, mode string) string {
	if mode == "md5" {
		return fmt.Sprintf("%v", md5.Sum([]byte(origin)))
	} else if mode == "sha1" {
		return fmt.Sprintf("%v", sha1.Sum([]byte(origin)))
	} else if mode == "sha256" {
		return fmt.Sprintf("%v", sha256.Sum256([]byte(origin)))
	} else if mode == "sha512" {
		return fmt.Sprintf("%v", sha512.Sum512([]byte(origin)))
	} else {
		return "Encryption does not supported"
	}
}

func VisitURL(link string) (string, error) {
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func ChatGPT(apiKey string, question string) (string, error) {
	url := "https://api.openai.com/v1/completions"

	requestBody := struct {
		Model     string `json:"model"`
		Prompt    string `json:"prompt"`
		MaxTokens int    `json:"max_tokens"`
		User      string `json:"user"`
	}{
		Model:     "text-davinci-003",
		Prompt:    question,
		MaxTokens: 3000,
		User:      "alkhanif",
	}

	jsonValue, _ := json.Marshal(requestBody)
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type Choice struct {
		Text     string  `json:"text"`
		Index    int     `json:"index"`
		Logprobs float64 `json:"logprobs"`
		Reason   string  `json:"finish_reason"`
	}

	type Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	type Response struct {
		ID      string   `json:"id"`
		Object  string   `json:"object"`
		Created int      `json:"created"`
		Model   string   `json:"model"`
		Choices []Choice `json:"choices"`
		Usage   Usage    `json:"usage"`
	}

	var response Response
	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response.Choices[0].Text), err
}
