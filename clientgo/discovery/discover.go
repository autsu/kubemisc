package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

// discoveryClient 可以查询集群里的 gvr 之类的信息
// 没太搞明白这玩意的使用场景，除了 kubectl api-resources
// 老实说一直没搞懂 resource 和 kind 的区别，为毛需要这两种定义，gvr 和 gvk 有啥区别？什么时候用 gvr，什么时候用 gvk？

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		panic(err)
	}

	apiGroup, apiResourceListSlice, err := client.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}

	//fmt.Printf("APIGroup :\n\n %v\n\n\n\n", apiGroup)

	fmt.Println("\n============ group list =============")
	for _, group := range apiGroup {
		fmt.Println(group.Name)
	}
	fmt.Printf("\n================================\n")

	for _, apiResourceList := range apiResourceListSlice {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			panic(err)
		}
		fmt.Printf("gv: %v\n", gv)

		// 获取该 gv 下的所有 resource
		for _, r := range apiResourceList.APIResources {
			fmt.Println(r.Name)
		}
		fmt.Printf("\n\n")
	}
}
