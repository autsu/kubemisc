package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"time"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// FILE_PATH=/Users/xxx/pj/GoProjects/kubemisc/ctrlruntime/source/testdata/test.txt
func main() {
	log.SetLogger(klogr.New())

	filepath := os.Getenv("FILE_PATH")
	if filepath == "" {
		panic("file path can't be nil")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		klog.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		klog.Error(err, "unable to set up manager")
		os.Exit(1)
	}

	mgr.GetScheme().AddKnownTypes(
		schema.GroupVersion{Group: "example.io", Version: "v1beta1"}, &FileWatch{})

	ctr := &Ctrl{
		fw: NewFileWatch(filepath, time.Second),
	}

	if err := (ctr).SetupWithManager(mgr); err != nil {
		klog.Error(err)
		os.Exit(1)
	}

	klog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
