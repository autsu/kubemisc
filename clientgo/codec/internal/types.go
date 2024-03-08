package internal

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ListOptions struct {
	metav1.TypeMeta

	Test string
}

func (i *ListOptions) DeepCopyObject() runtime.Object {
	return &ListOptions{Test: i.Test}
}

func (i *ListOptions) String() string {
	return `Test: ` + i.Test + `,`
}
