package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

type ServiceController struct {
	sharedInformerFactory informers.SharedInformerFactory
	workqueue             workqueue.RateLimitingInterface
	serviceLister         listercorev1.ServiceLister
	serviceSynced         cache.InformerSynced
	cli                   *kubernetes.Clientset
}

func NewServiceController(factory informers.SharedInformerFactory) *ServiceController {
	c := &ServiceController{
		sharedInformerFactory: factory,
		workqueue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "service"),
		serviceLister:         factory.Core().V1().Services().Lister(),
		serviceSynced:         factory.Core().V1().Services().Informer().HasSynced,
	}
	factory.Core().V1().Services().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
		DeleteFunc: c.deleteService,
	})
	return c
}

// service 创建事件
// 获取该 service 代理的 deployment，然后更新它们的 annotation，
func (s *ServiceController) addService(obj interface{}) {
	svc := obj.(*corev1.Service)
	s.enqueueService(svc)
}

// service 发生了更新事件
// 比较 oldObj 和 newObj selector 是否变更
// 变更的话，要：
// 获取 oldObj 代理的所有 deployment，把 oldObj 的 id 从它们的 slb annotation 中移除
// 获取 newObj 代理的所有 deployment，把 newObj 的 id 添加到它们的 slb annotation 中
//
// 比较 service type 是否变更
// 如果
func (s *ServiceController) updateService(oldObj, newObj interface{}) {
	//svc := newObj.(*corev1.Service)
	//s.enqueueService(svc)
}

// service 发生了删除事件
//  1. 判断是否是 SLB service
//     1.1 获取该 service 代理的所有 deployment
//     1.2 遍历 deployment，看有没有 "SLB_IDS" 这个 annotation，这个 annotation 的 value 保存的是这个 deployment 绑定的 SLB id
//     1.3 有的话，将当前 service（被删除的这个）从 deployment 的 annotation value 中删除，如果删除后 value 为空，
//     则移除这个 annotation
//     1.4 更新 deployment 的 annotation
func (s *ServiceController) deleteService(obj interface{}) {
	svc := obj.(*corev1.Service)
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		deploys, err := findServiceDeploys(context.TODO(), s.cli, svc)
		if err != nil {

		}
		for _, deploy := range deploys {
			if v, ok := deploy.Annotations["SLB_IDS"]; ok {
				slbIDS := strings.Split(v, ",")
				for idx, id := range slbIDS {
					if id == svc.Labels["SLB_ID"] {
						slbIDS = append(slbIDS[:idx], slbIDS[idx+1:]...)
					}
				}
				if len(slbIDS) == 0 {
					delete(deploy.Annotations, "SLB_IDS")
					continue
				}
				var val string
				for i, id := range slbIDS {
					val += id
					if i < len(slbIDS)-1 {
						val += ","
					}
				}
				deploy.Annotations["SLB_IDS"] = val
			}
		}
	}

	//s.enqueueService(svc)
}

func (s *ServiceController) enqueueService(svc *corev1.Service) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(svc)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	klog.Info(key)
	s.workqueue.Add(key)
}

func (s *ServiceController) worker(ctx context.Context) {
	for s.processNextWorkItem(ctx) {
	}
}

func (s *ServiceController) processNextWorkItem(ctx context.Context) bool {
	key, quit := s.workqueue.Get()
	klog.Info(key)
	if quit {
		return false
	}
	klog.Info("1")
	defer s.workqueue.Done(key)
	err := s.syncHandler(ctx, key.(string))
	if err != nil {
		klog.Error(err)
		return false
	}
	return true
}

func (s *ServiceController) Run(ctx context.Context, workers int) {
	defer runtime.HandleCrash()
	defer s.workqueue.ShutDown()

	if !cache.WaitForNamedCacheSync("service", ctx.Done(), s.serviceSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, s.worker, time.Second)
	}

	<-ctx.Done()
}

func (s *ServiceController) syncHandler(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	klog.Info(namespace, name)
	svc, err := s.serviceLister.Services(namespace).Get(name)
	// delete event
	if errors.IsNotFound(err) {
		klog.Error(err)
		return err
	}
	if err != nil {
		return err
	}
	fmt.Println(svc.Name)
	return nil
}

func main() {
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}
	factory := informers.NewSharedInformerFactory(cli, time.Minute*3)
	controller := NewServiceController(factory)

	factory.Start(stopCh.Done())

	controller.Run(context.TODO(), 2)
}
