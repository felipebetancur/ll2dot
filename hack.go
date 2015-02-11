// HACK: This entire file is a hack!
//
// LLVM IR has a notion of unnamed variables and basic blocks which are given
// function scoped IDs during assembly generation. The in-memory representation
// does not include this ID, so instead of reimplementing the logic of ID slots
// we capture the output of Value.Dump to locate the basic block names. Note
// that unnamed basic blocks are not given explicit labels during vanilla LLVM
// assembly generation, but rather comments which include the basic block ID.
// For this reason the "unnamed.patch" has been applied to the LLVM code base,
// which ensures that all basic blocks are given explicit labels.

package main

// #include <stdio.h>
//
// void fflush_stderr(void) {
// 	fflush(stderr);
// }
import "C"

import (
	"io/ioutil"

	"github.com/mewkiz/pkg/errutil"
	"golang.org/x/sys/unix"
	"llvm.org/llvm/bindings/go/llvm"
)

// hackDump returns the value dump as a string.
func hackDump(v llvm.Value) (string, error) {
	// Open temp file.
	// TODO: Use an in-memory file instead of /tmp/x.
	fd, err := unix.Open("/tmp/x", unix.O_WRONLY|unix.O_TRUNC|unix.O_CREAT, 0644)
	if err != nil {
		return "", errutil.Err(err)
	}

	// Store original stderr.
	stderr, err := unix.Dup(2)
	if err != nil {
		return "", errutil.Err(err)
	}

	// Capture stderr and redirect its output to the temp file.
	err = unix.Dup2(fd, 2)
	if err != nil {
		return "", errutil.Err(err)
	}
	err = unix.Close(fd)
	if err != nil {
		return "", errutil.Err(err)
	}

	// Dump value.
	v.Dump()
	C.fflush_stderr()

	// Restore stderr.
	err = unix.Dup2(stderr, 2)
	if err != nil {
		return "", errutil.Err(err)
	}
	err = unix.Close(stderr)
	if err != nil {
		return "", errutil.Err(err)
	}

	// Return content of temp file.
	buf, err := ioutil.ReadFile("/tmp/x")
	if err != nil {
		return "", errutil.Err(err)
	}
	return string(buf), nil
}
