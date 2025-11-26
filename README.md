go-z3 provides Go bindings for
the [Z3 SMT solver](https://github.com/Z3Prover/z3).

## Installation

First, follow the instructions to
[download and install](https://github.com/Z3Prover/z3/blob/master/README.md) the Z3 C library.

go-z3 requires Z3 version 4.15.4 or later.

If you installed the C library to a non-default location (such as a
directory under `$HOME`), set the following environment variables:

```sh
# For building:
export CGO_CFLAGS=-I$Z3PREFIX/include CGO_LDFLAGS=-L$Z3PREFIX/lib
# For running binaries (including tests):
export LD_LIBRARY_PATH=$Z3PREFIX/lib
```

Then download and build go-z3:

```sh
go get -u github.com/ralscha/go-z3/z3
```

## Example

Here's a simple example that solves a classic puzzle: *"A farmer has rabbits and pheasants. There are 9 animals total and 24 legs. How many of each?"*

```go
package main

import (
	"fmt"
	"github.com/ralscha/go-z3/z3"
)

func main() {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	rabbits := ctx.IntConst("rabbits")
	pheasants := ctx.IntConst("pheasants")

	// 9 animals total
	solver.Assert(rabbits.Add(pheasants).Eq(ctx.Int(9)))
	// 24 legs total (rabbits have 4 legs, pheasants have 2)
	solver.Assert(rabbits.Mul(ctx.Int(4)).Add(pheasants.Mul(ctx.Int(2))).Eq(ctx.Int(24)))
	// Non-negative counts
	solver.Assert(rabbits.GE(ctx.Int(0)))
	solver.Assert(pheasants.GE(ctx.Int(0)))

	if sat, _ := solver.Check(); sat {
		model := solver.Model()
		r, _, _ := model.EvalAsInt64(rabbits, true)
		p, _, _ := model.EvalAsInt64(pheasants, true)
		fmt.Printf("Rabbits: %d, Pheasants: %d\n", r, p)
		// Output: Rabbits: 3, Pheasants: 6
	}
}
```
