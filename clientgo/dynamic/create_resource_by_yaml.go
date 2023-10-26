package main

import (
	"context"
	"fmt"
	"io"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func applyYaml(file string) {
	namespace := metav1.NamespaceDefault

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	dc := client.Discovery()
	restMapperRes, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		panic(err)
	}

	restMapper := restmapper.NewDiscoveryRESTMapper(restMapperRes)
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	d := yaml.NewYAMLOrJSONDecoder(f, 4096)

	// 为什么这里要用 for 循环
	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Println("raw: ", string(ext.Raw))

		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("gvk: %+v\n", gvk)

		mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			panic(err)
		}
		fmt.Printf("mapping: %+v\n", mapping)

		// runtime.Object 转换为 unstructed
		unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			panic(err)
		}
		fmt.Printf("unstructuredObj: %+v", unstructuredObj)

		var unstruct unstructured.Unstructured
		unstruct.Object = unstructuredObj

		if md, ok := unstruct.Object["metadata"]; ok {
			metadata := md.(map[string]interface{})
			if internalns, ok := metadata["namespace"]; ok {
				namespace = internalns.(string)
			}
		}

		// 动态客户端
		dclient, err := dynamic.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		res, err := dclient.Resource(mapping.Resource).Namespace(namespace).Create(context.TODO(), &unstruct, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
}

func main() {
	applyYaml("")
}
