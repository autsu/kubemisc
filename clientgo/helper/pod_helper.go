package helper

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// PodHelperForClientSet wrap pod crud for clientset
type PodHelperForClientSet struct {
	cli *kubernetes.Clientset
	ns  string
}

func NewPodHelperForClientSet(cli *kubernetes.Clientset, namespace string) *PodHelperForClientSet {
	ph := &PodHelperForClientSet{}

	if cli == nil {
		panic("clientset is not set")
	}
	ph.cli = cli

	if namespace == "" {
		ph.ns = metav1.NamespaceDefault
	} else {
		ph.ns = namespace
	}

	return ph
}

func (p *PodHelperForClientSet) Create(ctx context.Context, po *corev1.Pod, opts metav1.CreateOptions) {
	pod, err := p.cli.CoreV1().Pods(p.ns).Create(ctx, po, opts)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return
		}
		panic(err)
	}
	slog.Info("pod create success", slog.String("name", pod.Name))
}

func (p *PodHelperForClientSet) Get(ctx context.Context, name string, opts metav1.GetOptions) *corev1.Pod {
	pod, err := p.cli.CoreV1().Pods(p.ns).Get(ctx, name, opts)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("get pod: %v %v\n", pod.Name, pod.Spec.Containers[0].Image)
	return pod
}

func (p *PodHelperForClientSet) Update(ctx context.Context, newPod *corev1.Pod, opts metav1.UpdateOptions) {
	_, err := p.cli.CoreV1().Pods(p.ns).Update(ctx, newPod, opts)
	if err != nil {
		panic(err)
	}
	slog.Info("pod update success")
}

func (p *PodHelperForClientSet) List(ctx context.Context, opts metav1.ListOptions, handleFunc func([]corev1.Pod)) {
	podList, err := p.cli.CoreV1().Pods(p.ns).List(ctx, opts)
	if err != nil {
		panic(err)
	}
	if handleFunc != nil {
		handleFunc(podList.Items)
	}
}

func (p *PodHelperForClientSet) Watch(ctx context.Context, opts metav1.ListOptions, handleFunc func(watch.Event)) {
	w, err := p.cli.CoreV1().Pods(p.ns).Watch(ctx, opts)
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

func (p *PodHelperForClientSet) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) {
	err := p.cli.CoreV1().Pods(p.ns).Delete(ctx, name, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("delete pod")
}

// PodHelperForRESTClient wrap pod crud for RESTClient
type PodHelperForRESTClient struct {
	cli *rest.RESTClient
	ns  string
}

func NewPodHelperForRESTClient(cli *rest.RESTClient, namespace string) *PodHelperForRESTClient {
	ph := &PodHelperForRESTClient{}

	if cli == nil {
		panic("RESTClient is not set")
	}
	ph.cli = cli

	if namespace == "" {
		ph.ns = metav1.NamespaceDefault
	} else {
		ph.ns = namespace
	}

	return ph
}

func (p *PodHelperForRESTClient) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) {
	err := p.cli.Post().
		Namespace(p.ns).
		Resource("pods").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(pod).
		Do(ctx).
		Error()
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return
		}
		panic(err)
	}
	slog.Info("pod create success", slog.String("name", pod.Name))
}

func (p *PodHelperForRESTClient) Get(ctx context.Context, name string, opts metav1.GetOptions) *corev1.Pod {
	result := &corev1.Pod{}
	err := p.cli.Get().
		Namespace(p.ns).
		Resource("pods").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("get pod: %v %v\n", pod.Name, pod.Spec.Containers[0].Image)
	return result
}

func (p *PodHelperForRESTClient) Update(ctx context.Context, newPod *corev1.Pod, opts metav1.UpdateOptions) {
	result := &corev1.Pod{}
	err := p.cli.Put().
		Namespace(p.ns).
		Resource("pods").
		Name(newPod.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(newPod).
		Do(ctx).
		Into(result)
	if err != nil {
		panic(err)
	}
	slog.Info("pod update success")
}

func (p *PodHelperForRESTClient) List(ctx context.Context, opts metav1.ListOptions, handleListFunc func(items []corev1.Pod)) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &corev1.PodList{}
	err := p.cli.Get().
		Namespace(p.ns).
		Resource("pods").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		panic(err)
	}
	if handleListFunc != nil {
		handleListFunc(result.Items)
	}
}

func (p *PodHelperForRESTClient) Watch(ctx context.Context, opts metav1.ListOptions, handleFunc func(watch.Event)) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	w, err := p.cli.Get().
		Namespace(p.ns).
		Resource("pods").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
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

func (p *PodHelperForRESTClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) {
	err := p.cli.Delete().
		Namespace(p.ns).
		Resource("pods").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
	if err != nil {
		panic(err)
	}
	slog.Info("delete pod success")
}
