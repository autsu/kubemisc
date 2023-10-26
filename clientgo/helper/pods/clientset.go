package pods

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// PodHelperForClientSet wrap pod crud for clientset
type PodHelperForClientSet struct {
	cli       *kubernetes.Clientset
	namespace string
}

func NewPodHelperForClientSet(cli *kubernetes.Clientset, namespace string) *PodHelperForClientSet {
	ph := &PodHelperForClientSet{}

	if cli == nil {
		panic("clientset is not set")
	}
	ph.cli = cli

	if namespace == "" {
		ph.namespace = metav1.NamespaceDefault
	} else {
		ph.namespace = namespace
	}

	return ph
}

func (p *PodHelperForClientSet) MustCreate(ctx context.Context, po *corev1.Pod, opts metav1.CreateOptions) {
	_, err := p.cli.CoreV1().Pods(p.namespace).Create(ctx, po, opts)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return
		}
		panic(err)
	}
	//slog.Info("pod create success", slog.String("name", pod.Name))
}

func (p *PodHelperForClientSet) MustGet(ctx context.Context, name string, opts metav1.GetOptions) *corev1.Pod {
	pod, err := p.cli.CoreV1().Pods(p.namespace).Get(ctx, name, opts)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("get pod: %v %v\n", pod.Name, pod.Spec.Containers[0].Image)
	return pod
}

func (p *PodHelperForClientSet) MustUpdate(ctx context.Context, newPod *corev1.Pod, opts metav1.UpdateOptions) {
	_, err := p.cli.CoreV1().Pods(p.namespace).Update(ctx, newPod, opts)
	if err != nil {
		panic(err)
	}
	//slog.Info("pod update success")
}

func (p *PodHelperForClientSet) MustList(ctx context.Context, opts metav1.ListOptions, handleFunc func([]corev1.Pod)) {
	podList, err := p.cli.CoreV1().Pods(p.namespace).List(ctx, opts)
	if err != nil {
		panic(err)
	}
	if handleFunc != nil {
		handleFunc(podList.Items)
	}
}

func (p *PodHelperForClientSet) MustWatch(ctx context.Context, opts metav1.ListOptions, handleFunc func(watch.Event)) {
	w, err := p.cli.CoreV1().Pods(p.namespace).Watch(ctx, opts)
	if err != nil {
		panic(err)
	}
	if handleFunc == nil {
		handleFunc = func(event watch.Event) {
			fmt.Printf("watch a pod event, name=%v, eventType=%v\n", event.Object.(*corev1.Pod).Name, event.Type)
		}
	}
	for event := range w.ResultChan() {
		handleFunc(event)
	}
}

func (p *PodHelperForClientSet) MustDelete(ctx context.Context, name string, opts metav1.DeleteOptions) {
	err := p.cli.CoreV1().Pods(p.namespace).Delete(ctx, name, opts)
	if err != nil {
		panic(err)
	}
}

func (p *PodHelperForClientSet) MustPath(ctx context.Context, name string, data []byte, opts metav1.PatchOptions) {
	_, err := p.cli.CoreV1().Pods(p.namespace).Patch(ctx, name, types.StrategicMergePatchType, data, opts)
	if err != nil {
		panic(err)
	}
}
