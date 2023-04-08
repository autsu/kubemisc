package main

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Ctrl struct {
	fw *FileWatch
	ch chan event.GenericEvent
}

func (c *Ctrl) Sync() {
	ticket := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticket.C:
			c.ch <- event.GenericEvent{Object: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "qwe",
					Namespace: metav1.NamespaceDefault,
				},
			}}
		}
	}
}

func (c *Ctrl) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("req: ", req.String())
	return ctrl.Result{}, nil
}

func (c *Ctrl) SetupWithManager(mgr ctrl.Manager) error {
	c.fw = NewFileWatch("/var/tmp/test.txt")
	c.ch = make(chan event.GenericEvent)
	go c.Sync()
	//fw.Sync()
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Watches(c.fw, &handler.EnqueueRequestForObject{}).
		Watches(&source.Channel{Source: c.ch}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(&predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return e.Object.GetNamespace() == corev1.NamespaceDefault
			},
			DeleteFunc:  nil,
			UpdateFunc:  nil,
			GenericFunc: nil,
		}).
		Complete(c)
}
