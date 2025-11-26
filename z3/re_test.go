// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestRESort(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()
	reSort := ctx.RESort(strSort)
	if reSort.Kind() != KindRE {
		t.Errorf("expected KindRE, got %v", reSort.Kind())
	}
}

func TestStringToRE(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello")
	re := s.ToRE()

	// String "hello" should match the regex for exactly "hello"
	solver := NewSolver(ctx)
	solver.Assert(s.InRE(re))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for string in its own regex")
	}
}

func TestREInRE(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello")
	re := s.ToRE()

	// Test that "hello" is in the regex
	solver := NewSolver(ctx)
	solver.Assert(s.InRE(re))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT")
	}

	// Test that "world" is NOT in the regex for "hello"
	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("world").InRE(re))
	if sat, _ := solver2.Check(); sat {
		t.Error("expected UNSAT")
	}
}

func TestREStar(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: a*
	a := ctx.FromString("a").ToRE()
	aStar := a.Star()

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("aaa").InRE(aStar))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'aaa' in a*")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("").InRE(aStar))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for empty string in a*")
	}
}

func TestREPlus(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: a+
	a := ctx.FromString("a").ToRE()
	aPlus := a.Plus()

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("aaa").InRE(aPlus))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'aaa' in a+")
	}

	// Empty string should NOT match a+
	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("").InRE(aPlus))
	if sat, _ := solver2.Check(); sat {
		t.Error("expected UNSAT for empty string in a+")
	}
}

func TestREOption(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: a?
	a := ctx.FromString("a").ToRE()
	aOpt := a.Option()

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("a").InRE(aOpt))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'a' in a?")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("").InRE(aOpt))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for empty string in a?")
	}
}

func TestREConcat(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: ab
	a := ctx.FromString("a").ToRE()
	b := ctx.FromString("b").ToRE()
	ab := a.Concat(b)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("ab").InRE(ab))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'ab' in ab")
	}
}

func TestREUnion(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: a|b
	a := ctx.FromString("a").ToRE()
	b := ctx.FromString("b").ToRE()
	aOrB := a.Union(b)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("a").InRE(aOrB))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'a' in a|b")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("b").InRE(aOrB))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for 'b' in a|b")
	}
}

func TestRERange(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: [a-z]
	azRange := ctx.RERange(ctx.FromString("a"), ctx.FromString("z"))

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("m").InRE(azRange))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'm' in [a-z]")
	}
}

func TestRELoop(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: a{2,4}
	a := ctx.FromString("a").ToRE()
	a2to4 := a.Loop(2, 4)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("aa").InRE(a2to4))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'aa' in a{2,4}")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("aaaa").InRE(a2to4))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for 'aaaa' in a{2,4}")
	}

	// "a" should not match a{2,4}
	solver3 := NewSolver(ctx)
	solver3.Assert(ctx.FromString("a").InRE(a2to4))
	if sat, _ := solver3.Check(); sat {
		t.Error("expected UNSAT for 'a' in a{2,4}")
	}
}

func TestREComplement(t *testing.T) {
	ctx := NewContext(nil)
	// Regex: not(a)
	a := ctx.FromString("a").ToRE()
	notA := a.Complement()

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("b").InRE(notA))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'b' in complement(a)")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(ctx.FromString("a").InRE(notA))
	if sat, _ := solver2.Check(); sat {
		t.Error("expected UNSAT for 'a' in complement(a)")
	}
}

func TestRESymbolic(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.StringSort()).(String)

	// x matches a*b
	a := ctx.FromString("a").ToRE()
	b := ctx.FromString("b").ToRE()
	pattern := a.Star().Concat(b)

	solver := NewSolver(ctx)
	solver.Assert(x.InRE(pattern))
	solver.Assert(x.Length().Eq(ctx.FromInt(4, ctx.IntSort()).(Int)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for symbolic regex match")
	}

	model := solver.Model()
	xVal := model.Eval(x, true)
	t.Logf("x = %v", xVal)
}
