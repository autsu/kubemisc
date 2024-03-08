package sets

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
)

func printSet[T comparable](set sets.Set[T], desc string) {
	fmt.Print(desc)
	for v := range set {
		fmt.Printf("%v ", v)
	}
	fmt.Printf("\n")
}

func TestBasic(t *testing.T) {
	set := sets.New(1, 2, 3, 3)
	set.Insert(5)

	printSet(set, "set values: ")

	fmt.Printf("set.Has(1): %v\n", set.Has(1))
	fmt.Printf("set.Has(10): %v\n", set.Has(10))

	set1 := sets.New(3, 3, 8)
	printSet(set1, "set1 values: ")

	set.HasAny()

	printSet(set.Difference(set1), "set.Difference(set1): ")

	fmt.Printf("set.Equal(set1): %v\n", set.Equal(set1))
	fmt.Printf("set.Equal(set): %v\n", set.Equal(set.Clone()))

	// set 和 set1 的交集
	printSet(set.Intersection(set1), "set.Intersection(set1): ")
	// set 是否是 set1 的超集
	fmt.Println("set.IsSuperset(set1): ", set.IsSuperset(set1))
	// SymmetricDifference 返回一组元素，这些元素位于任一集合中，但不在它们的交集中。
	printSet(set.SymmetricDifference(set1), "set.SymmetricDifference(set1): ")
	// set 和 set1 所有元素的集合
	printSet(set.Union(set1), "set.Union(set1): ")
}

func TestPtr(t *testing.T) {
	a := 1
	b := 1
	printSet(sets.New(&a, &b), "sets.New(&a, &b): ")
}
