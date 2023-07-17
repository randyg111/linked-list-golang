package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// List represents a singly-linked list that holds
// values of ordered type (string, int, float)
// Implements Stringer interface
type List[T constraints.Ordered] struct {
	next *List[T]
	val  T
}

// Index out of bounds error
// Implements error interface
type IndexError struct {
	index int
}

// Invalid operation error
// Implements error interface
type InvalidError struct {
	op string
}

// Iterator for list
// Efficient way to traverse linked list
type Iterator[T constraints.Ordered] struct {
	list *List[T]
	ret  *List[T]
}

// Returns whether iterator has next element
func (iter *Iterator[T]) hasNext() bool {
	return iter.list.next != nil
}

// Returns next element of iterator
func (iter *Iterator[T]) next() T {
	// Do not increment if return value not set
	if iter.ret != nil {
		iter.list = iter.list.next
	}
	iter.ret = iter.list.next
	return iter.list.next.val
}

// Removes value last returned by iterator
func (iter *Iterator[T]) remove() error {
	if iter.ret == nil {
		return &InvalidError{"remove"}
	}
	iter.list.next = iter.ret.next
	iter.ret = nil
	return nil
}

// Returns error message
func (e *IndexError) Error() string {
	return fmt.Sprintf("Index %v out of bounds in list", e.index)
}

