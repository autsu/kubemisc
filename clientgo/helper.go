package clientgo

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type WaitFunc = func() (bool, error)

func WaitPodCreate(ctx context.Context, cli *kubernetes.Clientset, podName string) WaitFunc {
	return func() (bool, error) {
		pod, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		return false, nil
	}
}

func WaitPodUpdate(ctx context.Context, cli *kubernetes.Clientset, podName string) WaitFunc {
	return WaitPodCreate(ctx, cli, podName)
}

func WaitPodDelete(ctx context.Context, cli *kubernetes.Clientset, podName string) WaitFunc {
	return func() (bool, error) {
		_, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Get(ctx, podName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		return false, nil
	}
}
