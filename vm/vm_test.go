package vm_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/bartekpacia/toyvm/vm"
)

type test struct {
	execFilename string
	output       string
}

var testCases = []test{
	{
		execFilename: "simple3",
		output:       "Hello\n",
	},
	{
		execFilename: "cmp_equal_true",
		output:       "T",
	},
	{
		execFilename: "cmp_equal_false",
		output:       "F",
	},
	{
		execFilename: "hello",
		output:       "Hello World\n",
	},
}

func TestVM(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.execFilename, func(t *testing.T) {
			virtualMachine := vm.NewVM()
			fmt.Println(os.Getwd())
			err := virtualMachine.LoadMemoryFromFile(0, "../examples/"+testCase.execFilename)
			if err != nil {
				t.Errorf("failed to load memory from file: %v", err)
				return
			}

			stdout := bytes.NewBuffer(nil)

			virtualMachine.Stdout = stdout
			err = virtualMachine.Run()
			if err != nil {
				t.Errorf("failed to run: %v", err)
				return
			}

			got := stdout.String()
			want := testCase.output

			if got != want {
				t.Errorf("stdout does not match, got: %#v, want: %#v", got, want)
				return
			}
		})
	}
}
