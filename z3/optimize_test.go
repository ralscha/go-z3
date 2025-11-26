// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"testing"
)

func TestOptimize(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	// Set pareto priority mode
	config := NewContextConfig()
	config.SetString("priority", "pareto")
	opt.SetParams(config)

	x := ctx.IntConst("x")
	y := ctx.IntConst("y")
	zero := ctx.FromInt(0, ctx.IntSort()).(Int)
	ten := ctx.FromInt(10, ctx.IntSort()).(Int)
	eleven := ctx.FromInt(11, ctx.IntSort()).(Int)

	opt.Assert(ten.GE(x).And(x.GE(zero)))
	opt.Assert(ten.GE(y).And(y.GE(zero)))
	opt.Assert(x.Add(y).LE(eleven))

	h1 := opt.Maximize(x)
	h2 := opt.Maximize(y)

	const TotalSolutions = 10
	var solutions int
	for {
		if sat, err := opt.Check(); sat {
			t.Log("x: ", h1.Lower(), ", y: ", h2.Lower())
			solutions++
		} else if err != nil {
			t.Fatalf("error: %s", err)
		} else if solutions > TotalSolutions {
			t.Fatalf("Too many solutions found (expected %d, found %d)\n",
				TotalSolutions, solutions)
		} else {
			break
		}
	}
}

// Based on an example from the z3 optimization tutorial
func TestOptimizeSoft(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	opt.AssertSoft(a, "1", "A")
	opt.AssertSoft(b, "2", "B")
	opt.AssertSoft(c, "3", "A")
	opt.Assert(a.Eq(c))
	opt.Assert(a.And(b).Not())

	if sat, err := opt.Check(); sat {
		model := opt.Model()
		if val, _ := model.Eval(c, false).(Bool).AsBool(); !val {
			t.Fatal("c has wrong value")
		}
		if val, _ := model.Eval(b, false).(Bool).AsBool(); val {
			t.Fatal("b has wrong value")
		}
		if val, _ := model.Eval(a, false).(Bool).AsBool(); !val {
			t.Fatal("a has wrong value")
		}
	} else if err != nil {
		t.Fatalf("error: %s", err)
	}
}

func TestOptimizeMinimize(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.IntConst("x")
	y := ctx.IntConst("y")
	zero := ctx.FromInt(0, ctx.IntSort()).(Int)
	ten := ctx.FromInt(10, ctx.IntSort()).(Int)

	// x >= 0, y >= 0, x + y >= 5
	opt.Assert(x.GE(zero))
	opt.Assert(y.GE(zero))
	opt.Assert(x.Add(y).GE(ctx.FromInt(5, ctx.IntSort()).(Int)))
	opt.Assert(x.LE(ten))
	opt.Assert(y.LE(ten))

	// Minimize x + y
	obj := opt.Minimize(x.Add(y))

	sat, err := opt.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}

	// The minimum value of x + y should be 5
	lower := obj.Lower()
	if lower.String() != "5" {
		t.Fatalf("expected minimum of 5, got %s", lower)
	}
}

func TestOptimizePushPop(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.IntConst("x")
	zero := ctx.FromInt(0, ctx.IntSort()).(Int)
	ten := ctx.FromInt(10, ctx.IntSort()).(Int)
	five := ctx.FromInt(5, ctx.IntSort()).(Int)

	opt.Assert(x.GE(zero))
	opt.Assert(x.LE(ten))

	obj := opt.Maximize(x)

	sat, err := opt.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}
	if obj.Upper().String() != "10" {
		t.Fatalf("expected max of 10, got %s", obj.Upper())
	}

	// Push and add more constraints
	opt.Push()
	opt.Assert(x.LE(five))
	obj2 := opt.Maximize(x)

	sat, err = opt.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}
	if obj2.Upper().String() != "5" {
		t.Fatalf("expected max of 5, got %s", obj2.Upper())
	}

	// Pop and check original constraints
	opt.Pop()
	obj3 := opt.Maximize(x)

	sat, err = opt.Check()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !sat {
		t.Fatal("expected satisfiable")
	}
	if obj3.Upper().String() != "10" {
		t.Fatalf("expected max of 10 after pop, got %s", obj3.Upper())
	}
}

func TestOptimizeString(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.IntConst("x")
	opt.Assert(x.GE(ctx.FromInt(0, ctx.IntSort()).(Int)))

	s := opt.String()
	if s == "" {
		t.Fatal("expected non-empty string representation")
	}
	t.Log("Optimize string representation:", s)
}

func TestOptimizeAssertions(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.IntConst("x")
	y := ctx.IntConst("y")
	zero := ctx.FromInt(0, ctx.IntSort()).(Int)

	opt.Assert(x.GE(zero))
	opt.Assert(y.GE(zero))

	assertions := opt.Assertions()
	if len(assertions) != 2 {
		t.Fatalf("expected 2 assertions, got %d", len(assertions))
	}
}
