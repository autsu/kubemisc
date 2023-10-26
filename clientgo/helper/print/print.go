package print

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func ResourceListItemName[T runtime.Object](list []T) {
	for _, r := range list {
		fmt.Println(meta.NewAccessor().Name(r))
	}
}
