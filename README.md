go-z3 provides Go bindings for
the [Z3 SMT solver](https://github.com/Z3Prover/z3).

## Installation

### Prerequisites

- **Go 1.21 or later**
- **Z3 4.15.4 or later** - Download from [Z3 Releases](https://github.com/Z3Prover/z3/releases)
- **C compiler** - GCC or compatible (CGO is required)

### Linux / macOS

1. Download and install the Z3 library. If you installed it to a non-default location, set the following environment variables:

```sh
export Z3PREFIX=/path/to/z3

# For building:
export CGO_CFLAGS="-I$Z3PREFIX/include"
export CGO_LDFLAGS="-L$Z3PREFIX/lib -lz3"

# For running binaries (including tests):
export LD_LIBRARY_PATH=$Z3PREFIX/lib
```

2. Install go-z3:

```sh
go get github.com/ralscha/go-z3/z3
```

### Windows

1. Download the Z3 release for Windows and extract it.

2. Install MinGW-w64 or another GCC distribution for Windows.

3. Set the required environment variables (adjust paths as needed):

**Command Prompt (cmd.exe):**
```cmd
SET CGO_ENABLED=1
SET CGO_CFLAGS=-IC:/path/to/z3/include
SET CGO_LDFLAGS=-LC:/path/to/z3/bin -lz3
SET PATH=C:/path/to/mingw64/bin;C:/path/to/z3/bin;%PATH%
SET CC=C:/path/to/mingw64/bin/gcc.exe
```

**PowerShell:**
```powershell
$env:CGO_ENABLED = "1"
$env:CGO_CFLAGS = "-IC:/path/to/z3/include"
$env:CGO_LDFLAGS = "-LC:/path/to/z3/bin -lz3"
$env:PATH = "C:/path/to/mingw64/bin;C:/path/to/z3/bin;$env:PATH"
$env:CC = "C:/path/to/mingw64/bin/gcc.exe"
```

**Git Bash / MSYS2:**
```sh
export CGO_ENABLED=1
export CGO_CFLAGS="-I/c/path/to/z3/include"
export CGO_LDFLAGS="-L/c/path/to/z3/bin -lz3"
export PATH="/c/path/to/mingw64/bin:/c/path/to/z3/bin:$PATH"
export CC="/c/path/to/mingw64/bin/gcc.exe"
```

4. Install go-z3:

```sh
go get github.com/ralscha/go-z3/z3
```

> **Note:** On Windows, the Z3 DLL (`libz3.dll`) must be in your PATH at runtime.

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

For more practical examples including Sudoku solving, N-Queens, and other constraint satisfaction problems, see [real_world_test.go](z3/real_world_test.go).
