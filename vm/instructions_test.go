package vm

import (
	"testing"
)

func TestVmov(t *testing.T) {
	testCases := []struct {
		desc   string
		seed   func(*VM)
		args   []byte
		verify func(*VM) bool
	}{
		{
			desc:   "copy value of R4 to R1",
			seed:   func(vm *VM) { vm.reg[4].value = 1 },
			args:   []byte{1, 4},
			verify: func(vm *VM) bool { return vm.reg[1].value == 1 },
		},
	}

	for _, tc := range testCases {
		vm := NewVM()
		tc.seed(vm)
		VMOV(vm, tc.args)
		if !tc.verify(vm) {
			t.Errorf("%s failed", tc.desc)
		}
	}
}

func TestVset(t *testing.T) {
	testCases := []struct {
		args []byte
		want uint32
	}{
		{
			args: []byte{1, 0x12, 0x34, 0x00, 0x00},
			want: 0x3412,
		},
	}

	for _, tc := range testCases {
		vm := NewVM()

		VSET(vm, tc.args)
		got := vm.reg[tc.args[0]].value
		if got != tc.want {
			t.Errorf("got %x, want %x", got, tc.want)
		}
	}
}

func TestVld(t *testing.T) {
	vm := NewVM()
	vm.memory.mem[3] = 0x12
	vm.memory.mem[4] = 0x34
	vm.memory.mem[5] = 0x56
	vm.memory.mem[6] = 0x78

	vm.reg[0].value = 3

	VLD(vm, []byte{5, 0})
	got := vm.reg[5].value
	var want uint32 = 0x78563412
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestVst(t *testing.T) {
	vm := NewVM()
	vm.memory.mem[3] = 0x12
	vm.memory.mem[4] = 0x34
	vm.memory.mem[5] = 0x56
	vm.memory.mem[6] = 0x78

	var want uint32 = 0x12345678
	vm.reg[9].value = 0x1234
	vm.reg[5].value = want

	VST(vm, []byte{9, 5})
	var got uint32
	got = got | uint32(vm.memory.mem[0x1234])
	got = got | uint32(vm.memory.mem[0x1235])<<8
	got = got | uint32(vm.memory.mem[0x1236])<<16
	got = got | uint32(vm.memory.mem[0x1237])<<24

	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func TestXor(t *testing.T) {
	vm := NewVM()
	vm.reg[0].value = 0b010101
	vm.reg[1].value = 0b100101
	var want uint32 = 0b110101

	VXOR(vm, []byte{0, 1})
	got := vm.reg[0].value
	if got != want {
		t.Errorf("got %b, want %b", got, want)
	}
}
