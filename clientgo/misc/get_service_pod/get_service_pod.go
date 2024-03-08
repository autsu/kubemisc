package get_service_pod

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func getServicePodNumByEndpoints(ctx context.Context, cli *kubernetes.Clientset, svc *corev1.Service) int64 {
	if cli == nil || svc == nil {
		return 0
	}
	ep, err := cli.CoreV1().Endpoints(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	//for _, subset := range ep.Subsets {
	//	fmt.Printf("%+v\n", subset)
	//}
	count := 0
	for _, subset := range ep.Subsets {
		count += len(subset.Addresses)
	}
	//return int64(len(ep.Subsets))
	return int64(count)
}

func getServicePodNum(ctx context.Context, cli *kubernetes.Clientset, svc *corev1.Service) int64 {
	if cli == nil || svc == nil {
		return 0
	}
	ls := labels.Set(svc.Spec.Selector).AsSelector().String()
	podList, err := cli.CoreV1().Pods(svc.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: ls,
	})
	if err != nil {
		panic(err)
	}
	return int64(len(podList.Items))
}
