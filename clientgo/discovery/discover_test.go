package main

import (
	"fmt"
	"testing"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGetResourceGV(t *testing.T) {
	resources := []string{"pizzas", "pods", "deployments"}
	for _, resource := range resources {
		group, version, find := getResourceGV(resource)
		if !find {
			fmt.Printf("not found resource %v\n", resource)
			continue
		}
		fmt.Printf("%v Group/Version is %v/%v\n", resource, group, version)
	}
}

func TestGetAllGVR(t *testing.T) {
	GetAllGVR()
}

func TestName(t *testing.T) {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		panic(err)
	}

	initialAPIGroupResources, err := restmapper.GetAPIGroupResources(client)
	if err != nil {
		panic(err)
	}

	for _, gr := range initialAPIGroupResources {
		fmt.Printf("%+v\n", gr.Group)
		fmt.Printf("%+v\n", gr.VersionedResources)
		fmt.Println("----------------")
	}
}
