package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"

	"void.io/kubemisc/clientgo/helper/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	once   sync.Once
	deploy = func() *appsv1.Deployment {
		deploy := resource.NewDeploymentSample()
		deploy.Name = testDeployName
		deploy.Namespace = testNamespace
		return deploy
	}()
	svc = func() *corev1.Service {
		svc := resource.NewServiceForDeployment(deploy)
		svc.Name = testServiceName
		svc.Namespace = testNamespace
		return svc
	}()
)

func cleanup() {
	cli.AppsV1().
		Deployments(testNamespace).
		Delete(context.TODO(), testDeployName, metav1.DeleteOptions{})
	cli.CoreV1().
		Services(testNamespace).
		Delete(context.TODO(), testServiceName, metav1.DeleteOptions{})
}

func initTestResource() error {
	_, err := cli.AppsV1().
		Deployments(testNamespace).
		Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = cli.CoreV1().
		Services(testNamespace).
		Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func Init() {
	once.Do(func() {
		initClientSet()
		if err := initTestResource(); err != nil {
			panic(err)
		}
	})
}

func TestFindDeployServiceByInformer(t *testing.T) {
	//Init()
	//t.Cleanup(cleanup)
	serviceLister := initInformer()

	service, err := __findDeployServiceByInformerWithCache__(context.TODO(), serviceLister, deploy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
	//findDeployService1()
}

func TestFindDeployServiceByClientSet(t *testing.T) {
	initClientSet()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil,
	})))

	ctx := context.TODO()
	searchKey := []*appsv1.Deployment{
		deploy,
		deploy,
		deploy,
	}
	for i, deploy := range searchKey {
		service, err := __findDeployServiceByClientSetWithCache__(ctx, cli, deploy)
		if err != nil {
			t.Fatal(err)
		}
		//t.Log(service)
		for _, svc := range service {
			key, _ := cache.MetaNamespaceKeyFunc(svc)
			t.Logf("%v %v\n", i, key)
		}
	}
}

func TestInitForFindDeployService(t *testing.T) {
	Init()
}

func TestFindDeployServiceErr2(t *testing.T) {
	//Init()
	//t.Cleanup(cleanup)
	initClientSet()
	services, err := __findDeployService2Err__(cli, testNamespace, deploy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(func() (ret []string) {
		for _, service := range services {
			ret = append(ret, service.Name)
		}
		return
	}())
	//findDeployService1()
}
