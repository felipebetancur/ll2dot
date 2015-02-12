# ll2dot

[![GoDoc](https://godoc.org/github.com/mewrev/ll2dot?status.svg)](https://godoc.org/github.com/mewrev/ll2dot)

ll2dot is a tool which creates control flow graphs of LLVM IR assembly files (e.g. *.ll -> *.dot). The output is a set of GraphViz DOT files, each representing the control flow graph of a function using one node per basic block.

For a source file "foo.ll" containing the functions "bar" and "baz" the following DOT files will be created:

   foo_graphs/bar.dot
   foo_graphs/baz.dot

## Installation

```shell
go get github.com/mewrev/ll2dot
```

## Examples

### if

Input:
* [if.ll](testdata/if.ll)

Output:
* [main.dot](testdata/)

![CFG for the main function of if.ll](https://raw.githubusercontent.com/mewrev/ll2dot/master/testdata/if_graphs/main.png)

## Dependencies

* [llvm.org/llvm/bindings/go/llvm](https://godoc.org/llvm.org/llvm/bindings/go/llvm) with [unnamed.patch](unnamed.patch)
* `llvm-as` from [LLVM](http://llvm.org/)
* `dot` from [GraphViz](http://www.graphviz.org/)

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
