/*
Copyright 2022.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AutoLabelReconciler reconciles a AutoLabel object
type AutoLabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=misc.io.io,resources=autolabels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=misc.io.io,resources=autolabels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=misc.io.io,resources=autolabels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutoLabel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AutoLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logz := log.FromContext(ctx)

	// TODO(user): your logic here
	// Reconcile 在第一次运行时会 list 全部的 deploy，同时执行 Reconcile logic，也就是会为当前 k8s system 中存在的
	// 的每一个 deploy 打上一个 XXX = DO_NOT_EDIT_IT 标签
	// 此外，当有 创建/更新/删除 事件发生时，也会触发 Reconcile，比如 edit 一个 deploy，将 label XXX = DO_NOT_EDIT_IT 改为
	// XXXX = DO_NOT_EDIT_IT，此时就会触发 Reconcile，再为其打一个 XXX = DO_NOT_EDIT_IT 标签
	var deploy appsv1.Deployment
	err := r.Get(ctx, req.NamespacedName, &deploy)
	if errors.IsNotFound(err) {
		logz.Error(nil, "Could not find deploy")
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not fetch deploy: %+v", err)
	}
	logz.Info("get a deploy", "name", deploy.Name)

	if deploy.Labels == nil {
		deploy.Labels = make(map[string]string)
	}
	deploy.Labels["XXX"] = "DO_NOT_EDIT_IT"

	if err := r.Client.Update(ctx, &deploy); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For 只能调用一次，也就是说只能对一种资源进行 reconcile
		//For(&corev1.Pod{}).
		//For(&corev1.ReplicationController{}).
		//For(&appsv1.ReplicaSet{}).
		Complete(r)
}
