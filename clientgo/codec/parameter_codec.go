package codec

import (
	"fmt"
	"net/url"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"void.io/kubemisc/clientgo/codec/internal"
)

var _ runtime.Object = &ListOptions{}

// ListOptions 要想成功调用 ParameterCodec.DecodeParameters，必须要将 Object 注册到 scheme，并且还要注册对应的 url.Values 到
// Object 的转换函数
type ListOptions struct {
	metav1.TypeMeta

	LabelSelector string
}

func (l *ListOptions) DeepCopyObject() runtime.Object {
	return &ListOptions{LabelSelector: l.LabelSelector}
}

// 坑：记得重写 String 方法，不然默认会调用 TypeMeta 的，导致 fmt.Print 输出异常
func (l *ListOptions) String() string {
	return `LabelSelector: ` + l.LabelSelector + `,`
}

type internalListOptions struct {
	metav1.TypeMeta

	Test string
}

func (i *internalListOptions) DeepCopyObject() runtime.Object {
	return &internalListOptions{Test: i.Test}
}

func (i *internalListOptions) String() string {
	return `Test: ` + i.Test + `,`
}

var (
	SchemeGroupVersion         = schema.GroupVersion{Group: "void.io", Version: "v1beta1"}
	SchemeGroupVersionInternal = schema.GroupVersion{Group: "void.io", Version: runtime.APIVersionInternal}
	scheme                     = runtime.NewScheme()
	Codecs                     = serializer.NewCodecFactory(scheme)
	ParameterCodec             = runtime.NewParameterCodec(scheme)
)

var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	Install            = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(func(s *runtime.Scheme) error {
		// 将 Object 注册进 scheme
		scheme.AddKnownTypes(SchemeGroupVersion, &ListOptions{})
		scheme.AddKnownTypes(SchemeGroupVersionInternal, &internalListOptions{})
		scheme.AddKnownTypes(SchemeGroupVersionInternal, &internal.ListOptions{})

		// 注册 url.Values -> ListOptions 的转换函数，必须要做这个操作，不然 ParameterCodec.DecodeParameters 会报错
		// 前两个参数必须是指针类型
		if err := scheme.AddConversionFunc((*url.Values)(nil), (*ListOptions)(nil), func(in, out interface{}, scope conversion.Scope) error {
			inVal := in.(*url.Values)
			outVal := out.(*ListOptions)
			if values, ok := (*inVal)["labelSelector"]; ok {
				if err := runtime.Convert_Slice_string_To_string(&values, &outVal.LabelSelector, scope); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		if err := scheme.AddConversionFunc((*ListOptions)(nil), (*internalListOptions)(nil), func(in, out interface{}, scope conversion.Scope) error {
			inVal := in.(*ListOptions)
			outVal := out.(*internalListOptions)
			outVal.Test = inVal.LabelSelector
			return nil
		}); err != nil {
			return err
		}
		if err := scheme.AddConversionFunc((*ListOptions)(nil), (*internal.ListOptions)(nil), func(in, out interface{}, scope conversion.Scope) error {
			inVal := in.(*ListOptions)
			outVal := out.(*internal.ListOptions)
			outVal.Test = inVal.LabelSelector
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	utilruntime.Must(Install(scheme))
}

func decodeParameters(u string, t *testing.T) {
	//into := &ListOptions{}
	// 我定义了 url.Values -> ListOptions 的转换函数，以及 ListOptions -> internal.ListOptions
	// 的转换函数，但是没有直接定义 url.Values -> internal.ListOptions 的转换函数
	// 那么 DecodeParameters 还能转换成功吗？
	// ~~验证了，不行~~
	// 是可以的，但有以下要求：
	// - 两个 struct 的名字相同，这样他们的 Kind 才会默认相同
	// - 都注册到了同一个 scheme，且注册的 gv 不同。在这里，ListOptions 的 gv 是 void.io/v1beta1，
	//   而 internal.ListOptions 的 gv 是 void.io/__internal
	into := &internal.ListOptions{}

	parse, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}
	query := parse.Query()
	if err := ParameterCodec.DecodeParameters(query, SchemeGroupVersion, into); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", into)

	//requirements, _ := into.LabelSelector.Requirements()
	//for _, requirement := range requirements {
	//	fmt.Println(requirement)
	//}
}
