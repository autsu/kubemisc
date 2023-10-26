package wait

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

//type WaitFunc = func() (bool, error)

type Waiter interface {
	WaitCreate(ctx context.Context, namespace, name string) error
	WaitUpdate(ctx context.Context, namespace, name string) error
	WaitDelete(ctx context.Context, namespace, name string) error
}

var _ Waiter = &PodWait{}

// TODO: 做成一个更通用的 waiter，而不是具体某个资源

type PodWait struct {
	cli      *kubernetes.Clientset
	interval time.Duration
}

func (p *PodWait) WaitCreate(ctx context.Context, namespace, name string) error {
	return wait.PollUntilContextCancel(ctx,
		p.interval, true,
		func(ctx context.Context) (done bool, err error) {
			pod, err := p.cli.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			return pod.Status.Phase == corev1.PodRunning, nil
		})
}

func (p *PodWait) WaitUpdate(ctx context.Context, namespace, name string) error {
	return wait.PollUntilContextCancel(ctx, p.interval, true, func(ctx context.Context) (done bool, err error) {
		pod, err := p.cli.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return pod.Status.Phase == corev1.PodRunning, nil
	})
}

func (p *PodWait) WaitDelete(ctx context.Context, namespace, name string) error {
	return wait.PollUntilContextCancel(ctx, p.interval, true, func(ctx context.Context) (done bool, err error) {
		_, err = p.cli.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		return false, nil
	})
}

func NewPodWait(cli *kubernetes.Clientset, interval time.Duration) *PodWait {
	pw := &PodWait{}

	if cli == nil {
		panic("clientset is not set")
	}
	pw.cli = cli

	//if namespace == "" {
	//	pw.namespace = metav1.NamespaceDefault
	//} else {
	//	pw.namespace = namespace
	//}

	return pw
}

var _ Waiter = &WaitFunc{}

type WaitFunc struct {
	cli      *kubernetes.Clientset
	interval time.Duration
}

func (w *WaitFunc) WaitCreate(ctx context.Context, namespace, name string) error {
	//TODO implement me
	panic("implement me")
}

func (w *WaitFunc) WaitUpdate(ctx context.Context, namespace, name string) error {
	//TODO implement me
	panic("implement me")
}

func (w *WaitFunc) WaitDelete(ctx context.Context, namespace, name string) error {
	//TODO implement me
	panic("implement me")
}

//func (p *PodWait) WaitPodCreate(ctx context.Context, podName string) WaitFunc {
//	return func() (bool, error) {
//		pod, err := p.cli.CoreV1().Pods(p.namespace).Get(ctx, podName, metav1.GetOptions{})
//		if err != nil {
//			return false, err
//		}
//		return pod.Status.Phase == corev1.PodRunning, nil
//	}
//}
//
//func (p *PodWait) WaitPodUpdate(ctx context.Context, podName string) WaitFunc {
//	return p.WaitPodCreate(ctx, podName)
//}
//
//func (p *PodWait) WaitPodDelete(ctx context.Context, podName string) WaitFunc {
//	return func() (bool, error) {
//		_, err := p.cli.CoreV1().Pods(p.namespace).Get(ctx, podName, metav1.GetOptions{})
//		if errors.IsNotFound(err) {
//			return true, nil
//		}
//		if err != nil {
//			return false, err
//		}
//		return false, nil
//	}
//}
