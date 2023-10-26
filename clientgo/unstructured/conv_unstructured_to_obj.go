package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// 将 unstructured.Unstructured 转换为具体的 k8s 对象

func main() {
	// Create a new dynamic client.
	restConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	assertNoError(err)

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	assertNoError(err)

	resp, err := dynamicClient.
		Resource(schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "deployments",
		}).
		Namespace(metav1.NamespaceDefault).
		// 你需要确保你连接的集群拥有这个名为 nginx-deployment 的 deployment
		Get(context.TODO(), "nginx-deployment", metav1.GetOptions{})
	assertNoError(err)

	// Convert the unstructured object to cluster.
	unstructured := resp.UnstructuredContent()
	var deploy appsv1.Deployment
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(unstructured, &deploy)
	assertNoError(err)

	// Use the typed object.
	fmt.Println(deploy)
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}
