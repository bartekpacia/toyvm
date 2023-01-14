package main

import (
	"fmt"
	"os"

	"github.com/bartekpacia/toyvm/vm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: vm <source file>")
		os.Exit(1)
	}

	machine := vm.NewVM()
	machine.LoadMemoryFromFile(0, os.Args[1])
	machine.Run()
}
