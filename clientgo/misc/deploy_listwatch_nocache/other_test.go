package main

import (
	"context"
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"void.io/kubemisc/clientgo/helper"
)

func TestTolerationsContain(t *testing.T) {
	cli := helper.NewClientSetOrDie()
	dep, err := cli.AppsV1().Deployments(metav1.NamespaceDefault).Get(context.TODO(), "nginx-deployment", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, toleration := range dep.Spec.Template.Spec.Tolerations {
		t.Logf("%+v\n", toleration)
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
	t.Log(slices.Contains(dep.Spec.Template.Spec.Tolerations, t1))
	t.Log(slices.Contains(dep.Spec.Template.Spec.Tolerations, t2))
}
