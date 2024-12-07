package vm

import "fmt"

//lint:file-ignore ST1020 documentation for instructions is in the book

type InstructionHandler func(vm *VM, args []byte)

//region Data copying instructions

type opcode struct {
	handler  InstructionHandler
	length   int
	mnemonic string
}

// mov
func VMOV(vm *VM, args []byte) {
	vm.reg[args[0]].value = vm.reg[args[1]].value
}

// set
func VSET(vm *VM, args []byte) {
	var value uint32
	value = uint32(args[1])
	value = value | uint32(args[2])<<8
	value = value | uint32(args[3])<<16
	value = value | uint32(args[4])<<24

	vm.reg[args[0]].value = value
}

// load
func VLD(vm *VM, args []byte) {
	addr := vm.reg[args[1]].value
	data, err := vm.memory.FetchDword(uint16(addr))
	if err != nil {
		vm.interrupt(IntMemoryError)
		return
	}
	vm.reg[args[0]].value = data
}

// store
func VST(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	err := vm.memory.StoreDword(uint16(rdst.value), rsrc.value)
	if err != nil {
		vm.interrupt(IntMemoryError)
	}
}

// load byte
func VLDB(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	b, err := vm.memory.FetchByte(uint16(rsrc.value))
	if err != nil {
		vm.interrupt(IntMemoryError)
		return
	}

	rdst.value = uint32(b)
}

// store byte
func VSTB(vm *VM, args []byte) {
}

//endregion

//region Arithmetic and logic instructions

// add
func VADD(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value + rsrc.value
}

// subtract
func VSUB(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value - rsrc.value
}

// multiply
func VMUL(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value * rsrc.value
}

// divide
func VDIV(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]

	if rsrc.value == 0 {
		vm.interrupt(IntDivisionError)
		return
	}

	rdst.value = rdst.value / rsrc.value
}

// modulo
func VMOD(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value % rsrc.value
}

// or
func VOR(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value | rsrc.value
}

// and
func VAND(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value & rsrc.value
}

// xor
func VXOR(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]
	rdst.value = rdst.value ^ rsrc.value
}

// not
func VNOT(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rdst.value = ^rdst.value
}

// shift left
func VSHL(vm *VM, args []byte) {
	// TODO: Implement VSHL
}

// shift right
func VSHR(vm *VM, args []byte) {
	// TODO: Implement VSHR
}

//endregion

//region Comparison and conditional jumps instructions

// compare
func VCMP(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]

	result := int(rdst.value) - int(rsrc.value)
	vm.fr &= 0xfffffffc

	if result == 0 {
		vm.fr |= FlagZF
	} else if result < 0 {
		vm.fr |= FlagCF
	}
}

// jump if zero
func VJZ(vm *VM, args []byte) {
	if len(args) != 2 {
		vm.interrupt(IntMemoryError)
		return
	}

	// 2 bytes as args is an imm16, little-endian
	// For example, imm16 is 2137, is 0x859, is 100001011001. In big endian it's written like:
	// 00001000 01011001
	// 87654321 87654321
	// But in little-endian (as we receive it) would be:
	// 01011001 00001000
	// 87654321 87654321

	// 01011001 00001000 -> 00001000 01011001 ???

	if vm.fr&FlagZF == FlagZF {
		diff := uint32(args[1]) + uint32(args[0])<<8
		vm.pc.value = vm.pc.value + diff
	}
}

// jump if equal
func VJE(vm *VM, args []byte) {
	// TODO: Implement VJE
}

// jump if not zero
func VJNZ(vm *VM, args []byte) {
	// TODO: Implement VJNZ
}

func VJNE(vm *VM, args []byte) {
	// TODO: implement VJNE
}

// jump if carry
func VJC(vm *VM, args []byte) {
	// TODO: Implement VJC
}

// jump if below
func VJB(vm *VM, args []byte) {
	// TODO: Implement VJB
}

// jump if not carry
func VJNC(vm *VM, args []byte) {
	// TODO: Implement VJNC
}

// jump if above or equal
func VJAE(vm *VM, args []byte) {
	// TODO: Implement VJAE
}

// jump if below or equal
func VJBE(vm *VM, args []byte) {
	// TODO: Implement VJBE
}

// jump if above
func VJA(vm *VM, args []byte) {
	// TODO: Implement VJA
}

//endregion

//region Stack manipulation instructions

// push
func VPUSH(vm *VM, args []byte) {
	// TODO: Implement VPUSH
}

// pop
func VPOP(vm *VM, args []byte) {
	// TODO: Implement VPOP
}

//endregion

//region Unconditional jumps instructions

// jump
func VJMP(vm *VM, args []byte) {
	if len(args) != 2 {
		vm.interrupt(IntMemoryError)
		return
	}

	diff := uint32(args[1]) | uint32(args[0])<<8
	// diff := uint32(args[0]) | uint32(args[1])<<8

	if vm.debug {
		fmt.Printf("jump by diff: %v\n", diff)
	}

	// Example: VJMP is at address 0x13, we want to jump to 0x30
	// * VJMP opcode: 0x40
	// * Address of instruction directly after VJMP is: 0x13 + 1 + 2 = 0x16
	// 0x30 - 0x16 = 48 - 22 = 26 = 0x1A

	// 40 1A 00
	// which means:
	// VJMP 0x13 + 3 + 0x1a

	vm.pc.value = vm.pc.value + 3 + diff
}

