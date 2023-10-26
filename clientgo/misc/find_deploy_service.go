package main

import (
	"context"
	"log/slog"

	"void.io/kubemisc/clientgo/helper/maps"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	informersappsv1 "k8s.io/client-go/informers/apps/v1"
	informerscorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

func findDeployServices(ctx context.Context, cli *kubernetes.Clientset, deploy *appsv1.Deployment) ([]*corev1.Service, error) {
	if cli == nil || deploy == nil {
		return nil, nil
	}
	var ret []*corev1.Service
	services, err := cli.CoreV1().Services(deploy.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, service := range services.Items {
		// if service.Spec.Type != corev1.ServiceTypeLoadBalancer {
		// 	continue
		// }
		if !maps.Contains(deploy.Spec.Selector.MatchLabels, service.Spec.Selector) {
			continue
		}
		ret = append(ret, &service)
	}
	return ret, nil
}

func serviceIsProxyDeploy(svc *corev1.Service, deploy *appsv1.Deployment) bool {
	if svc == nil || deploy == nil {
		return false
	}
	serviceSelector := svc.Spec.Selector
	for k, v := range serviceSelector {
		v1, ok := deploy.Spec.Selector.MatchLabels[k]
		if !ok || v != v1 {
			return false
		}
	}
	return true
}

// ========================= junk code =========================

var DeployMapServiceCache = make(map[string][]string)

// 有问题，无法保证缓存与 k8s 的一致性，这种方案废弃
func __findDeployServiceByClientSetWithCache__(ctx context.Context, cli *kubernetes.Clientset, deploy *appsv1.Deployment) ([]*corev1.Service, error) {
	if cli == nil || deploy == nil {
		return nil, nil
	}
	var ret []*corev1.Service
	deployKey, err := cache.MetaNamespaceKeyFunc(deploy)
	if v, ok := DeployMapServiceCache[deployKey]; ok {
		slog.Debug("hit cache")
		for i, vv := range v {
			namespace, name, err := cache.SplitMetaNamespaceKey(vv)
			if err != nil {
				return nil, err
			}
			svc, err := cli.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			// deploy 的标签可能更新，或者 service 已解绑，所以从缓存中拿到后需要二次检查
			if !serviceIsProxyDeploy(svc, deploy) {
				// 从缓存中移除
				DeployMapServiceCache[deployKey] = append(DeployMapServiceCache[deployKey][:i], DeployMapServiceCache[deployKey][i+1:]...)
				continue
			}
			ret = append(ret, svc)
		}
		return ret, nil
	}
	slog.Debug("cache miss")
	// cache miss
	services, err := cli.CoreV1().Services(deploy.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, service := range services.Items {
		//if service.Spec.Type != "" {
		//	continue
		//}
		if !serviceIsProxyDeploy(&service, deploy) {
			continue
		}
		cacheKey, err1 := cache.MetaNamespaceKeyFunc(deploy)
		cacheVal, err2 := cache.MetaNamespaceKeyFunc(&service)
		if err1 == nil && err2 == nil {
			DeployMapServiceCache[cacheKey] = append(DeployMapServiceCache[cacheKey], cacheVal)
		}
		ret = append(ret, &service)
	}
	return ret, nil
}

func __findDeployServiceByInformerWithCache__(ctx context.Context, lister listerscorev1.ServiceLister, deploy *appsv1.Deployment) ([]*corev1.Service, error) {
	if cli == nil || deploy == nil {
		return nil, nil
	}
	var ret []*corev1.Service
	deployKey, err := cache.MetaNamespaceKeyFunc(deploy)
	if v, ok := DeployMapServiceCache[deployKey]; ok {
		slog.Debug("hit cache")
		for i, vv := range v {
			namespace, name, err := cache.SplitMetaNamespaceKey(vv)
			if err != nil {
				return nil, err
			}
			svc, err := cli.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			// deploy 的标签可能更新，或者 service 已解绑，所以从缓存中拿到后需要二次检查
			if !serviceIsProxyDeploy(svc, deploy) {
				// 从缓存中移除
				DeployMapServiceCache[deployKey] = append(DeployMapServiceCache[deployKey][:i], DeployMapServiceCache[deployKey][i+1:]...)
				continue
			}
			ret = append(ret, svc)
		}
		return ret, nil
	}
	slog.Debug("cache miss")
	// cache miss
	services, err := lister.Services(deploy.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		//if service.Spec.Type != "" {
		//	continue
		//}
		if !serviceIsProxyDeploy(service, deploy) {
			continue
		}
		cacheKey, err1 := cache.MetaNamespaceKeyFunc(deploy)
		cacheVal, err2 := cache.MetaNamespaceKeyFunc(service)
		if err1 == nil && err2 == nil {
			DeployMapServiceCache[cacheKey] = append(DeployMapServiceCache[cacheKey], cacheVal)
		}
		ret = append(ret, service)
	}
	return ret, nil
}

// ChatGPT 生成的
// 不对，这个是用 deploy 的 selector label 作为查询条件，去找 metadata.Labels 里面
// 同样拥有此 label 的 service，我们期望的是查找 Service.spec.selector 中拥有的
func __findDeployService2Err__(cli *kubernetes.Clientset, namespace string, deploy *appsv1.Deployment) ([]corev1.Service, error) {
	selector := deploy.Spec.Selector.MatchLabels
	servicesClient := cli.CoreV1().Services(namespace)
	ls := labels.Set(selector).String()
	services, err := servicesClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: ls,
	})
	if err != nil {
		panic(err)
	}
	return services.Items, nil
}

type Informers struct {
	ServiceInformer    informerscorev1.ServiceInformer
	DeploymentInformer informersappsv1.DeploymentInformer
}

func InitInformers() {

}

func InitTest() {

}
