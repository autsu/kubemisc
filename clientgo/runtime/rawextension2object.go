package main

import (
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"void.io/kubemisc/clientgo/helper/resource"
)

// runtime.RawExtension 转换为具体的资源对象
func RawExtension2ObjectSample1() {
	ext := runtime.RawExtension{
		Raw:    []byte(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"deployment.kubernetes.io/revision":"1","shadow.clusterpedia.io/cluster-name":"k3s-1"},"creationTimestamp":"2023-11-16T02:51:22Z","generation":1,"labels":{"app":"nginx"},"name":"nginx-deployment","namespace":"default","resourceVersion":"1790","uid":"83397df9-603c-4041-babf-34d2453908c1"}}`),
		Object: nil,
	}
	deploy := &appsv1.Deployment{}
	_, _, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, &schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}, deploy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", deploy)
}

func RawExtension2ObjectSample2() {
	type WrapExt struct {
		ext runtime.RawExtension
		// gvk 和 obj 方便我们之后的 Decode 操作
		gvk schema.GroupVersionKind
		obj runtime.Object
	}

	var exts []*WrapExt
	var objs = []runtime.Object{resource.NewPodSample(), resource.NewDeploymentSample()}

	// 先把这两个资源序列化为 json 放到 runtime.RawExtension.Raw 中，同时填充 gvk 和 obj 信息
	for _, obj := range objs {
		data, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}

		var obj1 runtime.Object
		switch obj.GetObjectKind().GroupVersionKind().Kind {
		case "Deployment":
			obj1 = &appsv1.Deployment{}
		case "Pod":
			obj1 = &corev1.Pod{}
		}

		ext := &WrapExt{
			ext: runtime.RawExtension{Raw: data},
			gvk: obj.GetObjectKind().GroupVersionKind(),
			obj: obj1,
		}
		exts = append(exts, ext)
	}

	for _, ext := range exts {
		// 通过 gvk 和 obj 反序列出具体的对象
		_, _, err := unstructured.UnstructuredJSONScheme.Decode(ext.ext.Raw, &ext.gvk, ext.obj)
		if err != nil {
			panic(err)
		}
		switch ext.obj.GetObjectKind().GroupVersionKind().Kind {
		case "Deployment":
			deploy := ext.obj.(*appsv1.Deployment)
			fmt.Printf("%+v\n", deploy)
		case "Pod":
			po := ext.obj.(*corev1.Pod)
			fmt.Printf("%+v\n", po)
		}
	}

}
