// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestAtMost(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// At most 1 of a, b, c can be true
	solver := NewSolver(ctx)
	solver.Assert(ctx.AtMost([]Bool{a, b, c}, 1))
	solver.Assert(a)
	solver.Assert(b)

	// This should be UNSAT because we're asserting 2 bools are true
	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT when 2 bools are true with AtMost(1)")
	}
}

func TestAtMostSat(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// At most 2 of a, b, c can be true
	solver := NewSolver(ctx)
	solver.Assert(ctx.AtMost([]Bool{a, b, c}, 2))
	solver.Assert(a)
	solver.Assert(b)

	// This should be SAT because 2 <= 2
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT when 2 bools are true with AtMost(2)")
	}
}

func TestAtLeast(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// At least 2 of a, b, c must be true
	solver := NewSolver(ctx)
	solver.Assert(ctx.AtLeast([]Bool{a, b, c}, 2))
	solver.Assert(a.Not())
	solver.Assert(b.Not())

	// This should be UNSAT because at most 1 can be true
	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT when at most 1 bool can be true with AtLeast(2)")
	}
}

func TestAtLeastSat(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// At least 2 of a, b, c must be true
	solver := NewSolver(ctx)
	solver.Assert(ctx.AtLeast([]Bool{a, b, c}, 2))
	solver.Assert(a)
	solver.Assert(b)

	// This should be SAT
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT when 2 bools are true with AtLeast(2)")
	}
}

func TestPbEq(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// Exactly 2 of a, b, c must be true (with weights 1,1,1)
	solver := NewSolver(ctx)
	solver.Assert(ctx.PbEq([]Bool{a, b, c}, []int{1, 1, 1}, 2))
	solver.Assert(a)
	solver.Assert(b)
	solver.Assert(c.Not())

	// This should be SAT
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for exactly 2 true")
	}
}

func TestPbEqWeighted(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")

	// a*2 + b*3 = 5
	solver := NewSolver(ctx)
	solver.Assert(ctx.PbEq([]Bool{a, b}, []int{2, 3}, 5))
	solver.Assert(a)
	solver.Assert(b)

	// This should be SAT (2 + 3 = 5)
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 2+3=5")
	}
}

func TestPbLE(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// a*1 + b*2 + c*3 <= 4
	solver := NewSolver(ctx)
	solver.Assert(ctx.PbLE([]Bool{a, b, c}, []int{1, 2, 3}, 4))
	solver.Assert(a)
	solver.Assert(b)
	solver.Assert(c)

	// This should be UNSAT (1+2+3 = 6 > 4)
	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT for 1+2+3 <= 4")
	}
}

func TestPbGE(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	// a*1 + b*2 + c*3 >= 5
	solver := NewSolver(ctx)
	solver.Assert(ctx.PbGE([]Bool{a, b, c}, []int{1, 2, 3}, 5))
	solver.Assert(b)
	solver.Assert(c)

	// This should be SAT (2+3 = 5 >= 5)
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 2+3 >= 5")
	}
}
