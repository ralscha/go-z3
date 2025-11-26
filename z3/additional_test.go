// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestArrayNE(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	arrSort := ctx.ArraySort(intSort, intSort)

	x := ctx.Const("x", arrSort).(Array)
	y := ctx.Const("y", arrSort).(Array)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for array inequality")
	}
}

func TestArrayStore(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	arrSort := ctx.ArraySort(intSort, intSort)

	arr := ctx.Const("arr", arrSort).(Array)
	idx := ctx.FromInt(0, intSort)
	val := ctx.FromInt(42, intSort)

	arr2 := arr.Store(idx, val)
	selected := arr2.Select(idx).(Int)

	solver := NewSolver(ctx)
	solver.Assert(selected.Eq(ctx.FromInt(42, intSort).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for array store/select")
	}
}

func TestArrayDefault(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	defVal := ctx.FromInt(0, intSort)
	arr := ctx.ConstArray(intSort, defVal)

	def := arr.Default().(Int)
	solver := NewSolver(ctx)
	solver.Assert(def.Eq(ctx.FromInt(0, intSort).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for array default")
	}
}

func TestConstArray(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	defVal := ctx.FromInt(99, intSort)
	arr := ctx.ConstArray(intSort, defVal)

	// All elements should be 99
	idx := ctx.FromInt(123, intSort)
	selected := arr.Select(idx).(Int)

	solver := NewSolver(ctx)
	solver.Assert(selected.Eq(ctx.FromInt(99, intSort).(Int)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for const array")
	}
}

func TestArrayMap(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	arrSort := ctx.ArraySort(intSort, intSort)

	arr1 := ctx.Const("arr1", arrSort).(Array)
	arr2 := ctx.Const("arr2", arrSort).(Array)

	// Create a function declaration for addition
	addFunc := ctx.FuncDecl("+add", []Sort{intSort, intSort}, intSort)

	// Map the add function over the arrays
	result := ctx.ArrayMap(addFunc, arr1, arr2)

	solver := NewSolver(ctx)
	solver.Assert(result.Eq(result)) // Just check it works
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for array map")
	}
}

func TestAsArray(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()

	// Create a function and convert to array
	f := ctx.FuncDecl("f", []Sort{intSort}, intSort)
	arr := ctx.AsArray(f)

	solver := NewSolver(ctx)
	solver.Assert(arr.Eq(arr)) // Just check it works
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for as-array")
	}
}

func TestSolverReset(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	x := ctx.BoolConst("x")
	solver.Assert(x)
	solver.Assert(x.Not())

	if sat, _ := solver.Check(); sat {
		t.Error("expected UNSAT before reset")
	}

	solver.Reset()
	solver.Assert(x)

	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT after reset")
	}
}

func TestUninterpNE(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.UninterpretedSort("T")

	x := ctx.Const("x", sort).(Uninterpreted)
	y := ctx.Const("y", sort).(Uninterpreted)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for uninterpreted NE")
	}
}

func TestFiniteDomainSort(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.FiniteDomainSort("FD", 10)

	x := ctx.Const("x", sort).(FiniteDomain)
	y := ctx.Const("y", sort).(FiniteDomain)

	solver := NewSolver(ctx)
	solver.Assert(x.Eq(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for finite domain equality")
	}
}

func TestFiniteDomainNE(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.FiniteDomainSort("FD", 10)

	x := ctx.Const("x", sort).(FiniteDomain)
	y := ctx.Const("y", sort).(FiniteDomain)

	solver := NewSolver(ctx)
	solver.Assert(x.NE(y))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for finite domain NE")
	}
}

func TestASTContext(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.BoolConst("x")
	ast := x.AsAST()

	ctx2 := ast.Context()
	if ctx2 == nil {
		t.Error("expected non-nil context")
	}
}

func TestASTString(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.BoolConst("x")
	ast := x.AsAST()

	str := ast.String()
	if str == "" {
		t.Error("expected non-empty string")
	}
}

func TestASTKindString(t *testing.T) {
	kind := ASTKindNumeral
	str := kind.String()
	if str == "" {
		t.Error("expected non-empty string for AST kind")
	}
}

func TestSortKindString(t *testing.T) {
	kind := KindBool
	str := kind.String()
	if str == "" {
		t.Error("expected non-empty string for sort kind")
	}
}

