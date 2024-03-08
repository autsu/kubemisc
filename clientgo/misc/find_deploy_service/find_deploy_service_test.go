package find_deploy_service

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"void.io/kubemisc/clientgo/helper/resource"
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
	cli = func() *kubernetes.Clientset {
		cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			panic(err)
		}
		c, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			panic(err)
		}
		return c
	}()
)

const (
	testNamespace   = "test-1018"
	testDeployName  = "test-1018"
	testServiceName = "test-1018"
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
		if err := initTestResource(); err != nil {
			panic(err)
		}
	})
}

func initInformer() listerscorev1.ServiceLister {
	factory := informers.NewSharedInformerFactory(cli, time.Second*30)
	serviceInformer := factory.Core().V1().Services()
	stopCh := wait.NeverStop
	factory.Start(stopCh)
	//if !cache.WaitForCacheSync(stopCh, serviceInformer.Informer().HasSynced) {
	//	runtime.HandleError(errors.New("failed to sync"))
	//	return nil
	//}
	return serviceInformer.Lister()
}

func TestFindDeployServiceByInformer(t *testing.T) {
	//Init()
	//t.Cleanup(cleanup)
	serviceLister := initInformer()

	service, err := __findDeployServiceByInformerWithCache__(cli, context.TODO(), serviceLister, deploy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
	//findDeployService1()
}

func TestFindDeployServiceByClientSet(t *testing.T) {
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
