package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
	"time"
	"k8s.io/client-go/informers"
)
var controller cache.Controller
var store cache.Store

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./kubeconfig")
	if err != nil {
		log.Println(err)
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return
	}
	shardinformerfactory:=informers.NewSharedInformerFactory(clientset,time.Minute)
	shardinformerfactory.InformerFor()
	stop:=make(chan struct{})
	shardinformerfactory.Start(stop)
	podlister:=shardinformerfactory.Core().V1().Pods().Lister()
	log.Println(podlister.Pods("default").Get("busybox-58b5cb5bf5-xcj8z"))
	log.Println("begin watch")

}

func handlepodsAdd(obj interface{}){
	log.Println(obj.(*v1.Pod).Name)

}
func handlerpodsupdate(oldObj, newObj interface{}){
	log.Println(oldObj.(*v1.Pod),newObj.(*v1.Pod))
}