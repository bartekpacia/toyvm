// Package vm implements a very simple virtual machine.
package vm

import (
	"fmt"
	"os"
)

// gpRegister is a general-purpose register.
type gpRegister struct {
	value uint32
}

type VM struct {
	memory     *Memory
	registers  []gpRegister
	pc         *gpRegister // program counter
	sp         *gpRegister // stack pointer
	fr         uint32      // flag register
	terminated bool
}

func NewVM() *VM {
	var registers []gpRegister
	for i := 0; i < 16; i++ {
		registers = append(registers, gpRegister{})
	}

	return &VM{
		memory:     &Memory{mem: make([]byte, 64*1024)}, // 64KB
		registers:  registers,
		pc:         &registers[14],
		sp:         &registers[15],
		fr:         0,
		terminated: false,
	}
}

func (vm *VM) LoadMemoryFromFile(addr uint16, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read code from file %s: %w", filename, err)
	}

	err = vm.memory.StoreMany(addr, data)
	if err != nil {
		return fmt.Errorf("failed to load code to memory: %w", err)
	}

	return nil
}

func (vm *VM) Run() {
	for !vm.terminated {
		// vm.runSingleStep
	}

	// vm.devConsole.terminate()
	// vm.devPIT.terminate()

	// Simple (though ugly) method to make sure all threads exit.
	return
}
