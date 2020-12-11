/*
Copyright 2019 Jim Zhang (jim.zoumo@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package goset

import (
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
)

type interfaces []interface{}

func Test_typedAssert(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want typed
	}{
		{"", 1, typedInt},
		{"", "str", typedString},
		{"", Empty{}, typedAny},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := typedAssert(tt.args); got != tt.want {
				t.Errorf("typedAssert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Add(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args interfaces
	}{
		{"", newTypedSetGroup(), interfaces{1, "str", 1.2, Empty{}}},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Add(tt.args...)
			if !tt.s.ContainsAll(tt.args...) {
				t.Errorf("typedSet.Add() = %v, want %v", tt.s.Elements(), tt.args)
			}
		})
	}
}

func Test_typedSet_Remove(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args interfaces
		want interfaces
	}{
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{2, "str", 1.2}, interfaces{1, Empty{}}},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Remove(tt.args...)
			if !tt.s.ContainsAll(tt.want...) {
				t.Errorf("typedSet.Remove() = %v, want %v", tt.s.Elements(), tt.want)
			}
		})
	}
}

func Test_typedSet_Contains(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args interface{}
		want bool
	}{
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), 1, true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), "str", true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), Empty{}, true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), 3, false},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), "s", false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Contains(tt.args); got != tt.want {
				t.Errorf("typedSet.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Contains_Any(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args interfaces
		want bool
	}{
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{0, 1, 3}, true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{"str", -1}, true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{1.2}, true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{}, false},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{3}, false},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), interfaces{2.1}, false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ContainsAny(tt.args...); got != tt.want {
				t.Errorf("typedSet.ContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Copy(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		want typedSetGroup
	}{
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", 1.2, Empty{})},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Copy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typedSet.Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Len(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		want int
	}{
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), 5},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("typedSet.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Equal(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args typedSetGroup
		want bool
	}{
		{"", newTypedSetGroup(), newTypedSetGroup(), true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", 1.2, Empty{}), true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), newTypedSetGroup(1, "str", 1.2, Empty{}), false},

		// TODO: Add test cases.
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Equal(tt.args); got != tt.want {
				t.Errorf("typedSet.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_IsSubsetOf(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args typedSetGroup
		want bool
	}{
		{"", newTypedSetGroup(), newTypedSetGroup(), true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", 1.2, Empty{}), true},
		{"", newTypedSetGroup(1, "str", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", 1.2, Empty{}), true},
		{"", newTypedSetGroup(1, 2, "str", 1.2, Empty{}), newTypedSetGroup(1, "str", 1.2, Empty{}), false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsSubsetOf(tt.args); got != tt.want {
				t.Errorf("typedSet.IsSubsetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Range(t *testing.T) {
	s := newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{})
	mu := sync.Mutex{}
	got := []int{}
	want := rand.Perm(s.Len())
	s.Range(func(index int, elem interface{}) bool {
		mu.Lock()
		defer mu.Unlock()
		got = append(got, index)
		return true
	})

	sort.Ints(got)
	sort.Ints(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("typedSet.Range() = %v, want %v", got, want)
	}

	got = []int{}
	want = []int{0, 1, 2}
	s.Range(func(index int, elem interface{}) bool {
		mu.Lock()
		defer mu.Unlock()
		if index > 2 {
			return false
		}
		got = append(got, index)
		return true
	})

	sort.Ints(got)
	sort.Ints(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("typedSet.Range() = %v, want %v", got, want)
	}
}

func Test_typedSet_Diff(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args typedSetGroup
		want typedSetGroup
	}{
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup()},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(2, 3, "str2", Empty{}), newTypedSetGroup(1, "str", 1.2)},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{})},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Diff(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typedSet.Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_SymmetricDiff(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args typedSetGroup
		want typedSetGroup
	}{
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup()},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(2, 3, "str2", Empty{}), newTypedSetGroup(1, 3, "str", 1.2)},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{})},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.SymmetricDiff(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typedSet.SymmetricDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Unite(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args typedSetGroup
		want typedSetGroup
	}{
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{})},
		{"", newTypedSetGroup(1, "str", 1.2, Empty{}), newTypedSetGroup(2, 3, "str2", Empty{}), newTypedSetGroup(1, 2, 3, "str", "str2", 1.2, Empty{})},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{})},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Unite(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typedSet.Unite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Intersect(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		args typedSetGroup
		want typedSetGroup
	}{
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{})},
		{"", newTypedSetGroup(1, "str", 1.2, Empty{}), newTypedSetGroup(2, 3, "str2", Empty{}), newTypedSetGroup(Empty{})},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), newTypedSetGroup(), newTypedSetGroup()},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Intersect(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typedSet.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typedSet_Elements(t *testing.T) {
	tests := []struct {
		name string
		s    typedSetGroup
		want interfaces
	}{
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), interfaces{1, 2, "str", "str2", 1.2, Empty{}}},
		{"", newTypedSetGroup(1, "str", 1.2, Empty{}), interfaces{1, "str", 1.2, Empty{}}},
		{"", newTypedSetGroup(1, 2, "str", "str2", 1.2, Empty{}), interfaces{1, 2, "str", "str2", 1.2, Empty{}}},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Elements(); !tt.s.ContainsAll(got...) {
				t.Errorf("typedSet.Elements() = %v, want %v", got, tt.want)
			}
		})
	}
}
