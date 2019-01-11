package main

import (
	"net/http"
	"log"
)

func main(){
	http.HandleFunc("/",echohandler)
	err:=http.ListenAndServe(":9001",nil)
	if err!=nil {
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

func echohandler(w http.ResponseWriter,r *http.Request){
	for _,h:=range headersToCopy{
		log.Println(h,"  :",r.Header.Get(h))
	}
	w.Write([]byte("hellow"))
}