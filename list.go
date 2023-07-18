package main

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

// Return whether iterator has next element
func (iter *Iterator[T]) hasNext() bool {
	return iter.list != nil && iter.list.next != nil && (iter.ret == nil || iter.list.next.next != nil)
}

// Return next element of iterator
func (iter *Iterator[T]) next() (*List[T], error) {
	if !iter.hasNext() {
		return nil, &InvalidError{"next", "no next element"}
	}
	// Increment only if return value is set
	if iter.ret != nil {
		iter.list = iter.list.next
	}
	iter.ret = iter.list.next
	return iter.ret, nil
}

// Remove element last returned by iterator
func (iter *Iterator[T]) remove() error {
	if iter.ret == nil {
		return &InvalidError{"remove", "no element to remove"}
	}
	iter.list.next = iter.ret.next
	iter.ret = nil
	return nil
}

// Add element with iterator
func (iter *Iterator[T]) add(v T) {
	if iter.list.next == nil {
		iter.list.next = &List[T]{nil, v}
		iter.list = iter.list.next
	} else {
		iter.list.next.next = &List[T]{iter.list.next.next, v}
		iter.list = iter.list.next.next
		iter.ret = nil
	}
}

// Return error message
func (e *IndexError) Error() string {
	return fmt.Sprintf("Index %v out of bounds in list", e.index)
}

