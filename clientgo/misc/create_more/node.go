package main

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"void.io/kubemisc/clientgo/helper"
)

func main() {
	cli := helper.NewClientSetOrDie()

	for i := range 100 {
		node := &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-" + strconv.Itoa(i),
			},
			Spec: corev1.NodeSpec{},
		}
		_, err := cli.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
	}
}
