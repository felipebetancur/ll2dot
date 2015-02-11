package main

// #include <stdio.h>
//
// void myfflush(void) {
// 	fflush(stdout);
// }
import "C"

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"llvm.org/llvm/bindings/go/llvm"
)

// hackDump returns the value dump as a string.
func hackDump(v llvm.Value) string {
	// HACK!

	// Capture stdout and stderr.
	stdout, err := syscall.Dup(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalln("hackDump:", err)
	}
	stderr, err := syscall.Dup(int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalln("hackDump:", err)
	}
	// TODO: Use a memory file instead of /tmp/x.
	f, err := os.Create("/tmp/x")
	if err != nil {
		log.Fatalln("hackDump:", err)
	}
	syscall.Dup2(int(f.Fd()), 1)
	syscall.Dup2(int(f.Fd()), 2)

	// Dump value.
	v.Dump()
	C.myfflush()

	err = f.Close()
	if err != nil {
		log.Fatalln("hackDump:", err)
	}

	// Restore stdout and stderr.
	syscall.Dup2(stdout, 1)
	syscall.Dup2(stderr, 2)

	buf, err := ioutil.ReadFile("/tmp/x")
	if err != nil {
		log.Fatalln("hackDump:", err)
	}

	return string(buf)
}
