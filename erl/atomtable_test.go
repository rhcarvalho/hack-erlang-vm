package erl

import (
	"strconv"
	"strings"
	"testing"
)

func TestAddFirstManyAtomsIncreaseSize(t *testing.T) {
	at := NewAtomTable()
	checkSize(t, at, 0)
	for i := 0; i < 1<<oneByte; i++ {
		at.Add("atom" + strconv.Itoa(i))
		checkSize(t, at, uint32(i+1))
	}
}

func TestAddAtomTwice(t *testing.T) {
	at := NewAtomTable()
	at.Add("atom")
	at.Add("atom")
	checkSize(t, at, 1)
	checkOffset(t, at, "atom", 4)
}

func TestTableHasAtom(t *testing.T) {
	at := NewAtomTable()
	checkHas(t, at, "atom", false)
	at.Add("atom")
	checkHas(t, at, "atom", true)
}

func TestAtomOffset(t *testing.T) {
	at := NewAtomTable()
	at.Add("atom")
	at.Add("second")
	checkOffset(t, at, "atom", 4)
	checkOffset(t, at, "second", 9)
}

func TestAtomAt(t *testing.T) {
	at := NewAtomTable()
	at.Add("atom")
	at.Add("second")
	checkAt(t, at, 4, "atom")
	checkAt(t, at, 9, "second")
	checkAt(t, at, 127, "")
	checkAt(t, at, 7, "m\x06second")
}

func TestNthAtom(t *testing.T) {
	at := NewAtomTable()
	at.Add("atom")
	at.Add("second")
	checkNth(t, at, 0, "atom")
	checkNth(t, at, 1, "second")
	checkNth(t, at, 2, "")
	checkNth(t, at, -1, "second")
	checkNth(t, at, -2, "atom")
	checkNth(t, at, -3, "")
}

func TestString(t *testing.T) {
	at := NewAtomTable()
	for _, atom := range []string{"fac", "state", "erlang", "-", "*", "module_info", "get_module_info"} {
		at.Add(atom)
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

func checkSize(t *testing.T, at *AtomTable, expectedSize uint32) {
	if size := at.Size(); size != expectedSize {
		t.Errorf("AtomTable.Size() => %v, want %v", size, expectedSize)
	}
}

func checkHas(t *testing.T, at *AtomTable, atom string, expectedHas bool) {
	if has := at.Has(atom); has != expectedHas {
		t.Errorf("AtomTable.Has(%q) => %v, want %v", atom, has, expectedHas)
	}
}

func checkOffset(t *testing.T, at *AtomTable, atom string, expectedOffset uint32) {
	if offset := at.Offset(atom); offset != expectedOffset {
		t.Errorf("AtomTable.Offset(%q) => %v, want %v", atom, offset, expectedOffset)
	}
}

func checkAt(t *testing.T, at *AtomTable, offset uint32, expectedAtom string) {
	if atom := at.At(offset); atom != expectedAtom {
		t.Errorf("AtomTable.At(%v) => %q, want %q", offset, atom, expectedAtom)
	}
}

func checkNth(t *testing.T, at *AtomTable, n int32, expectedAtom string) {
	if atom := at.Nth(n); atom != expectedAtom {
		t.Errorf("AtomTable.Nth(%v) => %q, want %q", n, atom, expectedAtom)
	}
}

func BenchmarkSizeT(b *testing.B) {
	at := newBigAtomTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at.sizeT()
	}
}

func BenchmarkSizeA(b *testing.B) {
	at := newBigAtomTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at.sizeA()
	}
}

func BenchmarkSizeH(b *testing.B) {
	at := newBigAtomTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at.sizeH()
	}
}

func newBigAtomTable() *AtomTable {
	at := NewAtomTable()
	for i := 0; i < 1<<(2*oneByte); i++ {
		at.Add("atom" + strconv.Itoa(i))
	}
	return at
}
