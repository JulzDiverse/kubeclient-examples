package main

import (
	"fmt"
	"os"
	"time"

	v11 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

/*
  The informer package makes getting kubernetes events easy.
  It comes with a simple Factory for all kubernetes resources.
*/
func main() {
	//kubeconfig := os.Getenv("KUBECONFIG")
	kubeconfig := ""
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("POD_NAME:", os.Getenv("POD_NAME"))
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 20*time.Second, informers.WithNamespace("opi"), informers.WithTweakListOptions(listTweaker))
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//AddFunc: func(obj interface{}) {
		//mObj := obj.(v1.Object)

		//name := mObj.GetName()
		//log.Printf("\nNew Pod Added to Store: %s", name)
		//},
		UpdateFunc: func(obj interface{}, newObj interface{}) {
			fmt.Println("UPDATE")
			mObj := obj.(*v11.Pod)
			fmt.Printf("Pod IP updated:%s", mObj.Status.PodIP)
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)

			fmt.Printf("\nPod Removed from Store: %s", mObj.GetName())
			fmt.Printf("Labels: %s\n", mObj.GetLabels())
		},
	})

	informer.Run(stopper)
}

func listTweaker(listOptions *v1.ListOptions) {
	pod := os.Getenv("POD_NAME")
	listOptions.LabelSelector = fmt.Sprintf("statefulset.kubernetes.io/pod-name=%s", pod)
}
