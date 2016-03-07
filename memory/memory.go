package memory

import (
	"fmt"
	"github.com/niconoe/gameboy-emu/types"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// This method should be called to create a Mmu, it ensures it is set up correctly
func MakeMmu() Mmu {
	mmu := Mmu{}

	// We have to load/initialize the BIOS data:
	// TODO: get rid of this absolute path !!!
	data, err := ioutil.ReadFile("/Users/nicolasnoe/Dropbox/go/src/github.com/niconoe/gameboy-emu/memory/bios.bin")
	check(err)

	// At initialization time, the BIOS is mapped
	mmu.biosIsMapped = true
	mmu.biosData = data

	return mmu
}

// Mmu should not be instanciated directly.
// It should be instead instanciated with the MakeMmu() method
type Mmu struct {
	biosData     []byte
	romBank0     [16384]byte
	otherRomBank [16384]byte

	biosIsMapped bool
}

func (mmu *Mmu) LoadRom(romPath string) {
	rom, err := os.Open(romPath)
	check(err)
	defer rom.Close()

	// TODO: Throw error message is ROM type is not supported
	// (ROM largers than 32kb?)

	// Bank 0 of ROM is always available at 0000-3fff
	rom.Read(mmu.romBank0[:])

	// We currently only support 32k ROMS (without MBC chips), so Bank 1 is directly mapped at 4000-7fff
	rom.ReadAt(mmu.otherRomBank[:], 16384)
}

func (mmu Mmu) ReadByte(addr types.MemoryAddress) byte {
	if addr <= 0x3fff { // Rom Bank 0
		if addr <= 0x00ff && mmu.biosIsMapped {
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
	b2 := mmu.ReadByte(addr + 1)

	return types.WordFromBytes(b1, b2)
}

func (Mmu) WriteByte(addr types.MemoryAddress, val byte) {
	fmt.Printf("Write byte %.2x to Address: %x", val, addr)

	if addr >= 0x8000 && addr <= 0x9fff {
		fmt.Printf(" (It's VRAM!) \n")
	}

	if addr >= 0xff00 && addr <= 0xff7f {
		fmt.Printf(" (It's Memory-mapped IO!) \n")
	}
}

func (Mmu) WriteWord(addr types.MemoryAddress, val types.Word) {

}
