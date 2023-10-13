package main

import (
	"fmt"

	"void.io/kubemisc/clientgo/helper"

	"k8s.io/apimachinery/pkg/api/meta"
)

// 具体来说，MetadataAccessor 是 Kubernetes 中的一个 Go 语言接口，通常用于访问或操作 Kubernetes API 资源对象
// 的元数据（Metadata），如名称、命名空间、标签、注释等。这个接口定义了一些用于处理 Kubernetes 资源对象元数据的方法，
// 可以帮助开发者编写与资源对象元数据相关的操作代码。
func main() {
	accessor := meta.NewAccessor()
	name, err := accessor.Name(helper.NewPodSimple())
	if err != nil {
		panic(err)
	}
	namespace, err := accessor.Namespace(helper.NewPodSimple())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\\%v\n", namespace, name)
}
