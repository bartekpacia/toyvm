package vm

//lint:file-ignore ST1020 documentation for instructions is in the book

type InstructionHandler func(vm *VM, args []byte)

type opcode struct {
	handler InstructionHandler
	length  int
}

// data copying instructions

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
	dst := vm.reg[args[0]].value
	src := vm.reg[args[1]].value

	vm.memory.StoreDword(uint16(dst), src)
}

// load byte
func VLDB(vm *VM, args []byte) {
}

// store byte
func VSTB(vm *VM, args []byte) {
}

// arithmetic and logic instructions

// add
func VADD(vm *VM, args []byte) {
}

// subtract
func VSUB(vm *VM, args []byte) {
}

// multiply
func VMUL(vm *VM, args []byte) {
}

// divide
func VDIV(vm *VM, args []byte) {
}

// modulo
func VMOD(vm *VM, args []byte) {
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
}

// shift right
func VSHR(vm *VM, args []byte) {
}

// comparison and conditional jumps instructions

// compare
func VCMP(vm *VM, args []byte) {
}

// jump if zero
func VJZ(vm *VM, args []byte) {
}

// jump if equal
func VJE(vm *VM, args []byte) {
}

// jump if not zero
func VJNZ(vm *VM, args []byte) {
}

// jump if not equal
func VJNE(vm *VM, args []byte) {
}

// jump if carry
func VJC(vm *VM, args []byte) {
}

// jump if below
func VJB(vm *VM, args []byte) {
}

// jump if not carry
func VJNC(vm *VM, args []byte) {
}

// jump if above or equal
func VJAE(vm *VM, args []byte) {
}

// jump if below or equal
func VJBE(vm *VM, args []byte) {
}

// jump if above
func VJA(vm *VM, args []byte) {
}

// stack manipulation instructions

// push
func VPUSH(vm *VM, args []byte) {
}

// pop
func VPOP(vm *VM, args []byte) {
}

// unconditional jumps instructions

// jump
func VJMP(vm *VM, args []byte) {
}

// jump to address from register
func VJMPR(vm *VM, args []byte) {
}

// call
func VCALL(vm *VM, args []byte) {
}

// call an address from register
func VCALLR(vm *VM, args []byte) {
}

// return
func VRET(vm *VM, args []byte) {
}

// additional instructions

// control register load
func VCRL(vm *VM, args []byte) {
}

// control register store
func VCRS(vm *VM, args []byte) {
}

// output byte
func VOUTB(vm *VM, args []byte) {
}

// input byte
func VINB(vm *VM, args []byte) {
}

// interrupt return
func VIRET(vm *VM, args []byte) {
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
