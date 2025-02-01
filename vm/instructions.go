package vm

import "fmt"

//lint:file-ignore ST1020 documentation for instructions is in the book

type InstructionHandler func(vm *VM, args []byte)

type opcode struct {
	handler InstructionHandler
	length  int
}

// region Data copying instructions

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

// endregion

// region Arithmetic and logic instructions

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

// endregion

// region Comparison and conditional jumps instructions

// compare
func VCMP(vm *VM, args []byte) {
	rdst := &vm.reg[args[0]]
	rsrc := &vm.reg[args[1]]

	result := int(rdst.value) - int(rsrc.value)
	vm.fr &= 0xfffffffc

	vm.debugLog("==> VCMP: %v - %v = %v\n", rdst.value, rsrc.value, result)

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

	diff := uint32(args[0]) | uint32(args[1])<<8

	vm.debugLog("==> VJZ: conditional jump by diff: %x\n", diff)

	// 2 bytes as args is an imm16, little-endian
	// For example, imm16 is 2137, is 0x859, is 100001011001. In big endian it's written like:
	// 00001000 01011001
	// 87654321 87654321
	// But in little-endian (as we receive it) would be:
	// 01011001 00001000
	// 87654321 87654321

	// 01011001 00001000 -> 00001000 01011001 ???

	if vm.fr&FlagZF == FlagZF {
		vm.debugLog("==> VJZ: condition true, increased pc by", diff)
		vm.pc.value = vm.pc.value + diff
	} else {
		vm.debugLog("==> VJZ: condition false, no-op")
	}
}

// jump if equal
func VJE(vm *VM, args []byte) {
	VJZ(vm, args)
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

// endregion

// region Stack manipulation instructions

// push
func VPUSH(vm *VM, args []byte) {
	// TODO: Implement VPUSH
}

// pop
func VPOP(vm *VM, args []byte) {
	// TODO: Implement VPOP
}

// endregion

// region Unconditional jumps instructions

// jump
func VJMP(vm *VM, args []byte) {
	if len(args) != 2 {
		vm.interrupt(IntMemoryError)
		return
	}

	diff := uint32(args[0]) | uint32(args[1])<<8

	vm.debugLog("==> VJMP: unconditional jump by diff: %v\n", diff)

	// Example: VJMP is at address 0x13, we want to jump to 0x30
	// * VJMP opcode: 0x40
	// * Address of instruction directly after VJMP is: 0x13 + 1 + 2 = 0x16
	// 0x30 - 0x16 = 48 - 22 = 26 = 0x1A

	// 40 1A 00
	// which means:
	// VJMP 0x13 + 3 + 0x1a

	vm.pc.value = vm.pc.value + diff
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

// endregion

// region Additional instructions

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
	// TODO: Dumb implementation. A proper one should be interrupt-based.
	//  See also: https://github.com/gynvael/zrozumiec-programowanie/blob/master/007-Czesc_II-Rozdzial_3-Podstawy_architektury_komputerowe/vm_dev_con.py

	rsrc := args[0]
	value := vm.reg[rsrc].value & 0xff
	toWrite := fmt.Sprintf("%c", value)
	_, err := vm.Stdout.Write([]byte(toWrite))
	if err != nil {
		vm.interrupt(IntMemoryError)
	}
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
func VCRSH(vm *VM, _ []byte) {
	vm.crash()
}

// power off
func VOFF(vm *VM, _ []byte) {
	vm.terminated = true
}

func asd() opcode {
	var asdds = [][]string{
		[]string{"asd"},
	}
}

var opcodes = map[byte]opcode{
	// data copying instructions
	0x00: {handler: VMOV, length: 1 + 1},
	0x01: {handler: VSET, length: 1 + 4},
	0x02: {handler: VLD, length: 1 + 1},
	0x03: {handler: VST, length: 1 + 1},
	0x04: {handler: VLDB, length: 1 + 1},
	0x05: {handler: VSTB, length: 1 + 1},
	// arithmetic and logic instructions
	0x10: {handler: VADD, length: 1 + 1},
	0x11: {handler: VSUB, length: 1 + 1},
	0x12: {handler: VMUL, length: 1 + 1},
	0x13: {handler: VDIV, length: 1 + 1},
	0x14: {handler: VMOD, length: 1 + 1},
	0x15: {handler: VOR, length: 1 + 1},
	0x16: {handler: VAND, length: 1 + 1},
	0x17: {handler: VXOR, length: 1 + 1},
	0x18: {handler: VNOT, length: 1},
	0x19: {handler: VSHL, length: 1 + 1},
	0x1A: {handler: VSHR, length: 1 + 1},
	// comparison and conditional jumps instructions
	0x20: {handler: VCMP, length: 1 + 1},
	0x21: {handler: VJZ, length: 2},
	0x22: {handler: VJNZ, length: 2},
	0x23: {handler: VJC, length: 2},
	0x24: {handler: VJNC, length: 2},
	0x25: {handler: VJBE, length: 2},
	0x26: {handler: VJA, length: 2},
	// stack manipulation instructions
	0x30: {handler: VPUSH, length: 1},
	0x31: {handler: VPOP, length: 1},
	// unconditional jumps instructions
	0x40: {handler: VJMP, length: 2},
	0x41: {handler: VJMPR, length: 1},
	0x42: {handler: VCALL, length: 2},
	0x43: {handler: VCALLR, length: 1},
	0x44: {handler: VRET, length: 0},
	// additional instructions
	0xF0: {handler: VCRL, length: 1 + 2},
	0xF1: {handler: VCRS, length: 1 + 2},
	0xF2: {handler: VOUTB, length: 1 + 1},
	0xF3: {handler: VINB, length: 1 + 1},
	0xF4: {handler: VIRET, length: 0},
	0xFE: {handler: VCRSH, length: 0},
	0xFF: {handler: VOFF, length: 0},
}

// endregion
