package get_service_pod

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	
	"void.io/kubemisc/clientgo/helper"
)

const (
	testNamespace   = "test-1018"
	testDeployName  = "test-1018"
	testServiceName = "test-1018"
)

func TestGetServicePodNumByEndpoints(t *testing.T) {
	cli := helper.NewClientSetOrDie()

	ctx := context.TODO()
	svc, err := cli.CoreV1().Services(testNamespace).Get(ctx, testServiceName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	num := getServicePodNumByEndpoints(ctx, cli, svc)
	t.Log(num)
}

func TestGetServicePodNum(t *testing.T) {
	cli := helper.NewClientSetOrDie()

	ctx := context.TODO()
	svc, err := cli.CoreV1().Services(testNamespace).Get(ctx, testServiceName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	num := getServicePodNum(ctx, cli, svc)
	t.Log(num)
}
