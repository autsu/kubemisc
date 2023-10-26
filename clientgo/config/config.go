package config

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// TODO: config.APIPath 的作用

func getKubeConfig1() *rest.Config {
	defaultVal := filepath.Join(os.Getenv("HOME"), ".kube/config")
	kubeconfig := flag.String("kubeconfig", defaultVal, "kubeconfig file path")
	flag.Parse()

	// 从 master url 或者 kubeconfig 中获取集群配置
	// 下面这里是从用户指定的 -kubeconfig 参数中获取 kubeconfig 所在的位置
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// 如果想要指定为 kubeconfig 的默认位置（~/.kube/kubeconfig），那么可以直接用提供好的函数 clientcmd.RecommendedHomeFile
	// config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	return config
}

func getKubeConfig2() map[string]*rest.Config {
	// 获取 kubeconfig 文件
	cfg, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	configs := make(map[string]*rest.Config)

	// 获取 kubeconfig 中的所有 context 信息（kubeconfig 里面可以保存多个集群的信息）
	for ctx := range cfg.Contexts {
		restCfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*api.Config, error) {
			// 将 ctx 设置为 CurrentContext（对应到 kubeconfig 是 current-context 字段），
			cfg.CurrentContext = ctx
			// return 的是一个深拷贝对象
			return cfg.DeepCopy(), nil
		})
		if err != nil {
			panic(err)
		}
		configs[ctx] = restCfg
	}
	return configs
}

func testConfig(ctx context.Context, config *rest.Config) *wrapErr {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return newWrapErr(newClientSetError, err)
	}

	podList, err := clientset.CoreV1().Pods(metav1.NamespaceDefault).List(ctx, metav1.ListOptions{})
	if err != nil {
		//panic(err)
		return newWrapErr(listPodsError, err)
	}

	fmt.Println(func() []string {
		var s []string
		for _, item := range podList.Items {
			s = append(s, item.Name)
		}
		return s
	}())

	return nil
}

func NewRestConfig(server string, caData []byte, certData []byte, keyData []byte, timeout string) (*rest.Config, error) {
	config := clientcmdapi.Config{
		Preferences: *clientcmdapi.NewPreferences(),
		Clusters: map[string]*clientcmdapi.Cluster{
			"kubernetes": &clientcmdapi.Cluster{
				Server:                   server,
				CertificateAuthorityData: caData,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"kubernetes": &clientcmdapi.AuthInfo{
				ClientCertificateData: certData,
				ClientKeyData:         keyData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"kubernetes": &clientcmdapi.Context{
				Cluster:  "kubernetes",
				AuthInfo: "kubernetes",
			},
		},
		CurrentContext: "kubernetes",
	}
	return clientcmd.NewNonInteractiveClientConfig(config, config.CurrentContext, &clientcmd.ConfigOverrides{
		Timeout: timeout,
	}, nil).ClientConfig()
}
