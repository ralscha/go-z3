// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestStringSort(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.StringSort()
	if sort.Kind() != KindSeq {
		t.Errorf("expected KindSeq, got %v", sort.Kind())
	}
}

func TestStringConst(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.StringSort()).(String)
	if x.Sort().Kind() != KindSeq {
		t.Errorf("expected KindSeq, got %v", x.Sort().Kind())
	}
}

func TestStringLiteral(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello")
	solver := NewSolver(ctx)
	solver.Assert(s.Eq(ctx.FromString("hello")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT")
	}
}

func TestStringConcat(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.FromString("hello")
	b := ctx.FromString(" world")
	c := a.Concat(b)

	solver := NewSolver(ctx)
	solver.Assert(c.Eq(ctx.FromString("hello world")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for concat")
	}
}

func TestStringLength(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello")
	length := s.Length()

	solver := NewSolver(ctx)
	solver.Assert(length.Eq(ctx.FromInt(5, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for length")
	}
}

func TestStringContains(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello world")
	sub := ctx.FromString("world")

	solver := NewSolver(ctx)
	solver.Assert(s.Contains(sub))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for contains")
	}
}

func TestStringPrefixSuffix(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello world")
	prefix := ctx.FromString("hello")
	suffix := ctx.FromString("world")

	solver := NewSolver(ctx)
	solver.Assert(prefix.PrefixOf(s)) // "hello" is a prefix of "hello world"
	solver.Assert(suffix.SuffixOf(s)) // "world" is a suffix of "hello world"
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for prefix/suffix")
	}
}

func TestStringExtract(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello world")
	offset := ctx.FromInt(0, ctx.IntSort()).(Int)
	length := ctx.FromInt(5, ctx.IntSort()).(Int)
	extracted := s.Extract(offset, length)

	solver := NewSolver(ctx)
	solver.Assert(extracted.Eq(ctx.FromString("hello")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for extract")
	}
}

func TestStringAt(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello")
	idx := ctx.FromInt(1, ctx.IntSort()).(Int)
	ch := s.At(idx)

	solver := NewSolver(ctx)
	solver.Assert(ch.Eq(ctx.FromString("e")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for at")
	}
}

func TestStringIndexOf(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello world")
	sub := ctx.FromString("world")
	offset := ctx.FromInt(0, ctx.IntSort()).(Int)
	idx := s.IndexOf(sub, offset)

	solver := NewSolver(ctx)
	solver.Assert(idx.Eq(ctx.FromInt(6, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for indexOf")
	}
}

func TestStringReplace(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello world")
	src := ctx.FromString("world")
	dst := ctx.FromString("z3")
	replaced := s.Replace(src, dst)

	solver := NewSolver(ctx)
	solver.Assert(replaced.Eq(ctx.FromString("hello z3")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for replace")
	}
}

func TestStringToInt(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("42")
	i := s.ToInt()

	solver := NewSolver(ctx)
	solver.Assert(i.Eq(ctx.FromInt(42, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for toInt")
	}
}

func TestIntToString(t *testing.T) {
	ctx := NewContext(nil)
	i := ctx.FromInt(42, ctx.IntSort()).(Int)
	s := ctx.IntToString(i)

	solver := NewSolver(ctx)
	solver.Assert(s.Eq(ctx.FromString("42")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for intToString")
	}
}

func TestStringSymbolic(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.StringSort()).(String)
	y := ctx.Const("y", ctx.StringSort()).(String)

	solver := NewSolver(ctx)
	// x + y = "hello"
	solver.Assert(x.Concat(y).Eq(ctx.FromString("hello")))
	// length(x) = 2
	solver.Assert(x.Length().Eq(ctx.FromInt(2, ctx.IntSort()).(Int)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for symbolic string")
	}

	model := solver.Model()
	xVal := model.Eval(x, true)
	yVal := model.Eval(y, true)
	t.Logf("x = %v, y = %v", xVal, yVal)
}
