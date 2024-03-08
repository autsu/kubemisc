package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
}

// 只做 ListAndWatch 操作，不坐本地缓存，通过 Reflector + 队列实现
func listWatchDeploy(ctx context.Context) {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	lw := cache.NewListWatchFromClient(
		cli.AppsV1().RESTClient(), "deployments", metav1.NamespaceAll, fields.Everything())
	store := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction:           cache.MetaNamespaceKeyFunc,
		KnownObjects:          nil,
		EmitDeltaTypeReplaced: false,
		Transformer:           nil,
	})
	reflector := cache.NewReflector(lw, &appsv1.Deployment{}, store, 0)
	go reflector.Run(ctx.Done())

	for {
		if _, err := store.Pop(func(obj interface{}, isInInitialList bool) error {
			switch v := obj.(type) {
			case cache.Deltas:
				vv := v.Newest()
				dep := vv.Object.(*appsv1.Deployment)
				switch vv.Type {
				case cache.Sync, cache.Added:
					slog.Debug("", "event type", cache.Sync+", "+cache.Added, "namespace/name", fmt.Sprintf("%v/%v", dep.Namespace, dep.Name))
					if !changeRollingUpdateStrategy(dep) {
						return nil
					}
					_, err := cli.AppsV1().Deployments(dep.Namespace).Update(ctx, dep, metav1.UpdateOptions{})
					return err
				case cache.Updated:
					slog.Debug("", "event type", cache.Updated, "namespace/name", fmt.Sprintf("%v/%v", dep.Namespace, dep.Name))
					b1 := changeRollingUpdateStrategy(dep)
					b2 := addTolerationsAndLabel(dep)
					if !b1 && !b2 {
						return nil
					}
					_, err := cli.AppsV1().Deployments(dep.Namespace).Update(ctx, dep, metav1.UpdateOptions{})
					if errors.IsConflict(err) {
						slog.Debug("conflict")
					}
					return err
				}
			}
			return nil
		}); err != nil {
			slog.Error("handle queue obj error", "error", err)
			continue
		}
	}
}

func checkDeploymentNs(dep *appsv1.Deployment) bool {
	if dep == nil {
		return false
	}
	return dep.Namespace == metav1.NamespaceDefault
}

func changeRollingUpdateStrategy(dep *appsv1.Deployment) (needUpdate bool) {
	if dep == nil || dep.Spec.Replicas == nil || dep.Spec.Strategy.RollingUpdate == nil {
		return false
	}
	if !checkDeploymentNs(dep) {
		return false
	}
	oldMaxSurge := dep.Spec.Strategy.RollingUpdate.MaxSurge
	oldMaxUnavailable := dep.Spec.Strategy.RollingUpdate.MaxUnavailable

	newVal := intstr.FromString("25%")
	if *dep.Spec.Replicas >= 20 {
		newVal = intstr.FromInt32(10)
	}
	if *oldMaxSurge == newVal && *oldMaxUnavailable == newVal {
		return false
	}
	dep.Spec.Strategy.RollingUpdate.MaxSurge = &newVal
	dep.Spec.Strategy.RollingUpdate.MaxUnavailable = &newVal
	return true
}

func addTolerationsAndLabel(dep *appsv1.Deployment) (needUpdate bool) {
	if dep == nil {
		return false
	}
	if !checkDeploymentNs(dep) {
		return false
	}
	t1 := corev1.Toleration{
		Key:               "node.kubernetes.io/unreachable",
		Operator:          "Exists",
		Effect:            "NoExecute",
		TolerationSeconds: ptr.To(int64(60)),
	}
	t2 := corev1.Toleration{
		Key:               "node.kubernetes.io/not-ready",
		Operator:          "Exists",
		Effect:            "NoExecute",
		TolerationSeconds: ptr.To(int64(60)),
	}
	// 防止重复添加
	if !slices.ContainsFunc(dep.Spec.Template.Spec.Tolerations, func(toleration corev1.Toleration) bool {
		return toleration.MatchToleration(&t1)
	}) {
		dep.Spec.Template.Spec.Tolerations = append(dep.Spec.Template.Spec.Tolerations, t1)
		needUpdate = true
	}
	if !slices.ContainsFunc(dep.Spec.Template.Spec.Tolerations, func(toleration corev1.Toleration) bool {
		return toleration.MatchToleration(&t2)
	}) {
		dep.Spec.Template.Spec.Tolerations = append(dep.Spec.Template.Spec.Tolerations, t2)
		needUpdate = true
	}
	return false
}

func main() {
	listWatchDeploy(context.TODO())
}
