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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	voidiov1 "void.io/kubemisc/api/v1"
)

// FooReconciler reconciles a Foo object
type FooReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=void.io.void.io,resources=foos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=void.io.void.io,resources=foos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=void.io.void.io,resources=foos/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Foo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *FooReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var foo voidiov1.Foo
	if err := r.Client.Get(ctx, req.NamespacedName, &foo); err != nil {
		r.Recorder.Event(&foo, corev1.EventTypeWarning, "GetError", err.Error())
		return ctrl.Result{}, err
	}

	var deploy = &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: foo.Spec.DeploymentName}, deploy)
	if errors.IsNotFound(err) {
		r.Recorder.Eventf(&foo, corev1.EventTypeNormal, "CreateDeploy",
			"foo.spec.deploymentName %v is not found, will create it\n",
			foo.Spec.DeploymentName)
		deploy = newDeployment(&foo)
		err = r.Client.Create(ctx, deploy)
	}

	if err != nil {
		r.Recorder.Event(&foo, corev1.EventTypeWarning, "GetError", err.Error())
		return ctrl.Result{}, err
	}

	if !metav1.IsControlledBy(deploy, &foo) {
		msgFmt := "Deployment %v is not controller by Foo %v\n"
		args := []any{deploy.Name, foo.Name}
		r.Recorder.Eventf(&foo, corev1.EventTypeWarning, "Forbidden", msgFmt, args)
		return ctrl.Result{}, fmt.Errorf(msgFmt, args...)
	}

	if foo.Spec.Replicas != nil &&
		foo.Spec.Replicas != deploy.Spec.Replicas {
		r.Recorder.Eventf(&foo, corev1.EventTypeNormal,
			"Reconcile",
			"ready reconcile foo replicas (cur %v, want %v)\n",
			foo.Spec.Replicas, deploy.Spec.Replicas,
		)
		if err := r.Client.Update(ctx, newDeployment(&foo)); err != nil {
			r.Recorder.Event(&foo, corev1.EventTypeWarning, "UpdateError", err.Error())
			return ctrl.Result{}, err
		}
		r.Recorder.Eventf(&foo, corev1.EventTypeNormal,
			"ReconcileSuccess", "foo cur replicas: %v", foo.Spec.Replicas)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FooReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&voidiov1.Foo{}).
		Complete(r)
}

func newDeployment(foo *voidiov1.Foo) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": foo.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      foo.Spec.DeploymentName,
			Namespace: foo.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(foo, voidiov1.GroupVersion.WithKind("Foo")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: foo.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}
