# goset

Set is a useful collection but there is no built-in implementation in Go lang.

## Why use it

The only one pkg which provides set operations now is [golang-set](https://github.com/deckarep/golang-set)

Unfortunately, the api of golang-set is not good enough.

For example, I want to generate a set from a int slice

```go
import "github.com/deckarep/golang-set"

func main() {
	ints := []int{1, 2, 3, 4}
	mapset.NewSet(ints...)
	mapset.NewSetFromSlice(ints)
	mapset.NewSetWith(ints...)
}
```

the code above can not work, according to

>    cannot use ints (type []int) as type []interface{}

You can not assign any slice to an `[]interface{}`  in Go lang.

>   https://github.com/golang/go/wiki/InterfaceSlice

So you need to copy your elements from `[]int` to `[]interface` by a loop.

That means you must do this manually every time you want to generate a set from slice.

**That is ugly. So I made this one**

## Usage

```go
import "github.com/zoumo/goset"

func main() {
	ints := []int{1, 2, 3, 4}
	floats := []float64{1.0, 2.0, 3.0}
	strings := []string{"1", "2", "3"}
	goset.NewSetFrom(ints)
	goset.NewSafeSetFrom(ints)
	goset.NewSetFrom(floats)
	goset.NewSafeSetFrom(floats)
	goset.NewSetFrom(strings)
	goset.NewSafeSetFrom(strings)
}
```

Full API

```go
// Set provides a collection of oprations for sets
//
// The implementation of Set is base on hash table.
// So the elements must be hashable. functions, maps,
// slices are unhashable type, adding these elements
// will cause panic.
//
// There are two implementations of Set:
// 1. default is unsafe based on hash table
// 2. thread safe based on sync.RWMutex (maybe change to sync.Map in GO1.9)
//
// The two kinds of sets can easily convert to the
// other one. But you must know exactly what you are doing
// to avoid the concurrent race
type Set interface {
	// Add adds an element to the set.
	// If the elem is already exists, it will not return an ErrAlreadyExisted by default.
	// You can set package level var RaiseErrAlreadyExisted=true to raise it.
	Add(elem interface{}) error

	// Extend adds all elements in the given interface b to this set
	// the given interface must be array, slice or set.
	// Extend will ignore elements already existed error.
	Extend(b interface{}) error

	// Remove removes an element from the set.
	Remove(elem interface{})

	// Clear removes all elevemnts from the set.
	Clear()

	// Copy clones the set.
	Copy() Set

	// Len returns the size of set. aka Cardinality.
	Len() int

	// Elements returns all elements in this set.
	Elements() []interface{}

	// Containts checks whether the given item is in the set.
	Contains(item interface{}) bool

	// Equal checks whether this set is equal to the given one.
	// There are two constraints if set a is equal to set b.
	// the two set must have the same size and contain the same elements.
	Equal(b Set) bool

	// IsSubsetOf checks whether this set is the subset of the given set
	// In other words, all elements in this set are also the elements
	// of the given set.
	IsSubsetOf(b Set) bool

	// IsSupersetOf checks whether this set is the superset of the given set
	// In other words, all elements in the given set are also the elements
	// of this set.
	IsSupersetOf(b Set) bool

	// String returns the string representation of the set.
	String() string

	// ToThreadUnsafe returns a thread unsafe set.
	// Carefully use the method.
	ToThreadUnsafe() Set

	// ToThreadSafe returns a thread safe set.
	// Carefully use the method.
	ToThreadSafe() Set

	// Range calls f sequentially for each element present in the set.
	// If f returns false, range stops the iteration.
	//
	// Note: the iteration order is not specified and is not guaranteed
	// to be the same from one iteration to the next. The index only
	// means how many elements have been visited in the iteration, it not
	// specifies the index of an element in the set
	Range(foreach func(index int, elem interface{}) bool)

	// ---------------------------------------------------------------------
	// Set Oprations

	// Diff returns the difference between the set and this given
	// one, aka Difference Set
	// math formula: a - b
	Diff(b Set) Set

	// SymmetricDiff returns the symmetric difference between this set
	// and the given one. aka Symmetric Difference Set
	// math formula: (a - b) ∪ (b - a)
	SymmetricDiff(b Set) Set

	// Unite combines two sets into a new one, aka Union Set
	// math formula: a ∪ b
	Unite(b Set) Set

	// Intersect returns the intersection of two set, aka Intersection Set
	// math formula: a ∩ b
	Intersect(b Set) Set
}
```

