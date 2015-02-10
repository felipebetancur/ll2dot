// ll2dot is a tool which converts LLVM IR assembly files to GraphViz DOT files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/mewkiz/pkg/pathutil"
	"llvm.org/llvm/bindings/go/llvm"
)

func main() {
	flag.Parse()
	for _, llPath := range flag.Args() {
		err := createDOT(llPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// createDOT parses the provided LLVM IR assembly file and converts each of its
// defined functions to a directed graph with one node per basic block.
func createDOT(llPath string) error {
	// foo.ll -> foo.bc
	cmd := exec.Command("llvm-as", llPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	bcPath := pathutil.TrimExt(llPath) + ".bc"
	module, err := llvm.ParseBitcodeFile(bcPath)
	if err != nil {
		return err
	}
	defer module.Dispose()
	fmt.Println("=== [ module ] ===")
	spew.Dump(module)
	for f := module.FirstFunction(); ; f = llvm.NextFunction(f) {
		fmt.Println("--- [ function ] ---")
		spew.Dump(f)
		f.Dump()
		bbs := f.BasicBlocks()
		fmt.Println("--- [ basic blocks ] ---")
		spew.Dump(bbs)
		for _, bb := range bbs {
			fmt.Println("___ [ basic block ] ___")
			spew.Dump(bb)
			for inst := bb.FirstInstruction(); ; inst = llvm.NextInstruction(inst) {
				fmt.Println("~~~ [ instruction ] ~~~")
				spew.Dump(inst)
				inst.Dump()
				if inst == bb.LastInstruction() {
					break
				}
			}
		}
		if f == module.LastFunction() {
			break
		}
	}
	return nil
}
