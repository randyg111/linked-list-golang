package main

import "fmt"

// List represents a singly-linked list that holds
// values of comparable type
type List[T comparable] struct {
	next *List[T]
	val  T
}

// Index out of bounds error
// Implements error interface
type IndexError struct {
	index int
}

// Error message
func (e *IndexError) Error() string {
	return fmt.Sprintf("Index %v out of bounds in list", e.index)
}

// Convert list to a string
// Implements Stringer interface
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
	if index < 0 {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i <= index; i++ {
		if list.next == nil {
			var fail T
			return fail, &IndexError{index}
		}
		list = list.next
	}
	old := list.val
	list.val = v
	return old, nil
}

// Insert values at index
// Return error if index out of bounds
func (list *List[T]) insert(index int, vs ...T) error {
	if index < 0 {
		return &IndexError{index}
	}
	for i := 0; i < index; i++ {
		if list.next == nil {
			return &IndexError{index}
		}
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
	if index < 0 {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i <= index; i++ {
		if list.next == nil {
			var fail T
			return fail, &IndexError{index}
		}
		list = list.next
	}
	return list.val, nil
}

// Remove element at index
// Return error if index out of bounds
func (list *List[T]) remove(index int) (T, error) {
	if index < 0 {
		var fail T
		return fail, &IndexError{index}
	}
	for i := 0; i < index; i++ {
		if list.next == nil {
			var fail T
			return fail, &IndexError{index}
		}
		list = list.next
	}
	if list.next == nil {
		var fail T
		return fail, &IndexError{index}
	}
	old := list.next.val
	list.next = list.next.next
	return old, nil
}

func main() {
	// Initialize with dummy node
	fmt.Println("Initializing...")
	var list List[int]
	fmt.Println(list)
	fmt.Println()

	// Test append
	fmt.Println("Testing append...")
	list.add(1)
	list.add(2)
	fmt.Println(list)
	fmt.Println()

	// Test delete
	fmt.Println("Testing delete...")
	found := list.delete(1)
	fmt.Println("Found:", found)
	fmt.Println(list)
	fmt.Println()

	// Test set and error
	fmt.Println("Testing set...")
	elem, err := list.set(0, 1)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
	fmt.Println(list)

	elem, err = list.set(2, 1)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
	fmt.Println(list)
	fmt.Println()

	// Test insert
	// Note: Go does not support method overloading
	fmt.Println("Testing insert...")
	err = list.insert(0, 2)
	fmt.Println("Error:", err)
	fmt.Println(list)

	err = list.insert(2, 3)
	fmt.Println("Error:", err)
	fmt.Println(list)

	err = list.insert(4, 4)
	fmt.Println("Error:", err)
	fmt.Println(list)
	fmt.Println()

	// Test indexOf
	fmt.Println("Testing indexOf...")
	index := list.indexOf(1)
	fmt.Println("Index:", index)
	fmt.Println(list)

	index = list.indexOf(0)
	fmt.Println("Index:", index)
	fmt.Println(list)
	fmt.Println()

	// Test get (a bit late)
	fmt.Println("Testing get...")
	val, err := list.get(0)
	fmt.Println("Value:", val)
	fmt.Println("Error:", err)
	fmt.Println(list)

	val, err = list.get(4)
	fmt.Println("Value:", val)
	fmt.Println("Error:", err)
	fmt.Println(list)
	fmt.Println()

	// Test multiple arguments in add and insert
	fmt.Println("Testing multiple arguments...")
	list.add(4, 5)
	fmt.Println(list)

	list.insert(0, 6, 7)
	fmt.Println(list)
	fmt.Println()

	// Test remove
	fmt.Println("Testing remove...")
	elem, err = list.remove(3)
	fmt.Println(list)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)

	elem, err = list.remove(6)
	fmt.Println(list)
	fmt.Println("Removed element:", elem)
	fmt.Println("Error:", err)
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

	// TODO: implement merge sort
	// progress: immediately gave up
}
