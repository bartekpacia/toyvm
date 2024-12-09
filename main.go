package main

import (
	"log"
	"os"

	"github.com/bartekpacia/toyvm/vm"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalln("usage: vm <source file>")
	}

	debug := false
	for _, arg := range os.Args {
		if arg == "-debug" || arg == "--debug" {
			debug = true
		}
	}

	machine := vm.NewVM()
	err := machine.LoadMemoryFromFile(0, os.Args[1])
	if err != nil {
		log.Fatalln("failed to load memory from file:", err)
	}

	machine.SetDebug(debug)
	err = machine.Run()
	if err != nil {
		log.Fatalln("error while running virtual machine:", err)
	}
}
