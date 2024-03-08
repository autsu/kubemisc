package service_bind_limit

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"void.io/kubemisc/clientgo/helper"
)

// deploy1 找到 -> service 找到 -> deploy1, deploy2, deploy3
func TestCheckServiceLimit(t *testing.T) {
	cli := helper.NewClientSetOrDie()
	ctx := context.TODO()

	deploy, err := cli.AppsV1().Deployments(metav1.NamespaceDefault).Get(ctx, "slb-limit-test1", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	replicas := pointer.Int32Deref(deploy.Spec.Replicas, 0)
	deploy.Spec.Replicas = pointer.Int32(15)

	allow, err := CheckServiceLimit(ctx, helper.NewClientSetOrDie(), deploy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(allow)

	deploy.Spec.Replicas = pointer.Int32(replicas)
	allow, err = CheckServiceLimit(ctx, helper.NewClientSetOrDie(), deploy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(allow)

	deploy.Spec.Replicas = pointer.Int32(5)
	allow, err = CheckServiceLimit(ctx, helper.NewClientSetOrDie(), deploy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(allow)
}
