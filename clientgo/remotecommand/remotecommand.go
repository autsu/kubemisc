package remotecommand

import (
	"context"

	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type RemoteCommand struct {
	config *restclient.Config
	client clientset.Interface
}

func New(config *restclient.Config, client clientset.Interface) *RemoteCommand {
	return &RemoteCommand{config, client}
}

func NewForConfig(config *restclient.Config) (*RemoteCommand, error) {
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	rc := &RemoteCommand{
		config: config,
		client: client,
	}
	return rc, nil
}

func (rc *RemoteCommand) BuildExecOptions(ctx context.Context, namespace, pod, container string) *ExecOptions {
	options := &ExecOptions{
		config:    rc.config,
		client:    rc.client,
		namespace: namespace,
		pod:       pod,
		container: container,
		ctx:       ctx,
	}
	return options
}
