package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	// rest client 必须要设置 Group/Version 信息
	cfg.GroupVersion = &appsv1.SchemeGroupVersion
	// NegotiatedSerializer 也必须设置
	// NegotiatedSerializer is required when initializing a RESTClient
	cfg.NegotiatedSerializer = scheme.Codecs

	client, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err)
	}

	opts := metav1.ListOptions{
		LabelSelector: "test",
	}

	url := client.Get().
		Resource("collectionresources").
		Name("workloads").
		VersionedParams(&opts, scheme.ParameterCodec).
		URL()

	fmt.Println(url)
}
