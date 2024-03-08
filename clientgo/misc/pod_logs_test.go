package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPodLogStream(t *testing.T) {
	initGlobalClientSet()

	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancelFunc()

	stream, err := globalCliSet.
		CoreV1().
		Pods(metav1.NamespaceDefault).
		GetLogs("nginx-deployment-7c5ddbdf54-669hf", &corev1.PodLogOptions{
			Container: "",
			Follow:    true,
			TailLines: nil,
		}).
		Stream(ctx)
	if err != nil {
		panic(err)
	}
	//defer stream.Close()

	//go func() {
	//	<-time.After(time.Second * 3)
	//	fmt.Println("close stream")
	//	stream.Close()
	//}()

	bufioStream := bufio.NewReader(stream)

	//b := make([]byte, 4096)
	for {
		b, err := bufioStream.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				return
			}
			panic(err)
		}
		fmt.Println(string(b))
	}
}
