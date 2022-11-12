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
			desc:   "copy value of R5 to R2",
			seed:   func(vm *VM) { vm.reg[4].value = 1 },
			args:   []byte{2, 5},
			verify: func(vm *VM) bool { return vm.reg[1].value == 1 },
		},
	}

	for _, tc := range testCases {
		vm := NewVM()
		tc.seed(vm)
		VMOV(vm, tc.args)
		if !tc.verify(vm) {
			t.Fatalf("%s failed", tc.desc)
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
			t.Fatalf("want %x, got %x", tc.want, got)
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
		t.Fatalf("want %x, got %x", want, got)
	}
}
