/*
Copyright 2017 Jim Zhang (jim.zoumo@gmail.com)

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
	"testing"
)

func Test_set_Add(t *testing.T) {
	s := newSet()
	type temp struct {
		test int
	}
	type temp2 struct {
		test int
		m    map[string]int
	}

	tests := []struct {
		name    string
		elem    interface{}
		wantErr bool
	}{
		{"add nil", nil, false},
		{"add dup nil", nil, false},
		{"add int", 10, false},
		{"add float", 10.12, false},
		{"add dup float", 10.12, false},
		{"add struct", temp{}, false},
		{"add error struct", temp2{}, true},
		{"add slice error", []int{1, 2, 3}, true},
		{"add map error", map[string]int{"1": 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Add(tt.elem); (err != nil) != tt.wantErr {
				t.Errorf("set.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_set_Extend(t *testing.T) {
	s := newSet()

	tests := []struct {
		name    string
		b       interface{}
		wantErr bool
	}{
		{"extend int slice", []int{1, 2, 3}, false},
		{"extend string slice", []string{"1", "2"}, false},
		{"extend array", [3]int{3, 4, 5}, false},
		{"extend set", newSet(1, 2, 3, 4), false},
		{"extend err", 1, true},
		{"extend err", "test", true},
		{"extend nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Extend(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("set.Extend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_set_Remove(t *testing.T) {
	s := newSet()
	s.Extend([]interface{}{1, 2, "test", 3.0})

	tests := []struct {
		name string
		elem interface{}
	}{
		{"remove int", 1},
		{"remove string", "test"},
		{"remove float", 3.0},
		{"remove missing", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Remove(tt.elem)
			if s.Contains(tt.elem) != false {
				t.Error("set.Remove() element is not removed")
			}
		})
	}
}

func Test_set_Contains(t *testing.T) {
	s := newSet()
	s.Extend([]interface{}{1, 2, "test", 3.0})

	tests := []struct {
		name string
		elem interface{}
		want bool
	}{
		{"contains int", 1, true},
		{"contains string", "test", true},
		{"contains float", 3.0, true},
		{"contains missing", 3, false},
		{"contains missing", []int{1, 2, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Contains(tt.elem); got != tt.want {
				t.Errorf("set.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Equal(t *testing.T) {

	tests := []struct {
		name string
		a    *set
		b    *set
		want bool
	}{
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 3, 4),
			true,
		},
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 3),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Errorf("set.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_IsSubsetOf(t *testing.T) {

	tests := []struct {
		name string
		a    Set
		b    Set
		want bool
	}{
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 3, 4),
			true,
		},
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 3),
			false,
		},
		{
			"",
			newSet(1, 2, 3),
			newSet(1, 2, 3, 4),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsSubsetOf(tt.b); got != tt.want {
				t.Errorf("set.IsSubsetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_IsSupersetOf(t *testing.T) {

	tests := []struct {
		name string
		a    Set
		b    Set
		want bool
	}{
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 3, 4),
			true,
		},
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 3),
			true,
		},
		{
			"",
			newSet(1, 2, 3),
			newSet(1, 2, 3, 4),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsSupersetOf(tt.b); got != tt.want {
				t.Errorf("set.IsSupersetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_ToThreadUnsafe_And_Safe(t *testing.T) {
	tests := []struct {
		name string
		a    *set
		want bool
	}{
		{
			"",
			newSet(1, 2, 3, 4),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unsafe := tt.a.ToThreadUnsafe()
			if _, got := unsafe.(*set); got != tt.want {
				t.Errorf("set.ToThreadUnsafe() = %v, want %v", got, tt.want)
			}
			safe := tt.a.ToThreadSafe()
			if _, got := safe.(*threadSafeSet); got != tt.want {
				t.Errorf("set.ToThreadUnsafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Diff(t *testing.T) {
	tests := []struct {
		name string
		a    *set
		b    *set
		want Set
	}{
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2),
			newSet(3, 4),
		},
		{
			"",
			newSet(1, "2", 3, 4),
			newSet(1, 2),
			newSet("2", 3, 4),
		},
		{
			"",
			newSet(1),
			newSet(1, "2", 3, 42),
			newSet(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Diff(tt.b); !got.Equal(tt.want) {
				t.Errorf("set.Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_SymmetricDiff(t *testing.T) {
	tests := []struct {
		name string
		a    *set
		b    *set
		want Set
	}{
		{
			"",
			newSet(1, 2, 3, 4),
			newSet(1, 2, 5),
			newSet(3, 4, 5),
		},
		{
			"",
			newSet(1, "2", 3, 4),
			newSet(1, 2),
			newSet(2, "2", 3, 4),
		},
		{
			"",
			newSet(1),
			newSet(1, "2", 3, 4),
			newSet("2", 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.SymmetricDiff(tt.b); !got.Equal(tt.want) {
				t.Errorf("set.SymmetricDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Unite(t *testing.T) {
	tests := []struct {
		name string
		a    *set
		b    *set
		want Set
	}{
		{
			"",
			newSet(1, 2),
			newSet(2, 3, 4),
			newSet(1, 2, 3, 4),
		},
		{
			"",
			newSet(1, "2", 3, 4),
			newSet(1, 2),
			newSet(1, 2, "2", 3, 4),
		},
		{
			"",
			newSet(1),
			newSet(1, "2", 3, 4),
			newSet(1, "2", 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Unite(tt.b); !got.Equal(tt.want) {
				t.Errorf("set.Unite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Intersect(t *testing.T) {
	tests := []struct {
		name string
		a    *set
		b    *set
		want Set
	}{
		{
			"",
			newSet(1, 2),
			newSet(2, 3, 4),
			newSet(2),
		},
		{
			"",
			newSet(1, "2", 3, 4),
			newSet(1, 2, 3),
			newSet(1, 3),
		},
		{
			"",
			newSet(1),
			newSet(1, "2"),
			newSet(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Intersect(tt.b); !got.Equal(tt.want) {
				t.Errorf("set.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Range(t *testing.T) {
	tests := []struct {
		name string
		a    *set
	}{
		{
			"",
			newSet(1, 2, 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := make([]interface{}, 0, tt.a.Len())
			tt.a.Range(func(index int, elem interface{}) bool {
				seen = append(seen, elem)
				return true
			})
			if len(seen) != tt.a.Len() {
				t.Errorf("Range visited %v elements of %v-element Map", len(seen), tt.a.Len())
			}
		})
	}
}
