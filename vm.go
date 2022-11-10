// Package vm implements a very simple virtual machine.
package vm

type VM struct{}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) LoadMemoryFromFile(addr uint16, filename string) {
	// TODO
}

func (vm *VM) Run() {
	
}
