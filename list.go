package list

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/exp/constraints"
)

// List represents a singly-linked list that holds
// values of ordered type (string, int, float),
// Implements Stringer interface
type List[T constraints.Ordered] struct {
	next *List[T]
	val  T
}

// Index out of bounds error,
// Implements error interface
type IndexError struct {
	index int
}

// Invalid operation error,
// Implements error interface
type InvalidError struct {
	op     string
	reason string
}

// Iterator for efficient traversal of linked list
type Iterator[T constraints.Ordered] struct {
	list *List[T]
	ret  *List[T]
}

// Return error message for IndexError
func (e *IndexError) Error() string {
	return fmt.Sprintf("Index %v out of bounds in list", e.index)
}

// Return error message for InvalidError
func (e *InvalidError) Error() string {
	return fmt.Sprintf("Invalid operation: %v\nReason: %v", e.op, e.reason)
}

// Add element with iterator
func (iter *Iterator[T]) Add(v T) {
	if iter.ret == nil {
		iter.list.next = &List[T]{iter.list.next, v}
		iter.list = iter.list.next
	} else {
		iter.list.next.next = &List[T]{iter.list.next.next, v}
		iter.list = iter.list.next.next
		iter.ret = nil
	}
}

// Return whether iterator has next element
func (iter *Iterator[T]) HasNext() bool {
	return iter.list != nil && iter.list.next != nil && (iter.ret == nil || iter.list.next.next != nil)
}

// Return next element of iterator
func (iter *Iterator[T]) Next() (*List[T], error) {
	if !iter.HasNext() {
		return nil, &InvalidError{"next", "no next element"}
	}
	// Increment only if return value is Set
	if iter.ret != nil {
		iter.list = iter.list.next
	}
	iter.ret = iter.list.next
	return iter.ret, nil
}

// Remove element last returned by iterator
func (iter *Iterator[T]) Remove() error {
	if iter.ret == nil {
		return &InvalidError{"Remove", "no element to Remove"}
	}
	iter.list.next = iter.ret.next
	iter.ret = nil
	return nil
}

// Append values to the list
func (list *List[T]) Add(vs ...T) {
	for list.next != nil {
		list = list.next
	}
	for _, v := range vs {
		list.next = &List[T]{nil, v}
		list = list.next
	}
}

// Bogo sort
func (list *List[T]) Bogo() {
	for !list.Sorted() {
		list.Shuffle()
	}
}

// Clear list
func (list *List[_]) Clear() {
	list.next = nil
}

// Copy list to a new list
func (list *List[T]) Copy() List[T] {
	copy := List[T]{nil, list.val}
	next := &copy
	for list.next != nil {
		next.next = &List[T]{nil, list.next.val}
		next = next.next
		list = list.next
	}
	return copy
}

// Delete the first occurence of v,
// Return whether deletion succeeded
func (list *List[T]) Delete(v T) bool {
	for list.next != nil {
		n := list.next
		if n.val == v {
			list.next = n.next
			return true
		}
		list = n
	}
	return false
}

