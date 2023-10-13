package helper

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestPrintNameForResourceList(t *testing.T) {
	podL := []*corev1.Pod{
		NewPodSimple(),
		NewPodSimple(),
	}

	deployL := []*appsv1.Deployment{
		NewDeploymentSimple(),
		NewDeploymentSimple(),
	}

	PrintNameForResourceList(podL)
	PrintNameForResourceList(deployL)
}