// Returns error message
func (e *InvalidError) Error() string {
	return fmt.Sprintf("Invalid operation: %v", e.op)
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

// Append values to the list
func (list *List[T]) add(vs ...T) {
	for list.next != nil {
		list = list.next
	}
	for _, v := range vs {
		list.next = &List[T]{nil, v}
		list = list.next
	}
}

// Delete the first occurence of v
// Returns whether deletion succeeded
func (list *List[T]) delete(v T) bool {
	for list != nil {
		n := list.next
		if n.val == v {
			list.next = n.next
			return true
		}
		list = n
	}
	return false
}

// Set index to v
// Return error if index out of bounds
func (list *List[T]) set(index int, v T) (T, error) {
	if index < 0 || index >= list.length() {
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

// Insert values at index
// Return error if index out of bounds
func (list *List[T]) insert(index int, vs ...T) error {
	if index < 0 || index > list.length() {
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

// Return index of v, -1 if not found
func (list *List[T]) indexOf(v T) int {
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

// Get value at index
// Return error if index out of bounds
func (list *List[T]) get(index int) (T, error) {
	if index < 0 || index >= list.length() {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i <= index; i++ {
		list = list.next
	}
	return list.val, nil
}

// Remove element at index
// Return error if index out of bounds
func (list *List[T]) remove(index int) (T, error) {
	if index < 0 || index >= list.length() {
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

// Return length of list
func (list *List[_]) length() int {
	len := 0
	for list.next != nil {
		list = list.next
		len++
	}
	return len
}

// Splice list to get new list from lo to hi
// func (list *List[T]) splice(lo, hi int) (List[T], error) {
// 	len := list.length()
// 	if lo < 0 || lo > len {
// 		var fail List[T]
// 		return fail, &IndexError{lo}
// 	}
// 	if hi < 0 || hi > len {
// 		var fail List[T]
// 		return fail, &IndexError{lo}
// 	}
// 	if lo > hi {
// 		var fail List[T]
// 		return fail, &InvalidError{fmt.Sprintf("splice with indices %v and %v", lo, hi)}
// 	}
// 	var new List[T]

// 	for i := lo; i <= hi; i++ {
// 		if list.next == nil {
// 			var fail T
// 			return fail, &IndexError{index}
// 		}
// 		list = list.next
// 	}
// 	return list.val, nil
// }

// Binary search for v (inefficient in linked list)
// List must be sorted
// Return index, -1 if not found
func (list *List[T]) search(v T) int {
	hi := list.length()
	lo := 0
	for hi >= lo {
		mid := (hi + lo) / 2
		m, _ := list.get(mid)
		if v > m {
			lo = mid + 1
		} else if v < m {
			hi = mid - 1
		} else {
			return mid
		}
	}
	return -(lo + 1)
}

// Merge sort the list
func (list *List[T]) sort() {
	list.msort(0, list.length())
}

// Merge sort with indices
func (list *List[T]) msort(lo, hi int) {
	if hi > lo {
		mid := (hi + lo) / 2
		list.msort(lo, mid)
		list.msort(mid, hi)
		// list.merge(lo, hi)
	}
}

// Merge sort helper method
// func (list *List[T]) merge(lo, hi int) {
// 	mid := (hi + lo) / 2
// 	var new List[T]
// 	s1, _ := list.getNode(lo)
// 	e1, _ := list.getNode(mid)
// 	s2, _ := list.getNode(mid)
// 	e2, _ := list.getNode(hi)
// 	for s1 != e1 && s2 != e2 {
// 		if s1.val < s2.val {
// 			new.next = s1
// 			new = *new.next
// 		}
// 	}
// }

func main() {
	// Initialize with dummy node
	fmt.Println("Initializing...")
	var list List[int]
	fmt.Println(list)
	fmt.Println()

	// Test append
	fmt.Println("Testing append...")
	fmt.Println("Add 1 and 2")
	list.add(1)
	list.add(2)
	fmt.Println(list)
	fmt.Println()

	// Test delete
	fmt.Println("Testing delete...")
	fmt.Println("Delete 1")
	found := list.delete(1)
	fmt.Println("Found:", found)
	fmt.Println(list)
	fmt.Println()

	// Test set and error
	fmt.Println("Testing set...")
	fmt.Println("Set index 0 to 1")
	elem, err := list.set(0, 1)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
	fmt.Println(list)

	fmt.Println("Set index 1 to 1")
	elem, err = list.set(1, 1)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
	fmt.Println(list)
	fmt.Println()

	// Test insert
	// Note: Go does not support method overloading
	fmt.Println("Testing insert...")
	fmt.Println("Insert 2 at index 0")
	err = list.insert(0, 2)
	fmt.Println("Error:", err)
	fmt.Println(list)

	fmt.Println("Insert 3 at index 2")
	err = list.insert(2, 3)
	fmt.Println("Error:", err)
	fmt.Println(list)

	fmt.Println("Insert 4 at index 4")
	err = list.insert(4, 4)
	fmt.Println("Error:", err)
	fmt.Println(list)
	fmt.Println()

	// Test indexOf
	fmt.Println("Testing indexOf...")
	fmt.Println("Index of 1")
	index := list.indexOf(1)
	fmt.Println("Index:", index)
	fmt.Println(list)

	fmt.Println("Index of 0")
	index = list.indexOf(0)
	fmt.Println("Index:", index)
	fmt.Println(list)
	fmt.Println()

	// Test get (a bit late)
	fmt.Println("Testing get...")
	fmt.Println("Get index 0")
	val, err := list.get(0)
	fmt.Println("Value:", val)
	fmt.Println("Error:", err)
	fmt.Println(list)

	fmt.Println("Get index 4")
	val, err = list.get(4)
	fmt.Println("Value:", val)
	fmt.Println("Error:", err)
	fmt.Println(list)
	fmt.Println()

	// Test variable arguments in add and insert
	fmt.Println("Testing variable number of arguments...")
	fmt.Println("Add 4 and 5")
	list.add(4, 5)
	fmt.Println(list)

	fmt.Println("Insert 6 and 7 at index 0")
	list.insert(0, 6, 7)
	fmt.Println(list)
	fmt.Println()

	// Test remove
	fmt.Println("Testing remove...")
	fmt.Println("Remove index 3")
	elem, err = list.remove(3)
	fmt.Println(list)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)

	fmt.Println("Remove index 6")
	elem, err = list.remove(6)
	fmt.Println(list)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
	fmt.Println()

	// Test length
	fmt.Println("Testing length...")
	len := list.length()
	fmt.Println("Length of list: ", len)
	fmt.Println()

	// Test all with generic type
	var strings List[string]
	fmt.Println(strings)
	strings.add("goodbye", "cruel", "world")
	fmt.Println(strings)
	strings.delete("cruel")
	fmt.Println(strings)
	strings.set(0, "hello")
	fmt.Println(strings)
	strings.insert(1, "big", "beautiful")
	fmt.Println(strings)
	fmt.Println(strings.indexOf("world"))
	str, err := strings.get(0)
	fmt.Println(str)
	strings.remove(1)
	fmt.Println(strings)
	fmt.Println()

	// // Test merge sort
	// fmt.Println(list)

	// fmt.Println()

	// // Test binary search
	// fmt.Println(list)
	// fmt.Println("Search for 6")
	// fmt.Println(list.search(1))
	// fmt.Println("Search for 1")
	// fmt.Println(list.search(6))
	// fmt.Println()

	// fmt.Println(strings)
	// fmt.Println("Search for hello")
	// fmt.Println(strings.search("hello"))
	// fmt.Println("Search for big")
	// fmt.Println(strings.search("big"))
	// fmt.Println()
}
