// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestCharSort(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.CharSort()
	if sort.Kind() != KindChar {
		t.Errorf("expected KindChar, got %v", sort.Kind())
	}
}

func TestCharToInt(t *testing.T) {
	ctx := NewContext(nil)
	// Create a char constant and test ToInt
	c := ctx.Const("c", ctx.CharSort()).(Char)
	cInt := c.ToInt()

	solver := NewSolver(ctx)
	// ASCII 'A' = 65
	solver.Assert(cInt.Eq(ctx.Int(65)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT")
	}
}

func TestCharLE(t *testing.T) {
	ctx := NewContext(nil)
	a := ctx.Const("a", ctx.CharSort()).(Char)
	b := ctx.Const("b", ctx.CharSort()).(Char)

	solver := NewSolver(ctx)
	solver.Assert(a.LE(b))
	solver.Assert(a.ToInt().Eq(ctx.Int(65))) // 'A'
	solver.Assert(b.ToInt().Eq(ctx.Int(66))) // 'B'

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for A <= B")
	}
}

func TestIsDigit(t *testing.T) {
	ctx := NewContext(nil)
	c := ctx.Const("c", ctx.CharSort()).(Char)

	solver := NewSolver(ctx)
	solver.Assert(c.IsDigit())
	// Digits are '0'-'9' (48-57)
	solver.Assert(c.ToInt().Eq(ctx.Int(53))) // '5'

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for digit '5'")
	}

	// Non-digit should fail
	solver2 := NewSolver(ctx)
	solver2.Assert(c.IsDigit())
	solver2.Assert(c.ToInt().Eq(ctx.Int(65))) // 'A'

	if sat, _ := solver2.Check(); sat {
		t.Error("expected UNSAT for non-digit 'A'")
	}
}

func TestCharToBV(t *testing.T) {
	ctx := NewContext(nil)
	c := ctx.Const("c", ctx.CharSort()).(Char)
	bv := c.ToBV()

	solver := NewSolver(ctx)
	solver.Assert(c.ToInt().Eq(ctx.Int(65)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT")
	}

	// Z3 uses 18-bit BVs for character code points (Unicode BMP is 16-bit, but extended needs more)
	// The exact size may vary by Z3 version
	bvSize := bv.Sort().BVSize()
	if bvSize == 0 {
		t.Error("expected non-zero BV size")
	}
	t.Logf("Char BV size: %d bits", bvSize)
}

func TestToCode(t *testing.T) {
	ctx := NewContext(nil)
	// Single character string
	s := ctx.FromString("A")
	code := s.ToCode()

	solver := NewSolver(ctx)
	solver.Assert(code.Eq(ctx.Int(65)))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for ToCode('A') = 65")
	}
}

func TestStringFromCode(t *testing.T) {
	ctx := NewContext(nil)
	code := ctx.Int(65)
	s := ctx.StringFromCode(code)

	solver := NewSolver(ctx)
	solver.Assert(s.Eq(ctx.FromString("A")))

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for StringFromCode(65) = 'A'")
	}
}
