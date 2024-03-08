package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

// discoveryClient 可以查询集群里的 gvr 之类的信息
// 没太搞明白这玩意的使用场景，除了 kubectl api-resources
// 老实说一直没搞懂 resource 和 kind 的区别，为毛需要这两种定义，gvr 和 gvk 有啥区别？什么时候用 gvr，什么时候用 gvk？

func GetAllGVR() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		panic(err)
	}

	version, err := client.ServerVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("server version: ", version)

	apiGroup, apiResourceListSlice, err := client.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}

	func() {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Group"})

		for _, group := range apiGroup {
			table.Append([]string{group.Name})
		}

		table.Render()
	}()

	func() {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Group/Version", "Resources"})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)

		var datas [][]string

		for _, apiResourceList := range apiResourceListSlice {
			var data []string
			gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
			if err != nil {
				panic(err)
			}
			data = append(data, gv.String())

			// 获取该 gv 下的所有 resource
			var names string
			for _, r := range apiResourceList.APIResources {
				names += r.Name + ", "
			}
			names = names[:len(names)-2]
			data = append(data, names)
			datas = append(datas, data)
		}

		for _, data := range datas {
			table.Append(data)
		}

		table.Render()
	}()
}
