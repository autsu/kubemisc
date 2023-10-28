package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func NewGenericRESTClientWithGroupVersion(cfg *rest.Config, gv schema.GroupVersion) (*rest.RESTClient, error) {
	cfg.GroupVersion = &gv
	if len(gv.Group) == 0 {
		cfg.APIPath = "/api"
	} else {
		cfg.APIPath = "/apis"
	}
	return rest.RESTClientFor(cfg)
}
