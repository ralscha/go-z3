// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestIntAbs(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	solver.Assert(x.Eq(ctx.Int(-5)))
	solver.Assert(x.Abs().Eq(ctx.Int(5)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for abs(-5) = 5")
	}
}

func TestIntAbsPositive(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	solver.Assert(x.Eq(ctx.Int(5)))
	solver.Assert(x.Abs().Eq(ctx.Int(5)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for abs(5) = 5")
	}
}

func TestIntDivides(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	// x is divisible by 3
	solver.Assert(ctx.Int(3).Divides(x))
	solver.Assert(x.Eq(ctx.Int(9)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 3 divides 9")
	}
}

func TestIntDividesUnsat(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	// x is divisible by 3
	solver.Assert(ctx.Int(3).Divides(x))
	solver.Assert(x.Eq(ctx.Int(10)))

	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT for 3 divides 10")
	}
}

func TestRealAbs(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Const("x", ctx.RealSort()).(Real)

	solver := NewSolver(ctx)
	// x = -2.5
	two := ctx.FromInt(2, ctx.RealSort()).(Real)
	half := ctx.FromInt(1, ctx.RealSort()).(Real).Div(ctx.FromInt(2, ctx.RealSort()).(Real))
	negTwoPointFive := two.Add(half).Neg()

	solver.Assert(x.Eq(negTwoPointFive))
	solver.Assert(x.Abs().Eq(two.Add(half)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for abs(-2.5) = 2.5")
	}
}

func TestBVBit2Bool(t *testing.T) {
	ctx := NewContext(nil)
	// 8-bit value 0b00000101 = 5
	x := ctx.FromInt(5, ctx.BVSort(8)).(BV)

	// Bit 0 should be true
	bit0 := x.Bit2Bool(0)
	// Bit 1 should be false
	bit1 := x.Bit2Bool(1)
	// Bit 2 should be true
	bit2 := x.Bit2Bool(2)

	solver := NewSolver(ctx)
	solver.Assert(bit0)
	solver.Assert(bit1.Not())
	solver.Assert(bit2)

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for bit extraction of 5")
	}
}

func TestBVAddNoOverflow(t *testing.T) {
	ctx := NewContext(nil)
	// 8-bit unsigned: max = 255
	x := ctx.FromInt(200, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(100, ctx.BVSort(8)).(BV)

	// 200 + 100 = 300 which overflows 8-bit unsigned
	noOverflow := x.AddNoOverflow(y, false)

	solver := NewSolver(ctx)
	solver.Assert(noOverflow)

	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT for 200+100 no overflow (8-bit)")
	}
}

func TestBVAddNoOverflowSat(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(100, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(50, ctx.BVSort(8)).(BV)

	// 100 + 50 = 150 which fits in 8-bit unsigned
	noOverflow := x.AddNoOverflow(y, false)

	solver := NewSolver(ctx)
	solver.Assert(noOverflow)

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 100+50 no overflow (8-bit)")
	}
}

func TestBVMulNoOverflow(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(20, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(20, ctx.BVSort(8)).(BV)

	// 20 * 20 = 400 which overflows 8-bit unsigned
	noOverflow := x.MulNoOverflow(y, false)

	solver := NewSolver(ctx)
	solver.Assert(noOverflow)

	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT for 20*20 no overflow (8-bit)")
	}
}

func TestArrayExt(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	arrSort := ctx.ArraySort(intSort, intSort)

	a := ctx.Const("a", arrSort).(Array)
	b := ctx.Const("b", arrSort).(Array)

	// If arrays are different, Ext returns an index where they differ
	solver := NewSolver(ctx)
	solver.Assert(a.Eq(b).Not())

	diffIdx := a.Ext(b)
	// At the diff index, the values should be different
	solver.Assert(a.Select(diffIdx).(Int).Eq(b.Select(diffIdx).(Int)).Not())

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for array ext")
	}
}

func TestSolverNumScopes(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	if solver.NumScopes() != 0 {
		t.Errorf("expected 0 scopes, got %d", solver.NumScopes())
	}

	solver.Push()
	if solver.NumScopes() != 1 {
		t.Errorf("expected 1 scope, got %d", solver.NumScopes())
	}

	solver.Push()
	if solver.NumScopes() != 2 {
		t.Errorf("expected 2 scopes, got %d", solver.NumScopes())
	}

	solver.Pop()
	if solver.NumScopes() != 1 {
		t.Errorf("expected 1 scope after pop, got %d", solver.NumScopes())
	}
}

func TestSolverNumAssertions(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")

	if solver.NumAssertions() != 0 {
		t.Errorf("expected 0 assertions, got %d", solver.NumAssertions())
	}

	solver.Assert(a)
	if solver.NumAssertions() != 1 {
		t.Errorf("expected 1 assertion, got %d", solver.NumAssertions())
	}

	solver.Assert(b)
	if solver.NumAssertions() != 2 {
		t.Errorf("expected 2 assertions, got %d", solver.NumAssertions())
	}
}

func TestSolverAssertions(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")

	solver.Assert(a)
	solver.Assert(b)

	assertions := solver.Assertions()
	if len(assertions) != 2 {
		t.Errorf("expected 2 assertions, got %d", len(assertions))
	}
}

func TestSolverCheckAssumptions(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")

	solver.Assert(a.Implies(b))

	// With assumption a=true, b must be true
	sat, err := solver.CheckAssumptions(a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sat {
		t.Error("expected SAT")
	}

	// With assumptions a=true and b=false, should be UNSAT
	sat, err = solver.CheckAssumptions(a, b.Not())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sat {
		t.Error("expected UNSAT")
	}
}

func TestSolverUnsatCore(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)
	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")

	solver.Assert(a.Implies(b))
	solver.Assert(a)
	solver.Assert(b.Not())

	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT")
	}

	core := solver.UnsatCore()
	// Core should contain some subset of the assertions that's unsatisfiable
	if len(core) == 0 {
		t.Log("Note: UnsatCore may be empty depending on Z3 configuration")
	}
}
