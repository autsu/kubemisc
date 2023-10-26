package main

import (
	"context"
	"testing"

	"void.io/kubemisc/clientgo/helper"
	"void.io/kubemisc/clientgo/helper/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func initResourceForTestFindServiceDeploys() {
	ns := "test-1022-1"
	ctx := context.TODO()

	cli := helper.NewClientSetOrDie()
	_, err := cli.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
	if err != nil && !kubeerrors.IsAlreadyExists(err) {
		panic(err)
	}

	newDeployObj := func(name string, labels map[string]string) *appsv1.Deployment {
		deploy := resource.NewDeploymentSample()
		deploy.Namespace = ns
		deploy.Name = name
		deploy.Spec.Selector.MatchLabels = labels
		deploy.Spec.Template.ObjectMeta.Labels = labels
		return deploy
	}

	matchLabels := make(map[string]string)
	matchLabels["app"] = "nginx"
	matchLabels["version"] = "v1"
	matchLabels["test"] = "true"

	_, err = cli.AppsV1().Deployments(ns).Create(ctx, newDeployObj("test-1022-1", map[string]string{
		"app":     "nginx",
		"version": "v1",
		"test":    "true",
	}), metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	_, err = cli.AppsV1().Deployments(ns).Create(ctx, newDeployObj("test-1022-2", map[string]string{
		"app":     "nginx",
		"version": "v1",
	}), metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	_, err = cli.AppsV1().Deployments(ns).Create(ctx, newDeployObj("test-1022-3", map[string]string{
		"app":     "nginx",
		"version": "v2",
		"test":    "true",
	}), metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	svc := resource.NewServiceSample()
	svc.Namespace = ns
	svc.Name = "test-1022"
	svc.Spec.Selector = map[string]string{
		"app":     "nginx",
		"version": "v1",
	}

	_, err = cli.CoreV1().Services(ns).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

}

func TestInitForTestFindServiceDeploys(t *testing.T) {
	initResourceForTestFindServiceDeploys()
}

func TestFindServiceDeploys(t *testing.T) {
	svc := resource.NewServiceSample()
	svc.Namespace = "test-1022-1"
	svc.Name = "test-1022"
	svc.Spec.Selector = map[string]string{
		"app":     "nginx",
		"version": "v1",
	}
	deploys, err := findServiceDeploys(context.TODO(), helper.NewClientSetOrDie(), svc)
	if err != nil {
		panic(err)
	}
	//t.Log(deploys)
	for i := range deploys {
		t.Log(deploys[i].Name)
	}
}

func TestDeleteNamespace(t *testing.T) {
	if err := helper.NewClientSetOrDie().CoreV1().Namespaces().Delete(context.TODO(), "test-1022", metav1.DeleteOptions{}); err != nil {
		panic(err)
	}
}
