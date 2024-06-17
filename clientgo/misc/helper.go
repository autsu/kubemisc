package main

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	testNamespace   = "test-1018"
	testDeployName  = "test-1018"
	testServiceName = "test-1018"
)

var (
	globalCliSet     *kubernetes.Clientset
	globalDynamicCli *dynamic.DynamicClient

	onceForCliSet  sync.Once
	onceForDynamic sync.Once
)

func initGlobalClientSet() {
	onceForCliSet.Do(func() {
		cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			panic(err)
		}
		globalCliSet, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			panic(err)
		}
	})
}

func initGlobalDynamicClient() {
	onceForDynamic.Do(func() {
		cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			panic(err)
		}

		globalDynamicCli = dynamic.NewForConfigOrDie(cfg)

	})
}

func initInformer() listerscorev1.ServiceLister {
	if globalCliSet == nil {
		initGlobalClientSet()
	}
	factory := informers.NewSharedInformerFactory(globalCliSet, time.Second*30)
	serviceInformer := factory.Core().V1().Services()
	stopCh := wait.NeverStop
	factory.Start(stopCh)
	//if !cache.WaitForCacheSync(stopCh, serviceInformer.Informer().HasSynced) {
	//	runtime.HandleError(errors.New("failed to sync"))
	//	return nil
	//}
	return serviceInformer.Lister()
}
