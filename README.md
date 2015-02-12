# ll2dot

[![GoDoc](https://godoc.org/github.com/mewrev/ll2dot?status.svg)](https://godoc.org/github.com/mewrev/ll2dot)

ll2dot is a tool which creates control flow graphs of LLVM IR assembly files (e.g. *.ll -> *.dot). The output is a set of GraphViz DOT files, each representing the control flow graph of a function using one node per basic block.

For a source file "foo.ll" containing the functions "bar" and "baz" the following DOT files will be created:

   * foo_graphs/bar.dot
   * foo_graphs/baz.dot

## Installation

```shell
go get github.com/mewrev/ll2dot
```

## Examples

### funcs

Input:
* [funcs.ll](testdata/funcs.ll)

Output:
* [bar.dot](testdata/funcs_graphs/bar.dot)
* [main.dot](testdata/funcs_graphs/main.dot)

![CFG funcs the bar function of funcs.ll](https://raw.githubusercontent.com/mewrev/ll2dot/master/testdata/funcs_graphs/bar.png)
![CFG funcs the main function of funcs.ll](https://raw.githubusercontent.com/mewrev/ll2dot/master/testdata/funcs_graphs/main.png)

### switch

Input:
* [switch.ll](testdata/switch.ll)

Output:
* [main.dot](testdata/switch_graphs/main.dot)

![CFG switch the main function of switch.ll](https://raw.githubusercontent.com/mewrev/ll2dot/master/testdata/switch_graphs/main.png)

## Dependencies

* [llvm.org/llvm/bindings/go/llvm](https://godoc.org/llvm.org/llvm/bindings/go/llvm) with [unnamed.patch](unnamed.patch)
* `llvm-as` from [LLVM](http://llvm.org/)
* `dot` from [GraphViz](http://www.graphviz.org/)

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
