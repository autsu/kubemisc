package helper

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type WaitFunc = func() (bool, error)

type PodWait struct {
	cli       *kubernetes.Clientset
	namespace string
}

func NewPodWait(cli *kubernetes.Clientset, namespace string) *PodWait {
	pw := &PodWait{}

	if cli == nil {
		panic("clientset is not set")
	}
	pw.cli = cli

	if namespace == "" {
		pw.namespace = metav1.NamespaceDefault
	} else {
		pw.namespace = namespace
	}

	return pw
}

func (p *PodWait) WaitPodCreate(ctx context.Context, podName string) WaitFunc {
	return func() (bool, error) {
		pod, err := p.cli.CoreV1().Pods(p.namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return pod.Status.Phase == corev1.PodRunning, nil
	}
}

func (p *PodWait) WaitPodUpdate(ctx context.Context, podName string) WaitFunc {
	return p.WaitPodCreate(ctx, podName)
}

func (p *PodWait) WaitPodDelete(ctx context.Context, podName string) WaitFunc {
	return func() (bool, error) {
		_, err := p.cli.CoreV1().Pods(p.namespace).Get(ctx, podName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		return false, nil
	}
}
