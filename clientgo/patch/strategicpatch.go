package patch

import (
	"encoding/json"
	"fmt"

	"void.io/kubemisc/clientgo/helper/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func StrategicPatchDeploy(deploy *appsv1.Deployment, patchData []byte) *appsv1.Deployment {
	originalDeployJson, err := json.Marshal(deploy)
	if err != nil {
		panic(err)
	}
	patchDeploy := new(appsv1.Deployment)

	patchDeployByte, err := strategicpatch.StrategicMergePatch(originalDeployJson, patchData, &appsv1.Deployment{})
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(patchDeployByte, patchDeploy); err != nil {
		panic(err)
	}
	return patchDeploy
}

func CreateTwoWayMergePatch() {
	deploy := resource.NewDeploymentBy("test-1028", metav1.NamespaceDefault, map[string]string{"a": "b"})
	newDeploy := deploy.DeepCopy()
	newDeploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, corev1.Container{
		Name:  "patch-demo-ctr-2",
		Image: "redis",
	})
	newDeploy.Spec.Template.Spec.Tolerations = newDeploy.Spec.Template.Spec.Tolerations[:0]
	newDeploy.Spec.Template.Spec.Tolerations = append(newDeploy.Spec.Template.Spec.Tolerations, corev1.Toleration{
		Key:               "disktype",
		Operator:          "",
		Value:             "ssd",
		Effect:            "NoSchedule",
		TolerationSeconds: nil,
	})
	j, err := json.Marshal(deploy)
	if err != nil {
		panic(err)
	}
	j1, err := json.Marshal(newDeploy)
	if err != nil {
		panic(err)
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(j, j1, &appsv1.Deployment{})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(patch))
}
