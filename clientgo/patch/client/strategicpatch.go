package client

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"void.io/kubemisc/clientgo/helper"
	"void.io/kubemisc/clientgo/helper/printhelper"
)

// StrategicMergePatch2K8sUseJson 根据 JSON 进行 Patch
func StrategicMergePatch2K8sUseJson() {
	deploy, dropFunc := initTestResource(true)
	if dropFunc != nil {
		defer dropFunc()
	}

	fmt.Println("before patch: ")
	printhelper.ObjJSON(deploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(deploy.Spec.Template.Spec.Tolerations)

	cli := helper.NewClientSetOrDie()
	patchData := []byte(`
{
  "metadata": {
    "labels": {
      "new-label": "new-value",
      "patch": "true"
    }
  },
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "patch-demo-ctr-2",
            "image": "redis"
          }
        ],
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "disktype",
            "value": "ssd"
          }
        ]
      }
    }
  }
}
`)
	patchDeploy, err := cli.AppsV1().Deployments(metav1.NamespaceDefault).Patch(
		context.TODO(), deploy.Name, types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	//fmt.Println(patchDeploy.Labels)
	fmt.Println("after patch: ")
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Tolerations)
}

// StrategicMergePatch2K8sUseCreateTwoWayMergePatch 调用 strategicpatch.CreateTwoWayMergePatch 获取两个 obj 之间的 patch，
// 然后再用这个 patch 对对象执行 Patch 操作
func StrategicMergePatch2K8sUseCreateTwoWayMergePatch() {
	deploy, dropFunc := initTestResource(true)
	if dropFunc != nil {
		defer dropFunc()
	}
	fmt.Printf("before patch: \n\n")
	printhelper.ObjJSON(deploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(deploy.Spec.Template.Spec.Tolerations)

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

	origin, err := json.Marshal(deploy)
	if err != nil {
		panic(err)
	}

	mod, err := json.Marshal(newDeploy)
	if err != nil {
		panic(err)
	}

	patchData, err := strategicpatch.CreateTwoWayMergePatch(origin, mod, &appsv1.Deployment{})
	if err != nil {
		panic(err)
	}
	fmt.Println("patch data: ", string(patchData))

	cli := helper.NewClientSetOrDie()

	patchDeploy, err := cli.AppsV1().Deployments(metav1.NamespaceDefault).Patch(
		context.TODO(), deploy.Name, types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	//fmt.Println(patchDeploy.Labels)
	fmt.Println("after patch: ")
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Tolerations)
}
