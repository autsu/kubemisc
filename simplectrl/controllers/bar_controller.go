/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	voidiov1 "void.io/kubemisc/api/v1"
)

var finalizerKey = "void.io/foo"

// BarReconciler reconciles a Bar object
type BarReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=void.io.void.io,resources=bars,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=void.io.void.io,resources=bars/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=void.io.void.io,resources=bars/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bar object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *BarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logz := log.FromContext(ctx)
	logz.Info("bar reconcile in")

	bar := &voidiov1.Bar{}
	if err := r.Get(ctx, req.NamespacedName, bar); err != nil {
		logz.Error(err, "get Bar error")
		return ctrl.Result{}, err
	}

	// 资源未被删除
	if bar.DeletionTimestamp.IsZero() {
		// 如果存在同名的 Foo，则添加 finalizer
		foo := &voidiov1.Foo{}
		if err := r.Get(ctx, req.NamespacedName, foo); err != nil {
			return ctrl.Result{}, err
		}
		// 添加前需要检查一下是否已经存在，不存在的话才添加
		if !containsString(bar.Finalizers, finalizerKey) {
			bar.Finalizers = append(bar.Finalizers, finalizerKey)
		}
		if err := r.Update(ctx, bar); err != nil {
			return ctrl.Result{}, err
		}
	} else { // 资源被标记删除
		foo := &voidiov1.Foo{}
		err := r.Get(ctx, req.NamespacedName, foo)
		// 如果同名的 Foo 资源不存在，则移除 finalizer
		if errors.IsNotFound(err) {
			bar.Finalizers = removeString(bar.Finalizers, finalizerKey)
			if err := r.Update(ctx, bar); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		if err != nil {
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BarReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&voidiov1.Bar{}).
		Complete(r)
}

// 辅助函数用于检查并从字符串切片中删除字符串。
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
