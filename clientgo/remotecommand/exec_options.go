package remotecommand

import (
	"context"
	"fmt"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// IOStreams provides the standard names for iostreams.  This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

type ExecOptions struct {
	config    *rest.Config
	client    kubernetes.Interface
	namespace string
	pod       string
	container string
	command   []string
	err       error
	ctx       context.Context

	IOStreams
}

func (o *ExecOptions) execute() error {
	pod, err := o.client.CoreV1().Pods(o.namespace).Get(o.ctx, o.pod, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// pod 状态判断
	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}

	if o.container == "" {
		o.container = pod.Spec.Containers[0].Name
	}

	podOptions := &corev1.PodExecOptions{
		Stdin:     o.In != nil,
		Stdout:    o.Out != nil,
		Stderr:    o.ErrOut != nil,
		Container: o.container,
		Command:   o.command,
	}

	req := o.client.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(podOptions, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(o.config, http.MethodPost, req.URL())
	if err != nil {
		return err
	}

	streamOptions := remotecommand.StreamOptions{
		Stdin:  o.In,
		Stdout: o.Out,
		Stderr: o.ErrOut,
	}

	return exec.StreamWithContext(o.ctx, streamOptions)
}
