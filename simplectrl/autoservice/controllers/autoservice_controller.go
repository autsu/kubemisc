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
	"math/rand"
	"net"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	LabelBindService       = "bind-service"
	LabelEnableAutoService = "enable-auto-service"
	LabelAutoServiceUUID   = "auto-service-uuid"
)

// AutoServiceReconciler reconciles a AutoService object
type AutoServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=misc.lubenwei.io,resources=autoservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=misc.lubenwei.io,resources=autoservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=misc.lubenwei.io,resources=autoservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutoService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AutoServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logz := log.FromContext(ctx)
	_ = logz

	// TODO(user): your logic here
	//logs.Info("reconcile")
	var deploy = new(appsv1.Deployment)
	err := r.Client.Get(ctx, req.NamespacedName, deploy)
	if errors.IsNotFound(err) {
		logz.Error(nil, "Could not find ReplicaSet")
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	if deploy.Labels == nil {
		return ctrl.Result{}, nil
	}

	// 检查 deploy 是否设置了 LabelEnableAutoService 标签，
	// 如果设置了还需要检查该标签的值是否为 true
	enable := deploy.Labels[LabelEnableAutoService]
	if enable == "" {
		return ctrl.Result{}, nil
	}
	enableBool, err := strconv.ParseBool(enable)
	if err != nil {
		return ctrl.Result{}, nil
	}
	// 如果用户修改了 label，且已经创建了对应的 auto service，将这个 svc 删除
	if !enableBool {
		var svc = new(corev1.Service)
		err = r.Client.Get(ctx, req.NamespacedName, svc)
		if err != nil {
			return ctrl.Result{}, err
		}
		if metav1.IsControlledBy(svc, deploy) {
			err = r.Client.Delete(ctx, &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: req.Name, Namespace: req.Namespace}})
			if err != nil && !errors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		}
	}

	newService(deploy)
	// 为这个 deploy 设置一个 uuid label，用于 svc 的唯一 label selector
	deploy.Labels["auto-service-uuid"] = uuid.New().String()
	//r.Client.Update()

	//var svc = new(corev1.Service)

	return ctrl.Result{}, nil
}

func (r *AutoServiceReconciler) checkSvcExist(ctx context.Context, name string) (exist bool, err error) {
	//r.Get(ctx)
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&appsv1.Deployment{}).
		Complete(r)
}

func newService(deploy *appsv1.Deployment) *corev1.Service {
	labels := map[string]string{
		"auto-create-by": deploy.Name,
	}
	svcPorts := make([]corev1.ServicePort, 0)
	for _, container := range deploy.Spec.Template.Spec.Containers {
		for _, port := range container.Ports {
			svcPorts = append(svcPorts, corev1.ServicePort{
				Name:       port.Name,
				Protocol:   port.Protocol,
				Port:       GetRandomPort(),
				TargetPort: intstr.FromInt(int(port.ContainerPort)),
			})
		}
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploy.Name,
			Namespace: deploy.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(deploy, appsv1.SchemeGroupVersion.WithKind("Deployment")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: deploy.Spec.Selector.MatchLabels,
			Ports:    svcPorts,
		},
	}
	return svc
}

func GetRandomPort() int32 {
	rand.Seed(time.Now().Unix())
	for {
		port := rand.Intn(10000) + 40000
		_, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
		if err != nil {
			continue
		} else {
			return int32(port)
		}
	}
}
