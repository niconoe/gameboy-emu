package cpu

import (
	"github.com/niconoe/gameboy-emu/memory"
	"github.com/niconoe/gameboy-emu/types"
)

type clock struct {
	m uint64 // Machine cycles
	t uint64 // Clock cycles
}

type GameboyCPU struct {
	// Registers

	a, b, c, d, e, h, l, f byte   // 8-bit registers
	pc, sp                 uint16 // 16-bit registers: program counter and stack pointer

	lastInstructionClock clock

	mmu *memory.Mmu
}

func (cpu GameboyCPU) AttachMmu(mmu *memory.Mmu) {
	cpu.mmu = mmu
}

func (cpu GameboyCPU) Reset() {
	cpu.a = 0x00
	cpu.b = 0x00
	cpu.c = 0x00
	cpu.d = 0x00
	cpu.e = 0x00
	cpu.h = 0x00
	cpu.l = 0x00
	cpu.f = 0x00

	cpu.pc = 0x0000
	cpu.sp = 0x0000

	cpu.lastInstructionClock.m = 0
	cpu.lastInstructionClock.t = 0
}

func (cpu GameboyCPU) dispatch(opcode byte) {
	switch opcode {
	case 0x00:
		cpu.nop()
	case 0x01:
		//cpu.ldBCWord()
	}

}

// TODO: Move to types (method od Word type)
func wordToBytes(w types.Word) (msb, lsb byte) {
	msb = byte(w >> 8)
	lsb = byte(w & 0xff)
	return
}

// Instructions
func (cpu GameboyCPU) nop() {
	cpu.lastInstructionClock.m = 1
	cpu.lastInstructionClock.t = 4
}

func (cpu GameboyCPU) ldBCWord(w types.Word) {
	cpu.b, cpu.c = wordToBytes(w)

	cpu.lastInstructionClock.m = 3
	cpu.lastInstructionClock.t = 12
}
