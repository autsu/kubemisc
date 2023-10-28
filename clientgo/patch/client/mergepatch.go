package client

import (
	"context"
	"encoding/json"
	"fmt"

	"void.io/kubemisc/clientgo/helper"
	"void.io/kubemisc/clientgo/helper/printhelper"

	jsonpatch "github.com/evanphx/json-patch/v5"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func MergePatch2K8sUseJSON() {
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
            "new-labels": "new-value",
            "patch": "true"
        }
    },
    "spec": {
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "a": "b"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "nginx:1.12",
                        "imagePullPolicy": "IfNotPresent",
                        "name": "web",
                        "ports": [
                            {
                                "containerPort": 80,
                                "name": "http",
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    },
                    {
                        "image": "redis",
                        "imagePullPolicy": "Always",
                        "name": "patch-demo-ctr-2",
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
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
	patchDeploy, err := cli.AppsV1().Deployments(
		metav1.NamespaceDefault).Patch(context.TODO(), deploy.Name, types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("after patch: ")
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Tolerations)
}

func MergePatch2K8sUseCreateMergePatch() {
	deploy, dropFunc := initTestResource(true)
	if dropFunc != nil {
		defer dropFunc()
	}

	fmt.Println("before patch: ")
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

	cli := helper.NewClientSetOrDie()

	// patchData 根据 origin 和 mod 生成 patch
	patchData, err := jsonpatch.CreateMergePatch(origin, mod)
	if err != nil {
		panic(err)
	}
	fmt.Println("patch data: ", string(patchData))

	patchDeploy, err := cli.AppsV1().Deployments(
		metav1.NamespaceDefault).Patch(context.TODO(), deploy.Name, types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("after patch: ")
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Tolerations)
}
