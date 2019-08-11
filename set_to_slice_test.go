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
	"reflect"
	"sort"
	"testing"
)

func Test_set_ToStrings(t *testing.T) {
	tests := []struct {
		name string
		s    *set
		want []string
	}{
		{"", newSet(1, 2, "str", "str1"), []string{"str", "str1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.ToStrings()
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("set.ToStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_ToInts(t *testing.T) {
	tests := []struct {
		name string
		s    *set
		want []int
	}{
		{"", newSet(1, 2, "str", "str1"), []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.ToInts()
			sort.Ints(got)
			sort.Ints(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("set.ToInts() = %v, want %v", got, tt.want)
			}
		})
	}
}
