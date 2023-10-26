package patch

import (
	"encoding/json"
	"testing"

	"void.io/kubemisc/clientgo/helper/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func StrategicPatchDeploy(deploy *appsv1.Deployment, patchData []byte) *appsv1.Deployment {
	deployJson, err := json.Marshal(deploy)
	if err != nil {
		panic(err)
	}
	patchDeploy := new(appsv1.Deployment)

	patchDeployByte, err := strategicpatch.StrategicMergePatch(
		deployJson, patchData, &appsv1.Deployment{})
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(patchDeployByte, patchDeploy); err != nil {
		panic(err)
	}
	return patchDeploy
}

func TestStrategicPatch(t *testing.T) {
	tests := []struct {
		name            string
		deploy          *appsv1.Deployment
		patchData       []byte
		beforePatchFunc func(*appsv1.Deployment)
		afterPatchFunc  func(*appsv1.Deployment)
	}{
		{
			name: "patch labels",
			deploy: func() *appsv1.Deployment {
				deploy := resource.NewDeploymentSample()
				deploy.Labels["patch"] = "false"
				deploy.Labels["test"] = "true"
				return deploy
			}(),
			patchData:       []byte(`{"metadata": {"labels": {"new-label": "new-value", "patch": "true"}}}`),
			beforePatchFunc: func(deploy *appsv1.Deployment) { t.Log("before patch: ", deploy.Labels) },
			afterPatchFunc:  func(deploy *appsv1.Deployment) { t.Log("after patch: ", deploy.Labels) },
		},
		{
			name: "patch container",
			deploy: func() *appsv1.Deployment {
				deploy := resource.NewDeploymentSample()
				return deploy
			}(),
			patchData: []byte(`{"spec":{"template":{"spec":{"containers":[{"name":"patch-demo-ctr-2","image":"redis"}]}}}}`),
			beforePatchFunc: func(deploy *appsv1.Deployment) {
				t.Log("before patch: ")
				for _, container := range deploy.Spec.Template.Spec.Containers {
					t.Logf("container name: %v, image: %v\n", container.Name, container.Image)
				}
			},
			afterPatchFunc: func(deploy *appsv1.Deployment) {
				t.Log("after patch: ")
				for _, container := range deploy.Spec.Template.Spec.Containers {
					t.Logf("container name: %v, image: %v\n", container.Name, container.Image)
				}
			},
		},
		{
			name: "patch tolerations",
			deploy: func() *appsv1.Deployment {
				deploy := resource.NewDeploymentSample()
				deploy.Spec.Template.Spec.Tolerations = append(deploy.Spec.Template.Spec.Tolerations, corev1.Toleration{
					Key:               "dedicated",
					Operator:          "",
					Value:             "test-team",
					Effect:            "NoSchedule",
					TolerationSeconds: nil,
				})
				return deploy
			}(),
			patchData: []byte(`{"spec":{"template":{"spec":{"tolerations":[{"effect":"NoSchedule","key":"disktype","value":"ssd"}]}}}}`),
			beforePatchFunc: func(deploy *appsv1.Deployment) {
				t.Logf("before patch: %+v\n", deploy.Spec.Template.Spec.Tolerations)
			},
			afterPatchFunc: func(deploy *appsv1.Deployment) {
				t.Logf("after patch: %+v\n", deploy.Spec.Template.Spec.Tolerations)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforePatchFunc != nil {
				tt.beforePatchFunc(tt.deploy)
			}
			patchDeploy := StrategicPatchDeploy(tt.deploy, tt.patchData)
			if tt.afterPatchFunc != nil {
				tt.afterPatchFunc(patchDeploy)
			}
		})
	}
}
