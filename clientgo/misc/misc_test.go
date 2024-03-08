package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
	"void.io/kubemisc/clientgo/helper/printhelper"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func TestName(t *testing.T) {
	testLabels := map[string]string{"key": "value"}
	selector := labels.SelectorFromSet(labels.Set(testLabels))

	requirement, err := labels.NewRequirement("k", selection.In, []string{"1", "2"})
	if err != nil {
		t.Fatal(err)
	}
	selector = selector.Add(*requirement)
	into := &metainternal.ListOptions{LabelSelector: selector}
	requirements, _ := into.LabelSelector.Requirements()
	for _, r := range requirements {
		t.Log(r.Operator())
	}
}

func TestGetNode(t *testing.T) {
	initGlobalClientSet()
	nodeList, err := globalCliSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range nodeList.Items {
		printhelper.ObjJSON(node.Status)
		//fmt.Printf("%+v\n", node.Status)
	}
}

func TestDep(t *testing.T) {
	initGlobalClientSet()
	depList, err := globalCliSet.AppsV1().Deployments(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range depList.Items {
		printhelper.ObjJSON(item.Spec.Strategy.RollingUpdate)
		//fmt.Printf("%+v\n", node.Status)
	}

	dep, err := globalCliSet.
		AppsV1().
		Deployments(metav1.NamespaceDefault).
		Get(context.TODO(), "nginx-deployment", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dep.Spec.Strategy.RollingUpdate)

	newVal := intstr.FromInt32(10)
	//dep.Spec.Strategy.RollingUpdate.MaxSurge.IntVal = 10
	//dep.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal = 10
	dep.Spec.Strategy.RollingUpdate.MaxSurge = &newVal
	dep.Spec.Strategy.RollingUpdate.MaxUnavailable = &newVal

	_, err = globalCliSet.
		AppsV1().
		Deployments(metav1.NamespaceDefault).
		Update(context.TODO(), dep, metav1.UpdateOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func changeRollingUpdate(dep *appsv1.Deployment) {
	newVal := intstr.FromString("25%")
	dep.Spec.Strategy.RollingUpdate.MaxSurge = &newVal
	dep.Spec.Strategy.RollingUpdate.MaxUnavailable = &newVal
}

func TestChangeRollingUpdate(t *testing.T) {
	initGlobalClientSet()
	dep, err := globalCliSet.
		AppsV1().
		Deployments(metav1.NamespaceDefault).
		Get(context.TODO(), "nginx-deployment", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dep.Spec.Strategy.RollingUpdate)
	changeRollingUpdate(dep)
	fmt.Println(dep.Spec.Strategy.RollingUpdate)
}
