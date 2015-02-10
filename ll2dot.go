// ll2dot is a tool which converts LLVM IR assembly files to GraphViz DOT files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

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
// defined functions to directed graphs with one node per basic block.
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
	for f := module.FirstFunction(); !f.IsNil(); f = llvm.NextFunction(f) {
		fmt.Println("--- [ function ] ---")
		f.Dump()
		if f.IsDeclaration() {
			continue
		}
		bbs := f.BasicBlocks()
		for _, bb := range bbs {
			fmt.Println("___ [ basic block ] ___")
			term := bb.LastInstruction()
			fmt.Println("~~~ [ terminator instruction ] ~~~")
			term.Dump()
			for i := 0; i < term.OperandsCount(); i++ {
				op := term.Operand(i)
				fmt.Println("### [ operand ] ###")
				op.Dump()
			}
			opcode := term.InstructionOpcode()
			switch opcode {
			case llvm.Ret:
				// exit node.
				//    ret <type> <value>
				//    ret <void>
				fmt.Println("ret")
			case llvm.Br:
				// 2-way conditional branch.
				//    br i1 <cond>, label <target_true>, label <target_false>
				// unconditional branch.
				//    br label <target>
				fmt.Println("br")
			case llvm.Switch:
				// n-way conditional branch.
				//    switch <int_type> <value>, label <default_target> [
				//       <int_type> <case1>, label <case1_target>
				//       <int_type> <case2>, label <case2_target>
				//       ...
				//    ]
				fmt.Println("switch")
			case llvm.Unreachable:
				// unreachable node (similar to exit node).
				//    unreachable
				fmt.Println("unreachable")
			default:
				// Not yet supported:
				//    - indirectbr
				//    - invoke
				//    - resume
				panic(fmt.Sprintf("not yet implemented; support for terminator %v", opcode))
			}
		}
		if f == module.LastFunction() {
			break
		}
	}
	return nil
}
