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
import (
	"runtime"
	"unsafe"
)

// String is a symbolic value representing a string (sequence of characters).
//
// String implements Value.
type String value

func init() {
	kindWrappers[KindSeq] = func(x value) Value {
		return String(x)
	}
}

// StringSort returns the string sort.
func (ctx *Context) StringSort() Sort {
	var sort Sort
	ctx.do(func() {
		sort = wrapSort(ctx, C.Z3_mk_string_sort(ctx.c), KindSeq)
	})
	return sort
}

// SeqSort returns a sequence sort over the given element sort.
func (ctx *Context) SeqSort(elem Sort) Sort {
	var sort Sort
	ctx.do(func() {
		sort = wrapSort(ctx, C.Z3_mk_seq_sort(ctx.c, elem.c), KindSeq)
	})
	runtime.KeepAlive(elem)
	return sort
}

// IsSeqSort returns true if s is a sequence sort.
func (s Sort) IsSeqSort() bool {
	var result bool
	s.ctx.do(func() {
		result = z3ToBool(C.Z3_is_seq_sort(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return result
}

// IsStringSort returns true if s is a string sort.
func (s Sort) IsStringSort() bool {
	var result bool
	s.ctx.do(func() {
		result = z3ToBool(C.Z3_is_string_sort(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return result
}

// SeqSortBasis returns the element sort of a sequence sort.
func (s Sort) SeqSortBasis() Sort {
	var result Sort
	s.ctx.do(func() {
		result = wrapSort(s.ctx, C.Z3_get_seq_sort_basis(s.ctx.c, s.c), KindUnknown)
	})
	runtime.KeepAlive(s)
	return result
}

// StringConst returns a string constant named "name".
func (ctx *Context) StringConst(name string) String {
	return ctx.Const(name, ctx.StringSort()).(String)
}

// FromString returns a string literal with value val.
func (ctx *Context) FromString(val string) String {
	cstr := C.CString(val)
	defer C.free(unsafe.Pointer(cstr))
	return String(wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_string(ctx.c, cstr)
	}))
}

// AsString returns the value of lit as a Go string. If lit is not a
// string literal, it returns "", false.
func (lit String) AsString() (val string, isLiteral bool) {
	var result C.Z3_string
	var isStr bool
	lit.ctx.do(func() {
		isStr = z3ToBool(C.Z3_is_string(lit.ctx.c, lit.c))
		if isStr {
			result = C.Z3_get_string(lit.ctx.c, lit.c)
		}
	})
	runtime.KeepAlive(lit)
	if !isStr {
		return "", false
	}
	return C.GoString(result), true
}

// Empty returns an empty string/sequence of the given sort.
func (ctx *Context) EmptySeq(s Sort) String {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_empty(ctx.c, s.c)
	})
	runtime.KeepAlive(s)
	return String(val)
}

// Unit returns a unit sequence containing the single element elem.
func (ctx *Context) SeqUnit(elem Value) String {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_unit(ctx.c, elem.impl().c)
	})
	runtime.KeepAlive(elem)
	return String(val)
}

// Eq returns a Value that is true if l and r are equal.
func (l String) Eq(r String) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_eq(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// NE returns a Value that is true if l and r are not equal.
func (l String) NE(r String) Bool {
	return l.ctx.Distinct(l, r)
}

// Concat returns the concatenation of l and r.
func (l String) Concat(r ...String) String {
	ctx := l.ctx
	cargs := make([]C.Z3_ast, len(r)+1)
	cargs[0] = l.c
	for i, arg := range r {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_concat(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return String(val)
}

// Length returns the length of l.
func (l String) Length() Int {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_length(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Int(val)
}

// Contains returns true if l contains the substring sub.
func (l String) Contains(sub String) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_contains(ctx.c, l.c, sub.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(sub)
	return Bool(val)
}

// PrefixOf returns true if l is a prefix of s.
func (l String) PrefixOf(s String) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_prefix(ctx.c, l.c, s.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(s)
	return Bool(val)
}

// SuffixOf returns true if l is a suffix of s.
func (l String) SuffixOf(s String) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_suffix(ctx.c, l.c, s.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(s)
	return Bool(val)
}

// Extract returns the substring of l starting at offset with given length.
func (l String) Extract(offset, length Int) String {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_extract(ctx.c, l.c, offset.c, length.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(offset)
	runtime.KeepAlive(length)
	return String(val)
}

// At returns the unit sequence at position index in l.
// The sequence is empty if the index is out of bounds.
func (l String) At(index Int) String {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_at(ctx.c, l.c, index.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(index)
	return String(val)
}

// Nth returns the element at position index in l.
// The function is under-specified if the index is out of bounds.
func (l String) Nth(index Int) Value {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_nth(ctx.c, l.c, index.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(index)
	return val.lift(KindUnknown)
}

// IndexOf returns the index of the first occurrence of substr in l
// starting from offset. Returns -1 if not found.
func (l String) IndexOf(substr String, offset Int) Int {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_index(ctx.c, l.c, substr.c, offset.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(substr)
	runtime.KeepAlive(offset)
	return Int(val)
}

// LastIndexOf returns the index of the last occurrence of substr in l.
// Returns -1 if not found.
func (l String) LastIndexOf(substr String) Int {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_last_index(ctx.c, l.c, substr.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(substr)
	return Int(val)
}

// Replace returns l with the first occurrence of src replaced by dst.
func (l String) Replace(src, dst String) String {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_replace(ctx.c, l.c, src.c, dst.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(src)
	runtime.KeepAlive(dst)
	return String(val)
}

// LT returns l < r (lexicographic comparison).
func (l String) LT(r String) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_str_lt(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// LE returns l <= r (lexicographic comparison).
func (l String) LE(r String) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_str_le(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// ToInt converts string l to an integer.
// Returns -1 if the string does not represent an integer.
func (l String) ToInt() Int {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_str_to_int(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Int(val)
}

// IntToString converts an integer to a string.
func (ctx *Context) IntToString(i Int) String {
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_int_to_str(ctx.c, i.c)
	})
	runtime.KeepAlive(i)
	return String(val)
}

// ToRE converts string l to a regular expression that matches exactly l.
func (l String) ToRE() RE {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_to_re(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return RE(val)
}

// InRE returns true if l is in the language of regular expression re.
func (l String) InRE(re RE) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_seq_in_re(ctx.c, l.c, re.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(re)
	return Bool(val)
}
