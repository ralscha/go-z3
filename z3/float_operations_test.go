// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"testing"
)

func TestFloatNE(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(1.5)
	y := ctx.Float32FromFloat64(2.5)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1.5 != 2.5")
	}
}

func TestFloatAbs(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(-3.5)
	result := x.Abs()

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(3.5)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for abs(-3.5) = 3.5")
	}
}

func TestFloatNeg(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(3.5)
	result := x.Neg()

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(-3.5)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for neg(3.5) = -3.5")
	}
}

func TestFloatAdd(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(1.5)
	y := ctx.Float32FromFloat64(2.5)
	result := x.Add(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(4.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 1.5 + 2.5 = 4.0")
	}
}

func TestFloatSub(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(5.0)
	y := ctx.Float32FromFloat64(2.0)
	result := x.Sub(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(3.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5.0 - 2.0 = 3.0")
	}
}

func TestFloatMul(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(2.0)
	y := ctx.Float32FromFloat64(3.0)
	result := x.Mul(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(6.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 2.0 * 3.0 = 6.0")
	}
}

func TestFloatDiv(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(6.0)
	y := ctx.Float32FromFloat64(2.0)
	result := x.Div(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(3.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 6.0 / 2.0 = 3.0")
	}
}

func TestFloatMulAdd(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(2.0)
	y := ctx.Float32FromFloat64(3.0)
	z := ctx.Float32FromFloat64(1.0)
	result := x.MulAdd(y, z) // x*y + z

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(7.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 2.0 * 3.0 + 1.0 = 7.0")
	}
}

func TestFloatSqrt(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(4.0)
	result := x.Sqrt()

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(2.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for sqrt(4.0) = 2.0")
	}
}

func TestFloatRem(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(5.0)
	y := ctx.Float32FromFloat64(2.0)
	result := x.Rem(y)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(ctx.Float32FromFloat64(1.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 5.0 rem 2.0 = 1.0")
	}
}

func TestFloatMinMax(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(3.0)
	y := ctx.Float32FromFloat64(5.0)

	minResult := x.Min(y)
	solver := NewSolver(ctx)
	solver.Assert(minResult.Eq(ctx.Float32FromFloat64(3.0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for min(3.0, 5.0) = 3.0")
	}

	maxResult := x.Max(y)
	solver2 := NewSolver(ctx)
	solver2.Assert(maxResult.Eq(ctx.Float32FromFloat64(5.0)))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for max(3.0, 5.0) = 5.0")
	}
}

func TestFloatComparisons(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.Float32FromFloat64(3.0)
	y := ctx.Float32FromFloat64(5.0)

	solver := NewSolver(ctx)
	solver.Assert(x.LT(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 3.0 < 5.0")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(x.LE(y))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for 3.0 <= 5.0")
	}

	solver3 := NewSolver(ctx)
	solver3.Assert(y.GT(x))
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for 5.0 > 3.0")
	}

	solver4 := NewSolver(ctx)
	solver4.Assert(y.GE(x))
	if sat, _ := solver4.Check(); !sat {
		t.Error("expected SAT for 5.0 >= 3.0")
	}

	solver5 := NewSolver(ctx)
	z := ctx.Float32FromFloat64(3.0)
	solver5.Assert(x.IEEEEq(z))
	if sat, _ := solver5.Check(); !sat {
		t.Error("expected SAT for 3.0 IEEE= 3.0")
	}
}

func TestFloatPredicates(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.Float32Sort()

	// Test IsNormal
	x := ctx.Float32FromFloat64(1.0)
	solver := NewSolver(ctx)
	solver.Assert(x.IsNormal())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for IsNormal(1.0)")
	}

	// Test IsZero
	zero := ctx.FloatZero(sort, false)
	solver2 := NewSolver(ctx)
	solver2.Assert(zero.IsZero())
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for IsZero(0.0)")
	}

	// Test IsPositive
	solver3 := NewSolver(ctx)
	solver3.Assert(x.IsPositive())
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for IsPositive(1.0)")
	}

	// Test IsNegative
	negX := ctx.Float32FromFloat64(-1.0)
	solver4 := NewSolver(ctx)
	solver4.Assert(negX.IsNegative())
	if sat, _ := solver4.Check(); !sat {
		t.Error("expected SAT for IsNegative(-1.0)")
	}

	// Test IsSubnormal
	// A subnormal is a very small number close to zero
	subnormal := ctx.Const("sub", sort).(Float)
	solver5 := NewSolver(ctx)
	solver5.Assert(subnormal.IsSubnormal())
	if sat, _ := solver5.Check(); !sat {
		t.Error("expected SAT for some subnormal")
	}

	// Test IsInfinite
	inf := ctx.FloatInf(sort, false)
	solver6 := NewSolver(ctx)
	solver6.Assert(inf.IsInfinite())
	if sat, _ := solver6.Check(); !sat {
		t.Error("expected SAT for IsInfinite(+inf)")
	}

	// Test IsNaN
	nan := ctx.FloatNaN(sort)
	solver7 := NewSolver(ctx)
	solver7.Assert(nan.IsNaN())
	if sat, _ := solver7.Check(); !sat {
		t.Error("expected SAT for IsNaN(NaN)")
	}
}

func TestFloatConversions(t *testing.T) {
	ctx := NewContext(nil)
	sort64 := ctx.Float64Sort() // float64

	// ToFloat (different precision)
	x := ctx.Float32FromFloat64(2.5)
	y := x.ToFloat(sort64)
	solver := NewSolver(ctx)
	solver.Assert(y.Eq(ctx.Float64(2.5)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for float32(2.5) to float64")
	}

	// ToUBV
	z := ctx.Float32FromFloat64(42.0)
	bvu := z.ToUBV(8)
	solver2 := NewSolver(ctx)
	solver2.Assert(bvu.Eq(ctx.FromInt(42, ctx.BVSort(8)).(BV)))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for float to unsigned BV")
	}

	// ToSBV
	w := ctx.Float32FromFloat64(-5.0)
	bvs := w.ToSBV(8)
	solver3 := NewSolver(ctx)
	solver3.Assert(bvs.Eq(ctx.FromInt(-5, ctx.BVSort(8)).(BV)))
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for float to signed BV")
	}

	// ToIEEEBV
	v := ctx.Float32FromFloat64(1.5)
	ieeeBV := v.ToIEEEBV()
	solver4 := NewSolver(ctx)
	// 1.5 in IEEE 754 float32 is 0x3FC00000
	solver4.Assert(ieeeBV.Eq(ctx.FromInt(0x3FC00000, ctx.BVSort(32)).(BV)))
	if sat, _ := solver4.Check(); !sat {
		t.Error("expected SAT for float to IEEE BV")
	}
}

func TestFloatNaN(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.Float32Sort()
	nan := ctx.FloatNaN(sort)

	solver := NewSolver(ctx)
	solver.Assert(nan.IsNaN())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for NaN.IsNaN()")
	}
}

