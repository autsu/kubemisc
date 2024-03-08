package main

import (
	"time"
	"void.io/kubemisc/clientgo/helper/resource"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type InformerManager struct {
	cli                     kubernetes.Interface
	gvrsMapGenericInformers map[schema.GroupVersionResource]informers.GenericInformer
	//gvrs             []schema.GroupVersionResource
	//genericInformers []informers.GenericInformer
}

func NewInformerManager(cli kubernetes.Interface, gvrs []schema.GroupVersionResource) *InformerManager {
	ret := &InformerManager{cli: cli}
	for _, gvr := range gvrs {
		genericInformer, err := informers.NewSharedInformerFactory(cli, time.Second*30).ForResource(gvr)
		if err != nil {
			panic(err)
		}
		ret.gvrsMapGenericInformers[gvr] = genericInformer
	}
	return ret
}

func (i *InformerManager) Run(stopCh <-chan struct{}) {
	for _, informer := range i.gvrsMapGenericInformers {
		go informer.Informer().Run(stopCh)
	}
}

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	manager := NewInformerManager(cli, []schema.GroupVersionResource{resource.GVR.Pod(), resource.GVR.Service(), resource.GVR.Deployment()})
	manager.Run(wait.NeverStop)

}
