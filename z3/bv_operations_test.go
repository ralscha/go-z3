// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"testing"
)

func TestBVNot(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	notX := x.Not()

	solver := NewSolver(ctx)
	solver.Assert(notX.Eq(ctx.FromInt(0, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for NOT 0xFF = 0x00")
	}
}

func TestBVNE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(10, ctx.BVSort(8)).(BV)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5 != 10")
	}
}

func TestBVAllBits(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	allBits := x.AllBits()

	solver := NewSolver(ctx)
	// AllBits of 0xFF should be 1 (all bits are 1)
	solver.Assert(allBits.Eq(ctx.FromInt(1, ctx.BVSort(1)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for AllBits(0xFF) = 1")
	}
}

func TestBVAnyBits(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0x80, ctx.BVSort(8)).(BV)
	anyBits := x.AnyBits()

	solver := NewSolver(ctx)
	// AnyBits of 0x80 should be 1 (at least one bit is 1)
	solver.Assert(anyBits.Eq(ctx.FromInt(1, ctx.BVSort(1)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for AnyBits(0x80) = 1")
	}
}

func TestBVAnd(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xF0, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0x0F, ctx.BVSort(8)).(BV)
	result := x.And(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 0xF0 & 0x0F = 0")
	}
}

func TestBVOr(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xF0, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0x0F, ctx.BVSort(8)).(BV)
	result := x.Or(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 0xF0 | 0x0F = 0xFF")
	}
}

func TestBVXor(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0x0F, ctx.BVSort(8)).(BV)
	result := x.Xor(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0xF0, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 0xFF ^ 0x0F = 0xF0")
	}
}

func TestBVNand(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	result := x.Nand(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for NAND(0xFF, 0xFF) = 0")
	}
}

func TestBVNor(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0, ctx.BVSort(8)).(BV)
	result := x.Nor(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for NOR(0, 0) = 0xFF")
	}
}

func TestBVXnor(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	result := x.Xnor(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for XNOR(0xFF, 0xFF) = 0xFF")
	}
}

