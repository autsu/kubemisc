package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Ctrl struct {
	fw *FileWatch
}

func (c *Ctrl) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("new modify event, file path: ", req.String())
	return ctrl.Result{}, nil
}

func (c *Ctrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// 必须要指定一种监听的资源，这里设置为 Pod
		For(&corev1.Pod{}).
		Watches(c.fw, &handler.EnqueueRequestForObject{}).
		WithEventFilter(&predicate.Funcs{
			// 在这里设置过滤事件，只监听指定 namespace 的 pod 的创建事件，
			// 这里的 namespace 我随便设置了一个，进而可以达到忽略 Pod 事件的目的
			// 因为我们的主要目的是观察自定义 Source 的效果，所以尽量避免其他资源的干扰
			CreateFunc: func(e event.CreateEvent) bool {
				return e.Object.GetNamespace() == "UNKNOWN_NAMESPACE"
			},
		}).
		Complete(c)
}
