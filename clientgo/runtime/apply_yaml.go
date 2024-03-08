package main

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// TODO：未完成
func main() {
	f, err := os.Open("../testdata/nginx-deploy.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewYAMLOrJSONDecoder(f, 4096)
	ext := runtime.RawExtension{}
	if err := decoder.Decode(&ext); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", ext)

	obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(gvk)

	switch obj.(type) {
	case *unstructured.Unstructured:
		fmt.Println("obj type is unstructured")
	}
}
