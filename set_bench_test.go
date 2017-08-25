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
	"math/rand"
	"testing"
)

// the following benchmark is no need to do,
// because finding an elements in a set using O(1) time.
// the operating time depends on how many
// elements need to be processed per operation
//
// s.IsSubsetOf(b) depends on len(s)
// s.IsSupersetOf(b) depends on len(b)
// s.Diff(b) depends on len(s)
// s.SymmetricDiff(b) depends on len(s) + len(b)
// s.Unit(b) depends on len(s) + len(b)
// s.Intersect(b) depends on min(len(s), len(b))

func benchmarkAdd(b *testing.B, s Set, clear bool) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
		if clear {
			s.Clear()
		}
	}
}
func BenchmarkUnsafeAddWithClear(b *testing.B) {
	benchmarkAdd(b, makeSet(), true)
}

func BenchmarkSafeAddWithClear(b *testing.B) {
	benchmarkAdd(b, makeThreadSafeSet(), true)
}

func BenchmarkUnsafeAddWithoutClear(b *testing.B) {
	benchmarkAdd(b, makeSet(), false)
}

func BenchmarkSafeAddWithoutClear(b *testing.B) {
	benchmarkAdd(b, makeThreadSafeSet(), false)
}

func benchmarkExtend(b *testing.B, s Set, clear bool) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Extend([]int{i})
		if clear {
			s.Clear()
		}
	}
}

func BenchmarkUnSafeExtendWithClear(b *testing.B) {
	benchmarkExtend(b, makeSet(), true)
}
func BenchmarkSafeExtendWithClear(b *testing.B) {
	benchmarkExtend(b, makeThreadSafeSet(), true)
}

func BenchmarkUnSafeExtendWithoutClear(b *testing.B) {
	benchmarkExtend(b, makeSet(), false)
}

func BenchmarkSafeExtendWithoutClear(b *testing.B) {
	benchmarkExtend(b, makeThreadSafeSet(), false)
}

func benchmarkCopy(b *testing.B, s Set, scale int) {
	s.Extend(rand.Perm(scale))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Copy()
	}
}

func BenchmarkUnsafeCopy1(b *testing.B) {
	benchmarkCopy(b, makeSet(), 1)
}

func BenchmarkSfaeCopy1(b *testing.B) {
	benchmarkCopy(b, makeThreadSafeSet(), 1)
}

func BenchmarkUnsafeCopy10(b *testing.B) {
	benchmarkCopy(b, makeSet(), 10)
}

func BenchmarkSfaeCopy10(b *testing.B) {
	benchmarkCopy(b, makeThreadSafeSet(), 10)
}

func BenchmarkUnsafeCopy100(b *testing.B) {
	benchmarkCopy(b, makeSet(), 100)
}

func BenchmarkSfaeCopy100(b *testing.B) {
	benchmarkCopy(b, makeThreadSafeSet(), 100)
}

func benchmarkLen(b *testing.B, s Set) {
	s.Extend(rand.Perm(1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Len()
	}
}

func BenchmarkUnsafeLen(b *testing.B) {
	benchmarkLen(b, makeSet())
}
func BenchmarkSafeLen(b *testing.B) {
	benchmarkLen(b, makeThreadSafeSet())
}

func benchmarkContains(b *testing.B, s Set) {
	s.Extend(rand.Perm(100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i)
	}
}

func BenchmarkUnsafeContains(b *testing.B) {
	benchmarkContains(b, makeSet())
}

func BenchmarkSafeContains(b *testing.B) {
	benchmarkContains(b, makeThreadSafeSet())
}

func benchmarkEqual(b *testing.B, x Set, y Set, scale int) {
	ints := rand.Perm(scale)
	x.Extend(ints)
	y.Extend(ints)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Equal(y)
	}
}

func BenchmarkUnsafeEqual1(b *testing.B) {
	benchmarkEqual(b, makeSet(), makeSet(), 1)
}

func BenchmarkSafeEqual1(b *testing.B) {
	benchmarkEqual(b, makeThreadSafeSet(), makeThreadSafeSet(), 1)
}

func BenchmarkUnsafeEqual10(b *testing.B) {
	benchmarkEqual(b, makeSet(), makeSet(), 10)
}

func BenchmarkSafeEqual10(b *testing.B) {
	benchmarkEqual(b, makeThreadSafeSet(), makeThreadSafeSet(), 10)
}
