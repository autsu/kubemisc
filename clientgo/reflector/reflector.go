package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
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

var _ cache.ListerWatcher = &LW{}

type LW struct {
	ListFunc  func(options metav1.ListOptions) (runtime.Object, error)
	WatchFunc func(options metav1.ListOptions) (watch.Interface, error)
}

func (l *LW) List(options metav1.ListOptions) (runtime.Object, error) {
	return l.ListFunc(options)
}

func (l *LW) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return l.WatchFunc(options)
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

func nginxPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-" + rand.String(5),
			Namespace: metav1.NamespaceDefault,
			Labels:    map[string]string{"test": "true"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}
}

type eventManage struct {
	alreadyCreatedPod []string
	cli               *kubernetes.Clientset
}

func (e *eventManage) newPod(ctx context.Context) error {
	pod, err := e.cli.
		CoreV1().
		Pods(metav1.NamespaceDefault).
		Create(ctx, busyboxPod(), metav1.CreateOptions{ /*DryRun: []string{"All"}*/ })
	if err != nil {
		return err
	}
	klog.Info("[eventManage] create new pod", "pod name", pod.Name)
	e.alreadyCreatedPod = append(e.alreadyCreatedPod, pod.Name)
	return nil
}

func (e *eventManage) deletePod(ctx context.Context) error {
	if len(e.alreadyCreatedPod) == 0 {
		return nil
	}
	err := e.cli.
		CoreV1().
		Pods(metav1.NamespaceDefault).
		Delete(ctx, e.alreadyCreatedPod[0], metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	klog.Info("[eventManage] delete pod", "pod name", e.alreadyCreatedPod[0])
	e.alreadyCreatedPod = e.alreadyCreatedPod[1:]
	return nil
}

func (e *eventManage) randomPodEvents(ctx context.Context) error {
	var t = time.NewTicker(time.Second * 5)
	var err error
	var errNum int

	for range t.C {
		if errNum == 5 {
			t.Stop()
			return fmt.Errorf("[eventManage] two many errors, last error: %v", err)
		}
		i := rand.Intn(10)
		switch {
		case i%2 == 0:
			klog.Info("[eventManage] delete event")
			if err = e.deletePod(ctx); err != nil {
				klog.Error(err)
				errNum++
			}
		default:
			klog.Info("[eventManage] create event")
			if err = e.newPod(ctx); err != nil {
				klog.Error(err)
				errNum++
			}
		}
	}

	return err
}

func listStoreState(s cache.Store) {
	var t = time.NewTicker(time.Second * 3)
	for range t.C {
		klog.Infof("\n")
		klog.Info("====================================")
		klog.Infof("[store] current store size: %v\n", len(s.ListKeys()))
		klog.Infof("[store] current store keys: %v\n", s.ListKeys())
		klog.Info("====================================")
		klog.Infof("\n")
	}
}

func main() {
	ctx := context.TODO()
	wg := errgroup.Group{}
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	em := &eventManage{
		alreadyCreatedPod: []string{},
		cli:               cli,
	}

	lw := &LW{
		ListFunc: func(opt metav1.ListOptions) (runtime.Object, error) {
			klog.Info("start list...")
			return cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (watch.Interface, error) {
			klog.Info("start watch...")
			return cli.CoreV1().Pods(metav1.NamespaceDefault).Watch(ctx, opt)
		},
	}

	s := cache.NewStore(cache.MetaNamespaceKeyFunc)
	r := cache.NewReflector(lw, &corev1.Pod{}, s, 0)
	go r.Run(wait.NeverStop)

	wg.Go(func() error {
		return em.randomPodEvents(ctx)
	})
	wg.Go(func() error {
		listStoreState(s)
		return nil
	})

	// clear test pod when exist
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, os.Kill)
		<-sigch
		if err := cli.CoreV1().
			Pods(metav1.NamespaceDefault).
			DeleteCollection(ctx, metav1.DeleteOptions{},
				metav1.ListOptions{
					LabelSelector: "test=true",
				}); err != nil {
			panic(err)
		}
		os.Exit(1)
	}()

	if err := wg.Wait(); err != nil {
		panic(err)
	}
}
