package main

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var _ cache.ListerWatcher = &LWLimit{}

type LWLimit struct {
	ListFunc  func(options metav1.ListOptions) (runtime.Object, error)
	WatchFunc func(options metav1.ListOptions) (watch.Interface, error)
}

func (l *LWLimit) List(options metav1.ListOptions) (runtime.Object, error) {
	options.Limit = 2
	return l.ListFunc(options)
}

func (l *LWLimit) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return l.WatchFunc(options)
}

func InitPod(ctx context.Context, cli *kubernetes.Clientset, nums int, waitChan chan struct{}) error {
	for i := 0; i < nums; i++ {
		_, err := cli.
			CoreV1().
			Pods(metav1.NamespaceDefault).
			Create(ctx, busyboxPod(), metav1.CreateOptions{ /*DryRun: []string{"All"}*/ })
		if err != nil {
			return err
		}
	}
	waitChan <- struct{}{}
	return nil
}

func busyboxPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "busybox-" + rand.String(5),
			Namespace: metav1.NamespaceDefault,
			Labels:    map[string]string{"test": "true"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func main() {
	ctx := context.TODO()
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	// ch := make(chan struct{}, 1)
	// if err := InitPod(ctx, cli, 20, ch); err != nil {
	// 	panic(err)
	// }
	// <-ch

	// 莫名其妙，这里设置了 limit 为 2，最终获取的结果也是 2，是正确的
	testLimit := func() {
		podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, metav1.ListOptions{Limit: 2})
		if err != nil {
			panic(err)
		}
		podNameList := func(podList *corev1.PodList) (ret []string) {
			for _, item := range podList.Items {
				ret = append(ret, item.Name)
			}
			return
		}(podList)
		klog.Infof("podList: %v, podList size: %v, continue: %v\n", podNameList, len(podNameList), podList.Continue)
	}
	_ = testLimit

	ListFunc := func(opt metav1.ListOptions) (runtime.Object, error) {
		klog.Info("start list...")
		opt.Limit = 2
		// 这里的 ResourceVersion 必须要设置为空，因为 Reflector 内部的 ListAndWatch 会将其设置为 "0"，如果这里不覆盖掉，
		// 会导致忽略 Limit 拉取全量数据，具体可以查一下 ResourceVersion 的相关文档资料
		opt.ResourceVersion = ""
		pods := &corev1.PodList{}
		for {
			podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opt)
			if err != nil {
				panic(err)
			}
			pods.Items = append(pods.Items, podList.Items...)
			podNameList := func(podList *corev1.PodList) (ret []string) {
				for _, item := range podList.Items {
					ret = append(ret, item.Name)
				}
				return
			}(podList)
			klog.Infof("podList: %v, podList size: %v, continue: %v\n", podNameList, len(podNameList), podList.Continue)
			if podList.Continue == "" {
				break
			} else {
				opt.Continue = podList.Continue
			}
		}
		return pods, nil
	}
	//ListFunc(metav1.ListOptions{})

	lw := &LWLimit{
		// ListFunc: func(opt metav1.ListOptions) (runtime.Object, error) {
		// 	klog.Info("start list...")
		// 	// 但是这里设置为 2，后面的 List 无效，还是会拉取全量数据
		// 	opt.Limit = 2
		// 	//klog.Infof("%+v\n", opt)
		// 	podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opt)
		// 	podNameList := func(podList *corev1.PodList) (ret []string) {
		// 		for _, item := range podList.Items {
		// 			ret = append(ret, item.Name)
		// 		}
		// 		return
		// 	}(podList)
		// 	klog.Infof("podList: %v, podList size: %v, continue: %v\n", podNameList, len(podNameList), podList.Continue)
		// 	return podList, err
		// },
		ListFunc: ListFunc,
		WatchFunc: func(opt metav1.ListOptions) (watch.Interface, error) {
			klog.Info("start watch...")
			return cli.CoreV1().Pods(metav1.NamespaceDefault).Watch(ctx, opt)
		},
	}
	s := cache.NewStore(cache.MetaNamespaceKeyFunc)
	r := cache.NewReflector(lw, &corev1.Pod{}, s, 0)
	go r.Run(wait.NeverStop)
	go func() {
		for {
			klog.Infof("store size: %v\n", len(s.List()))
			time.Sleep(time.Second * 3)
		}
	}()

	select {}
}
