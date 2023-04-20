package main

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// 错误示例
func __main() {
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

	// 创建一个 Pod 的 ListerWatcher
	// 坑爹的 chatGPT 生成的代码，最后一个参数传了 nil 导致 panic
	//podLW := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", metav1.NamespaceDefault, nil)

	podLW := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", metav1.NamespaceDefault, fields.Everything())
	list, err := podLW.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	podList, ok := list.(*corev1.PodList)
	if !ok {
		panic(err)
	}
	for _, item := range podList.Items {
		fmt.Println(item.Name)
	}

	//fmt.Printf("%+v\n", podLW)

	// 创建一个 Pod 的 Store
	podStore, podController := cache.NewInformer(podLW, &corev1.Pod{}, time.Second*30, cache.ResourceEventHandlerFuncs{
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
	//fmt.Printf("podStorage: %+v\n, podController: %+v\n", podStore, podController)

	// 创建一个 Reflector，用于同步 k8s 中的 Pod 资源到本地 Pod 的 Store 中
	reflector := cache.NewReflector(podLW, &corev1.Pod{}, podStore, time.Second*30)
	//fmt.Printf("reflector: %+v\n", reflector)

	// 启动 Reflector
	stopCh := make(chan struct{})
	defer close(stopCh)
	go reflector.Run(stopCh)

	// 等待 Reflector 同步 k8s 中的 Pod 资源到本地 Pod 的 Store 中
	if !cache.WaitForCacheSync(stopCh, podController.HasSynced) {
		panic("Timed out waiting for caches to sync")
	}

	// 打印本地 Pod 的数量
	fmt.Printf("Pods in the Store: %d\n", len(podStore.List()))

	// 等待 10s，然后退出程序
	wait.PollImmediateUntil(time.Second*1, func() (bool, error) { return false, nil }, wait.NeverStop)
}