func TestSortContext(t *testing.T) {
	ctx := NewContext(nil)
	sort := ctx.BoolSort()

	ctx2 := sort.Context()
	if ctx2 == nil {
		t.Error("expected non-nil context")
	}
}

func TestSortDomainAndRange(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	arrSort := ctx.ArraySort(intSort, intSort)

	domain, rng := arrSort.DomainAndRange()
	if domain.Context() == nil || rng.Context() == nil {
		t.Error("expected non-nil domain and range")
	}
}

func TestFuncDeclContext(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	f := ctx.FuncDecl("f", []Sort{intSort}, intSort)

	ctx2 := f.Context()
	if ctx2 == nil {
		t.Error("expected non-nil context")
	}
}

func TestFuncDeclString(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	f := ctx.FuncDecl("f", []Sort{intSort}, intSort)

	str := f.String()
	if str == "" {
		t.Error("expected non-empty string")
	}
}

func TestFreshFuncDecl(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	f := ctx.FreshFuncDecl("fresh", []Sort{intSort}, intSort)

	// Apply it
	x := ctx.FromInt(42, intSort)
	result := f.Apply(x)

	// Just check the result has the expected sort kind
	if result.Sort().Kind() != KindInt {
		t.Error("expected result to have KindInt")
	}
}

func TestFuncDeclMap(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	arrSort := ctx.ArraySort(intSort, intSort)

	f := ctx.FuncDecl("neg", []Sort{intSort}, intSort)
	arr := ctx.Const("arr", arrSort).(Array)

	result := f.Map(arr)
	solver := NewSolver(ctx)
	solver.Assert(result.Eq(result)) // Just check it works
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for func decl map")
	}
}

func TestModelString(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)
	x := ctx.BoolConst("x")
	solver.Assert(x)

	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		str := m.String()
		if str == "" {
			t.Error("expected non-empty model string")
		}
	}
}

func TestSolverString(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)
	x := ctx.BoolConst("x")
	solver.Assert(x)

	str := solver.String()
	if str == "" {
		t.Error("expected non-empty solver string")
	}
}

func TestSolverScopesAdditional(t *testing.T) {
	ctx := NewContext(nil)
	solver := NewSolver(ctx)

	// Check initial scopes
	if solver.NumScopes() != 0 {
		t.Error("expected 0 scopes initially")
	}
	solver.Push()
	if solver.NumScopes() != 1 {
		t.Error("expected 1 scope after push")
	}
	solver.Pop()
	if solver.NumScopes() != 0 {
		t.Error("expected 0 scopes after pop")
	}
}

func TestExprContext(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.BoolConst("x")

	ctx2 := x.Context()
	if ctx2 == nil {
		t.Error("expected non-nil context")
	}
}

func TestContextConfig(t *testing.T) {
	cfg := NewContextConfig()
	cfg.SetUint("timeout", 1000)
	ctx := NewContext(cfg)
	if ctx == nil {
		t.Error("expected non-nil context with config")
	}
}

func TestContextInterrupt(t *testing.T) {
	ctx := NewContext(nil)
	// Just test that it doesn't panic
	ctx.Interrupt()
}

func TestContextExtra(t *testing.T) {
	ctx := NewContext(nil)

	// SetExtra and Extra require a key and value
	ctx.SetExtra("mykey", 42)
	val := ctx.Extra("mykey")
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestRealConst(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.RealConst("x")

	solver := NewSolver(ctx)
	solver.Assert(x.Eq(ctx.FromBigRat(nil))) // This will create a 0
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for real const")
	}
}

func TestStringNE(t *testing.T) {
	ctx := NewContext(nil)
	s1 := ctx.FromString("hello")
	s2 := ctx.FromString("world")

	solver := NewSolver(ctx)
	solver.Assert(s1.NE(s2))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for string NE")
	}
}

func TestStringLT(t *testing.T) {
	ctx := NewContext(nil)
	s1 := ctx.FromString("a")
	s2 := ctx.FromString("b")

	solver := NewSolver(ctx)
	solver.Assert(s1.LT(s2))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'a' < 'b'")
	}
}

func TestStringLE(t *testing.T) {
	ctx := NewContext(nil)
	s1 := ctx.FromString("a")
	s2 := ctx.FromString("a")

	solver := NewSolver(ctx)
	solver.Assert(s1.LE(s2))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for 'a' <= 'a'")
	}
}

