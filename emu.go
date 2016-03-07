package main

import (
	"fmt"
	"github.com/niconoe/gameboy-emu/cpu"
	"github.com/niconoe/gameboy-emu/memory"
	"os"
)

func main() {
	fmt.Println("Niconoe's experimental Gameboy emulator...")
	if len(os.Args) != 2 {
		fmt.Println("No ROM path given, aborting.")
	} else {
		romPath := os.Args[1]

		var cpu cpu.GameboyCPU
		var mmu = memory.MakeMmu()
		mmu.LoadRom(romPath)

		cpu.AttachMmu(&mmu)
		cpu.Run()
	}

}
