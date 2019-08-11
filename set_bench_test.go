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
	"strconv"
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

func fill(s Set, scal int) {
	for i := 0; i < scal; i++ {
		s.Add(i)
		s.Add(strconv.Itoa(i))
		s.Add(float64(i))
	}
}

func benchmarkAdd(b *testing.B, s Set) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
		s.Add(strconv.Itoa(i))
		s.Add(float64(i))
	}
}

func BenchmarkUnsafeAdd(b *testing.B) {
	benchmarkAdd(b, newSet())
}

func BenchmarkSafeAdd(b *testing.B) {
	benchmarkAdd(b, newThreadSafeSet())
}

func benchmarkExtend(b *testing.B, s Set) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Extend([]interface{}{i, strconv.Itoa(i)})
	}
}

func BenchmarkUnSafeExtend(b *testing.B) {
	benchmarkExtend(b, newSet())
}

func BenchmarkSafeExtend(b *testing.B) {
	benchmarkExtend(b, newThreadSafeSet())
}

func benchmarkCopy(b *testing.B, s Set, scale int) {
	fill(s, scale)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Copy()
	}
}

func BenchmarkUnsafeCopy1(b *testing.B) {
	benchmarkCopy(b, newSet(), 1)
}

func BenchmarkSafeCopy1(b *testing.B) {
	benchmarkCopy(b, newThreadSafeSet(), 1)
}

func BenchmarkUnsafeCopy10(b *testing.B) {
	benchmarkCopy(b, newSet(), 10)
}

func BenchmarkSafeCopy10(b *testing.B) {
	benchmarkCopy(b, newThreadSafeSet(), 10)
}

func BenchmarkUnsafeCopy100(b *testing.B) {
	benchmarkCopy(b, newSet(), 100)
}

func BenchmarkSafeCopy100(b *testing.B) {
	benchmarkCopy(b, newThreadSafeSet(), 100)
}

func benchmarkLen(b *testing.B, s Set) {
	fill(s, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Len()
	}
}

func BenchmarkUnsafeLen(b *testing.B) {
	benchmarkLen(b, newSet())
}
func BenchmarkSafeLen(b *testing.B) {
	benchmarkLen(b, newThreadSafeSet())
}

func benchmarkContains(b *testing.B, s Set) {
	fill(s, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(i)
	}
}

func BenchmarkUnsafeContains(b *testing.B) {
	benchmarkContains(b, newSet())
}

func BenchmarkSafeContains(b *testing.B) {
	benchmarkContains(b, newThreadSafeSet())
}

func benchmarkEqual(b *testing.B, x Set, y Set, scale int) {
	fill(x, scale)
	fill(y, scale)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Equal(y)
	}
}

func BenchmarkUnsafeEqual1(b *testing.B) {
	benchmarkEqual(b, newSet(), newSet(), 1)
}

func BenchmarkSafeEqual1(b *testing.B) {
	benchmarkEqual(b, newThreadSafeSet(), newThreadSafeSet(), 1)
}

func BenchmarkUnsafeEqual11(b *testing.B) {
	benchmarkEqual(b, newSet(), newSet(), 10)
}

func BenchmarkSafeEqual10(b *testing.B) {
	benchmarkEqual(b, newThreadSafeSet(), newThreadSafeSet(), 10)
}

func benchmarkRange(b *testing.B, s Set, scale int) {
	fill(s, scale)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Range(func(_ int, elem interface{}) bool {
			return true
		})
	}
}

func BenchmarkUnsafeRange10(b *testing.B) {
	benchmarkRange(b, newSet(), 10)
}

func BenchmarkSafeRange10(b *testing.B) {
	benchmarkRange(b, newThreadSafeSet(), 10)
}
