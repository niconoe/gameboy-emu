package memory

import (
	"github.com/niconoe/gameboy-emu/types"
    "io/ioutil"
)


// This method should be called to create a Mmu, it ensures it is set up correctly
func MakeMmu() Mmu {
    mmu := Mmu{}

    // We have to load/initialize the BIOS data:

    // TODO: get rid of this absolute path !!!
    data, err := ioutil.ReadFile("/Users/nicolasnoe/Dropbox/go/src/github.com/niconoe/gameboy-emu/memory/bios.bin")
    if err != nil {
        // TODO: look how to use panic to display proper error message
        panic(err)
    }

    mmu.bios_data = data

    return mmu
}


// Mmu should not be instanciated directly.
// It should be instead instanciated with the MakeMmu() method
type Mmu struct {
    bios_data []byte
}

func (Mmu) readByte(addr types.MemoryAddress) byte {
	return 0x00
}

func (Mmu) readWord(addr types.MemoryAddress) types.Word {
	return 0xffff
}

func (Mmu) writeByte(addr types.MemoryAddress, val byte) {

}

func (Mmu) writeWord(addr types.MemoryAddress, val types.Word) {

}