func TestBVNeg(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(1, ctx.BVSort(8)).(BV)
	result := x.Neg()

	solver := NewSolver(ctx)
	// -1 in 8-bit two's complement is 0xFF
	solver.Assert(result.Eq(ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for -1 = 0xFF")
	}
}

func TestBVAdd(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(10, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(20, ctx.BVSort(8)).(BV)
	result := x.Add(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(30, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 10 + 20 = 30")
	}
}

func TestBVSub(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(30, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(10, ctx.BVSort(8)).(BV)
	result := x.Sub(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(20, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 30 - 10 = 20")
	}
}

func TestBVMul(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(4, ctx.BVSort(8)).(BV)
	result := x.Mul(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(20, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5 * 4 = 20")
	}
}

func TestBVUDiv(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(20, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(4, ctx.BVSort(8)).(BV)
	result := x.UDiv(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(5, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 20 / 4 = 5")
	}
}

func TestBVSDiv(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(-20, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(4, ctx.BVSort(8)).(BV)
	result := x.SDiv(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(-5, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for -20 / 4 = -5")
	}
}

func TestBVURem(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(23, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(5, ctx.BVSort(8)).(BV)
	result := x.URem(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(3, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 23 % 5 = 3")
	}
}

func TestBVSRem(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(-23, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(5, ctx.BVSort(8)).(BV)
	result := x.SRem(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(-3, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for -23 % 5 = -3")
	}
}

func TestBVSMod(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(-23, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(5, ctx.BVSort(8)).(BV)
	result := x.SMod(y)

	solver := NewSolver(ctx)
	// SMod follows the sign of the divisor
	solver.Assert(result.Eq(ctx.FromInt(2, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for -23 mod 5 = 2")
	}
}

func TestBVComparisons(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(5, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(10, ctx.BVSort(8)).(BV)

	solver := NewSolver(ctx)
	solver.Assert(x.ULT(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5 < 10 (unsigned)")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(x.SLT(y))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for 5 < 10 (signed)")
	}

	solver3 := NewSolver(ctx)
	solver3.Assert(x.ULE(y))
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for 5 <= 10 (unsigned)")
	}

	solver4 := NewSolver(ctx)
	solver4.Assert(x.SLE(y))
	if sat, _ := solver4.Check(); !sat {
		t.Error("expected SAT for 5 <= 10 (signed)")
	}

	solver5 := NewSolver(ctx)
	solver5.Assert(y.UGT(x))
	if sat, _ := solver5.Check(); !sat {
		t.Error("expected SAT for 10 > 5 (unsigned)")
	}

	solver6 := NewSolver(ctx)
	solver6.Assert(y.SGT(x))
	if sat, _ := solver6.Check(); !sat {
		t.Error("expected SAT for 10 > 5 (signed)")
	}

	solver7 := NewSolver(ctx)
	solver7.Assert(y.UGE(x))
	if sat, _ := solver7.Check(); !sat {
		t.Error("expected SAT for 10 >= 5 (unsigned)")
	}

	solver8 := NewSolver(ctx)
	solver8.Assert(y.SGE(x))
	if sat, _ := solver8.Check(); !sat {
		t.Error("expected SAT for 10 >= 5 (signed)")
	}
}

func TestBVConcat(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xAB, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(0xCD, ctx.BVSort(8)).(BV)
	result := x.Concat(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0xABCD, ctx.BVSort(16)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for concat(0xAB, 0xCD) = 0xABCD")
	}
}

func TestBVExtract(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xABCD, ctx.BVSort(16)).(BV)
	high := x.Extract(15, 8)
	low := x.Extract(7, 0)

	solver := NewSolver(ctx)
	solver.Assert(high.Eq(ctx.FromInt(0xAB, ctx.BVSort(8)).(BV)))
	solver.Assert(low.Eq(ctx.FromInt(0xCD, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for extract operations")
	}
}

func TestBVSignExtend(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(-1, ctx.BVSort(8)).(BV) // 0xFF
	result := x.SignExtend(8)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(-1, ctx.BVSort(16)).(BV))) // 0xFFFF
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for sign extend -1:8 to -1:16")
	}
}

func TestBVZeroExtend(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xFF, ctx.BVSort(8)).(BV)
	result := x.ZeroExtend(8)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0x00FF, ctx.BVSort(16)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for zero extend 0xFF:8 to 0x00FF:16")
	}
}

func TestBVRepeat(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0xAB, ctx.BVSort(8)).(BV)
	result := x.Repeat(2)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.FromInt(0xABAB, ctx.BVSort(16)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for repeat(0xAB, 2) = 0xABAB")
	}
}

func TestBVShifts(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(1, ctx.BVSort(8)).(BV)
	shift := ctx.FromInt(4, ctx.BVSort(8)).(BV)

	// Left shift
	lshResult := x.Lsh(shift)
	solver := NewSolver(ctx)
	solver.Assert(lshResult.Eq(ctx.FromInt(16, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1 << 4 = 16")
	}

	// Unsigned right shift
	y := ctx.FromInt(16, ctx.BVSort(8)).(BV)
	urshResult := y.URsh(shift)
	solver2 := NewSolver(ctx)
	solver2.Assert(urshResult.Eq(ctx.FromInt(1, ctx.BVSort(8)).(BV)))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for 16 >> 4 = 1 (unsigned)")
	}

	// Signed right shift
	z := ctx.FromInt(-16, ctx.BVSort(8)).(BV)
	shift2 := ctx.FromInt(2, ctx.BVSort(8)).(BV)
	srshResult := z.SRsh(shift2)
	solver3 := NewSolver(ctx)
	solver3.Assert(srshResult.Eq(ctx.FromInt(-4, ctx.BVSort(8)).(BV)))
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for -16 >> 2 = -4 (signed)")
	}
}

func TestBVRotate(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(0x81, ctx.BVSort(8)).(BV) // 10000001
	shift := ctx.FromInt(1, ctx.BVSort(8)).(BV)

	rotLeft := x.RotateLeft(shift)
	solver := NewSolver(ctx)
	// 10000001 rotated left by 1 = 00000011
	solver.Assert(rotLeft.Eq(ctx.FromInt(0x03, ctx.BVSort(8)).(BV)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for rotate left")
	}

	rotRight := x.RotateRight(shift)
	solver2 := NewSolver(ctx)
	// 10000001 rotated right by 1 = 11000000
	solver2.Assert(rotRight.Eq(ctx.FromInt(0xC0, ctx.BVSort(8)).(BV)))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for rotate right")
	}
}

func TestBVToInt(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(42, ctx.BVSort(8)).(BV)

	uintResult := x.UToInt()
	solver := NewSolver(ctx)
	solver.Assert(uintResult.Eq(ctx.FromInt(42, ctx.IntSort()).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for UToInt(42) = 42")
	}

	y := ctx.FromInt(-1, ctx.BVSort(8)).(BV)
	sintResult := y.SToInt()
	solver2 := NewSolver(ctx)
	solver2.Assert(sintResult.Eq(ctx.FromInt(-1, ctx.IntSort()).(Int)))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for SToInt(-1) = -1")
	}
}

func TestBVOverflowChecks(t *testing.T) {
	ctx := NewContext(nil)

	// AddNoUnderflow
	x := ctx.FromInt(-100, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(-50, ctx.BVSort(8)).(BV)
	addNoUnder := x.AddNoUnderflow(y)
	solver := NewSolver(ctx)
	solver.Assert(addNoUnder.Not())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for add underflow case")
	}

	// SubNoOverflow
	a := ctx.FromInt(100, ctx.BVSort(8)).(BV)
	b := ctx.FromInt(-50, ctx.BVSort(8)).(BV)
	subNoOver := a.SubNoOverflow(b)
	solver2 := NewSolver(ctx)
	solver2.Assert(subNoOver.Not())
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for sub overflow case")
	}

	// MulNoUnderflow
	c := ctx.FromInt(-100, ctx.BVSort(8)).(BV)
	d := ctx.FromInt(3, ctx.BVSort(8)).(BV)
	mulNoUnder := c.MulNoUnderflow(d)
	solver3 := NewSolver(ctx)
	solver3.Assert(mulNoUnder.Not())
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for mul underflow case")
	}

	// SDivNoOverflow (e.g., MIN_INT / -1)
	minInt := ctx.FromInt(-128, ctx.BVSort(8)).(BV)
	negOne := ctx.FromInt(-1, ctx.BVSort(8)).(BV)
	sdivNoOver := minInt.SDivNoOverflow(negOne)
	solver4 := NewSolver(ctx)
	solver4.Assert(sdivNoOver.Not())
	if sat, _ := solver4.Check(); !sat {
		t.Error("expected SAT for sdiv overflow case")
	}

	// NegNoOverflow
	negNoOver := minInt.NegNoOverflow()
	solver5 := NewSolver(ctx)
	solver5.Assert(negNoOver.Not())
	if sat, _ := solver5.Check(); !sat {
		t.Error("expected SAT for neg overflow case")
	}
}

func TestSubNoUnderflow(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.FromInt(-100, ctx.BVSort(8)).(BV)
	y := ctx.FromInt(50, ctx.BVSort(8)).(BV)
	noUnder := x.SubNoUnderflow(y, true)

	solver := NewSolver(ctx)
	solver.Assert(noUnder.Not())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for sub underflow case")
	}
}
