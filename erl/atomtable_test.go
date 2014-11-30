package erl

import (
	"strconv"
	"strings"
	"testing"
)

func TestAtomTablePutFirstManyAtomsIncreaseSize(t *testing.T) {
	at := NewAtomTable()
	checkSize(t, at, 0)
	for i := 0; i < 1<<oneByte; i++ {
		checkAtomTablePut(t, at, "atom"+strconv.Itoa(i), uint32(i))
		checkSize(t, at, uint32(i+1))
	}
}

func TestAtomTablePutAtomTwice(t *testing.T) {
	at := NewAtomTable()
	checkAtomTablePut(t, at, "atom", 0)
	checkAtomTablePut(t, at, "atom", 0)
	checkSize(t, at, 1)
}

func TestAtomTableHasAtom(t *testing.T) {
	at := NewAtomTable()
	checkAtomTableHas(t, at, "atom", false)
	at.Put("atom")
	checkAtomTableHas(t, at, "atom", true)
}

func TestAtomTableAtomAt(t *testing.T) {
	at := NewAtomTable()
	at.Put("atom")
	at.Put("second")
	checkAtomTableAt(t, at, 0, "atom")
	checkAtomTableAt(t, at, 1, "second")
	checkAtomTableAt(t, at, 2, "")
}

func TestAtomTableString(t *testing.T) {
	at := NewAtomTable()
	for _, atom := range []string{"fac", "state", "erlang", "-", "*", "module_info", "get_module_info"} {
		at.Put(atom)
	}
	expected := strings.Join(strings.Fields(`<<0,0,0,7,
	3,102,97,99,
	5,115,116,97,116,101,
	6,101,114,108,97,110,103,
	1,45,
	1,42,
	11,109,111,100,117,108,101,95,105,110,102,111,
	15,103,101,116,95,109,111,100,117,108,101,95,105,110,102,111>>`), "")
	if str := at.String(); str != expected {
		t.Errorf("AtomTable.String() => %q, want %q", str, expected)
	}
}

func checkAtomTablePut(t *testing.T, at *AtomTable, atom string, expectedIndex uint32) {
	if index := at.Put(atom); index != expectedIndex {
		t.Errorf("AtomTable.Put(%q) => %v, want %v", atom, index, expectedIndex)
	}
}

func checkAtomTableHas(t *testing.T, at *AtomTable, atom string, expectedHas bool) {
	if has := at.Has(atom); has != expectedHas {
		t.Errorf("AtomTable.Has(%q) => %v, want %v", atom, has, expectedHas)
	}
}

func checkAtomTableAt(t *testing.T, at *AtomTable, index uint32, expectedAtom string) {
	if atom := at.At(index); atom != expectedAtom {
		t.Errorf("AtomTable.At(%v) => %q, want %q", index, atom, expectedAtom)
	}
}

func BenchmarkAtomTableSizeT(b *testing.B) {
	at := newBigAtomTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at.sizeT()
	}
}

func BenchmarkAtomTableSizeA(b *testing.B) {
	at := newBigAtomTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at.sizeA()
	}
}

func BenchmarkAtomTableSizeH(b *testing.B) {
	at := newBigAtomTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at.sizeH()
	}
}

func newBigAtomTable() *AtomTable {
	at := NewAtomTable()
	for i := 0; i < 1<<(2*oneByte); i++ {
		at.Put("atom" + strconv.Itoa(i))
	}
	return at
}
