/*
/*
@Time : 2019/12/13 4:58 下午
@Author : tianpeng.du
@File : upload.go
@Software: GoLand
*/
package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
)

var (
	accessKey = os.Getenv("QINIU_ACCESS_KEY")
	secretKey = os.Getenv("QINIU_SECRET_KEY")
	bucket    = os.Getenv("QINIU_WECHAT_BUCKET")
	proxyURL  = os.Getenv("PROXY_URL")
	downUrl   string
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	if len(os.Args) < 1 {
		panic("no parm")
	}
	urlname := os.Args[1]

	uri, err := url.ParseRequestURI(urlname)
	if err != nil {
		panic(err)
	}
	filename := path.Base(uri.Path)

	dclient := http.Client{}

	if proxyURL != "" {
		proxyURI, _ := url.Parse(proxyURL)
		dclient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURI),
		}
	}

	req, err := http.NewRequest("GET", urlname, nil)
	if err != nil {
		panic(err)
	}

	resp, err := dclient.Do(req)
	if err != nil {
		log.Print(err)
	}

	tmpfile, err := ioutil.TempFile("", "tmpfile")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(tmpfile.Name(), os.O_RDWR|os.O_CREATE, 0644)
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Println(err)
	}
	client := http.Client{}
	upload(tmpfile.Name(), filename, client)
	os.Remove(tmpfile.Name())
}

func upload(file string, filename string, client http.Client) {
	dir := time.Now().Format("20060102")

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(accessKey, secretKey)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	resumeUploader := storage.NewResumeUploaderEx(&cfg, &storage.Client{Client: &client})
	upToken := putPolicy.UploadToken(mac)

	ret := storage.PutRet{}

	err := resumeUploader.PutFile(context.Background(), &ret, upToken, fmt.Sprintf("%s/%s", dir, filename), file, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("upload %s success \n", ret.Key)
}
