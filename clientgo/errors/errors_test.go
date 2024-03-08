package errors

import (
	stderrors "errors"
	"testing"

	"k8s.io/apimachinery/pkg/util/errors"
)

func TestAggregateError(t *testing.T) {
	e1 := stderrors.New("error1")
	e2 := stderrors.New("error2")
	e3 := stderrors.New("error3")
	e4 := stderrors.New("error4")

	aggregate := errors.NewAggregate([]error{e1, e2, e3})

	t.Log(aggregate.Error())
	t.Log(aggregate.Is(e1))
	t.Log(aggregate.Is(e4))
}
