package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

// https://github.com/kubernetes/client-go/tree/master/examples/leader-election
//
// # first terminal
// go run main.go -kubeconfig=/path/to/kubeconfig -logtostderr=true -lease-lock-name=example -lease-lock-namespace=default -id=1
//
// # second terminal
// go run main.go -kubeconfig=/path/to/kubeconfig -logtostderr=true -lease-lock-name=example -lease-lock-namespace=default -id=2
//
// # third terminal
// go run main.go -kubeconfig=/path/to/kubeconfig -logtostderr=true -lease-lock-name=example -lease-lock-namespace=default -id=3

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	klog.InitFlags(nil)

	var kubeconfig string
	var leaseLockName string
	var leaseLockNamespace string
	var id string

	flag.StringVar(&kubeconfig, "kubeconfig", clientcmd.RecommendedHomeFile, "absolute path to the kubeconfig file")
	flag.StringVar(&id, "id", uuid.New().String(), "the holder identity name")
	flag.StringVar(&leaseLockName, "lease-lock-name", "", "the lease lock resource name")
	flag.StringVar(&leaseLockNamespace, "lease-lock-namespace", "", "the lease lock resource namespace")
	flag.Parse()

	if leaseLockName == "" {
		klog.Fatal("unable to get lease lock resource name (missing lease-lock-name flag).")
	}
	if leaseLockNamespace == "" {
		klog.Fatal("unable to get lease lock resource namespace (missing lease-lock-namespace flag).")
	}

	// leader election uses the Kubernetes API by writing to a
	// lock object, which can be a LeaseLock object (preferred),
	// a ConfigMap, or an Endpoints (deprecated) object.
	// Conflicting writes are detected and each client handles those actions
	// independently.
	config, err := buildConfig(kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
	client := clientset.NewForConfigOrDie(config)

	run := func(ctx context.Context) {
		// complete your controller loop here
		klog.Info("Controller loop...")

		select {}
	}

	// use a Go context so we can tell the leaderelection code when we
	// want to step down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		cancel()
	}()

	// we use the Lease lock type since edits to Leases are less common
	// and fewer objects in the cluster watch "all Leases".
	// 用 Lease 这个 k8s 资源来做选主
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseLockName,
			Namespace: leaseLockNamespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				run(ctx)
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				klog.Infof("leader lost: %s", id)
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				// we're notified when new leader elected
				if identity == id {
					// I just got the lock
					return
				}
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})
}
