package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createPod(ctx context.Context, rc *rest.RESTClient, gvr schema.GroupVersionResource, obj runtime.Object) {
	err := rc.Post().
		Namespace(obj.(*corev1.Pod).Namespace).
		Resource(gvr.Resource).
		Body(obj).
		Do(ctx).
		Error()

	if err != nil && !errors.IsAlreadyExists(err) {
		panic(err.Error())
	}

	fmt.Println("Pod created")
}

func listPods(ctx context.Context, rc *rest.RESTClient, gvr schema.GroupVersionResource) {
	// 创建一个 PodList 对象
	podList := &corev1.PodList{}

	err := rc.Get().
		Namespace(metav1.NamespaceDefault).
		Resource(gvr.Resource).
		VersionedParams(&metav1.ListOptions{LabelSelector: "app=example"}, metav1.ParameterCodec).
		Do(ctx).
		Into(podList)

	if err != nil {
		panic(err.Error())
	}

	fmt.Print("list pod: ")
	for _, item := range podList.Items {
		fmt.Println(item.Name, item.Spec.Containers[0].Image)
	}
}

func getPod(ctx context.Context, rc *rest.RESTClient, gvr schema.GroupVersionResource, name string) *corev1.Pod {
	pod := &corev1.Pod{}

	err := rc.Get().
		Namespace(metav1.NamespaceDefault).
		Resource(gvr.Resource).
		Name(name).
		Do(ctx).
		Into(pod)

	if err != nil {
		panic(err)
	}
	fmt.Println("get pod: ", pod.Name, pod.Spec.Containers[0].Image)
	return pod
}

func updatePod(ctx context.Context, rc *rest.RESTClient, gvr schema.GroupVersionResource, newPod *corev1.Pod) {
	res := rc.Put().
		Namespace(metav1.NamespaceDefault).
		Resource(gvr.Resource).
		Name(newPod.Name).
		Body(newPod).
		Do(ctx)
	if res.Error() != nil {
		makeSureDeletePodWhenPanic(ctx, rc, gvr, newPod.Name)
		panic(res.Error())
	}
	fmt.Println("update pod")
}

func deletePod(ctx context.Context, rc *rest.RESTClient, gvr schema.GroupVersionResource, name string) {
	err := rc.Delete().
		Namespace(metav1.NamespaceDefault).
		Resource(gvr.Resource).
		Name(name).
		Do(ctx).
		Error()
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	fmt.Println("delete pod")
}

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	//clientcmd.BuildConfigFromKubeconfigGetter()

	cfg.APIPath = "/api"
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs

	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err)
	}

	// Pod 的 API 资源类型是 core，API 资源版本是 metav1
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	globalPod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: metav1.NamespaceDefault,
			Labels: map[string]string{
				"app": "example",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	ctx := context.Background()

	createPod(ctx, rc, podGVR, globalPod)
	// k8s 的操作是异步的，它不会按照顺序，等前一个操作执行结束了再执行下一个，比如在 main 函数里通过 client-go 依次顺序调用了 create, get,
	// update, delete 函数，我们期望的是 k8s 也按照代码逻辑的顺序执行相应的操作，create 完再 get，get 完再 update， update 完再 delete，
	// 但是实际在 k8s 层面是不会等 create 完再 get 的，这些操作之间没有前后顺序性，可能 delete 执行完了，update 才开始运行，此时就会失败，因为资源
	// 已经被删除了
	//
	// wait 直到 create 已经完成，再执行下面的代码
	err = wait.PollImmediate(time.Second, time.Second*30, func() (bool, error) {
		localPod := new(corev1.Pod)
		err := rc.Get().Namespace(metav1.NamespaceDefault).Resource(podGVR.Resource).Name(globalPod.Name).Do(ctx).Into(localPod)
		if err != nil {
			return false, err
		}
		if localPod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		makeSureDeletePodWhenPanic(ctx, rc, podGVR, globalPod.Name)
		panic(err)
	}

	po := getPod(ctx, rc, podGVR, globalPod.Name)
	//globalPod.Spec.Containers[0].Image = "nginx:1.21.1"

	// 更新操作，需要先 get 拿到对应的 pod object，再修改字段，最后将这个 pod object 传入 update()
	// 不能直接修改 globalPod 然后传入这个对象，因为 globalPod 是用于 create 的，部分字段无需要提供，
	// 而创建出来的 pod 会自动填充一些字段，比如 ResourceVersion，如果传入 globalPod，那么 globalPod
	// 的 ResourceVesion 是 ""，而实际创建出来的 pod 的这个字段却不为空，k8s 就会认为你要对这个字段进行更新
	// 操作，但是 pod 是不允许更新这个字段的，所以 update 会 error
	newpo := po.DeepCopy()
	newpo.Spec.Containers[0].Image = "nginx:1.21.1"
	updatePod(ctx, rc, podGVR, newpo)
	err = wait.PollImmediate(time.Second, time.Second*30, func() (bool, error) {
		localPod := new(corev1.Pod)
		err := rc.Get().Namespace(metav1.NamespaceDefault).Name(globalPod.Name).Resource(podGVR.Resource).Do(ctx).Into(localPod)
		if err != nil {
			return false, err
		}
		if localPod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		makeSureDeletePodWhenPanic(ctx, rc, podGVR, globalPod.Name)
		panic(err)
	}

	listPods(ctx, rc, podGVR)
	deletePod(ctx, rc, podGVR, globalPod.Name)
	// 等待 delete 完成
	err = wait.PollImmediate(time.Second, time.Second*30, func() (bool, error) {
		localPod := new(corev1.Pod)
		err := rc.Get().Namespace(metav1.NamespaceDefault).Name(globalPod.Name).Resource(podGVR.Resource).Do(ctx).Into(localPod)
		// globalPod 不存在则说明删除完成了，返回 true 结束 wait
		if errors.IsNotFound(err) {
			return true, nil
		}
		// 其他情况都需要继续等待
		if err != nil {
			return false, err
		}
		return false, nil
	})
	// 如果上面不等待，那么这里可能会依然输出 globalPod，因为 delete 操作是异步的，k8s 那边还没删除完成，但是这里的代码却不会阻塞，会继续往下执行到这里，
	// 然后就会输出 globalPod，这是不符合预期结果的
	listPods(ctx, rc, podGVR)
}

func makeSureDeletePodWhenPanic(ctx context.Context, rc *rest.RESTClient, gvr schema.GroupVersionResource, name string) {
	deletePod(ctx, rc, gvr, name)
}
