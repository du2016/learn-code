package main

import (
	"bytes"
	"encoding/json"
	"io"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api/v1"
	"log"
	"net/http"
)

var (
	kubeconfig string = "xxx"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("hellowrold"))
	})
	http.HandleFunc("/predicates/test", testPredicateHandler)
	http.HandleFunc("/prioritize/test", testPrioritizeHandler)
	http.HandleFunc("/bind/test", BindHandler)
	http.ListenAndServe(":8880", nil)
}

func testPredicateHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	log.Println(buf.String())

	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult

	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Nodes:       nil,
			FailedNodes: nil,
			Error:       err.Error(),
		}
	} else {
		extenderFilterResult = predicateFunc(extenderArgs)
	}
	if resultBody, err := json.Marshal(extenderFilterResult); err != nil {
		panic(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultBody)
	}

}

func testPrioritizeHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		panic(err)
	}
	if list, err := prioritizeFunc(extenderArgs); err != nil {
		panic(err)
	} else {
		hostPriorityList = list
	}
	if resultBody, err := json.Marshal(hostPriorityList); err != nil {
		panic(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultBody)
	}
}

func predicateFunc(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	pod := args.Pod
	canSchedule := make([]v1.Node, 0, len(args.Nodes.Items))
	canNotSchedule := make(map[string]string)
	for _, node := range args.Nodes.Items {
		result, err := func(pod v1.Pod, node v1.Node) (bool, error) {
			return true, nil
		}(pod, node)
		if err != nil {
			canNotSchedule[node.Name] = err.Error()
		} else {
			if result {
				canSchedule = append(canSchedule, node)
			}
		}
	}
	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: canSchedule,
		},
		FailedNodes: canNotSchedule,
		Error:       "",
	}
	return &result
}

func prioritizeFunc(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	nodes := args.Nodes.Items
	var priorityList schedulerapi.HostPriorityList
	priorityList = make([]schedulerapi.HostPriority, len(nodes))
	for i, node := range nodes {
		priorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: 0,
		}
	}
	return &priorityList, nil
}

func BindHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderBindingArgs schedulerapi.ExtenderBindingArgs
	if err := json.NewDecoder(body).Decode(&extenderBindingArgs); err != nil {
		panic(err)
	}
	b := &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{Namespace: extenderBindingArgs.PodNamespace, Name: extenderBindingArgs.PodName, UID: extenderBindingArgs.PodUID},
		Target: v1.ObjectReference{
			Kind: "Node",
			Name: extenderBindingArgs.Node,
		},
	}
	bind(b)

}

func bind(b *v1.Binding) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset.CoreV1().Pods(b.Namespace).Bind(b)
}
