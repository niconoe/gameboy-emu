package main

import (
	"fmt"
	"github.com/niconoe/gameboy-emu/cpu"
	"github.com/niconoe/gameboy-emu/memory"
)

func main() {
	var cpu cpu.GameboyCPU
	var mmu = memory.MakeMmu()

	fmt.Println("Niconoe's experimental Gameboy emulator...")

	cpu.AttachMmu(&mmu)
	cpu.Reset()
}