func TestStringLastIndexOf(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello hello")
	sub := ctx.FromString("hello")
	idx := s.LastIndexOf(sub)

	solver := NewSolver(ctx)
	solver.Assert(idx.Eq(ctx.Int(6)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for LastIndexOf")
	}
}

func TestStringConstAdditional(t *testing.T) {
	ctx := NewContext(nil)
	x := ctx.StringConst("x")

	solver := NewSolver(ctx)
	solver.Assert(x.Eq(ctx.FromString("test")))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for string const")
	}
}

func TestStringNth(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("hello")
	idx := ctx.Int(1)
	ch := s.Nth(idx).(Char)

	solver := NewSolver(ctx)
	// Nth returns a Char
	solver.Assert(ch.ToInt().Eq(ctx.Int(101))) // 'e' = 101
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for Nth")
	}
}

func TestSeqSort(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	seqSort := ctx.SeqSort(intSort)

	if !seqSort.IsSeqSort() {
		t.Error("expected IsSeqSort to return true")
	}

	basis := seqSort.SeqSortBasis()
	if basis.Kind() != KindInt {
		t.Error("expected int basis")
	}
}

func TestIsStringSort(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()

	if !strSort.IsStringSort() {
		t.Error("expected IsStringSort to return true")
	}
}

func TestEmptySeq(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	seqSort := ctx.SeqSort(intSort)
	empty := ctx.EmptySeq(seqSort)

	solver := NewSolver(ctx)
	solver.Assert(empty.Length().Eq(ctx.Int(0)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for empty seq")
	}
}

func TestSeqUnit(t *testing.T) {
	ctx := NewContext(nil)
	intSort := ctx.IntSort()
	val := ctx.FromInt(42, intSort)
	unit := ctx.SeqUnit(val)

	solver := NewSolver(ctx)
	solver.Assert(unit.Length().Eq(ctx.Int(1)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for seq unit")
	}
}

func TestAsString(t *testing.T) {
	ctx := NewContext(nil)
	s := ctx.FromString("test")
	str, ok := s.AsString()
	if !ok {
		t.Error("expected ok")
	}
	if str != "test" {
		t.Errorf("expected 'test', got %s", str)
	}
}

func TestREEq(t *testing.T) {
	ctx := NewContext(nil)
	re1 := ctx.FromString("a").ToRE()
	re2 := ctx.FromString("a").ToRE()

	solver := NewSolver(ctx)
	solver.Assert(re1.Eq(re2))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE equality")
	}
}

func TestRENE(t *testing.T) {
	ctx := NewContext(nil)
	re1 := ctx.FromString("a").ToRE()
	re2 := ctx.FromString("b").ToRE()

	solver := NewSolver(ctx)
	solver.Assert(re1.NE(re2))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE NE")
	}
}

func TestREIntersect(t *testing.T) {
	ctx := NewContext(nil)
	re1 := ctx.FromString("a").ToRE().Star()
	re2 := ctx.FromString("a").ToRE().Plus()
	inter := re1.Intersect(re2)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("a").InRE(inter))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE intersect")
	}
}

func TestREDiff(t *testing.T) {
	ctx := NewContext(nil)
	re1 := ctx.FromString("a").ToRE().Star()
	re2 := ctx.FromString("a").ToRE()
	diff := re1.Diff(re2)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("aa").InRE(diff))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE diff")
	}
}

func TestREPower(t *testing.T) {
	ctx := NewContext(nil)
	re := ctx.FromString("a").ToRE()
	pow := re.Power(3)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("aaa").InRE(pow))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE power")
	}
}

func TestREEmpty(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()
	reSort := ctx.RESort(strSort)
	empty := ctx.REEmpty(reSort)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("").InRE(empty).Not())
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE empty")
	}
}

func TestREFull(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()
	reSort := ctx.RESort(strSort)
	full := ctx.REFull(reSort)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("anything").InRE(full))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE full")
	}
}

func TestREAllChar(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()
	reSort := ctx.RESort(strSort)
	allChar := ctx.REAllChar(reSort)

	solver := NewSolver(ctx)
	solver.Assert(ctx.FromString("x").InRE(allChar))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for RE all char")
	}
}

