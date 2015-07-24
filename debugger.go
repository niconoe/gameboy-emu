package main

import (
	"fmt"
	"github.com/niconoe/gameboy-emu/cpu"
	"github.com/niconoe/gameboy-emu/memory"
	"os"
)

func main() {
	var cpu cpu.GameboyCPU
	var mmu = memory.MakeMmu()

	fmt.Println("Niconoe's experimental Gameboy debugger...")
	cpu.AttachMmu(&mmu)

	for {
		fmt.Print(">> ")
		var input string
		fmt.Scanln(&input)

		switch input {
		case "n", "next":
			opcode, extended_opcode := cpu.FetchNextOpcode()

			fmt.Printf("Executing opcode: %.2x", opcode)
			if opcode == 0xcb {
				fmt.Printf(" -- Extended opcode: %.2x", extended_opcode)
			}
			fmt.Printf("\n")

			cpu.Execute(opcode, extended_opcode)

		case "s", "show":
			// Show important state (CPU, ...)
			fmt.Println(cpu)

		case "q", "quit":
			os.Exit(0)

		default:
			fmt.Println("Unknown command.")
		}

	}
}
