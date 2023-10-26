package main

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetServicePodNumByEndpoints(t *testing.T) {
	initClientSet()

	ctx := context.TODO()
	svc, err := cli.CoreV1().Services(testNamespace).Get(ctx, testServiceName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	num := getServicePodNumByEndpoints(ctx, cli, svc)
	t.Log(num)
}

func TestGetServicePodNum(t *testing.T) {
	initClientSet()

	ctx := context.TODO()
	svc, err := cli.CoreV1().Services(testNamespace).Get(ctx, testServiceName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	num := getServicePodNum(ctx, cli, svc)
	t.Log(num)
}
