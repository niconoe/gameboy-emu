package cpu

import (
	"github.com/niconoe/gameboy-emu/memory"
	"github.com/niconoe/gameboy-emu/types"

	"fmt"
)

type clock struct {
	t uint64 // Clock cycles
}

type GameboyCPU struct {
	// Registers

	a, b, c, d, e, h, l, f byte                // 8-bit registers
	pc, sp                 types.MemoryAddress // 16-bit registers: program counter and stack pointer

	lastInstructionClock clock

	mmu *memory.Mmu
}

func (cpu GameboyCPU) String() string {
	return fmt.Sprintf("--------------------------\n"+
		"AF: 0x%.2x%.2x \n"+
		"BC: 0x%.2x%.2x \n"+
		"DE: 0x%.2x%.2x \n"+
		"HL: 0x%.2x%.2x \n"+
		"\n"+
		"SP: 0x%.4x\n"+
		"PC: 0x%.4x\n"+
		"--------------------------\n", cpu.a, cpu.f, cpu.b, cpu.c, cpu.d, cpu.e, cpu.h, cpu.l, cpu.sp, cpu.pc)
}

// This function will need a Mmu attached
func (cpu GameboyCPU) Run() {
	cpu.Reset()

	for {
		//fmt.Println(cpu)

		// Get the opcode of the next instruction
		opcode, extended_opcode := cpu.FetchNextOpcode()

		/*fmt.Printf("\nOpcode: %.2x", opcode)
		if opcode == 0xcb {
			fmt.Printf(" -- Extended opcode: %.2x", extended_opcode)
		}
		fmt.Printf("\n")*/
		cpu.Execute(opcode, extended_opcode)

		// For debugging purposes...
		//time.Sleep(100 * time.Millisecond)
	}
}

func (cpu *GameboyCPU) AttachMmu(mmu *memory.Mmu) {
	cpu.mmu = mmu
}

func (cpu *GameboyCPU) Reset() {
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

	cpu.lastInstructionClock.t = 0
}

func (cpu *GameboyCPU) Execute(opcode, extended_opcode byte) {
	if opcode != 0xcb {
		cpu.dispatch(opcode)
	} else {
		// This is an extended instruction, we pass the second part of the opcode
		cpu.dispatchExtended(extended_opcode)
	}

}

// If normal instruction, only the first byte is significant and the second one should be ignored
// If extended instruction, both bytes are returned
func (cpu *GameboyCPU) FetchNextOpcode() (first, second byte) {
	opcode := cpu.mmu.ReadByte(cpu.pc)
	if opcode == 0xcb {
		return 0xcb, cpu.mmu.ReadByte(cpu.pc + 1)
	} else {
		return opcode, 0x00
	}
}

func (cpu *GameboyCPU) dispatch(opcode byte) {
	// When we enter here, PC is always pointing right before the current instruction
	switch opcode {
	case 0x00:
		cpu.nop()
	case 0x01:
		cpu.ldBCd16()
	case 0x06:
		cpu.ldBd8()
	case 0x0c:
		cpu.incC()
	case 0x0e:
		cpu.ldCd8()
	case 0x11:
		cpu.ldDEd16()
	case 0x17:
		cpu.rlA()
	case 0x1a:
		cpu.ldAParDEPar()
	case 0x20:
		cpu.jrNzR8()
	case 0x21:
		cpu.ldHLd16()
	case 0x31:
		cpu.ldSPd16()
	case 0x32:
		cpu.LDParHL_ParA()
	case 0x3e:
		cpu.ldAd8()
	case 0x4f:
		cpu.ldCA()
	case 0x77:
		cpu.ldParHLParA()
	case 0xaf:
		cpu.xorA()
	case 0xc1:
		cpu.popBC()
	case 0xc5:
		cpu.pushBC()
	case 0xcd:
		cpu.callA16()
	case 0xe0:
		cpu.LDHPara8ParA()
	case 0xe2:
		cpu.ldParCParA()

	default:
		panic(fmt.Sprintf("Opcode not found: %.2x. CPU state: %s", opcode, cpu))
	}
}

