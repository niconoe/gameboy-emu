package memory

import (
	"github.com/niconoe/gameboy-emu/types"
    "io/ioutil"

    "fmt"
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

    // At initialization time, the BIOS is mapped
    mmu.biosIsMapped = true
    mmu.biosData = data

    return mmu
}


// Mmu should not be instanciated directly.
// It should be instead instanciated with the MakeMmu() method
type Mmu struct {
    biosData []byte
    romBank0 [16384]byte   

    biosIsMapped    bool
}

func (mmu Mmu) ReadByte(addr types.MemoryAddress) byte {
    if addr <= 0x3fff { // Rom Bank 0
        if addr <=0x00ff && mmu.biosIsMapped {
                // If BIOS is mapped, it shadows the cartridge ROM
                return mmu.biosData[addr]
            } else {
                return mmu.romBank0[addr]
            }
    }
    return 0x00
}

func (mmu Mmu) ReadWord(addr types.MemoryAddress) types.Word {
	b1 := mmu.ReadByte(addr)
    b2 := mmu.ReadByte(addr+1)

    return types.WordFromBytes(b1, b2)
}

func (Mmu) WriteByte(addr types.MemoryAddress, val byte) {
    if addr >= 0x8000 && addr <= 0x9fff {
        fmt.Printf("Write value %.2x to VRAM!\n", val)
    }
}

func (Mmu) WriteWord(addr types.MemoryAddress, val types.Word) {

}
