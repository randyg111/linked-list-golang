package main

import "testing"

var cases = 1000

func handle(s string, t *testing.T) {
	if r := recover(); r != nil {
		t.Error(s+":", r)
	}
}

func TestSet(t *testing.T) {
	defer handle("Set", t)
	var list List[int]
	_, err := list.Set(0, 0)
	if err == nil {
		t.Error("Set: expected error but got", err)
	}
}

func TestDelete(t *testing.T) {
	defer handle("Delete", t)
	var list List[int]
	b := list.Delete(0)
	if b {
		t.Error("Delete: expected false but got", b)
	}
}

func TestAdd(t *testing.T) {
	defer handle("Add", t)
	var list List[int]
	list.Add(0)
}

func TestLength(t *testing.T) {
	defer handle("Length", t)
	var list List[int]
	l := list.Length()
	if l != 0 {
		t.Error("Length: expected 0 but got", l)
	}
}

func TestIndexOf(t *testing.T) {
	defer handle("IndexOf", t)
	var list List[int]
	i := list.IndexOf(0)
	if i != -1 {
		t.Error("IndexOf: expected -1 but got", i)
	}
	for i := 1; i <= cases; i++ {
		list.Add(i)
	}
	list.Shuffle()
	for i := 1; i <= cases; i++ {
		index := list.IndexOf(i)
		v, err := list.Get(index)
		if err != nil || i != v {
			t.Error("IndexOf: expected", i, "at index", index, "but got", v)
		}
	}
}

func TestInsert(t *testing.T) {
	defer handle("Insert", t)
	var list List[int]
	err := list.Insert(0, 0)
	if err != nil {
		t.Error("Insert:", err)
	}
}

func TestGet(t *testing.T) {
	defer handle("Get", t)
	var list List[int]
	_, err := list.Get(0)
	if err == nil {
		t.Error("Get: expected error but got", err)
	}
}

func TestRemove(t *testing.T) {
	defer handle("Remove", t)
	var list List[int]
	_, err := list.Remove(0)
	if err == nil {
		t.Error("Remove: expected error but got", err)
	}
}

func TestCopy(t *testing.T) {
	defer handle("Copy", t)
	var list List[int]
	l := list.Copy()
	if l.String() != list.String() {
		t.Error("Copy: expected", list.String(), "but got", l.String())
	}
}

func TestClear(t *testing.T) {
	defer handle("Clear", t)
	var list List[int]
	list.Clear()
}

func TestSublist(t *testing.T) {
	defer handle("Sublist", t)
	var list List[int]
	l, err := list.Sublist(0)
	if err != nil {
		t.Error("Sublist:", err)
	}
	if l.String() != list.String() {
		t.Error("Sublist: expected", list.String(), "but got", l.String())
	}
}

func TestInsertList(t *testing.T) {
	defer handle("InsertList", t)
	var list List[int]
	err := list.InsertList(0, &List[int]{})
	if err != nil {
		t.Error("InsertList:", err)
	}
}

func TestSort(t *testing.T) {
	defer handle("Sort", t)
	var list List[int]
	list.Sort()
	for i := 1; i <= cases; i++ {
		list.Add(i)
	}
	list.Shuffle()
	l := list.Copy()
	list.Sort()
	if !list.Sorted() {
		t.Error("Sort: list not sorted with list", l)
	}
}

func TestSearch(t *testing.T) {
	defer handle("Search", t)
	var list List[int]
	i := list.Search(0)
	if i != -1 {
		t.Error("Search: expected -1 but got", i)
	}
	for i := 1; i <= cases; i++ {
		list.Add(i)
	}
	for i := 1; i <= cases; i++ {
		index := list.Search(i)
		v, err := list.Get(index)
		if err != nil || i != v {
			t.Error("Search: expected", i, "at index", index, "but got", v)
		}
	}
}

func TestShuffle(t *testing.T) {
	defer handle("Shuffle", t)
	var list List[int]
	list.Shuffle()
}

func TestSorted(t *testing.T) {
	defer handle("Sorted", t)
	var list List[int]
	b := list.Sorted()
	if !b {
		t.Error("Sorted: expected false but got", b)
	}
}

func TestBogo(t *testing.T) {
	defer handle("Bogo", t)
	var list List[int]
	list.Bogo()
	for i := 1; i <= cases; i *= 4 {
		list.Add(i)
	}
	list.Shuffle()
	l := list.Copy()
	list.Bogo()
	if !list.Sorted() {
		t.Error("Bogo: list not sorted with list", l)
	}
}