// Return error message
func (e *InvalidError) Error() string {
	return fmt.Sprintf("Invalid operation: %v\nReason: %v", e.op, e.reason)
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

// Delete the first occurence of v,
// Return whether deletion succeeded
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

// Set index to v,
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

// Insert values at index,
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

// Get value at index,
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

// Remove element at index,
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

// Clear list
func (list *List[_]) clear() {
	list.next = nil
}

// Return sublist starting at index
func (list *List[T]) sublist(index int) (*List[T], error) {
	len := list.length()
	if index < 0 || index >= len {
		return nil, &IndexError{index}
	}
	// List starts 1 before index
	index--
	for i := 0; i <= index; i++ {
		list = list.next
	}
	return list, nil
}

// Insert list at index
func (list *List[T]) insertList(index int, other *List[T]) error {
	if index < 0 || index > list.length() {
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

// Merge sort the list
func (list *List[_]) sort() {
	list.msort(0, list.length()-1)
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
	s1, _ := list.sublist(lo)
	s2, _ := list.sublist(mid + 1)

	// End
	e, _ := list.sublist(hi + 2)

	// Iterators
	i1 := Iterator[T]{s1, nil}
	i2 := Iterator[T]{s2, nil}

	// Values
	v1, _ := i1.next()
	v2, _ := i2.next()

	// Temporary list
	temp := List[T]{}
	iter := Iterator[T]{&temp, nil}

	for v1 != v2 && v2 != e {
		if v1.val < v2.val {
			iter.add(v1.val)
			i1.remove()
			v1, _ = i1.next()
		} else {
			iter.add(v2.val)
			i2.remove()
			v2, _ = i2.next()
		}
	}

	for v1 != e {
		iter.add(v1.val)
		i1.remove()
		v1, _ = i1.next()
	}
	s1.insertList(0, &temp)
}

// Binary search for v (inefficient in linked list)
// List must be sorted
// Return index, (-insertion point-1) if not found
func (list *List[T]) search(v T) int {
	hi := list.length() - 1
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

// Fisher-Yates shuffle
func (list *List[T]) shuffle() {
	rand.Seed(time.Now().UnixNano())
	len := list.length()
	for i := 0; i < len-1; i++ {
		randi := rand.Intn(len-i) + i
		val, _ := list.get(i)
		swap, _ := list.set(randi, val)
		list.set(i, swap)
	}
}

// search, sort, import package from github, helper package
func main() {
	// Initialize with dummy node
	fmt.Println(red + "Initializing..." + reset)
	var list List[int]
	fmt.Println(list)
	fmt.Println()

	// Test append
	appendTest(&list)

	// Test delete
	deleteTest(&list)

	// Test set and error
	setTest(&list)

	// Test insert
	// Note: Go does not support method overloading
	insertTest(&list)

	// Test indexOf
	indexOfTest(&list)

	// Test get (a bit late)
	getTest(&list)

	// Test variable arguments in add and insert
	variableTest(&list)

	// Test remove
	removeTest(&list)

	// Test length
	lengthTest(&list)

	// Test iterator
	iteratorTest(&list)

	// Test sublist
	sublistTest(&list)

	// Test insertList
	insertListTest(&list)

	// Test shuffle
	shuffleTest(&list)

	// Test merge sort
	// sortTest(&list)

	// Test binary search
	searchTest(&list)

	// Test all with generic type
	genericTest()
}

// Colors
var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"
	white  = "\033[97m"
)

func appendTest(list *List[int]) {
	fmt.Println(red + "Testing append..." + reset)
	fmt.Println(list)
	fmt.Println("Add 1 and 2")
	list.add(1)
	list.add(2)
	fmt.Println(list)
	fmt.Println()
}

func deleteTest(list *List[int]) {
	fmt.Println(red + "Testing delete..." + reset)
	fmt.Println(list)
	fmt.Println("Delete 1")
	found := list.delete(1)
	fmt.Println("Found:", found)
	fmt.Println(list)
	fmt.Println()
}

func setTest(list *List[int]) {
	fmt.Println(red + "Testing set..." + reset)
	fmt.Println(list)
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
}

func insertTest(list *List[int]) {
	fmt.Println(red + "Testing insert..." + reset)
	fmt.Println(list)
	fmt.Println("Insert 2 at index 0")
	err := list.insert(0, 2)
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
}

func indexOfTest(list *List[int]) {
	fmt.Println(red + "Testing indexOf..." + reset)
	fmt.Println(list)
	fmt.Println("Index of 1")
	index := list.indexOf(1)
	fmt.Println("Index:", index)
	fmt.Println(list)

	fmt.Println("Index of 0")
	index = list.indexOf(0)
	fmt.Println("Index:", index)
	fmt.Println(list)
	fmt.Println()
}

func getTest(list *List[int]) {
	fmt.Println(red + "Testing get..." + reset)
	fmt.Println(list)
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
}

func variableTest(list *List[int]) {
	fmt.Println(red + "Testing variable number of arguments..." + reset)
	fmt.Println(list)
	fmt.Println("Add 4 and 5")
	list.add(4, 5)
	fmt.Println(list)

	fmt.Println("Insert 6 and 7 at index 0")
	list.insert(0, 6, 7)
	fmt.Println(list)
	fmt.Println()
}

func removeTest(list *List[int]) {
	fmt.Println(red + "Testing remove..." + reset)
	fmt.Println(list)
	fmt.Println("Remove index 3")
	elem, err := list.remove(3)
	fmt.Println(list)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)

	fmt.Println("Remove index 6")
	elem, err = list.remove(6)
	fmt.Println(list)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
	fmt.Println()
}

func lengthTest(list *List[int]) {
	fmt.Println(red + "Testing length..." + reset)
	fmt.Println(list)
	len := list.length()
	fmt.Println("Length of list: ", len)
	fmt.Println()
}

func iteratorTest(list *List[int]) {
	fmt.Println(red + "Testing iterator..." + reset)
	fmt.Println(list)
	iter := Iterator[int]{list, nil}
	fmt.Println("Get the next element")
	val, err := iter.next()
	fmt.Println(val.val)
	fmt.Println("Error: ", err)
	err = iter.remove()
	fmt.Println("Remove the last returned element")
	fmt.Println("Error: ", err)
	fmt.Println(list)
	err = iter.remove()
	fmt.Println("Remove the last returned element")
	fmt.Println("Error: ", err)
	fmt.Println(list)
	fmt.Println("Iterate through the list")
	for iter.hasNext() {
		val, err := iter.next()
		fmt.Println(val.val)
		fmt.Println("Error: ", err)
	}
	val, err = iter.next()
	fmt.Println(val)
	fmt.Println("Error: ", err)
	fmt.Println()
}

func sublistTest(list *List[int]) {
	fmt.Println(red + "Testing sublist..." + reset)
	fmt.Println(list)
	fmt.Println("Sublist starting from 0")
	li, err := list.sublist(0)
	fmt.Println(li)
	fmt.Println("Error: ", err)
	fmt.Println("Sublist starting from 1")
	li, err = list.sublist(1)
	fmt.Println(li)
	fmt.Println("Error: ", err)
	fmt.Println("Sublist starting from 4")
	li, err = list.sublist(4)
	fmt.Println(li)
	fmt.Println("Error: ", err)
	fmt.Println("Sublist starting from 5")
	li, err = list.sublist(5)
	fmt.Println(li)
	fmt.Println("Error: ", err)
	fmt.Println()
}

func insertListTest(list *List[int]) {
	fmt.Println(red + "Testing insertList..." + reset)
	fmt.Println(list)
	list2 := List[int]{}
	list2.add(8, 9, 1)
	fmt.Println("Insert", list2, "at index 3")
	err := list.insertList(3, &list2)
	fmt.Println(list)
	fmt.Println("Error: ", err)
	list3 := List[int]{}
	list3.add(6)
	fmt.Println("Insert", list3, "at index 8")
	err = list.insertList(8, &list3)
	fmt.Println(list)
	fmt.Println("Error: ", err)
	fmt.Println()
}

func shuffleTest(list *List[int]) {
	fmt.Println(red + "Testing shuffle..." + reset)
	fmt.Println(list)
	fmt.Println("Shuffle list")
	list.shuffle()
	fmt.Println(list)
	fmt.Println()
}

func sortTest(list *List[int]) {
	fmt.Println(red + "Testing merge sort..." + reset)
	fmt.Println(list)
	fmt.Println("Sort list")
	list.sort()
	fmt.Println(list)
	fmt.Println()
}

func searchTest(list *List[int]) {
	fmt.Println(red + "Testing binary search..." + reset)
	fmt.Println(list)
	for i := 1; i <= 9; i++ {
		fmt.Println("Search for", i)
		fmt.Println(list.search(i))
	}
	fmt.Println("Search for 0")
	fmt.Println(list.search(0))
	fmt.Println("Search for 10")
	fmt.Println(list.search(10))
	fmt.Println()
}

func genericTest() {
	fmt.Println(red + "Testing generic types..." + reset)
	fmt.Println(red + "Initialize" + reset)
	var strings List[string]
	fmt.Println(strings)
	fmt.Println(red + "Append" + reset)
	strings.add("goodbye", "cruel", "world")
	fmt.Println(strings)
	fmt.Println(red + "Delete" + reset)
	strings.delete("cruel")
	fmt.Println(strings)
	fmt.Println(red + "Set" + reset)
	strings.set(0, "hello")
	fmt.Println(strings)
	fmt.Println(red + "Insert" + reset)
	strings.insert(1, "big", "beautiful")
	fmt.Println(strings)
	fmt.Println(red + "IndexOf" + reset)
	fmt.Println(strings.indexOf("world"))
	fmt.Println(red + "Get" + reset)
	str, err := strings.get(0)
	fmt.Println(str)
	fmt.Println("Error: ", err)
	fmt.Println(red + "Remove" + reset)
	strings.remove(1)
	fmt.Println(strings)
	fmt.Println(red + "Shuffle" + reset)
	for c := 'A'; c <= 'Z'; c++ {
		strings.add(string(c))
	}
	for c := 'a'; c <= 'z'; c++ {
		strings.add(string(c))
	}
	strings.shuffle()
	fmt.Println(strings)
	fmt.Println(red + "Sort" + reset)
	// strings.sort()
	fmt.Println(strings)
	fmt.Println()
	strings.clear()
	strings.add("Eandy", "Kandy", "Dandy", "Oandy", "Handy", "beautiful", "Gandy", "Landy", "Xandy", "Pandy", "Yandy", "Sandy", "hello", "Vandy", "Andy", "Nandy", "Candy", "world", "Iandy", "Mandy", "Wandy", "Randy", "Bandy", "Fandy", "Jandy", "Zandy", "Uandy", "Tandy", "Qandy")
	strings.sort()
	fmt.Println(strings)
}
