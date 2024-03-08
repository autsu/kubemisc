package main

//import (
//	"context"
//	"testing"
//
//	corev1 "k8s.io/api/core/v1"
//	rbacv1 "k8s.io/api/rbac/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
//
//	"void.io/kubemisc/clientgo/helper"
//)
//
//func createRBAC(cli *kubernetes.Clientset) {
//	commonName := "try"
//	commonNamespace := "kube-test"
//
//	sa := &corev1.ServiceAccount{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "ServiceAccount",
//			APIVersion: "v1",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      commonName,
//			Namespace: commonNamespace,
//		},
//	}
//
//	role := &rbacv1.Role{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "Role",
//			APIVersion: "rbac.authorization.k8s.io/v1",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      commonName,
//			Namespace: commonNamespace,
//		},
//		Rules: []rbacv1.PolicyRule{
//			{
//				Verbs:     []string{"get", "list", "watch"},
//				APIGroups: []string{"", "apps"},
//				Resources: []string{"pods", "deployments"},
//			},
//		},
//	}
//
//	roleBinding := &rbacv1.RoleBinding{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "RoleBinding",
//			APIVersion: "rbac.authorization.k8s.io/v1",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      commonName,
//			Namespace: commonNamespace,
//		},
//		Subjects: []rbacv1.Subject{
//			{
//				Kind:      "ServiceAccount",
//				Name:      commonName,
//				Namespace: commonNamespace,
//			},
//		},
//		RoleRef: rbacv1.RoleRef{
//			APIGroup: "rbac.authorization.k8s.io",
//			Kind:     "Role",
//			Name:     commonName,
//		},
//	}
//
//	if _, err := cli.CoreV1().ServiceAccounts(commonNamespace).Create(context.TODO(), sa, metav1.CreateOptions{}); err != nil {
//		panic(err)
//	}
//	if _, err := cli.RbacV1().Roles(commonNamespace).Create(context.TODO(), role, metav1.CreateOptions{}); err != nil {
//		panic(err)
//	}
//	if _, err := cli.RbacV1().RoleBindings(commonNamespace).Create(context.TODO(), roleBinding, metav1.CreateOptions{}); err != nil {
//		panic(err)
//	}
//}
//
//func TestServiceAccount(t *testing.T) {
//	cli := helper.NewClientSetOrDie()
//	createRBAC(cli)
//
//
//}
