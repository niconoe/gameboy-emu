package types

import "encoding/binary"

type MemoryAddress uint16
type Word uint16

// Returns a MemoryAddress whose initial value is the base address of the memory-mapped IOs
// (0xff00)
func MakeMemoryMappedIOAddress() MemoryAddress {
	var a MemoryAddress

	a = 0xff00
	return a
}

func (m MemoryAddress) AddSignedOffset(offset byte) MemoryAddress {
	// Returns a new address with offset added
	// Since offset is a signed 8-bit integer, the new adress is in the range
	// [oldaddress-128, oldaddress+127]
	return m + MemoryAddress(int8(offset))
}

func (m MemoryAddress) AddUnsignedOffset(offset byte) MemoryAddress {
	return m + MemoryAddress(uint8(offset))
}

func (w Word) ToBytes() (msb, lsb byte) {
	msb = byte(w >> 8)
	lsb = byte(w & 0xff)
	return
}

// Byte order: little-endian, as they appear in memory
func WordFromBytes(b1, b2 byte) Word {
	array := []byte{b1, b2}
	return Word(binary.LittleEndian.Uint16(array))
}
