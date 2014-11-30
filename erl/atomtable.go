package erl

import (
	"bytes"
	"fmt"
)

const oneByte = 8

// An AtomTable represents an Erlang VM Atom Table.
type AtomTable struct {
	t []byte
	a []uint32
	h map[string]uint32
}

// String returns the AtomTable as an Erlang Bit String.
func (at *AtomTable) String() string {
	var b bytes.Buffer
	fmt.Fprint(&b, "<<")
	for i, v := range at.t {
		if i > 0 {
			fmt.Fprint(&b, ",")
		}
		fmt.Fprint(&b, v)
	}
	fmt.Fprint(&b, ">>")
	return b.String()
}

// NewAtomTable returns a properly initialized AtomTable.
func NewAtomTable() *AtomTable {
	return &AtomTable{t: []byte{0, 0, 0, 0}, h: make(map[string]uint32)}
}

// Put adds an atom to the AtomTable and returns its index.
//
// Putting the same atom again returns the index of the original atom.
func (at *AtomTable) Put(atom string) uint32 {
	if at.Has(atom) {
		return at.h[atom]
	}
	index := at.incSize() - 1
	offset := uint32(len(at.t))
	at.h[atom] = index
	at.a = append(at.a, offset)
	at.t = append(append(at.t, byte(len(atom))), []byte(atom)...)
	return index
}

func (at *AtomTable) incSize() uint32 {
	newSize := at.Size() + 1
	copy(at.t[:4], []byte{
		byte(newSize >> (3 * oneByte)),
		byte(newSize >> (2 * oneByte)),
		byte(newSize >> (1 * oneByte)),
		byte(newSize >> (0 * oneByte))})
	return newSize
}

// Size returns the number of atoms in the AtomTable.
func (at *AtomTable) Size() uint32 {
	return at.sizeA()
}

func (at *AtomTable) sizeT() uint32 {
	return uint32(at.t[0])<<(3*oneByte) +
		uint32(at.t[1])<<(2*oneByte) +
		uint32(at.t[2])<<(1*oneByte) +
		uint32(at.t[3])<<(0*oneByte)
}

func (at *AtomTable) sizeA() uint32 {
	return uint32(len(at.a))
}

func (at *AtomTable) sizeH() uint32 {
	return uint32(len(at.h))
}

// Has returns true if the AtomTable contains the atom, false otherwise.
func (at *AtomTable) Has(atom string) bool {
	_, has := at.h[atom]
	return has
}

// Offset returns the offset of the atom in the AtomTable.
//
// Returns 0 when the AtomTable does not contain the atom.
func (at *AtomTable) Offset(atom string) uint32 {
	index, ok := at.h[atom]
	if ok {
		return at.a[index]
	}
	return 0
}

// At returns the atom at the given offset of the AtomTable.
//
// Returns garbage if offset is not valid.
func (at *AtomTable) At(offset uint32) string {
	m := uint32(len(at.t))
	offset = offset % m
	atomSize := uint32(at.t[offset])
	return string(at.t[min(offset+1, m):min(offset+1+atomSize, m)])
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// Nth returns the nth atom in the AtomTable.
//
// Returns the empty string when n is out of bounds.
func (at *AtomTable) Nth(n int32) string {
	if n < 0 {
		n += int32(len(at.a))
	}
	if n >= int32(len(at.a)) || n < 0 {
		return ""
	}
	return at.At(at.a[n])
}