// jump to address from register
func VJMPR(vm *VM, args []byte) {
	// TODO: Implement VJMPR
}

// call
func VCALL(vm *VM, args []byte) {
	// TODO: Implement VCALL
}

// call an address from register
func VCALLR(vm *VM, args []byte) {
	// TODO: Implement VCALLR
}

// return
func VRET(vm *VM, args []byte) {
	// TODO: Implement VRET
}

//endregion

//region Additional instructions

// control register load
func VCRL(vm *VM, args []byte) {
	// TODO: Implement VCRL
}

// control register store
func VCRS(vm *VM, args []byte) {
	// TODO: Implement VCRS
}

// output byte
func VOUTB(vm *VM, args []byte) {
	// TODO: Implement VOUTB
}

// input byte
func VINB(vm *VM, args []byte) {
	// TODO: Implement VINB
}

// interrupt return
func VIRET(vm *VM, args []byte) {
	// TODO: Implement VIRET
}

// crash
func VCRSH(vm *VM, args []byte) {
	vm.crash()
}

// power off
func VOFF(vm *VM, args []byte) {
	vm.terminated = true
}

var opcodes = map[byte]opcode{
	// data copying instructions
	0x00: {handler: VMOV, length: 1 + 1, mnemonic: "MOV"},
	0x01: {handler: VSET, length: 1 + 4, mnemonic: "SET"},
	0x02: {handler: VLD, length: 1 + 1, mnemonic: "LD"},
	0x03: {handler: VST, length: 1 + 1, mnemonic: "ST"},
	0x04: {handler: VLDB, length: 1 + 1, mnemonic: "LDB"},
	0x05: {handler: VSTB, length: 1 + 1, mnemonic: "STB"},
	// arithmetic and logic instructions
	0x10: {handler: VADD, length: 1 + 1, mnemonic: "ADD"},
	0x11: {handler: VSUB, length: 1 + 1, mnemonic: "SUB"},
	0x12: {handler: VMUL, length: 1 + 1, mnemonic: "MUL"},
	0x13: {handler: VDIV, length: 1 + 1, mnemonic: "DIV"},
	0x14: {handler: VMOD, length: 1 + 1, mnemonic: "MOD"},
	0x15: {handler: VOR, length: 1 + 1, mnemonic: "OR"},
	0x16: {handler: VAND, length: 1 + 1, mnemonic: "AND"},
	0x17: {handler: VXOR, length: 1 + 1, mnemonic: "XOR"},
	0x18: {handler: VNOT, length: 1, mnemonic: "NOT"},
	0x19: {handler: VSHL, length: 1 + 1, mnemonic: "SHL"},
	0x1A: {handler: VSHR, length: 1 + 1, mnemonic: "SHR"},
	// comparison and conditional jumps instructions
	0x20: {handler: VCMP, length: 1 + 1, mnemonic: "CMP"},
	0x21: {handler: VJZ, length: 2, mnemonic: "JZ"},
	0x22: {handler: VJNZ, length: 2, mnemonic: "JNZ"},
	0x23: {handler: VJC, length: 2, mnemonic: "JC"},
	0x24: {handler: VJNC, length: 2, mnemonic: "JNC"},
	0x25: {handler: VJBE, length: 2, mnemonic: "JBE"},
	0x26: {handler: VJA, length: 2, mnemonic: "JA"},
	// stack manipulation instructions
	0x30: {handler: VPUSH, length: 1, mnemonic: "PUSH"},
	0x31: {handler: VPOP, length: 1, mnemonic: "POP"},
	// unconditional jumps instructions
	0x40: {handler: VJMP, length: 2, mnemonic: "JMP"},
	0x41: {handler: VJMPR, length: 1, mnemonic: "JMPR"},
	0x42: {handler: VCALL, length: 2, mnemonic: "CALL"},
	0x43: {handler: VCALLR, length: 1, mnemonic: "CALLR"},
	0x44: {handler: VRET, length: 0, mnemonic: "RET"},
	// additional instructions
	0xF0: {handler: VCRL, length: 1 + 2, mnemonic: "CRL"},
	0xF1: {handler: VCRS, length: 1 + 2, mnemonic: "CRS"},
	0xF2: {handler: VOUTB, length: 1 + 1, mnemonic: "OUTB"},
	0xF3: {handler: VINB, length: 1 + 1, mnemonic: "INB"},
	0xF4: {handler: VIRET, length: 0, mnemonic: "IRET"},
	0xFE: {handler: VCRSH, length: 0, mnemonic: "CRSH"},
	0xFF: {handler: VOFF, length: 0, mnemonic: "OFF"},
}

//endregion
