package print

import (
	"testing"

	"void.io/kubemisc/clientgo/helper/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestPrintNameForResourceList(t *testing.T) {
	podL := []*corev1.Pod{
		resource.NewPodSample(),
		resource.NewPodSample(),
	}

	deployL := []*appsv1.Deployment{
		resource.NewDeploymentSample(),
		resource.NewDeploymentSample(),
	}

	ResourceListItemName(podL)
	ResourceListItemName(deployL)
}