func (cpu *GameboyCPU) dispatchExtended(secondByteOfOpcode byte) {
	// We receive only the second byte of the opcode, but like for the normal dispatch,
	// PC is till pointing before the whole (2 bytes) instruction
	switch secondByteOfOpcode {
	case 0x11:
		cpu.rl_C()
	case 0x7c:
		cpu.bit7H()
	default:
		panic(fmt.Sprintf("EXTENDED Opcode not found: %.2x. CPU state: %s", secondByteOfOpcode, cpu))
	}

}

// Instructions
func (cpu *GameboyCPU) rlA() {
	cpu.a = cpu.rlN(cpu.a)

	cpu.lastInstructionClock.t = 4
 	cpu.pc += 1
}

func (cpu *GameboyCPU) rl_C() {
	cpu.c = cpu.rlN(cpu.c)

	cpu.lastInstructionClock.t = 8
 	cpu.pc += 2
}

func (cpu *GameboyCPU) popBC() {
	cpu.setBC(cpu.popWordFromStack())

	cpu.lastInstructionClock.t = 12
	cpu.pc += 1

	_ = "breakpoint"
}

func (cpu *GameboyCPU) pushBC() {
	cpu.pushWordOnStack(cpu.getBC())

	cpu.lastInstructionClock.t = 16
	cpu.pc += 1
}

func (cpu *GameboyCPU) ldCA() {
	cpu.c = cpu.a

	cpu.lastInstructionClock.t = 4
	cpu.pc += 1
}

// Each instruction manipulates PC appropriately
func (cpu *GameboyCPU) callA16() {
	// CALL nn
	// Description: Push address of next instruction onto stack and then jump to address nn.
	// Use with: two byte immediate value (LS byte first).
	// Cycles: 12
	cpu.lastInstructionClock.t = 12
	
	cpu.pushWordOnStack(types.Word(cpu.pc) + 3)
	cpu.pc = types.MemoryAddress(cpu.mmu.ReadWord(cpu.pc + 1))
	// TODO: untested instruction
}

func (cpu *GameboyCPU) nop() {
	cpu.lastInstructionClock.t = 4

	cpu.pc += 1
}

func (cpu *GameboyCPU) ldBCd16() {
	w := cpu.mmu.ReadWord(cpu.pc + 1)
	cpu.setBC(w)

	cpu.lastInstructionClock.t = 12

	cpu.pc += 3
}

func (cpu *GameboyCPU) ldSPd16() {
	cpu.sp = types.MemoryAddress(cpu.mmu.ReadWord(cpu.pc + 1))

	cpu.lastInstructionClock.t = 12

	cpu.pc += 3
}

func (cpu *GameboyCPU) LDParHL_ParA() {
	// Put A into memory address HL. Decrement HL.
	// Known as: LD (HL-),A
	// Known as: LD (HLD),A or LDD (HL),A
	dest := types.MemoryAddress(cpu.getHL())
	cpu.mmu.WriteByte(dest, cpu.a)

	// Decrement HL
	cpu.decHL()

	cpu.lastInstructionClock.t = 8

	cpu.pc += 1
}

func (cpu *GameboyCPU) ldParHLParA() {
	dest := types.MemoryAddress(cpu.getHL())
	cpu.mmu.WriteByte(dest, cpu.a)

	cpu.lastInstructionClock.t = 8

	cpu.pc += 1
}

func (cpu *GameboyCPU) ldDEd16() {
	w := cpu.mmu.ReadWord(cpu.pc + 1)
	cpu.setDE(w)

	cpu.lastInstructionClock.t = 12

	cpu.pc += 3
}

