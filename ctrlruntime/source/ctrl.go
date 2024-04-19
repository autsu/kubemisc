package main

import (
	"context"
	"path/filepath"

	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

type Ctrl struct {
	fw *FileWatch
}

func (c *Ctrl) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("new modify event, file path: ", filepath.Clean(req.String()))
	return ctrl.Result{}, nil
}

func (c *Ctrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("watchFileChange").
		//Watches(c.fw, &handler.EnqueueRequestForObject{}).
		// 新版本用下面这个函数
		WatchesRawSource(c.fw, &handler.EnqueueRequestForObject{}).
		Complete(c)
}
