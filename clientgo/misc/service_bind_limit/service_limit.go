package service_bind_limit

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"

	"void.io/kubemisc/clientgo/misc/find_deploy_service"
	"void.io/kubemisc/clientgo/misc/find_service_deploy"
)

var limit int32 = 50

func CheckServiceLimit(ctx context.Context, cli *kubernetes.Clientset, deploy *appsv1.Deployment) (bool, error) {
	if deploy == nil {
		return true, nil
	}
	services, err := find_deploy_service.FindDeployServices(ctx, cli, deploy)
	if err != nil {
		return false, err
	}
	if len(services) == 0 {
		return true, nil
	}
	deploys, err := find_service_deploy.FindServiceDeploys(ctx, cli, services[0])
	if err != nil {
		return false, err
	}
	var totalReplicas int32
	if deploy.Spec.Replicas != nil {
		totalReplicas += *deploy.Spec.Replicas
	}
	for _, d := range deploys {
		if d.UID == deploy.UID {
			continue
		}
		if d.Spec.Replicas != nil {
			totalReplicas += *d.Spec.Replicas
		}
	}
	if totalReplicas > limit {
		return false, nil
	}
	return true, nil
}