func (cpu *GameboyCPU) ldHLd16() {
	w := cpu.mmu.ReadWord(cpu.pc + 1)
	cpu.setHL(w)

	cpu.lastInstructionClock.t = 12

	cpu.pc += 3
}

func (cpu *GameboyCPU) xorA() {
	// Xor A, with itself, effectively setting it to zero
	cpu.a = 0
	cpu.setZeroFlag()

	cpu.lastInstructionClock.t = 4

	cpu.pc += 1
}

func (cpu *GameboyCPU) bit7H() {
	if !hasBit(cpu.h, 7) {
		cpu.setZeroFlag()
	} else {
		cpu.clearZeroFlag()
	}

	cpu.clearSubstractFlag()
	cpu.setHalfCarryFlag()

	cpu.lastInstructionClock.t = 8

	cpu.pc += 2
}

func (cpu *GameboyCPU) jrNzR8() {
	cpu.pc += 2 // We advance it before jump, since it is relative to the next instruction

	if !cpu.hasZeroFlag() {
		cpu.pc = cpu.pc.AddSignedOffset(cpu.mmu.ReadByte(cpu.pc - 1))

		cpu.lastInstructionClock.t = 12
	} else {
		cpu.lastInstructionClock.t = 8
	}

}

func (cpu *GameboyCPU) ldBd8() {
	cpu.b = cpu.mmu.ReadByte(cpu.pc + 1)

	cpu.lastInstructionClock.t = 8

	cpu.pc += 2
}

func (cpu *GameboyCPU) ldCd8() {
	cpu.c = cpu.mmu.ReadByte(cpu.pc + 1)

	cpu.lastInstructionClock.t = 8

	cpu.pc += 2
}

func (cpu *GameboyCPU) ldAd8() {
	cpu.a = cpu.mmu.ReadByte(cpu.pc + 1)

	cpu.lastInstructionClock.t = 8

	cpu.pc += 2
}

func (cpu *GameboyCPU) ldParCParA() {
	// Put A into address $FF00 + register C.
	// Also known as LD (C),A and LD ($FF00+C),A

	dest := types.MakeMemoryMappedIOAddress()
	dest = dest.AddUnsignedOffset(cpu.c)
	cpu.mmu.WriteByte(dest, cpu.a)

	cpu.lastInstructionClock.t = 8

	cpu.pc += 1
}

func (cpu *GameboyCPU) LDHPara8ParA() {
	//Put A into memory address $FF00+n.
	dest := types.MakeMemoryMappedIOAddress()
	dest = dest.AddUnsignedOffset(cpu.mmu.ReadByte(cpu.pc + 1))
	cpu.mmu.WriteByte(dest, cpu.a)

	cpu.lastInstructionClock.t = 12

	cpu.pc += 2
}

func (cpu *GameboyCPU) incC() {
	previousVal := cpu.c
	newVal := cpu.c + 1

	cpu.c = newVal

	if cpu.c == 0 {
		cpu.setZeroFlag()
	} else {
		cpu.clearZeroFlag()
	}

	cpu.clearSubstractFlag()

	if (newVal^0x01^previousVal)&0x10 == 0x10 {
		cpu.setHalfCarryFlag()
	} else {
		cpu.clearHalfCarryFlag()
	}

	cpu.lastInstructionClock.t = 4

	cpu.pc += 1
}

func (cpu *GameboyCPU) ldAParDEPar() {
	source := types.MemoryAddress(cpu.getDE())

	cpu.a = cpu.mmu.ReadByte(source)

	cpu.lastInstructionClock.t = 8

	cpu.pc += 1
}


// Helpers

