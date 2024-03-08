package resource

import "k8s.io/apimachinery/pkg/runtime/schema"

var GVR = &GVRGetter{}

type GVRGetter struct{}

func (g *GVRGetter) Pod() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
}

func (g *GVRGetter) Deployment() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
}

func (g *GVRGetter) Service() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}
}

func (g *GVRGetter) Namespace() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}
}

func (g *GVRGetter) ConfigMap() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
}

func (g *GVRGetter) Secret() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}
}

func (g *GVRGetter) StatefulSet() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "statefulsets",
	}
}

//var (
//	GvrPods = schema.GroupVersionResource{
//		Group:    "",
//		Version:  "v1",
//		Resource: "pods",
//	}
//
//	GvrDeployments = schema.GroupVersionResource{
//		Group:    "apps",
//		Version:  "v1",
//		Resource: "deployments",
//	}
//
//	GvrServices = schema.GroupVersionResource{
//		Group:    "",
//		Version:  "v1",
//		Resource: "services",
//	}
//
//	GvrNamespaces = schema.GroupVersionResource{
//		Group:    "",
//		Version:  "v1",
//		Resource: "namespaces",
//	}
//
//	GvrConfigMaps = schema.GroupVersionResource{
//		Group:    "",
//		Version:  "v1",
//		Resource: "configmaps",
//	}
//
//	GvrSecrets = schema.GroupVersionResource{
//		Group:    "",
//		Version:  "v1",
//		Resource: "secrets",
//	}
//
//	GvrStatefulSets = schema.GroupVersionResource{
//		Group:    "apps",
//		Version:  "v1",
//		Resource: "statefulsets",
//	}
//)
