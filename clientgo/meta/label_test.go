package main

import (
	"fmt"
	"testing"

	"void.io/kubemisc/clientgo/helper/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSample(t *testing.T) {
	deploy := resource.NewDeploymentSample()
	labels := deploy.Spec.Selector.MatchLabels
	labels["app"] = "nginx"
	labels["test"] = "true"
	fmt.Println(deploy.Spec.Selector.MatchLabels)

	selector, err := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	if err != nil {
		panic(err)
	}
	_ = selector
}
