/*
/*
@Time : 2019/12/16 3:51 下午
@Author : tianpeng.du
@File : file-parse.go
@Software: GoLand
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	f, err := os.OpenFile(os.Args[1], os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	res := parseFile(string(content))

	_, err = f.WriteAt([]byte(res), 0)
	if err != nil {
		log.Println(err)
	}
}

func parseFile(content string) string {
	exp := regexp.MustCompile(`!\[.*\]\((.*)\)`)
	res := exp.FindAllStringSubmatch(content, -1)
	dir := time.Now().Format("20060102")
	for _, v := range res {
		log.Printf("process %s", v[1])
		uri, err := url.ParseRequestURI(v[1])
		if err != nil {
			panic(err)
		}
		filename := path.Base(uri.Path)
		downloadAndUpload(v[1], filename)
		newname := fmt.Sprintf("https://rocdu.top/%s/%s", dir, filename)
		content = strings.ReplaceAll(content, v[1], newname)
		log.Printf("process finish %s", newname)
	}
	return content
}