func TestFloatInf(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.Float32Sort()
	posInf := ctx.FloatInf(sort, false)
	negInf := ctx.FloatInf(sort, true)

	solver := NewSolver(ctx)
	solver.Assert(posInf.IsInfinite())
	solver.Assert(posInf.IsPositive())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for +inf")
	}

	solver2 := NewSolver(ctx)
	solver2.Assert(negInf.IsInfinite())
	solver2.Assert(negInf.IsNegative())
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for -inf")
	}
}

func TestFromFloat32(t *testing.T) {
	ctx := NewContext(nil)
	f := ctx.Float32(2.5)

	solver := NewSolver(ctx)
	solver.Assert(f.Eq(ctx.Float32FromFloat64(2.5)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for FromFloat32")
	}
}

func TestBVToFloat(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.Float32Sort()

	// IEEEToFloat
	bv := ctx.FromInt(0x3FC00000, ctx.BVSort(32)).(BV) // 1.5 in IEEE 754
	f := bv.IEEEToFloat(sort)
	solver := NewSolver(ctx)
	solver.Assert(f.Eq(ctx.Float32FromFloat64(1.5)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for IEEEToFloat")
	}

	// SToFloat
	sBV := ctx.FromInt(-5, ctx.BVSort(8)).(BV)
	sf := sBV.SToFloat(sort)
	solver2 := NewSolver(ctx)
	solver2.Assert(sf.Eq(ctx.Float32FromFloat64(-5.0)))
	if sat, _ := solver2.Check(); !sat {
		t.Error("expected SAT for SToFloat")
	}

	// UToFloat
	uBV := ctx.FromInt(42, ctx.BVSort(8)).(BV)
	uf := uBV.UToFloat(sort)
	solver3 := NewSolver(ctx)
	solver3.Assert(uf.Eq(ctx.Float32FromFloat64(42.0)))
	if sat, _ := solver3.Check(); !sat {
		t.Error("expected SAT for UToFloat")
	}
}
