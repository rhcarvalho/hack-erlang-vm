package erl

import "testing"

type Sizer interface {
	Size() uint32
}

func checkSize(t *testing.T, s Sizer, expectedSize uint32) {
	if size := s.Size(); size != expectedSize {
		t.Errorf("%T.Size() => %v, want %v", s, size, expectedSize)
	}
}
