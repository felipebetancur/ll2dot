// ll2dot is a tool which converts LLVM IR assembly files to GraphViz DOT files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewlang/llvm/asm/lexer"
	"github.com/mewlang/llvm/asm/token"
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

	// Parse foo.bc
	module, err := llvm.ParseBitcodeFile(bcPath)
	if err != nil {
		return err
	}
	defer module.Dispose()

	// Create one graph for each function using basic blocks as nodes.
	fmt.Println("=== [ module ] ===")
	for f := module.FirstFunction(); !f.IsNil(); f = llvm.NextFunction(f) {
		fmt.Println("--- [ function ] ---")
		f.Dump()
		if f.IsDeclaration() {
			continue
		}
		graph := dot.NewGraph()
		graphName := f.Name()
		fmt.Println("graph name:", graphName)
		graph.SetName(graphName)
		bbs := f.BasicBlocks()
		for _, bb := range bbs {
			// TODO: Mark the entry basic block.
			//    if bb == f.EntryBasicBlock()

			// Add node (i.e. basic block) to graph.
			fmt.Println("___ [ basic block ] ___")
			nodeName, err := getBBName(bb.AsValue())
			if err != nil {
				return err
			}
			fmt.Println("node name:", nodeName)
			graph.AddNode(graphName, nodeName, nil)

			// Add edges from node (i.e. target basic blocks) to graph.
			term := bb.LastInstruction()
			fmt.Println("~~~ [ terminator instruction ] ~~~")
			term.Dump()
			nops := term.OperandsCount()
			for i := 0; i < nops; i++ {
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
				// unconditional branch.
				//    br label <target>
				// 2-way conditional branch.
				//    br i1 <cond>, label <target_true>, label <target_false>
				switch nops {
				case 1:
					// unconditional branch
					target := term.Operand(0)
					targetName, err := getBBName(target)
					if err != nil {
						return err
					}
					fmt.Println("target name:", targetName)
					graph.AddEdge(nodeName, targetName, true, nil)
				case 3:
					// 2-way conditional branch
					targetTrue, targetFalse := term.Operand(1), term.Operand(2)
					targetTrueName, err := getBBName(targetTrue)
					if err != nil {
						return err
					}
					targetFalseName, err := getBBName(targetFalse)
					if err != nil {
						return err
					}
					fmt.Println("target true name:", targetTrueName)
					fmt.Println("target false name:", targetFalseName)
					graph.AddEdge(nodeName, targetTrueName, true, nil)  // TODO: Add "true" to attrs?
					graph.AddEdge(nodeName, targetFalseName, true, nil) // TODO: Add "false" to attrs?
				default:
					return fmt.Errorf("invalid number of parameters (%d) for br", nops)
				}
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
		fmt.Println("### [ graph ] ###")
		fmt.Println(graph)
		if f == module.LastFunction() {
			break
		}
	}
	return nil
}

// getBBName returns the name (or ID if unnamed) of a basic block.
func getBBName(v llvm.Value) (string, error) {
	if !v.IsBasicBlock() {
		return "", errutil.Newf("invalid value type; expected basic block, got %v", v.Type())
	}

	// Return name of named basic block.
	if name := v.Name(); len(name) > 0 {
		return name, nil
	}

	// Return ID of unnamed basic block.

	// Search for the basic block label in the value dump.
	//
	// Example value dump:
	//    0:
	//      br i1 true, label %1, label %2
	//
	// Each basic block is expected to have a label, which requires the
	// unnamed.patch to be applied to the llvm.org/llvm/bindings/go/llvm code
	// base.
	s, err := hackDump(v)
	if err != nil {
		return "", errutil.Err(err)
	}
	fmt.Println("s:", s)
	tokens := lexer.ParseString(s)
	if len(tokens) < 1 {
		return "", errutil.Newf("unable to locate basic block name in %q", s)
	}
	tok := tokens[0]
	if tok.Kind != token.Label {
		return "", errutil.Newf("invalid token; expected %v, got %v", token.Label, tok.Kind)
	}
	return tok.Val, nil
}
