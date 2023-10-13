package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	pod       = flag.String("pod", "", "pod name")
	namespace = flag.String("n", metav1.NamespaceDefault, "namespace")
	shellType = flag.String("t", "bash", "bash | sh")
	container = flag.String("c", "", "container name")
)

// go run std_example.go -pod busybox -n dsp-test
func main() {
	flag.Parse()
	if *pod == "" {
		panic("miss -pod flag param")
	}

	ctx := context.TODO()
	_ = ctx

	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// 构建出原始的 HTTP 请求 URL
	req := clientset.
		CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(*pod).
		Namespace(*namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: *container,
		Command:   []string{*shellType},
	}, scheme.ParameterCodec)

	// https://10.21.5.74:16443/api/v1/namespaces/default/pods/busybox/exec?command=sh&stderr=true&stdin=true&stdout=true&tty=true
	fmt.Println(req.URL())

	executor, err := remotecommand.NewSPDYExecutor(cfg, http.MethodPost, req.URL())
	if err != nil {
		panic(err)
	}
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		// 全部绑定到 std
		// 这样从本地 stdin 输入的内容会发送给容器内的 stdin
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: nil,
	}); err != nil {
		panic(err)
	}
}
