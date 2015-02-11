package main

// #include <stdio.h>
//
// void myfflush(void) {
// 	fflush(stdout);
// 	fflush(stderr);
// }
import "C"

import (
	"fmt"
	"io/ioutil"

	"github.com/mewkiz/pkg/errutil"
	"golang.org/x/sys/unix"
	"llvm.org/llvm/bindings/go/llvm"
)

// hackDump returns the value dump as a string.
func hackDump(v llvm.Value) (string, error) {
	// HACK!

	// Open temp file.
	// TODO: Use an in-memory file instead of /tmp/x.
	fd, err := unix.Open("/tmp/x", unix.O_WRONLY|unix.O_TRUNC|unix.O_CREAT, 0644)
	if err != nil {
		return "", errutil.Err(err)
	}
	fmt.Println("fd:", fd)

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
	C.myfflush()

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
