package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	s1 := NewSet([]interface{}{1, 3.14, true})
	s1.Add("Sam")
	fmt.Println(s1.Get(), s1.Contains("Sam"))
	s1.Remove("Sam")
	fmt.Println(s1.Get(), s1.Contains("Sam"))
	s2 := NewSet([]interface{}{3.14, "Jack"})
	intersect := s1.Intersect(s2)
	union := s1.Union(s2)
	diff := s1.Diff(s2)
	fmt.Println(intersect.Get())
	fmt.Println(union.Get())
	fmt.Println(diff.Get())
	assert.Equal(t, true, reflect.DeepEqual(union, intersect.Union(diff)))
}
