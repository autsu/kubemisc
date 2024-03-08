package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/pager"

	"void.io/kubemisc/clientgo/helper"
)

// TODO：代码还不能用，不知道这个库的作用和使用方法
func main() {
	p := pager.New(func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
		cli := helper.NewClientSetOrDie()
		podList, err := cli.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, opts)
		if err != nil {
			return nil, err
		}
		return podList, nil
	})

	object, _, err := p.List(context.TODO(), metav1.ListOptions{Limit: 5})
	if err != nil {
		panic(err)
	}

	list, err := meta.ExtractList(object)
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		fmt.Println(item.GetObjectKind().GroupVersionKind().String())
	}
}
