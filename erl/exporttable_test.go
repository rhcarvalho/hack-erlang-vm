package erl

import "testing"

func TestExportTablePutFirstManyFunctionsIncreaseSize(t *testing.T) {
	et := NewExportTable()
	checkSize(t, et, 0)
	for i := uint32(0); i < 1<<oneByte; i++ {
		et.Put(i, 1, 0)
		checkSize(t, et, uint32(i+1))
	}
}

func TestExportTableAddFunctionTwice(t *testing.T) {
	et := NewExportTable()
	et.Put(0, 1, 0)
	et.Put(0, 1, 1)
	checkSize(t, et, 1)
	checkExportTableLabel(t, et, 0, 1, 0)
}

func TestExportTableHasFunction(t *testing.T) {
	et := NewExportTable()
	checkExportTableHas(t, et, 0, 1, false)
	et.Put(0, 1, 0)
	checkExportTableHas(t, et, 0, 1, true)
}

func TestExportTableLabel(t *testing.T) {
	et := NewExportTable()
	et.Put(0, 1, 42)
	checkExportTableLabel(t, et, 0, 1, 42)
	checkExportTableLabel(t, et, 0, 0, 0)
}

func TestExportTableString(t *testing.T) {
	et := NewExportTable()
	et.Put(0x2021, 0x1, 0xc1f2)
	expected := "<<0,0,0,1,0,0,32,33,0,0,0,1,0,0,193,242>>"
	if str := et.String(); str != expected {
		t.Errorf("ExportTable.String() => %q, want %q", str, expected)
	}
}

func checkExportTableHas(t *testing.T, et *ExportTable, name, arity uint32, expectedHas bool) {
	if has := et.Has(name, arity); has != expectedHas {
		t.Errorf("ExportTable.Has(%v, %v) => %v, want %v", name, arity, has, expectedHas)
	}
}

func checkExportTableLabel(t *testing.T, et *ExportTable, name, arity, expectedLabel uint32) {
	if label := et.Label(name, arity); label != expectedLabel {
		t.Errorf("ExportTable.Label(%v, %v) => %v, want %v", name, arity, label, expectedLabel)
	}
}
