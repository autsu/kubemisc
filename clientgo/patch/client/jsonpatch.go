package client

import (
	"context"
	"fmt"

	"void.io/kubemisc/clientgo/helper"
	"void.io/kubemisc/clientgo/helper/printhelper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func JsonPatch2K8s() {
	deploy, dropFunc := initTestResource(true)
	if dropFunc != nil {
		defer dropFunc()
	}

	fmt.Println("before patch: ")
	printhelper.ObjJSON(deploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(deploy.Spec.Template.Spec.Tolerations)

	cli := helper.NewClientSetOrDie()
	patchData := []byte(`
[
  {
    "op": "add",
    "path": "/metadata/labels",
    "value": {
		"new-labels": "new-value",
		"patch": "true"
	}
  },
  {
    "op": "add",
    "path": "/spec/template/spec/containers/1",
    "value": 
      {
        "name": "patch-demo-ctr-2",
        "image": "redis"
      }
  },
  {
    "op": "add",
    "path": "/spec/template/spec/tolerations",
    "value": [
      {
        "effect": "NoSchedule",
        "key": "disktype",
        "value": "ssd"
      }
    ]
  }
]

`)
	patchDeploy, err := cli.AppsV1().Deployments(
		metav1.NamespaceDefault).Patch(context.TODO(), deploy.Name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("after patch: ")
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Containers)
	printhelper.ObjJSON(patchDeploy.Spec.Template.Spec.Tolerations)
}
