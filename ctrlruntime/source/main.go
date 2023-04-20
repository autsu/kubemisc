package main

import (
	"os"

	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
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
	if err := (&Ctrl{
		fw: NewFileWatch(filepath),
	}).SetupWithManager(mgr); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	klog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
