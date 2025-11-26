// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"math/big"
	"testing"
)

func TestIntNE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.IntSort()).(Int)
	y := ctx.FromInt(10, ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5 != 10")
	}
}

func TestIntDiv(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(23, ctx.IntSort()).(Int)
	y := ctx.FromInt(5, ctx.IntSort()).(Int)
	result := x.Div(y)

	solver := NewSolver(ctx)
	// Z3 integer division floors, so 23 / 5 = 4
	solver.Assert(result.Eq(ctx.FromInt(4, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 23 div 5 = 4")
	}
}

func TestIntMod(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(23, ctx.IntSort()).(Int)
	y := ctx.FromInt(5, ctx.IntSort()).(Int)
	result := x.Mod(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(3, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 23 mod 5 = 3")
	}
}

func TestIntRem(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(-23, ctx.IntSort()).(Int)
	y := ctx.FromInt(5, ctx.IntSort()).(Int)
	result := x.Rem(y)

	solver := NewSolver(ctx)
	// Z3's Rem is based on floored division, not truncated division like Go's %.
	// For floored division: -23 / 5 = -5, so -23 rem 5 = -23 - (-5 * 5) = 2
	solver.Assert(result.Eq(ctx.FromInt(2, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for -23 rem 5 = 2")
	}
}

func TestIntMul(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.IntSort()).(Int)
	y := ctx.FromInt(4, ctx.IntSort()).(Int)
	result := x.Mul(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(20, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5 * 4 = 20")
	}
}

func TestIntSub(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(10, ctx.IntSort()).(Int)
	y := ctx.FromInt(3, ctx.IntSort()).(Int)
	result := x.Sub(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(7, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 10 - 3 = 7")
	}
}

func TestIntNeg(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.IntSort()).(Int)
	result := x.Neg()

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(-5, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for -5 = -5")
	}
}

func TestIntExp(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(2, ctx.IntSort()).(Int)
	y := ctx.FromInt(3, ctx.IntSort()).(Int)
	result := x.Exp(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(8, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 2^3 = 8")
	}
}

func TestIntLT(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.IntSort()).(Int)
	y := ctx.FromInt(10, ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	solver.Assert(x.LT(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5 < 10")
	}
}

func TestIntGT(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(10, ctx.IntSort()).(Int)
	y := ctx.FromInt(5, ctx.IntSort()).(Int)

	solver := NewSolver(ctx)
	solver.Assert(x.GT(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 10 > 5")
	}
}

func TestIntToBV(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(42, ctx.IntSort()).(Int)
	bv := x.ToBV(8)

	solver := NewSolver(ctx)
	solver.Assert(bv.Eq(ctx.FromInt(42, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for ToBV(42) = 42:8")
	}
}

func TestIntAsUint64(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(42, ctx.IntSort()).(Int)
	val, _, ok := x.AsUint64()
	if !ok {
		t.Error("expected ok for AsUint64")
	}
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

// Real tests

func TestRealNE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(1, 2))
	y := ctx.FromBigRat(big.NewRat(1, 3))

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1/2 != 1/3")
	}
}

func TestRealMul(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(1, 2))
	y := ctx.FromBigRat(big.NewRat(2, 3))
	result := x.Mul(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromBigRat(big.NewRat(1, 3))))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1/2 * 2/3 = 1/3")
	}
}

func TestRealSub(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(3, 4))
	y := ctx.FromBigRat(big.NewRat(1, 4))
	result := x.Sub(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromBigRat(big.NewRat(1, 2))))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 3/4 - 1/4 = 1/2")
	}
}

func TestRealLT(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(1, 3))
	y := ctx.FromBigRat(big.NewRat(1, 2))

	solver := NewSolver(ctx)
	solver.Assert(x.LT(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1/3 < 1/2")
	}
}

func TestRealLE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(1, 2))
	y := ctx.FromBigRat(big.NewRat(1, 2))

	solver := NewSolver(ctx)
	solver.Assert(x.LE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1/2 <= 1/2")
	}
}

func TestRealGT(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(2, 3))
	y := ctx.FromBigRat(big.NewRat(1, 2))

	solver := NewSolver(ctx)
	solver.Assert(x.GT(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 2/3 > 1/2")
	}
}

func TestRealGE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(1, 2))
	y := ctx.FromBigRat(big.NewRat(1, 2))

	solver := NewSolver(ctx)
	solver.Assert(x.GE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1/2 >= 1/2")
	}
}

func TestRealIsInt(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(4, 2)) // 2

	solver := NewSolver(ctx)
	solver.Assert(x.IsInt())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 4/2 is integer")
	}
}

func TestRealToFloatExp(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBigRat(big.NewRat(5, 2)) // 2.5
	floatSort := ctx.FloatSort(8, 24)     // float32
	exp := ctx.FromInt(0, ctx.IntSort()).(Int)
	f := x.ToFloatExp(exp, floatSort)

	solver := NewSolver(ctx)
	solver.Assert(f.Eq(ctx.FromFloat64(2.5, floatSort)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for real to float")
	}
}

// Logic tests

func TestBoolNE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBool(true)
	y := ctx.FromBool(false)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for true != false")
	}
}

func TestBoolIfThenElse(t *testing.T) {
	ctx := NewContext(nil)
	cond := ctx.FromBool(true)
	x := ctx.FromInt(1, ctx.IntSort()).(Int)
	y := ctx.FromInt(2, ctx.IntSort()).(Int)
	result := cond.IfThenElse(x, y).(Int)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(x))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for if true then 1 else 2 = 1")
	}
}

func TestBoolXor(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromBool(true)
	y := ctx.FromBool(false)
	result := x.Xor(y)

	solver := NewSolver(ctx)
	solver.Assert(result)
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for true XOR false = true")
	}
}
