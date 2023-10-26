package wait

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func TestPollImmediateUntil(t *testing.T) {
	count := 0
	stopCh := make(chan struct{})
	go func() {
		time.AfterFunc(time.Second*5, func() {
			stopCh <- struct{}{}
		})
	}()
	// 每隔 time.Second 执行一次 func，除非遇到下面几种情况才会停止：
	// 1. done 返回 true
	// 2. err != nil
	// 3. 能从 stopCh 读取到值
	if err := wait.PollImmediateUntil(time.Second, func() (done bool, err error) {
		if count < 10 {
			fmt.Println(count)
			count++
		} else {
			//done = true
			return false, errors.New("test error")
		}
		return
	}, stopCh); err != nil {
		t.Fatal(err)
	}
}

func TestWaitUntil(t *testing.T) {
	// wait.Until 和 PollImmediateUntil 有点类似，不过它的停止情况只有一种，
	// 就是能从 stopCh 中读取到值
	// 它的回调函数是没有返回值的
	go wait.Until(func() {
		fmt.Println("1")
	}, time.Second, wait.NeverStop)
}

func TestPollUntilContextCancel(t *testing.T) {
	wait.PollUntilContextCancel(
		context.TODO(),
		time.Second*5, // interval
		true,          // immediate，如果为 true，则立即执行一次 func，否则等待一次 interval 再执行
		func(ctx context.Context) (done bool, err error) {
			return false, err
		})
}
