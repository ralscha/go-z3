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

// Char is a symbolic value representing a character.
//
// Char implements Value.
type Char value

func init() {
	kindWrappers[KindChar] = func(x value) Value {
		return Char(x)
	}
}

// CharSort returns the character sort.
func (ctx *Context) CharSort() Sort {
	var sort Sort
	ctx.do(func() {
		sort = wrapSort(ctx, C.Z3_mk_char_sort(ctx.c), KindChar)
	})
	return sort
}

// IsCharSort returns true if s is a character sort.
func (s Sort) IsCharSort() bool {
	var result bool
	s.ctx.do(func() {
		result = z3ToBool(C.Z3_is_char_sort(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return result
}

// CharConst returns a character constant named "name".
func (ctx *Context) CharConst(name string) Char {
	return ctx.Const(name, ctx.CharSort()).(Char)
}

// Eq returns a Value that is true if l and r are equal.
func (l Char) Eq(r Char) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_eq(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// NE returns a Value that is true if l and r are not equal.
func (l Char) NE(r Char) Bool {
	return l.ctx.Distinct(l, r)
}

// LE returns l <= r.
func (l Char) LE(r Char) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_char_le(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// ToInt returns the code point of character l.
func (l Char) ToInt() Int {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_char_to_int(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Int(val)
}

// ToBV returns the character l as a bit-vector.
// The bit-vector size depends on Z3's encoding setting (default: 18 bits for Unicode).
func (l Char) ToBV() BV {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_char_to_bv(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// CharFromBV creates a character from a bit-vector.
// The bit-vector size must match Z3's encoding setting (default: 18 bits for Unicode).
func (ctx *Context) CharFromBV(bv BV) Char {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_char_from_bv(ctx.c, bv.c)
	})
	runtime.KeepAlive(bv)
	return Char(val)
}

// IsDigit returns true if character l is a digit (0-9).
func (l Char) IsDigit() Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_char_is_digit(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Bool(val)
}

// StringToCode returns the code point of the first character in s,
// or -1 if s is empty.
func (s String) ToCode() Int {
	ctx := s.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_string_to_code(ctx.c, s.c)
	})
	runtime.KeepAlive(s)
	return Int(val)
}

// StringFromCode creates a string from a code point.
func (ctx *Context) StringFromCode(code Int) String {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_string_from_code(ctx.c, code.c)
	})
	runtime.KeepAlive(code)
	return String(val)
}
