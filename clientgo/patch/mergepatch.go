package patch

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func JsonMergePatch() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	target := []byte(`{"name": "Jane", "age": 24}`)

	// 这里的 patch 是 original 和 target 之间差异的集合
	// 比如对于上面的内容，返回的 patch 将会是 {"height":null,"name":"Jane"}
	// height 为 null 代表要删除这个字段，name 为 Jane 代表我们要将 name 更新为 Jane
	patch, err := jsonpatch.CreateMergePatch(original, target)
	if err != nil {
		panic(err)
	}

	// Output: {"height":null,"name":"Jane"}
	fmt.Printf("patch document:   %s\n", patch)

	alternative := []byte(`{"name": "Tina", "age": 28, "height": 3.75}`)
	modifiedAlternative, err := jsonpatch.MergePatch(alternative, patch)
	fmt.Printf("updated alternative doc: %s\n", modifiedAlternative)
}
