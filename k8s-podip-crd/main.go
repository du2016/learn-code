package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func main() {
	kubeconfig := flag.String("kubeconfig", "./kubeconfig", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()
	pic := newPodipcontroller(*kubeconfig)
	var stopCh <-chan struct{}
	pic.Run(2, stopCh)
}

type Podipcontroller struct {
	kubeClient *kubernetes.Clientset
	crdclient  *Podipclient

	podStore cache.Store
	//cache.AppendFunc
	podController cache.Controller
	podtoip       *PodToIp
	podsQueue     workqueue.RateLimitingInterface
}

func (slm *Podipcontroller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	fmt.Println("Starting Controller")
	//slm.registerTPR()
	go slm.podController.Run(stopCh)
	//go slm.endpointController.Run(stopCh)
	// wait for the controller to List. This help avoid churns during start up.
	if !cache.WaitForCacheSync(stopCh, slm.podController.HasSynced) {
		return
	}
	for i := 0; i < workers; i++ {
		go wait.Until(slm.podWorker, time.Second, stopCh)
	}

	<-stopCh
	fmt.Printf("Shutting down Controller")
	slm.podsQueue.ShutDown()
}

func (slm *Podipcontroller) podWorker() {
	workFunc := func() bool {
		key, quit := slm.podsQueue.Get()
		log.Println(key)
		if quit {
			return true
		}
		defer slm.podsQueue.Done(key)
		slm.podStore.Resync()
		obj, exists, err := slm.podStore.GetByKey(key.(string))
		log.Printf("%#v", obj)
		if !exists {
			fmt.Printf("Pod has been deleted %v\n", key)
			return false
		}
		if err != nil {
			fmt.Printf("cannot get pod: %v\n", key)
			return false
		}
		pod := obj.(*v1.Pod)
		if pod.DeletionTimestamp != nil {
			log.Println(slm.crdclient.Delete(pod.Name, pod.Namespace))
			return false
		}
		log.Println(slm.crdclient.Create(&PodToIp{
			Metadata: metav1.ObjectMeta{
				Name:      pod.ObjectMeta.Name,
				Namespace: pod.Namespace,
			},
			PodName:    pod.ObjectMeta.Name,
			PodAddress: pod.Status.PodIP,
		}))
		return false
	}
	for {
		if quit := workFunc(); quit {
			fmt.Printf("pod worker shutting down")
			return
		}
	}
}
func newPodipcontroller(kubeconfig string) *Podipcontroller {
	slm := &Podipcontroller{
		kubeClient: getClientsetOrDie(kubeconfig),
		crdclient:  getCRDClientOrDie(kubeconfig),
		podsQueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pods"),
	}
	watchList := cache.NewListWatchFromClient(slm.kubeClient.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())
	slm.podStore, slm.podController = cache.NewInformer(
		watchList,
		&v1.Pod{},
		time.Second*30,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    slm.enqueuePod,
			UpdateFunc: slm.updatePod,
		},
	)
	return slm
}

func (slm *Podipcontroller) enqueuePod(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		fmt.Printf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	slm.podsQueue.Add(key)
}

func (slm *Podipcontroller) updatePod(oldObj, newObj interface{}) {
	oldPod := oldObj.(*v1.Pod)
	newPod := newObj.(*v1.Pod)

	if newPod.Status.PodIP == oldPod.Status.PodIP {
		return
	}
	slm.enqueuePod(newObj)
}

type Podipclient struct {
	rest *rest.RESTClient
}

type PodToIp struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`

	PodName    string `json:"podName"`
	PodAddress string `json:"podAddress"`
}

func getClientsetOrDie(kubeconfig string) *kubernetes.Clientset {
	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	//clientset.CoreV1().Pods("111").Get("",nil)
	return clientset
}

func (c *Podipclient) Create(body *PodToIp) (*PodToIp, error) {
	var ret PodToIp
	err := c.rest.Post().
		Resource("podtoips").
		Namespace(body.Metadata.Namespace).
		Body(body).
		Do().Into(&ret)
	return &ret, err
}

func (c *Podipclient) Update(body *PodToIp) (*PodToIp, error) {
	var ret PodToIp
	err := c.rest.Put().
		Resource("podtoips").
		Namespace(body.Metadata.Namespace).
		Name(body.Metadata.Name).
		Body(body).
		Do().Into(&ret)
	return &ret, err
}

func (c *Podipclient) Get(name string, namespace string) (*PodToIp, error) {
	var ret PodToIp
	err := c.rest.Get().
		Resource("podtoips").
		Namespace(namespace).
		Name(name).
		Do().Into(&ret)
	return &ret, err
}

func (c *Podipclient) Delete(name string, namespace string) (*PodToIp, error) {
	var ret PodToIp
	err := c.rest.Delete().
		Resource("podtoips").
		Namespace(namespace).
		Name(name).
		Do().Into(&ret)
	return &ret, err
}

func getCRDClientOrDie(kubeconfig string) *Podipclient {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	configureClient(config)
	rest, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	return &Podipclient{rest}
}

func configureClient(config *rest.Config) {
	groupversion := schema.GroupVersion{
		Group:   "example.com",
		Version: "v1",
	}

	config.GroupVersion = &groupversion
	config.APIPath = "/apis"
	// Currently TPR only supports JSON
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				groupversion,
				&PodToIp{},
				&PodToIpList{},
				&metav1.ListOptions{},
				&metav1.DeleteOptions{},
			)
			return nil
		})
	schemeBuilder.AddToScheme(scheme.Scheme)
}

type PodToIpList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata"`

	Items []PodToIp `json:"items"`
}

func (in *PodToIp) DeepCopy() *PodToIp {
	if in == nil {
		return nil
	}
	out := new(PodToIp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodToIp) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodToIp) DeepCopyInto(out *PodToIp) {
	*out = *in
	return
}

func (in *PodToIpList) DeepCopy() *PodToIpList {
	if in == nil {
		return nil
	}
	out := new(PodToIpList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodToIpList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodToIpList) DeepCopyInto(out *PodToIpList) {
	*out = *in
	return
}