// Get value at index,
// Return error if index out of bounds
func (list *List[T]) Get(index int) (T, error) {
	if index < 0 || index >= list.Length() {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i <= index; i++ {
		list = list.next
	}
	return list.val, nil
}

// Return index of v, or -1 if not found
func (list *List[T]) IndexOf(v T) int {
	list = list.next
	i := 0
	for list != nil {
		if list.val == v {
			return i
		}
		list = list.next
		i++
	}
	return -1
}

// Insert values at index,
// Return error if index out of bounds
func (list *List[T]) Insert(index int, vs ...T) error {
	if index < 0 || index > list.Length() {
		return &IndexError{index}
	}
	for i := 0; i < index; i++ {
		list = list.next
	}
	if list.next == nil {
		for _, v := range vs {
			list.next = &List[T]{nil, v}
			list = list.next
		}
	} else {
		for _, v := range vs {
			list.next = &List[T]{list.next, v}
			list = list.next
		}
	}
	return nil
}

// Insert list at index
func (list *List[T]) InsertList(index int, other *List[T]) error {
	if index < 0 || index > list.Length() {
		return &IndexError{index}
	}
	for i := 0; i < index; i++ {
		list = list.next
	}
	other = other.next
	if list.next == nil {
		list.next = other
	} else if other != nil {
		next := list.next
		list.next = other
		for other.next != nil {
			other = other.next
		}
		other.next = next
	}
	return nil
}

// Return length of list
func (list *List[_]) Length() int {
	len := 0
	for list.next != nil {
		list = list.next
		len++
	}
	return len
}

// Remove element at index,
// Return error if index out of bounds
func (list *List[T]) Remove(index int) (T, error) {
	if index < 0 || index >= list.Length() {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i < index; i++ {
		list = list.next
	}
	old := list.next.val
	list.next = list.next.next
	return old, nil
}

// Binary search for v (inefficient in linked list),
// List must be sorted,
// Return index, or (-insertion point-1) if not found
func (list *List[T]) Search(v T) int {
	hi := list.Length() - 1
	lo := 0
	prev := Iterator[T]{list, nil}
	iter := Iterator[T]{list, nil}
	for hi >= lo {
		mid := (hi + lo) / 2
		for i := lo; i < mid; i++ {
			iter.Next()
		}
		ml, _ := iter.Next()
		m := ml.val
		if v > m {
			lo = mid + 1
			prev = iter
		} else if v < m {
			hi = mid - 1
			iter = prev
		} else {
			return mid
		}
	}
	return -(lo + 1)
}

// Set index to v,
// Return error if index out of bounds
func (list *List[T]) Set(index int, v T) (T, error) {
	if index < 0 || index >= list.Length() {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i <= index; i++ {
		list = list.next
	}
	old := list.val
	list.val = v
	return old, nil
}

// Fisher-Yates shuffle
func (list *List[T]) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	len := list.Length()
	for i := 0; i < len-1; i++ {
		randi := rand.Intn(len-i) + i
		val, _ := list.Get(i)
		swap, _ := list.Set(randi, val)
		list.Set(i, swap)
	}
}

// Merge sort the list
func (list *List[_]) Sort() {
	list.msort(0, list.Length()-1)
}

// Merge sort with recursion
func (list *List[_]) msort(lo, hi int) {
	if hi > lo {
		mid := (hi + lo) / 2
		list.msort(lo, mid)
		list.msort(mid+1, hi)
		list.merge(lo, hi)
	}
}

// Merge sort helper method
func (list *List[T]) merge(lo, hi int) {
	mid := (hi + lo) / 2

	// Start
	s1, _ := list.Sublist(lo)
	s2, _ := list.Sublist(mid + 1)

	// End
	e, _ := list.Sublist(hi + 2)

	// Iterators
	i1 := Iterator[T]{s1, nil}
	i2 := Iterator[T]{s2, nil}

	// Values
	v1, _ := i1.Next()
	v2, _ := i2.Next()

	// Temporary list
	temp := List[T]{}
	iter := Iterator[T]{&temp, nil}

	for v1 != v2 && v2 != e {
		if v1.val < v2.val {
			iter.Add(v1.val)
			i1.Remove()
			v1, _ = i1.Next()
		} else {
			iter.Add(v2.val)
			i2.Remove()
			v2, _ = i2.Next()
		}
	}

	for v1 != e {
		iter.Add(v1.val)
		i1.Remove()
		v1, _ = i1.Next()
	}
	s1.InsertList(0, &temp)
}

// Check if list sorted
func (list *List[T]) Sorted() bool {
	list = list.next
	for list != nil && list.next != nil {
		if list.next.val < list.val {
			return false
		}
		list = list.next
	}
	return true
}

// Convert list to a string
func (list List[_]) String() string {
	if list.next == nil {
		return "[]"
	}
	list = *list.next
	s := fmt.Sprintf("[%v", list.val)
	for list.next != nil {
		list = *list.next
		s += fmt.Sprintf(", %v", list.val)
	}
	s += "]"
	return s
}

// Return sublist starting at index
func (list *List[T]) Sublist(index int) (*List[T], error) {
	len := list.Length()
	if index < 0 || index > len {
		return nil, &IndexError{index}
	}
	// List starts 1 before index
	index--
	for i := 0; i <= index; i++ {
		list = list.next
	}
	return list, nil
}
