package main

import (
	"log"
	"os"

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
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 5000, informers.WithNamespace("eirini"))
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)

			name := mObj.GetName()
			log.Printf("\nNew Pod Added to Store: %s", name)
		},
		UpdateFunc: func(obj interface{}, newObj interface{}) {
			mObj := obj.(*v11.Pod)
			log.Printf("Pod IP:%s", mObj.Status.PodIP)
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)

			log.Printf("\nPod Removed to Store: %s", mObj.GetName())
			log.Printf("Labels: %s\n", mObj.GetLabels())
		},
	})

	informer.Run(stopper)
}
