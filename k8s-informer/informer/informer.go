package informer

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	//"k8s.io/client-go/pkg/api"
	"time"
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
	watchList := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())

	store, controller = cache.NewInformer(
		watchList,
		&v1.Pod{},
		time.Second*30,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handlepodsAdd,
			UpdateFunc: handlerpodsUpdate,
			DeleteFunc: handlerpodsDelete,
		},
	)
	stop := make(chan struct{})
	log.Println("begin watch")
	controller.Run(stop)
}

func handlepodsAdd(obj interface{}) {
	log.Println(obj.(*v1.Pod).Name, obj.(*v1.Pod).Namespace)
}
func handlerpodsUpdate(oldObj, newObj interface{}) {
	log.Println(oldObj.(*v1.Pod))
	log.Println(newObj.(*v1.Pod))
	log.Println()
}

func handlerpodsDelete(obj interface{}) {
	log.Println(obj.(*v1.Pod).Name, obj.(*v1.Pod).Namespace)
}
