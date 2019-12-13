/*
/*
@Time : 2019/11/27 7:00 下午
@Author : tianpeng.du
@File : meta.go
@Software: GoLand
*/
package main

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/metadata/metadatainformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "./kubeconfig")
	if err != nil {
		log.Println(err)
		return
	}
	client, err := metadata.NewForConfig(config)
	//watcher, err := client.Resource(schema.GroupVersionResource{
	//	Version:  "v1",
	//	Resource: "pods",
	//}).Watch(v1.ListOptions{})
	//if err != nil {
	//	log.Println(err)
	//}
	//for v := range watcher.ResultChan() {
	//	log.Println(v)
	//}

	informer := metadatainformer.NewFilteredMetadataInformer(client, schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}, v1.NamespaceAll, 0, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, nil)
	var stop chan struct{}
	go informer.Informer().Run(stop)
	ret, err := informer.Lister().ByNamespace("production").List(nil)
	if err != nil {
		log.Println(err)
	}
	for _, v := range ret {
		log.Println(v)
	}
}