func (cpu *GameboyCPU) rlN(register byte) byte {
	// Helper for the various "rotate left through carry (8-bit)" instructions
	// Sets the flag and return new register value.

	// assigning new value to register, updating the clock and PC is the caller responsability.

	// Rotate n left through Carry flag.
	// Flags affected:
 	//  Z - Set if result is zero.
 	//  N - Reset.
 	//  H - Reset.
 	//  C - Contains old bit 7 data.

 	bit7_was_set := false
 	if hasBit(register, 7){
 		bit7_was_set = true
 	}

 	register = register << 1

 	if cpu.hasCarryFlag() { 
 		register ^= 0x01
 	}

 	if bit7_was_set {
 		cpu.setCarryFlag()
 	} else {
 		cpu.clearCarryFlag()
 	}

 	if register == 0 {
 		cpu.setZeroFlag()
 	} else {
 		cpu.clearZeroFlag()
 	}

 	cpu.clearSubstractFlag()
 	cpu.clearHalfCarryFlag()

 	return register
}

func (cpu *GameboyCPU) popWordFromStack() types.Word {
  // Push the given word on the stack and update SP
  w := cpu.mmu.ReadWord(cpu.sp)
  cpu.sp += 2

  return w
}

func (cpu *GameboyCPU) pushWordOnStack(w types.Word) {
	// Push the given word on the stack and update SP
	cpu.mmu.WriteWord(cpu.sp, w)
	cpu.sp -= 2
}

// Instruction helpers to manipulate the flags register
func (cpu *GameboyCPU) hasZeroFlag() bool {
	return hasBit(cpu.f, 7)
}

func (cpu *GameboyCPU) setZeroFlag() {
	cpu.f = setBit(cpu.f, 7)
}

func (cpu *GameboyCPU) clearZeroFlag() {
	cpu.f = clearBit(cpu.f, 7)
}

func (cpu *GameboyCPU) clearSubstractFlag() {
	cpu.f = clearBit(cpu.f, 6)
}

func (cpu *GameboyCPU) setHalfCarryFlag() {
	cpu.f = setBit(cpu.f, 5)
}

func (cpu *GameboyCPU) clearHalfCarryFlag() {
	cpu.f = clearBit(cpu.f, 5)
}

func (cpu GameboyCPU) hasCarryFlag() bool {
	return hasBit(cpu.f, 4)
}

func (cpu *GameboyCPU) setCarryFlag() {
	cpu.f = setBit(cpu.f, 4)
}

func (cpu *GameboyCPU) clearCarryFlag() {
	cpu.f = clearBit(cpu.f, 4)
}

// Instruction helpers to manipulate 8-bit registers as pairs
func (cpu GameboyCPU) getHL() types.Word {
	return types.WordFromBytes(cpu.l, cpu.h) // !WordFrom bytes expect little endian!
}

func (cpu GameboyCPU) getDE() types.Word {
	return types.WordFromBytes(cpu.e, cpu.d) // !WordFrom bytes expect little endian!
}

func (cpu GameboyCPU) getBC() types.Word {
	return types.WordFromBytes(cpu.c, cpu.b) // !WordFrom bytes expect little endian!
}

func (cpu *GameboyCPU) setHL(w types.Word) {
	cpu.h, cpu.l = w.ToBytes()
}

func (cpu *GameboyCPU) setDE(w types.Word) {
	cpu.d, cpu.e = w.ToBytes()
}

func (cpu *GameboyCPU) decHL() {
	// Decrements HL
	w := cpu.getHL()
	w = w - 1
	cpu.h, cpu.l = w.ToBytes()
}

func (cpu *GameboyCPU) setBC(w types.Word) {
	cpu.b, cpu.c = w.ToBytes()
}

// Generic instruction helpers to manipulate individual bits

// Sets the bit at pos in the byte b.
func setBit(b byte, pos uint) byte {
	b |= (1 << pos)
	return b
}

// Clears the bit at pos in b.
func clearBit(b byte, pos uint) byte {
	mask := ^(1 << pos)
	b &= byte(mask)
	return b
}

func hasBit(b byte, pos uint) bool {
	val := b & (1 << pos)
	return (val > 0)
}
