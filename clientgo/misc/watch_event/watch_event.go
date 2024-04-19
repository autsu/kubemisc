package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"void.io/kubemisc/clientgo/helper"
)

func main() {
	cli := helper.NewClientSetOrDie()
	watch, err := cli.CoreV1().Events(metav1.NamespaceDefault).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for e := range watch.ResultChan() {
		event := e.Object.(*corev1.Event)
		if event.InvolvedObject.Kind == "Deployment" {
			fmt.Printf("Reason: %+v, Message: %v\n", event.Reason, event.Message)
		}
	}
}
