package main

import (
	"errors"
	"fmt"
	"time"

	"void.io/kubemisc/clientgo/helper/print"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	cfg.APIPath = "/api"
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := factory.Core().V1().Pods()
	_, err = podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("New Pod Added to Store: %s/%s\n", pod.Namespace, pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			oldPod := oldObj.(*corev1.Pod)
			if newPod.ResourceVersion == oldPod.ResourceVersion {
				return
			}
			fmt.Printf("Pod Updated in Store: %s/%s\n", newPod.Namespace, newPod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod Deleted from Store: %s/%s\n", pod.Namespace, pod.Name)
		},
	})
	if err != nil {
		panic(err)
	}

	stopCh := wait.NeverStop
	factory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced) {
		runtime.HandleError(errors.New("failed to sync"))
		return
		//panic("failed to sync")
	}

	// Lister 会直接从本地 cache 中获取，而不是请求 apiserver
	podLister := podInformer.Lister()
	pods, err := podLister.List(labels.Everything())
	if err != nil {
		panic(err)
	}
	print.ResourceListItemName(pods)

	pods, err = podLister.Pods(metav1.NamespaceDefault).List(labels.Everything())
	if err != nil {
		panic(err)
	}
	print.ResourceListItemName(pods)
}
