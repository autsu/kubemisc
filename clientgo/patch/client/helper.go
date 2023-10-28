package client

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"void.io/kubemisc/clientgo/helper"
	"void.io/kubemisc/clientgo/helper/resource"
)

func initTestResource(deleteAfterTest bool) (deploy *appsv1.Deployment, dropFunc func()) {
	name := "test-1027"
	ctx := context.TODO()
	deploy = resource.NewDeploymentBy(name, metav1.NamespaceDefault, map[string]string{"a": "b"})
	deploy.Spec.Template.Spec.Tolerations = append(deploy.Spec.Template.Spec.Tolerations, corev1.Toleration{
		Key:               "dedicated",
		Operator:          "",
		Value:             "test-team",
		Effect:            corev1.TaintEffectNoSchedule,
		TolerationSeconds: nil,
	})

	cli := helper.NewClientSetOrDie()
	deploy, err := cli.AppsV1().Deployments(metav1.NamespaceDefault).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	if deleteAfterTest {
		dropFunc = func() { cli.AppsV1().Deployments(metav1.NamespaceDefault).Delete(ctx, name, metav1.DeleteOptions{}) }
	}
	return deploy, dropFunc
}
