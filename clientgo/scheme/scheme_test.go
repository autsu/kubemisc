package scheme

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ runtime.Object = &Object{}

type Object struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func (o *Object) GetObjectKind() schema.ObjectKind {
	return o
}

func (o *Object) DeepCopyObject() runtime.Object {
	return &Object{}
}

func TestScheme(t *testing.T) {
	runtime.NewSchemeBuilder()
	//meta.RESTMapper()

	internalGV := schema.GroupVersion{Group: "test.group", Version: runtime.APIVersionInternal}
	internalGVK := internalGV.WithKind("Simple")
	externalGV := schema.GroupVersion{Group: "test.group", Version: "testExternal"}
	externalGVK := externalGV.WithKind("Simple")

	scheme_ := runtime.NewScheme()
	scheme_.AddKnownTypeWithName(internalGVK, &Object{})
	scheme_.AddKnownTypeWithName(externalGVK, &Object{})
	//utilruntime.Must(runtimetesting.RegisterConversions(scheme_))

	kinds, _, err := scheme_.ObjectKinds(&Object{})
	if err != nil {
		panic(err)
	}
	// [test.group/__internal, Kind=Simple test.group/testExternal, Kind=Simple]
	t.Log(kinds)

	mapper := meta.NewDefaultRESTMapper(func(gvk []schema.GroupVersionKind) (ret []schema.GroupVersion) {
		for _, v := range gvk {
			ret = append(ret, v.GroupVersion())
		}
		return
	}(kinds))

	mapping, err := mapper.RESTMapping(schema.GroupKind{Group: kinds[0].Group, Kind: kinds[0].Kind}, kinds[0].Version)
	if err != nil {
		panic(err)
	}
	resource := mapping.Resource.Resource
	t.Log(resource)
}

func TestName(t *testing.T) {
	groups := runtime.NewScheme().PreferredVersionAllGroups()
	for _, g := range groups {
		t.Log(g.String())
	}
}
