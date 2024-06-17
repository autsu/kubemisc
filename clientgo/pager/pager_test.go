package main

import (
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/pager"

	"void.io/kubemisc/clientgo/helper"
)

func TestListPage(t *testing.T) {
	cli := helper.NewClientSetOrDie()
	offset := 1
	limit := 5
	listOpt := metav1.ListOptions{Limit: int64(limit)}

	for {
		list, err := cli.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), listOpt)
		if err != nil {
			t.Fatal(err)
		}

		if list.Continue == "" {
			break
		}

		listOpt.Continue = list.Continue

		t.Logf("offset: %d, limit: %d", offset, limit)
		for _, pod := range list.Items {
			t.Logf("namespace/name: %s/%s", pod.Namespace, pod.Name)
		}
		offset++
		t.Logf("\n")
	}
}

// 貌似这个 pager 只是在内部请求 apiserver 时会分页，但最终还是会把拿到的所有数据一起返回给我们
// 如果是这样的话，那使用这个包只是可以降低 apiserver 的压力
// 如果调用方（比如 Paas 平台）想避免内存压力（比如某个集群有 100 万个 pod，我不想一次拿到所有的 pod，不然内存压力会很大），
// 这个时候用这个包就没用了，得用 TestListPage 里的方法去分页获取

func Test1(t *testing.T) {
	cli := helper.NewClientSetOrDie()

	p := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opts)
		if err != nil {
			return nil, err
		}
		return podList, nil
	})

	// 如果指定了 Limit，那么返回的 list 类型是 internalversion.List
	// 而不是 corev1.PodList
	//
	// 这里虽然指定了 Limit，但最终返回的 list 依然是全量数据
	list, hasMore, err := p.List(context.TODO(), metav1.ListOptions{Limit: 5, ResourceVersion: ""})
	if err != nil {
		t.Fatal(err)
	}
	_ = hasMore

	podList, err := meta.ExtractList(list)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range podList {
		pod := item.(*corev1.Pod)
		t.Logf("namespace/name: %s/%s", pod.Namespace, pod.Name)
	}

}

func Test2(t *testing.T) {
	cli := helper.NewClientSetOrDie()

	p := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opts)
		if err != nil {
			return nil, err
		}
		return podList, nil
	})

	list, hasMore, err := p.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	_ = hasMore

	pods, ok := list.(*corev1.PodList)
	if !ok {
		fmt.Printf("unexpected type %T\n", list)
		return
	}

	for _, item := range pods.Items {
		t.Logf("namespace/name: %s/%s", item.Namespace, item.Name)
	}
}

// 使用 EachListItem
func Test3(t *testing.T) {
	cli := helper.NewClientSetOrDie()

	p := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opts)
		if err != nil {
			return nil, err
		}
		return podList, nil
	})

	err := p.EachListItem(context.TODO(), metav1.ListOptions{Limit: 5, ResourceVersion: ""}, func(obj runtime.Object) error {
		pod, ok := obj.(*corev1.Pod)
		if !ok {
			return fmt.Errorf("expected Pod, got %T", obj)
		}
		t.Logf("namespace/name: %s/%s", pod.Namespace, pod.Name)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
