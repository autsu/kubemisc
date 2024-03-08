package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/tools/clientcmd"

	"void.io/kubemisc/clientgo/helper/resource"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	cli, err := metadata.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	data, err := cli.Resource(resource.GVR.Deployment()).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, item := range data.Items {
		fmt.Println(item)
	}
}
