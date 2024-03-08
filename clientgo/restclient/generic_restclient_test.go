package main

import (
	"testing"

	"void.io/kubemisc/clientgo/helper/resource"

	"k8s.io/client-go/tools/clientcmd"
)

func TestNewGenericRESTClientWithGroupVersion(t *testing.T) {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		t.Fatal(err)
	}
	cli, err := NewGenericRESTClientWithGroupVersion(cfg, resource.GVR.Pod().GroupVersion())
	if err != nil {
		t.Fatal(err)
	}
	_ = cli
	//cli.Get().NamespaceIfScoped().Resource(resource.GvrPods.Resource).Name().Do().Get()
}
