// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

// Pseudo-Boolean constraints are cardinality constraints over Boolean variables.

// AtMost returns a constraint that at most k of the args are true.
// This is equivalent to: args[0] + args[1] + ... + args[n-1] <= k
// where true is treated as 1 and false as 0.
func (ctx *Context) AtMost(args []Bool, k uint) Bool {
	cargs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cargs[i] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_atmost(ctx.c, C.uint(len(cargs)), &cargs[0], C.uint(k))
	})
	runtime.KeepAlive(&cargs[0])
	return Bool(val)
}

// AtLeast returns a constraint that at least k of the args are true.
// This is equivalent to: args[0] + args[1] + ... + args[n-1] >= k
// where true is treated as 1 and false as 0.
func (ctx *Context) AtLeast(args []Bool, k uint) Bool {
	cargs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cargs[i] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_atleast(ctx.c, C.uint(len(cargs)), &cargs[0], C.uint(k))
	})
	runtime.KeepAlive(&cargs[0])
	return Bool(val)
}

// PbLE returns a constraint that the weighted sum is at most k.
// This is equivalent to: coeffs[0]*args[0] + coeffs[1]*args[1] + ... <= k
// where true is treated as 1 and false as 0.
func (ctx *Context) PbLE(args []Bool, coeffs []int, k int) Bool {
	if len(args) != len(coeffs) {
		panic("args and coeffs must have the same length")
	}
	cargs := make([]C.Z3_ast, len(args))
	ccoeffs := make([]C.int, len(coeffs))
	for i, arg := range args {
		cargs[i] = arg.c
		ccoeffs[i] = C.int(coeffs[i])
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_pble(ctx.c, C.uint(len(cargs)), &cargs[0], &ccoeffs[0], C.int(k))
	})
	runtime.KeepAlive(&cargs[0])
	runtime.KeepAlive(&ccoeffs[0])
	return Bool(val)
}

// PbGE returns a constraint that the weighted sum is at least k.
// This is equivalent to: coeffs[0]*args[0] + coeffs[1]*args[1] + ... >= k
// where true is treated as 1 and false as 0.
func (ctx *Context) PbGE(args []Bool, coeffs []int, k int) Bool {
	if len(args) != len(coeffs) {
		panic("args and coeffs must have the same length")
	}
	cargs := make([]C.Z3_ast, len(args))
	ccoeffs := make([]C.int, len(coeffs))
	for i, arg := range args {
		cargs[i] = arg.c
		ccoeffs[i] = C.int(coeffs[i])
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_pbge(ctx.c, C.uint(len(cargs)), &cargs[0], &ccoeffs[0], C.int(k))
	})
	runtime.KeepAlive(&cargs[0])
	runtime.KeepAlive(&ccoeffs[0])
	return Bool(val)
}

// PbEq returns a constraint that the weighted sum equals k.
// This is equivalent to: coeffs[0]*args[0] + coeffs[1]*args[1] + ... = k
// where true is treated as 1 and false as 0.
func (ctx *Context) PbEq(args []Bool, coeffs []int, k int) Bool {
	if len(args) != len(coeffs) {
		panic("args and coeffs must have the same length")
	}
	cargs := make([]C.Z3_ast, len(args))
	ccoeffs := make([]C.int, len(coeffs))
	for i, arg := range args {
		cargs[i] = arg.c
		ccoeffs[i] = C.int(coeffs[i])
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_pbeq(ctx.c, C.uint(len(cargs)), &cargs[0], &ccoeffs[0], C.int(k))
	})
	runtime.KeepAlive(&cargs[0])
	runtime.KeepAlive(&ccoeffs[0])
	return Bool(val)
}
