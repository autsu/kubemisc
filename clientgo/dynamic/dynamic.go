package main

import (
	"encoding/json"
	"fmt"

	"void.io/kubemisc/clientgo/helper/resource"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var objs = []any{resource.NewPodSample(), resource.NewDeploymentSample()}

func main() {
	var (
		err  error
		exts []runtime.RawExtension
	)

	for _, obj := range objs {
		ext := runtime.RawExtension{}
		ext.Raw, err = json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		exts = append(exts, ext)
	}

	codec := unstructured.UnstructuredJSONScheme
	for _, ext := range exts {
		obj, gvk, err := codec.Decode(ext.Raw, nil, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(gvk)
		accessor := meta.NewAccessor()
		name, _ := accessor.Name(obj)
		namespace, _ := accessor.Namespace(obj)
		fmt.Printf("%v\\%v\n", namespace, name)
	}
}
