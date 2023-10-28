package patch

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func JsonPatch() {
	original := []byte(`
{
	"title": "Goodbye!",
	"author": {
		"givenName": "John",
		"familyName": "Doe"
	},
	"tags": [
		"example",
		"sample"
	],
	"content": "This will be unchanged"
}
`)

	// 貌似 add 的策略是覆盖（同 path），比如下面的例子中，后一个 add 会覆盖掉前一个 add 的值，最终
	// phoneNumber 的值将会是 +01-123-456-7891
	patchJSON := []byte(`[
		{ "op": "replace", "path": "/title", "value": "Hello!"},
  		{ "op": "remove", "path": "/author/familyName"},
  		{ "op": "add", "path": "/phoneNumber", "value": "+01-123-456-7890"},
  		{ "op": "add", "path": "/phoneNumber", "value": "+01-123-456-7891"},
  		{ "op": "replace", "path": "/tags", "value": ["example"]}
	]`)
	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original document: %s\n", original)
	fmt.Printf("Modified document: %s\n", modified)
}
