package main

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
)

func NewPodSample() *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: metav1.NamespaceDefault,
			Labels: map[string]string{
				"app": "example",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
	return pod
}

func NewDeploymentSample() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "demo-deployment",
			Labels: make(map[string]string),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func main() {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	_, err = cli.CoreV1().Pods(metav1.NamespaceDefault).Get(context.TODO(), "busybox", metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get pod error: %v\n", err)
	}
	_, err = cli.CoreV1().Pods(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("list pod error: %v\n", err)
	}
	_, err = cli.AppsV1().Deployments(metav1.NamespaceDefault).Get(context.TODO(), "nginx", metav1.GetOptions{})
	if _, err := cli.CoreV1().Pods(metav1.NamespaceDefault).Create(context.TODO(), NewPodSample(), metav1.CreateOptions{}); err != nil {
		klog.Errorf("create pod error: %v\n", err)
	}
	_, err = cli.AppsV1().Deployments(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("create pod error: %v\n", err)
	}
	_, err = cli.AppsV1().Deployments(metav1.NamespaceDefault).Create(context.TODO(), NewDeploymentSample(), metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("create pod error: %v\n", err)
	}
}

// 运行后结果：
// E0308 17:02:12.164152       1 pod.go:108] create pod error: pods is forbidden: User "system:serviceaccount:default:test-20240308" cannot create resource "pods" in API group "" in the namespace "default"
//E0308 17:02:12.166747       1 pod.go:116] create pod error: deployments.apps is forbidden: User "system:serviceaccount:default:test-20240308" cannot create resource "deployments" in API group "apps" in the namespace "default"
// 可以看到 pod 的 create 操作和 deployment 的 create 操作都失败了，这符合我们的预期，因为 sa 绑定的 role 并没有提供这两种操作的权限
