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
	"sync"
	"testing"
)

func Test_ThreadSafeSet_Add(t *testing.T) {
	s := newThreadSafeSet()
	is := []interface{}{1, 2, "3", 4.0, 5.1}
	var wg sync.WaitGroup
	for _, i := range is {
		wg.Add(1)
		go func(i interface{}) {
			defer wg.Done()
			s.Add(i)
		}(i)
	}

	wg.Wait()
	for _, i := range is {
		if !s.Contains(i) {
			t.Errorf("threadSafeSet.Add() missing element: %v", i)
		}
	}
}

func Test_threadSafeSet_Extend(t *testing.T) {
	s := newThreadSafeSet()
	is := []interface{}{
		[]int{1, 2},
		[]string{"3", "4"},
		NewSafeSet(1, "3", 5),
	}
	var wg sync.WaitGroup
	for _, i := range is {
		wg.Add(1)
		go func(i interface{}) {
			defer wg.Done()
			s.Extend(i)
		}(i)
	}

	wg.Wait()

	elems := []interface{}{1, 2, "3", "4", 5}
	for _, i := range elems {
		if !s.Contains(i) {
			t.Errorf("threadSafeSet.Extend() missing element: %v", i)
		}
	}
}

func Test_threadSafeSet_Remove(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)
	is := []interface{}{1, 2, "test", 4}

	var wg sync.WaitGroup
	for _, i := range is {
		wg.Add(1)
		go func(i interface{}) {
			defer wg.Done()
			s.Remove(i)
		}(i)
	}

	wg.Wait()

	for _, i := range is {
		if s.Contains(i) {
			t.Errorf("threadSafeSet.Remove() element: %v can not be removed", i)
		}
	}
}

func Test_threadSafeSet_Contains(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)
	is := []interface{}{1, 2, "3", 4.0, 5.1}

	var wg sync.WaitGroup
	for _, i := range is {
		wg.Add(1)
		go func(i interface{}) {
			defer wg.Done()
			if s.Contains(i) != true {
				t.Errorf("threadSafeSet.Contains() missing elements: %v", i)
			}
		}(i)
	}

	wg.Wait()
}

func Test_threadSafeSet_Equal(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want bool
	}

	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			true,
		},
		{
			newThreadSafeSet(1, 2, 3),
			false,
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.Equal(tt.b); got != tt.want {
				t.Errorf("threadSafeSet.Equal() == %v, want %v", got, tt.b)
			}
		}(tt)
	}

	wg.Wait()
}

func Test_threadSafeSet_IsSubsetOf(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want bool
	}

	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			true,
		},
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1, "6"),
			true,
		},
		{
			newThreadSafeSet(1, 2, 3),
			false,
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.IsSubsetOf(tt.b); got != tt.want {
				t.Errorf("threadSafeSet.IsSubsetOf() == %v, want %v", got, tt.b)
			}
		}(tt)
	}

	wg.Wait()
}

func Test_threadSafeSet_IsSupersetOf(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want bool
	}

	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			true,
		},
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1, "6"),
			false,
		},
		{
			newThreadSafeSet(1, 2, 3),
			false,
		},
		{
			newThreadSafeSet(1, 2),
			true,
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.IsSupersetOf(tt.b); got != tt.want {
				t.Errorf("threadSafeSet.IsSupersetOf() == %v, want %v", got, tt.want)
			}
		}(tt)
	}

	wg.Wait()
}

func Test_threadSafeSet_ToThreadUnsafe_And_Safe(t *testing.T) {
	tests := []struct {
		name string
		a    *threadSafeSet
		want bool
	}{
		{
			"",
			newThreadSafeSet(1, 2, 3, 4),
			true,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			unsafe := tt.a.ToThreadUnsafe()
			if _, got := unsafe.(*set); got != tt.want {
				t.Errorf("threadSafeSet.ToThreadUnsafe() = %v, want %v", got, tt.want)
			}
			safe := tt.a.ToThreadSafe()
			if _, got := safe.(*threadSafeSet); got != tt.want {
				t.Errorf("threadSafeSet.ToThreadUnsafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threadSafeSet_Diff(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want Set
	}
	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			newThreadSafeSet(),
		},
		{
			newThreadSafeSet("3", 4.0, 5.1, "6.6"),
			newThreadSafeSet(1, 2),
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.Diff(tt.b); !got.Equal(tt.want) {
				t.Errorf("threadSafeSet.Diff() == %v, want %v", got, tt.want)
			}
		}(tt)
	}

	wg.Wait()
}

func Test_threadSafeSet_SymmetricDiff(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want Set
	}
	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			newThreadSafeSet(),
		},
		{
			newThreadSafeSet("3", 4.0, 5.1, "6.6"),
			newThreadSafeSet(1, 2, "6.6"),
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.SymmetricDiff(tt.b); !got.Equal(tt.want) {
				t.Errorf("threadSafeSet.SymmetricDiff() == %v, want %v", got, tt.want)
			}
		}(tt)
	}

	wg.Wait()
}

func Test_threadSafeSet_Unite(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want Set
	}
	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
		},
		{
			newThreadSafeSet("3", 4.0, 5.1, "6.6"),
			newThreadSafeSet(1, 2, "3", 4.0, 5.1, "6.6"),
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.Unite(tt.b); !got.Equal(tt.want) {
				t.Errorf("threadSafeSet.Unite() == %v, want %v", got, tt.want)
			}
		}(tt)
	}
}

func Test_threadSafeSet_Intersect(t *testing.T) {
	s := newThreadSafeSet(1, 2, "3", 4.0, 5.1)

	type test struct {
		b    *threadSafeSet
		want Set
	}
	tests := []test{
		{
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
			newThreadSafeSet(1, 2, "3", 4.0, 5.1),
		},
		{
			newThreadSafeSet("3", 4.0, 5.1, "6.6"),
			newThreadSafeSet("3", 4.0, 5.1),
		},
	}
	var wg sync.WaitGroup
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		go func(tt test) {
			defer wg.Done()
			if got := s.Intersect(tt.b); !got.Equal(tt.want) {
				t.Errorf("threadSafeSet.Intersect() == %v, want %v", got, tt.want)
			}
		}(tt)
	}
}
