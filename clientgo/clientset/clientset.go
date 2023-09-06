package main

import (
	"context"
	"fmt"
	"time"

	"zh.cargo.io/goclient/helper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
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

	ctx := context.Background()
	pod := helper.NewPodSimple()
	ph := helper.NewPodHelperForClientSet(clientset, metav1.NamespaceDefault)
	pw := helper.NewPodWait(clientset, metav1.NamespaceDefault)

	ph.Create(ctx, pod, metav1.CreateOptions{})
	// 等待 pod 创建完成
	if err := wait.PollImmediate(time.Second, time.Second*30, pw.WaitPodCreate(ctx, pod.Name)); err != nil {
		panic(err)
	}

	podFromGet := ph.Get(ctx, pod.Name, metav1.GetOptions{})

	podFromGet.Spec.Containers[0].Image = "nginx:1.21.1"
	ph.Update(ctx, podFromGet, metav1.UpdateOptions{})
	if err := wait.PollImmediate(time.Second, time.Second*30, pw.WaitPodUpdate(ctx, pod.Name)); err != nil {
		panic(err)
	}

	ph.List(ctx, metav1.ListOptions{LabelSelector: "app=example"}, func(items []corev1.Pod) {
		fmt.Print("list pod: ")
		for _, item := range items {
			fmt.Println(item.Name, item.Spec.Containers[0].Image)
		}
	})

	ph.Delete(ctx, pod.Name, metav1.DeleteOptions{})
	if err := wait.PollImmediate(time.Second, time.Second*30, pw.WaitPodDelete(ctx, pod.Name)); err != nil {
		panic(err)
	}

	ph.List(ctx, metav1.ListOptions{LabelSelector: "app=example"}, func(items []corev1.Pod) {
		fmt.Print("after delete, the pod list: ")
		for _, item := range items {
			fmt.Println(item.Name, item.Spec.Containers[0].Image)
		}
	})
}
