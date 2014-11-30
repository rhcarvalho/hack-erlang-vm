package erl

// An ExportTable represents an Erlang VM Export Table.
type ExportTable struct {
	// t holds the export table in binary format.
	// The first four bytes represent the number of exported functions.
	// The next bytes are a sequence of 12 bytes per exported function:
	// * 4 bytes for the function name in the atom index
	// * 4 bytes for the function arity
	// * 4 bytes for the function entry point label in the code table
	//
	// Empty table:              <<0,0,0,0>>
	// Put(0x2021, 0x1, 0xc1f2): <<0,0,0,1,0,0,32,33,0,0,0,1,0,0,193,242>>
	t []byte

	// h maps (name, arity) to label.
	h map[[2]uint32]uint32
}

// String returns the ExportTable as an Erlang Bit String.
func (et *ExportTable) String() string {
	return bytesToBitString(et.t)
}

// NewExportTable returns a properly initialized ExportTable.
func NewExportTable() *ExportTable {
	return &ExportTable{t: []byte{0, 0, 0, 0}, h: make(map[[2]uint32]uint32)}
}

// Put adds a function to the ExportTable.
//
// Putting the same function again has no effect.
func (et *ExportTable) Put(name, arity, label uint32) {
	if et.Has(name, arity) {
		return
	}
	et.incSize()
	et.h[[2]uint32{name, arity}] = label
	et.t = append(et.t, []byte{
		byte(name >> (3 * oneByte)),
		byte(name >> (2 * oneByte)),
		byte(name >> (1 * oneByte)),
		byte(name >> (0 * oneByte)),
		byte(arity >> (3 * oneByte)),
		byte(arity >> (2 * oneByte)),
		byte(arity >> (1 * oneByte)),
		byte(arity >> (0 * oneByte)),
		byte(label >> (3 * oneByte)),
		byte(label >> (2 * oneByte)),
		byte(label >> (1 * oneByte)),
		byte(label >> (0 * oneByte))}...)
}

func (et *ExportTable) incSize() {
	newSize := et.Size() + 1
	copy(et.t[:4], []byte{
		byte(newSize >> (3 * oneByte)),
		byte(newSize >> (2 * oneByte)),
		byte(newSize >> (1 * oneByte)),
		byte(newSize >> (0 * oneByte))})
}

// Size returns the number of exported functions in the ExportTable.
func (et *ExportTable) Size() uint32 {
	return et.sizeH()
}

func (et *ExportTable) sizeH() uint32 {
	return uint32(len(et.h))
}

// Has returns true if the ExportTable contains the function, false otherwise.
func (et *ExportTable) Has(name, arity uint32) bool {
	_, has := et.h[[2]uint32{name, arity}]
	return has
}

// Label returns the entry point label in the code table for the given function.
//
// Returns 0 when the ExportTable does not contain the pair function/arity.
func (et *ExportTable) Label(name, arity uint32) uint32 {
	return et.h[[2]uint32{name, arity}]
}
