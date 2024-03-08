package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"void.io/kubemisc/clientgo/helper/pods"
	resourcehelper "void.io/kubemisc/clientgo/helper/resource"
	waithelper "void.io/kubemisc/clientgo/helper/wait"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	//cfg.APIPath = "/api"
	//cfg.GroupVersion = &corev1.SchemeGroupVersion
	//cfg.NegotiatedSerializer = scheme.Codecs

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pod := resourcehelper.NewPodSample()
	ph := pods.NewPodHelperForClientSet(clientset, metav1.NamespaceDefault)
	pw := waithelper.NewPodWait(clientset, time.Second)

	ph.MustCreate(ctx, pod, metav1.CreateOptions{})
	// 等待 pod 创建完成
	if err := pw.WaitCreate(ctx, pod.Namespace, pod.Name); err != nil {
		panic(err)
	}
	slog.Info("pod create success", slog.String("name", pod.Name))

	podFromGet := ph.MustGet(ctx, pod.Name, metav1.GetOptions{})
	slog.Info("get pod success", slog.String("pod name", pod.Name), slog.String("pod image", pod.Spec.Containers[0].Image))

	podFromGet.Spec.Containers[0].Image = "nginx:1.21.1"
	ph.MustUpdate(ctx, podFromGet, metav1.UpdateOptions{})
	if err := pw.WaitUpdate(ctx, pod.Namespace, pod.Name); err != nil {
		panic(err)
	}
	slog.Info("pod update success", slog.String("name", pod.Name))

	ph.MustList(ctx, metav1.ListOptions{LabelSelector: "app=example"}, func(items []corev1.Pod) {
		fmt.Print("list pod: ")
		for _, item := range items {
			fmt.Println(item.Name, item.Spec.Containers[0].Image)
		}
	})

	ph.MustDelete(ctx, pod.Name, metav1.DeleteOptions{})
	if err := pw.WaitDelete(ctx, pod.Namespace, pod.Name); err != nil {
		panic(err)
	}
	slog.Info("pod update success", slog.String("name", pod.Name))

	ph.MustList(ctx, metav1.ListOptions{LabelSelector: "app=example"}, func(items []corev1.Pod) {
		fmt.Print("after delete, the pod list: ")
		for _, item := range items {
			fmt.Println(item.Name, item.Spec.Containers[0].Image)
		}
	})
}
