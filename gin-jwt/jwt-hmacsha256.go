package main

import (
	"crypto/hmac"
	"crypto"
	"log"
	_ "crypto/sha256"
	"strings"
	"encoding/base64"
)

func main(){
	header:=`{"alg":"HS256","typ":"JWT"}`
	claims:=`{"sub":"1234567890","name":"John Doe","admin":true}`
	encodeheader:=base64.StdEncoding.EncodeToString([]byte(header))
	encodeclaims:= base64.StdEncoding.EncodeToString([]byte(claims))
	log.Println(encodeheader+"."+encodeclaims)
	hasher :=hmac.New(crypto.SHA256.New,[]byte("secret"))
	hasher.Write([]byte(encodeheader+"."+encodeclaims))
	log.Println(strings.TrimRight(base64.URLEncoding.EncodeToString(hasher.Sum(nil)), "="))
}