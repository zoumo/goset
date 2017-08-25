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
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

var (
	// RaiseErrAlreadyExisted : if true, set will
	// raise an ErrAlreadyExisted when you attempt
	// to add an already existed element to the set
	RaiseErrAlreadyExisted = false

	// ErrAlreadyExisted represents that an element is already existed in the set
	ErrAlreadyExisted = errors.New("element already exists")
)

type set struct {
	// map[interface{}]bool is 2 times faster for add
	// than map[interface{}]interface{},
	// and bool occupies only a byte.
	data map[interface{}]bool
}

func newSet() *set {
	return &set{
		data: make(map[interface{}]bool),
	}
}

func (s *set) Add(elem interface{}) (err error) {

	defer func() {
		// recover unhashable error
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	_, ok := s.data[elem]
	if ok {
		if RaiseErrAlreadyExisted {
			return ErrAlreadyExisted
		}
		return nil
	}

	s.data[elem] = true

	return nil
}

func (s *set) Extend(b interface{}) error {

	if b == nil {
		return nil
	}

	var elements []interface{}

	setb, ok := b.(Set)
	if !ok {
		v := reflect.ValueOf(b)
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			elements = make([]interface{}, 0)
			for i := 0; i < v.Len(); i++ {
				vv := v.Index(i)
				elements = append(elements, vv.Interface())
			}
		default:
			return fmt.Errorf("error extend set with kind: %v, only support array and slice and Set", v.Kind())
		}
	} else {
		elements = setb.Elements()
	}

	for _, e := range elements {
		err := s.Add(e)
		if err != nil && err != ErrAlreadyExisted {
			return err
		}
	}

	return nil
}

func (s *set) Remove(elem interface{}) {
	delete(s.data, elem)
}

func (s *set) Clear() {
	s.data = make(map[interface{}]bool)
}

func (s *set) Copy() Set {
	c := newSet()

	for k, v := range s.data {
		// TODO: deep copy
		c.data[k] = v
	}

	return c
}

func (s *set) Len() int {
	return len(s.data)
}

func (s *set) Elements() []interface{} {
	elems := make([]interface{}, 0, len(s.data))
	for k := range s.data {
		elems = append(elems, k)
	}
	return elems
}

func (s *set) Contains(elem interface{}) (ret bool) {
	defer func() {
		// recover unhashable error
		if e := recover(); e != nil {
			ret = false
			return
		}
	}()
	_, ok := s.data[elem]
	return ok
}

func (s *set) Equal(b Set) bool {
	b = b.ToThreadUnsafe()

	sLen := s.Len()
	bLen := b.Len()
	if sLen != bLen {
		return false
	}

	if sLen == 0 {
		// empty set is always equal to empty set
		return true
	}

	for key := range s.data {
		if !b.Contains(key) {
			return false
		}
	}

	return true
}

func (s *set) IsSubsetOf(b Set) bool {
	b = b.ToThreadUnsafe()

	sLen := s.Len()
	bLen := b.Len()

	if sLen > bLen {
		return false
	}

	if sLen == 0 {
		return true
	}

	for key := range s.data {
		if !b.Contains(key) {
			return false
		}
	}
	return true

}

func (s *set) IsSupersetOf(b Set) bool {
	b = b.ToThreadUnsafe()

	sLen := s.Len()
	bLen := b.Len()

	if sLen < bLen {
		return false
	}

	if bLen == 0 {
		return true
	}

	for _, e := range b.Elements() {
		if !s.Contains(e) {
			return false
		}
	}
	return true

}

func (s *set) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("Set[")
	first := true
	for key := range s.data {
		if first {
			buf.WriteString(fmt.Sprintf("%+v", key))
			first = false
		} else {
			buf.WriteString(fmt.Sprintf(" %+v", key))
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (s *set) ToThreadUnsafe() Set {
	return s
}

func (s *set) ToThreadSafe() Set {
	return &threadSafeSet{unsafe: s}

}

func (s *set) Diff(b Set) Set {
	b = b.ToThreadUnsafe()
	diff := newSet()
	for key := range s.data {
		if !b.Contains(key) {
			diff.data[key] = true
		}
	}
	return diff
}

func (s *set) SymmetricDiff(b Set) Set {
	b = b.ToThreadUnsafe()
	adiff := s.Diff(b)
	bdiff := b.Diff(s)
	return adiff.Unite(bdiff)
}

func (s *set) Unite(b Set) Set {
	b = b.ToThreadUnsafe()
	union := s.Copy()
	for _, e := range b.Elements() {
		union.Add(e)
	}
	return union

}

func (s *set) Intersect(b Set) Set {
	c := (b.ToThreadUnsafe()).(*set)

	var x, y *set

	// find the smaller one
	if s.Len() <= c.Len() {
		x = s
		y = c
	} else {
		x = c
		y = s
	}

	intersection := newSet()
	for key := range x.data {
		if y.Contains(key) {
			intersection.Add(key)
		}
	}
	return intersection
}

func (s *set) Range(foreach func(int, interface{}) bool) {
	i := 0
	for key := range s.data {
		if !foreach(i, key) {
			break
		}
		i++
	}
}
