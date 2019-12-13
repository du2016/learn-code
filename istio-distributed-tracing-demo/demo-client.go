package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var target string

func main() {
	target = os.Args[1]
	http.HandleFunc("/", gethandler)
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		log.Println(err)
	}
}

var headersToCopy = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
}

func gethandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", target, nil)
	for _, h := range headersToCopy {
		log.Println(h, "  :", r.Header.Get(h))
		val := r.Header.Get(h)
		if val != "" {
			req.Header.Set(h, val)
		}
	}
	log.Println(req.Header)
	if user_cookie := r.Header.Get("user"); user_cookie != "" {
		req.Header.Set("Cookie", "user="+user_cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.Write(nil)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.Write(nil)
		return
	}
	w.Write([]byte(result))
}
