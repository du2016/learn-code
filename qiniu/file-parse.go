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
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func main() {
	//chrome := svg2png.NewChrome().SetHeight(600).SetWith(600)
	//filepath := "Soccerball_mask_transparent_background.png"
	//if err := chrome.Screenshoot("https://upload.wikimedia.org/wikipedia/commons/8/84/Example.svg", filepath); err != nil {
	//	panic(err)
	//}
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
		var issvg bool = false
		uri, err := url.ParseRequestURI(v[1])
		if err != nil {
			panic(err)
		}
		filename := path.Base(uri.Path)
		if strings.HasSuffix(filename, "svg") {
			filename = strings.Replace(filename, "svg", "png", 1)
			issvg = true
		}
		log.Println(filename)
		if strings.HasPrefix(v[1], "file:///") {
			client := http.Client{}
			log.Println(strings.Replace(v[1], "file://", "", 1))
			upload(strings.Replace(v[1], "file://", "", 1), filename, client)
		} else {
			downloadAndUpload(v[1], filename, issvg)
		}
		newname := fmt.Sprintf("%s/%s/%s", os.Getenv("QINIU_WECHAT_DOMAIN"), dir, filename)
		content = strings.ReplaceAll(content, v[1], newname)
		log.Printf("process finish %s", newname)
	}
	return content
}