func TestIsRESort(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()
	reSort := ctx.RESort(strSort)

	if !reSort.IsRESort() {
		t.Error("expected IsRESort to return true")
	}
}

func TestRESortBasis(t *testing.T) {
	ctx := NewContext(nil)
	strSort := ctx.StringSort()
	reSort := ctx.RESort(strSort)

	basis := reSort.RESortBasis()
	if basis.Kind() != KindSeq {
		t.Error("expected seq basis")
	}
}

func TestCharEq(t *testing.T) {
	ctx := NewContext(nil)
	c1 := ctx.Const("c1", ctx.CharSort()).(Char)
	c2 := ctx.Const("c2", ctx.CharSort()).(Char)

	solver := NewSolver(ctx)
	solver.Assert(c1.Eq(c2))
	solver.Assert(c1.ToInt().Eq(ctx.Int(65)))
	solver.Assert(c2.ToInt().Eq(ctx.Int(65)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for char equality")
	}
}

func TestCharNE(t *testing.T) {
	ctx := NewContext(nil)
	c1 := ctx.Const("c1", ctx.CharSort()).(Char)
	c2 := ctx.Const("c2", ctx.CharSort()).(Char)

	solver := NewSolver(ctx)
	solver.Assert(c1.NE(c2))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for char NE")
	}
}

func TestIsCharSort(t *testing.T) {
	ctx := NewContext(nil)
	charSort := ctx.CharSort()

	if !charSort.IsCharSort() {
		t.Error("expected IsCharSort to return true")
	}
}

func TestCharConst(t *testing.T) {
	ctx := NewContext(nil)
	c := ctx.CharConst("c")

	solver := NewSolver(ctx)
	solver.Assert(c.ToInt().Eq(ctx.Int(65)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for char const")
	}
}

func TestCharFromBV(t *testing.T) {
	ctx := NewContext(nil)
	// Create a BV with value 65 (ASCII 'A')
	// Z3 uses 18-bit BVs for characters in Unicode mode (default)
	bv := ctx.FromInt(65, ctx.BVSort(18)).(BV)
	c := ctx.CharFromBV(bv)

	solver := NewSolver(ctx)
	solver.Assert(c.ToInt().Eq(ctx.Int(65)))
	if sat, _ := solver.Check(); !sat {
		t.Error("expected SAT for CharFromBV")
	}
}

func TestOptimizeAssertAndTrack(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.BoolConst("x")
	tracker := ctx.BoolConst("track")

	opt.AssertAndTrack(x, tracker)

	if sat, _ := opt.Check(); !sat {
		t.Error("expected SAT")
	}
}

func TestOptimizeCheckAssumptions(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.BoolConst("x")
	y := ctx.BoolConst("y")
	opt.Assert(x)
	opt.Assert(y.Not())

	if sat, _ := opt.CheckAssumptions(); !sat {
		t.Error("expected SAT")
	}
}

func TestOptimizeUnsatCore(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	x := ctx.BoolConst("x")
	t1 := ctx.BoolConst("t1")
	t2 := ctx.BoolConst("t2")

	opt.AssertAndTrack(x, t1)
	opt.AssertAndTrack(x.Not(), t2)

	sat, _ := opt.Check()
	if sat {
		t.Error("expected UNSAT")
	}

	core := opt.UnsatCore()
	if len(core) == 0 {
		t.Error("expected non-empty unsat core")
	}
}

func TestOptimizeFromString(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	opt.FromString("(declare-const x Int)\n(assert (> x 0))\n(minimize x)")
}

func TestOptimizeHelp(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	help := opt.Help()
	if help == "" {
		t.Error("expected non-empty help")
	}
}

func TestSimplifyConfig(t *testing.T) {
	ctx := NewContext(nil)
	cfg := NewSimplifyConfig(ctx)
	if cfg == nil {
		t.Error("expected non-nil config")
	}
}

func TestConfigSetBool(t *testing.T) {
	cfg := NewContextConfig()
	cfg.SetBool("proof", true)
	ctx := NewContext(cfg)
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestConfigSetUint(t *testing.T) {
	cfg := NewContextConfig()
	cfg.SetUint("timeout", 1000)
	ctx := NewContext(cfg)
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestConfigSetFloat(t *testing.T) {
	cfg := NewContextConfig()
	cfg.SetFloat("pp.decimal_precision", 10.0)
	ctx := NewContext(cfg)
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}
