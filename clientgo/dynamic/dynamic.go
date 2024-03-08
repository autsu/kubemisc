package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	resourcehelper "void.io/kubemisc/clientgo/helper/resource"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	dclient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	for _, gvr := range []schema.GroupVersionResource{resourcehelper.GVR.Pod(), resourcehelper.GVR.Deployment()} {
		// 根据 GVR 来 list 对象
		res, err := dclient.Resource(gvr).Namespace(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, item := range res.Items {
			fmt.Printf("kind: %v, namespace/name: %v/%v\n", item.GetKind(), item.GetNamespace(), item.GetName())
		}
	}
}
