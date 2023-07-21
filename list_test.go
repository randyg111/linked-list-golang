package main

import "testing"

var cases = 1000

func handle(s string, t *testing.T) {
	if r := recover(); r != nil {
		t.Error(s+":", r)
	}
}

func TestIterator_Add(t *testing.T) {
	defer handle("Iterator Add", t)
	var list List[int]
	iter := Iterator[int]{&list, nil}
	iter.Add(0)
}

func TestHasNext(t *testing.T) {
	defer handle("HasNext", t)
	var list List[int]
	iter := Iterator[int]{&list, nil}
	b := iter.HasNext()
	if b {
		t.Error("HasNext: expected false but got", b)
	}
}

func TestNext(t *testing.T) {
	defer handle("Next", t)
	var list List[int]
	iter := Iterator[int]{&list, nil}
	_, err := iter.Next()
	if err == nil {
		t.Error("Next: expected error but got", err)
	}
}

func TestIterator_Remove(t *testing.T) {
	defer handle("Iterator Remove", t)
	var list List[int]
	iter := Iterator[int]{&list, nil}
	err := iter.Remove()
	if err == nil {
		t.Error("Iterator Remove: expected error but got", err)
	}
}

func TestList_Add(t *testing.T) {
	defer handle("List Add", t)
	var list List[int]
	list.Add(0)
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

func TestClear(t *testing.T) {
	defer handle("Clear", t)
	var list List[int]
	list.Clear()
}

func TestCopy(t *testing.T) {
	defer handle("Copy", t)
	var list List[int]
	l := list.Copy()
	if l.String() != list.String() {
		t.Error("Copy: expected", list.String(), "but got", l.String())
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

func TestGet(t *testing.T) {
	defer handle("Get", t)
	var list List[int]
	_, err := list.Get(0)
	if err == nil {
		t.Error("Get: expected error but got", err)
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

func TestInsertList(t *testing.T) {
	defer handle("InsertList", t)
	var list List[int]
	err := list.InsertList(0, &List[int]{})
	if err != nil {
		t.Error("InsertList:", err)
	}
}

func TestLength(t *testing.T) {
	defer handle("Length", t)
	var list List[int]
	l := list.Length()
	if l != 0 {
		t.Error("Length: expected 0 but got", l)
	}
}

func TestList_Remove(t *testing.T) {
	defer handle("List Remove", t)
	var list List[int]
	_, err := list.Remove(0)
	if err == nil {
		t.Error("List Remove: expected error but got", err)
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

func TestSet(t *testing.T) {
	defer handle("Set", t)
	var list List[int]
	_, err := list.Set(0, 0)
	if err == nil {
		t.Error("Set: expected error but got", err)
	}
}

func TestShuffle(t *testing.T) {
	defer handle("Shuffle", t)
	var list List[int]
	list.Shuffle()
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

func TestSorted(t *testing.T) {
	defer handle("Sorted", t)
	var list List[int]
	b := list.Sorted()
	if !b {
		t.Error("Sorted: expected false but got", b)
	}
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
