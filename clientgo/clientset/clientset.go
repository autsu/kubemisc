package main

import (
	"context"
	"fmt"
	"time"

	clientgo "zh.cargo.io/goclient"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func createPod(ctx context.Context, cli *kubernetes.Clientset, po *corev1.Pod) {
	pod, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Create(ctx, po, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return
		}
		panic(err)
	}
	fmt.Printf("create pod %v\n", pod.Name)
}

func getPod(ctx context.Context, cli *kubernetes.Clientset, name string) *corev1.Pod {
	pod, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("get pod: %v %v\n", pod.Name, pod.Spec.Containers[0].Image)
	return pod
}

func updatePod(ctx context.Context, cli *kubernetes.Clientset, newPod *corev1.Pod) {
	_, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Update(ctx, newPod, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("update pod")
}

func listPods(ctx context.Context, cli *kubernetes.Clientset) {
	podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, metav1.ListOptions{LabelSelector: "app=example"})
	if err != nil {
		panic(err)
	}
	fmt.Print("list pod: ")
	for _, item := range podList.Items {
		fmt.Println(item.Name, item.Spec.Containers[0].Image)
	}
}

func deletePod(ctx context.Context, cli *kubernetes.Clientset, name string) {
	err := cli.CoreV1().Pods(metav1.NamespaceDefault).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("delete pod")
}

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
	pod := clientgo.NewPod()
	createPod(ctx, clientset, pod)

	if err := wait.PollImmediate(time.Second, time.Second*30, clientgo.WaitPodCreate(ctx, clientset, pod.Name)); err != nil {
		panic(err)
	}

	podFromGet := getPod(ctx, clientset, pod.Name)
	podFromGet.Spec.Containers[0].Image = "nginx:1.21.1"
	updatePod(ctx, clientset, podFromGet)
	if err := wait.PollImmediate(time.Second, time.Second*30, clientgo.WaitPodUpdate(ctx, clientset, pod.Name)); err != nil {
		panic(err)
	}

	listPods(ctx, clientset)
	deletePod(ctx, clientset, pod.Name)
	if err := wait.PollImmediate(time.Second, time.Second*30, clientgo.WaitPodDelete(ctx, clientset, pod.Name)); err != nil {
		panic(err)
	}

	listPods(ctx, clientset)
}
