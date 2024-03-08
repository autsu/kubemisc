package printhelper

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func ResourceListItemName[T runtime.Object](list []T) {
	for _, r := range list {
		fmt.Println(meta.NewAccessor().Name(r))
	}
}

func ObjJSON(obj any) {
	if obj == nil {
		return
	}
	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
