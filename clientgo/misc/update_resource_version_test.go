package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"void.io/kubemisc/clientgo/helper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 更新会导致对象的 ResourceVersion 发生变化吗
func TestUpdateResourceVersion(t *testing.T) {
	cli := helper.NewClientSetOrDie()
	ctx := context.TODO()
	deployName := "slb-limit-test1"

	deploy, err := cli.AppsV1().Deployments(metav1.NamespaceDefault).Get(ctx, deployName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("before update: ", deploy.ResourceVersion)

	deploy.Spec.Template.Spec.Containers[0].Env = append(
		deploy.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{Name: "key", Value: "value"})
	deploy, err = cli.AppsV1().Deployments(metav1.NamespaceDefault).Update(ctx, deploy, metav1.UpdateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("after update: ", deploy.ResourceVersion)

	// Output:
	// before update:  451683
	// after update:  705596
}

func TestUpdateConflict(t *testing.T) {
	ctx := context.TODO()
	cli := helper.NewClientSetOrDie()
	name := "busybox"
	wg := &sync.WaitGroup{}

	// goroutine A 先拿到该 pod
	po, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	wg.Add(1)

	// goroutine B 会对同一个 pod 进行更新操作，这将导致该 pod 的 resourceVersion 发生变更
	go func() {
		defer wg.Done()
		po, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		if po.Labels == nil {
			po.Labels = make(map[string]string)
		}
		po.Labels["test"] = strconv.Itoa(int(time.Now().Unix()))
		_, err = cli.CoreV1().Pods(metav1.NamespaceDefault).Update(ctx, po, metav1.UpdateOptions{})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 3) // wait pod update finish
	}()

	// 等待 goroutine B 更新完成
	wg.Wait()

	// 此时 goroutine A 再去更新 pod，会发现当前 pod 的 resourceVersion 和之前 get 拿到时的 resourceVersion 不一样，说明
	// 这个 pod 已经被其他人更新过了，会产生一个 conflict 错误
	_, err = cli.CoreV1().Pods(metav1.NamespaceDefault).Update(ctx, po, metav1.UpdateOptions{})
	if err != nil {
		if errors.IsConflict(err) {
			fmt.Println("conflict error")
		}
		// Operation cannot be fulfilled on pods "busybox": the object has been modified; please apply your changes to the latest version and try again
		panic(err)
	}

}
