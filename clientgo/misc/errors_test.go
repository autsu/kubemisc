package main

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
)

// errors.IsNotFound 传 nil 会不会 panic
func TestNotFoundErrParamNil(t *testing.T) {
	t.Log(errors.IsNotFound(nil))
}
