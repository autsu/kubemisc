package main

import (
	"context"

	"void.io/kubemisc/clientgo/helper/maps"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func findServiceDeploys(ctx context.Context, cli *kubernetes.Clientset, svc *corev1.Service) ([]*appsv1.Deployment, error) {
	if cli == nil || svc == nil {
		return nil, nil
	}
	var ret []*appsv1.Deployment
	deployList, err := cli.AppsV1().Deployments(svc.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range deployList.Items {
		deploy := deployList.Items[i]
		if !maps.Contains(deploy.Spec.Selector.MatchLabels, svc.Spec.Selector) {
			continue
		}
		ret = append(ret, &deploy)
	}
	return ret, nil
}
