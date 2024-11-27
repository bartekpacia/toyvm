// Package vm implements a very simple virtual machine.
package vm

import (
	"fmt"
	"os"
	"sync"
)

// gpRegister is a general-purpose register.
type gpRegister struct {
	value uint32
}

const (
	FlagZF = iota + 1 // zero flag
	FlagCF            // carry flag
)

const (
	IntMemoryError   = iota
	IntDivisionError = iota
	IntGeneralError  = iota

	IntPit     = 8 // generated by programmable timer
	IntConsole = 9 // generated by console

	CregIntFirst  = 0x100
	CregIntLast   = 0x10f
	CregIntContrl = 0x110
)

var MaskableInterrupts = []int{IntPit, IntConsole}

type VM struct {
	memory     *Memory
	reg        []gpRegister // general-purpose registers
	creg       map[int]int  // control registers
	pc         *gpRegister  // program counter
	sp         *gpRegister  // stack pointer
	fr         uint32       // flag register
	terminated bool
	opcodes    map[byte]opcode

	interruptQueue      []int
	interruptQueueMutex sync.Mutex

	deferredQueue []func()
}

func NewVM() *VM {
	var registers []gpRegister
	for i := 0; i < 16; i++ {
		registers = append(registers, gpRegister{})
	}

	vm := VM{
		memory:     &Memory{mem: make([]byte, 64*1024)}, // 64KB
		reg:        registers,
		creg:       make(map[int]int),
		pc:         &registers[14],
		sp:         &registers[15],
		fr:         0,
		terminated: false,
		opcodes:    opcodes,

		interruptQueue:      make([]int, 0),
		interruptQueueMutex: sync.Mutex{},

		deferredQueue: make([]func(), 0),
	}

	vm.sp.value = 0x10000

	// Interrupt registers.
	for creg := CregIntFirst; creg < CregIntLast+1; creg++ {
		vm.creg[creg] = 0xffffffff
	}
	vm.creg[CregIntContrl] = 0 // Maskable interrupts disabled.

	return &vm
}

// crash terminates the virtual machine on critical error
func (vm *VM) crash() {
	vm.terminated = true

	fmt.Println("the virtual machine entered an erroneous state and is terminating")
	fmt.Println("register values at termination:")
	for ri, r := range vm.reg {
		fmt.Printf("\tr%d = %x\n", ri, r.value)
	}
}

func (vm *VM) interrupt(interrupt int) {
	vm.interruptQueueMutex.Lock()
	defer vm.interruptQueueMutex.Unlock()

	vm.interruptQueue = append(vm.interruptQueue, interrupt)
}

// fetchPendingInterrupt returns a pending interrupt to be processed. If
// maskable interrupts are disabled, returns a non-maskable interrupt (NMI) if
// available. If no interrupts are available for processing, returns None.
func (vm *VM) fetchPendingInterrupt() *int {
	vm.interruptQueueMutex.Lock()
	defer vm.interruptQueueMutex.Unlock()

	if len(vm.interruptQueue) == 0 {
		return nil
	}

	// In disable-interrupts state, we can process only non-maskable interrupts
	// (faults).
	if vm.creg[CregIntContrl]&1 == 0 {
		// Maskable interrupts disabled. Find a non-maskable one.
		for i, interrupt := range vm.interruptQueue {
			for _, maskableInterrupt := range MaskableInterrupts {
				if interrupt == maskableInterrupt {
					q := vm.interruptQueue

					// pop element at index i
					_interrupt := q[i]
					q = append(q[:i], q[i+1:]...)
					vm.interruptQueue = q

					return &_interrupt
				}
			}
		}

		return nil
	}

	// No non-maskable interrupts found.

	// Return the first interrupt available (pop)
	q := vm.interruptQueue
	i, queue := q[len(q)-1], q[:len(q)-1]
	vm.interruptQueue = queue

	return &i
}

// processInterruptQueue processes an interrupt (if one is available).
func (vm *VM) processInterruptQueue() error {
	i := vm.fetchPendingInterrupt()

	if i == nil {
		return nil
	}

	// Save context
	tmpSp := vm.sp.value
	registerValues := make([]uint32, 0)
	for _, register := range vm.reg {
		registerValues = append(registerValues, register.value)
	}
	registerValues = append(registerValues, vm.fr)

	for _, val := range registerValues {
		tmpSp -= 4
		err := vm.memory.StoreDword(uint16(tmpSp), val)
		if err != nil {
			// Since there is no way to save the state, and therefore no way to
			// recover, crash the machine.
			vm.crash()
			return fmt.Errorf("failed to store dword: %w", err)
		}
	}

	vm.sp.value = tmpSp
	vm.pc.value = uint32(vm.creg[CregIntFirst+(*i&0xf)])
	return nil
}

func (vm *VM) runSingleStep() error {
	// If there is any interrupt on the queue, we need to know about it now.
	err := vm.processInterruptQueue()
	if err != nil {
		return fmt.Errorf("something failed hard: %w", err)
	}

	// Check if there is anything in the deferred queue. If so, process it now.
	for len(vm.deferredQueue) != 0 {
		// pop
		q := vm.deferredQueue
		action, queue := q[len(q)-1], q[:len(q)-1]
		vm.deferredQueue = queue

		action()
	}

	// Proceed with normal execution
	opcodeByte, err := vm.memory.FetchByte(uint16(vm.pc.value))
	if err != nil {
		vm.interrupt(IntMemoryError)
		return fmt.Errorf("failed to fetch opcode: %v", err)
	}

	opcode, ok := vm.opcodes[opcodeByte]
	if !ok {
		vm.interrupt(IntGeneralError)
	}

	length := opcode.length
	argBytes, err := vm.memory.FetchMany(uint16(vm.pc.value+1), length)
	if err != nil {
		vm.interrupt(IntMemoryError)
		return fmt.Errorf("failed to fetch arg bytes: %v", err)
	}

	handler := opcode.handler
	vm.pc.value = vm.pc.value + 1 + uint32(length)
	handler(vm, argBytes)
	return nil
}

func (vm *VM) LoadMemoryFromFile(addr uint16, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read code from file %s: %w", filename, err)
	}

	err = vm.memory.StoreMany(addr, data)
	if err != nil {
		return fmt.Errorf("store data at address %d: %w", addr, err)
	}

	return nil
}

func (vm *VM) Run() error {
	for !vm.terminated {
		err := vm.runSingleStep()
		if err != nil {
			return fmt.Errorf("run single step: %v", err)
		}
	}

	// vm.devConsole.terminate() vm.devPIT.terminate()
	return nil
}
