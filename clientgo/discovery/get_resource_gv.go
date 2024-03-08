package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

// 假如现在我们有一个自定义资源 pizzas，并且该 resource 已经注册进了 k8s（通过 APIService），我们只需要执行
// kubectl get pizzas 即可查找该资源，问题是我们没有提供任何 group/version 的信息，那 kubectl 是如何知道
// pizzas 这个资源在哪个 gv 下，从而构建请求 URL 的呢？
//
// 通过下面这个函数就可以找到，使用 discovery client
// 通过这个例子有助于理解 discovery client 的作用
func getResourceGV(wantResource string) (group, version string, find bool) {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		panic(err)
	}

	_, resourceLists, err := client.ServerGroupsAndResources()
	for _, list := range resourceLists {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			panic(err)
		}
		for _, resource := range list.APIResources {
			if resource.Name == wantResource {
				fmt.Printf("%v Group/Version is %v/%v\n", wantResource, gv.Group, gv.Version)
				return gv.Group, gv.Version, true
			}
		}
	}
	return "", "", false
}
