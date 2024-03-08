package main

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"void.io/kubemisc/clientgo/helper/resource"
)

func objToUnstructured() {
	obj := resource.NewPodSample()
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	ut := unstructured.Unstructured{}
	if err := json.Unmarshal(b, &ut.Object); err != nil {
		panic(err)
	}

	fmt.Println(ut)
}
