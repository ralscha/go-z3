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

// RE is a symbolic value representing a regular expression.
//
// RE implements Value.
type RE value

func init() {
	kindWrappers[KindRE] = func(x value) Value {
		return RE(x)
	}
}

// RESort returns a regular expression sort over the given sequence sort.
func (ctx *Context) RESort(seq Sort) Sort {
	var sort Sort
	ctx.do(func() {
		sort = wrapSort(ctx, C.Z3_mk_re_sort(ctx.c, seq.c), KindRE)
	})
	runtime.KeepAlive(seq)
	return sort
}

// IsRESort returns true if s is a regular expression sort.
func (s Sort) IsRESort() bool {
	var result bool
	s.ctx.do(func() {
		result = z3ToBool(C.Z3_is_re_sort(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return result
}

// RESortBasis returns the sequence sort that this RE sort is over.
func (s Sort) RESortBasis() Sort {
	var result Sort
	s.ctx.do(func() {
		result = wrapSort(s.ctx, C.Z3_get_re_sort_basis(s.ctx.c, s.c), KindSeq)
	})
	runtime.KeepAlive(s)
	return result
}

// Eq returns a Value that is true if l and r are equal.
func (l RE) Eq(r RE) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_eq(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// NE returns a Value that is true if l and r are not equal.
func (l RE) NE(r RE) Bool {
	return l.ctx.Distinct(l, r)
}

// Plus returns re+ (one or more occurrences of re).
func (re RE) Plus() RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_plus(ctx.c, re.c)
	})
	runtime.KeepAlive(re)
	return RE(val)
}

// Star returns re* (zero or more occurrences of re).
func (re RE) Star() RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_star(ctx.c, re.c)
	})
	runtime.KeepAlive(re)
	return RE(val)
}

// Option returns [re] (zero or one occurrence of re).
func (re RE) Option() RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_option(ctx.c, re.c)
	})
	runtime.KeepAlive(re)
	return RE(val)
}

// Complement returns the complement of re.
func (re RE) Complement() RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_complement(ctx.c, re.c)
	})
	runtime.KeepAlive(re)
	return RE(val)
}

// Union returns the union of re with others.
func (re RE) Union(others ...RE) RE {
	ctx := re.ctx
	cargs := make([]C.Z3_ast, len(others)+1)
	cargs[0] = re.c
	for i, arg := range others {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_union(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return RE(val)
}

// Concat returns the concatenation of re with others.
func (re RE) Concat(others ...RE) RE {
	ctx := re.ctx
	cargs := make([]C.Z3_ast, len(others)+1)
	cargs[0] = re.c
	for i, arg := range others {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_concat(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return RE(val)
}

// Intersect returns the intersection of re with others.
func (re RE) Intersect(others ...RE) RE {
	ctx := re.ctx
	cargs := make([]C.Z3_ast, len(others)+1)
	cargs[0] = re.c
	for i, arg := range others {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_intersect(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return RE(val)
}

// Diff returns the set difference of re minus other.
func (re RE) Diff(other RE) RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_diff(ctx.c, re.c, other.c)
	})
	runtime.KeepAlive(re)
	runtime.KeepAlive(other)
	return RE(val)
}

// Range returns a regular expression that matches any character
// in the range [lo, hi] (inclusive).
func (ctx *Context) RERange(lo, hi String) RE {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_range(ctx.c, lo.c, hi.c)
	})
	runtime.KeepAlive(lo)
	runtime.KeepAlive(hi)
	return RE(val)
}

// Loop returns a regular expression that matches between lo and hi
// occurrences of re.
func (re RE) Loop(lo, hi uint) RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_loop(ctx.c, re.c, C.uint(lo), C.uint(hi))
	})
	runtime.KeepAlive(re)
	return RE(val)
}

// Power returns a regular expression that matches exactly n occurrences of re.
func (re RE) Power(n uint) RE {
	ctx := re.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_power(ctx.c, re.c, C.uint(n))
	})
	runtime.KeepAlive(re)
	return RE(val)
}

// REEmpty returns the empty regular expression (matches nothing).
func (ctx *Context) REEmpty(s Sort) RE {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_empty(ctx.c, s.c)
	})
	runtime.KeepAlive(s)
	return RE(val)
}

// REFull returns the full regular expression (matches everything).
func (ctx *Context) REFull(s Sort) RE {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_full(ctx.c, s.c)
	})
	runtime.KeepAlive(s)
	return RE(val)
}

// REAllChar returns a regular expression that matches any single character.
func (ctx *Context) REAllChar(s Sort) RE {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_re_allchar(ctx.c, s.c)
	})
	runtime.KeepAlive(s)
	return RE(val)
}
