package config

type errKind int

const (
	newClientSetError errKind = iota
	listPodsError
)

//var _ error = &wrapErr{}

type wrapErr struct {
	kind errKind
	err  error
}

func newWrapErr(kind errKind, err error) *wrapErr {
	return &wrapErr{kind: kind, err: err}
}

//func (w *wrapErr) Error() string {
//	if w == nil || w.err == nil {
//		return ""
//	}
//	return w.err.Error()
//}
