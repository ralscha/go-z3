// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package st

import (
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/ralscha/go-z3/z3"
)

// TestEquivInt8 tests Int8 operations
func TestEquivInt8(t *testing.T) {
	testEquiv(t, reflect.TypeOf(Int8{}), Int8.sym,
		int8(0), int8(1), int8(-1), int8(7), int8(8),
		int8(math.MaxInt8), int8(math.MinInt8))
}

// TestEquivInt16 tests Int16 operations
func TestEquivInt16(t *testing.T) {
	testEquiv(t, reflect.TypeOf(Int16{}), Int16.sym,
		int16(0), int16(1), int16(-1), int16(15), int16(16),
		int16(math.MaxInt16), int16(math.MinInt16))
}

// TestEquivInt64 tests Int64 operations
func TestEquivInt64(t *testing.T) {
	testEquiv(t, reflect.TypeOf(Int64{}), Int64.sym,
		int64(0), int64(1), int64(-1), int64(63), int64(64),
		int64(math.MaxInt32), int64(math.MinInt32))
}

// TestEquivUint8 tests Uint8 operations
func TestEquivUint8(t *testing.T) {
	testEquiv(t, reflect.TypeOf(Uint8{}), Uint8.sym,
		uint8(0), uint8(1), uint8(7), uint8(8), uint8(255))
}

// TestEquivUint16 tests Uint16 operations
func TestEquivUint16(t *testing.T) {
	testEquiv(t, reflect.TypeOf(Uint16{}), Uint16.sym,
		uint16(0), uint16(1), uint16(15), uint16(16), uint16(65535))
}

// TestEquivUint64 tests Uint64 operations
func TestEquivUint64(t *testing.T) {
	testEquiv(t, reflect.TypeOf(Uint64{}), Uint64.sym,
		uint64(0), uint64(1), uint64(63), uint64(64), uint64(1<<32-1))
}

// TestBoolString tests Bool.String method
func TestBoolString(t *testing.T) {
	// Test concrete
	b1 := Bool{C: true}
	if b1.String() != "true" {
		t.Errorf("expected 'true', got %s", b1.String())
	}

	// Test symbolic
	ctx := z3.NewContext(nil)
	cache := getCache(ctx)
	b2 := Bool{S: ctx.FreshConst("b", cache.sortBool).(z3.Bool)}
	str := b2.String()
	if str == "" {
		t.Error("expected non-empty string for symbolic bool")
	}
}

// TestBoolEval tests Bool.Eval method
func TestBoolEval(t *testing.T) {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	x := AnyBool(ctx, "x")
	solver.Assert(x.S)

	if sat, _ := solver.Check(); !sat {
		t.Fatal("expected SAT")
	}

	m := solver.Model()
	val := x.Eval(m)
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

// TestInt32String tests Int32.String method
func TestInt32String(t *testing.T) {
	// Test concrete
	i1 := Int32{C: 42}
	if i1.String() != "42" {
		t.Errorf("expected '42', got %s", i1.String())
	}

	// Test symbolic
	ctx := z3.NewContext(nil)
	i2 := AnyInt32(ctx, "x")
	str := i2.String()
	if str == "" {
		t.Error("expected non-empty string for symbolic int32")
	}
}

// TestInt32Eval tests Int32.Eval method
func TestInt32Eval(t *testing.T) {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	x := AnyInt32(ctx, "x")
	solver.Assert(x.S.Eq(ctx.FromInt(42, ctx.BVSort(32)).(z3.BV)))

	if sat, _ := solver.Check(); !sat {
		t.Fatal("expected SAT")
	}

	m := solver.Model()
	val := x.Eval(m)
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

// TestUint32String tests Uint32.String method
func TestUint32String(t *testing.T) {
	// Test concrete
	i1 := Uint32{C: 42}
	if i1.String() != "42" {
		t.Errorf("expected '42', got %s", i1.String())
	}

	// Test symbolic
	ctx := z3.NewContext(nil)
	i2 := AnyUint32(ctx, "x")
	str := i2.String()
	if str == "" {
		t.Error("expected non-empty string for symbolic uint32")
	}
}

// TestUint32Eval tests Uint32.Eval method
func TestUint32Eval(t *testing.T) {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	x := AnyUint32(ctx, "x")
	solver.Assert(x.S.Eq(ctx.FromInt(42, ctx.BVSort(32)).(z3.BV)))

	if sat, _ := solver.Check(); !sat {
		t.Fatal("expected SAT")
	}

	m := solver.Model()
	val := x.Eval(m)
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

// TestIntegerString tests Integer.String method
func TestIntegerString(t *testing.T) {
	// Test concrete
	i1 := Integer{C: big.NewInt(42)}
	if i1.String() != "42" {
		t.Errorf("expected '42', got %s", i1.String())
	}

	// Test symbolic
	ctx := z3.NewContext(nil)
	i2 := AnyInteger(ctx, "x")
	str := i2.String()
	if str == "" {
		t.Error("expected non-empty string for symbolic integer")
	}
}

// TestIntegerEval tests Integer.Eval method
func TestIntegerEval(t *testing.T) {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	x := AnyInteger(ctx, "x")
	solver.Assert(x.S.Eq(ctx.FromInt(42, ctx.IntSort()).(z3.Int)))

	if sat, _ := solver.Check(); !sat {
		t.Fatal("expected SAT")
	}

	m := solver.Model()
	val := x.Eval(m)
	if val.Cmp(big.NewInt(42)) != 0 {
		t.Errorf("expected 42, got %v", val)
	}
}

// TestRealString tests Real.String method
func TestRealString(t *testing.T) {
	// Test concrete
	r1 := Real{C: big.NewRat(1, 2)}
	if r1.String() != "1/2" {
		t.Errorf("expected '1/2', got %s", r1.String())
	}

	// Test symbolic
	ctx := z3.NewContext(nil)
	r2 := AnyReal(ctx, "x")
	str := r2.String()
	if str == "" {
		t.Error("expected non-empty string for symbolic real")
	}
}

// TestRealEval tests Real.Eval method
func TestRealEval(t *testing.T) {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	x := AnyReal(ctx, "x")
	solver.Assert(x.S.Eq(ctx.FromBigRat(big.NewRat(1, 2))))

	if sat, _ := solver.Check(); !sat {
		t.Fatal("expected SAT")
	}

	m := solver.Model()
	val := x.Eval(m)
	if val.Cmp(big.NewRat(1, 2)) != 0 {
		t.Errorf("expected 1/2, got %v", val)
	}
}

// TestAnyFunctions tests all Any* functions for different types
func TestAnyFunctions(t *testing.T) {
	ctx := z3.NewContext(nil)

	// AnyInt
	intVal := AnyInt(ctx, "i")
	if intVal.IsConcrete() {
		t.Error("AnyInt should return symbolic value")
	}

	// AnyInt8
	int8Val := AnyInt8(ctx, "i8")
	if int8Val.IsConcrete() {
		t.Error("AnyInt8 should return symbolic value")
	}

	// AnyInt16
	int16Val := AnyInt16(ctx, "i16")
	if int16Val.IsConcrete() {
		t.Error("AnyInt16 should return symbolic value")
	}

	// AnyInt64
	int64Val := AnyInt64(ctx, "i64")
	if int64Val.IsConcrete() {
		t.Error("AnyInt64 should return symbolic value")
	}

	// AnyUint
	uintVal := AnyUint(ctx, "u")
	if uintVal.IsConcrete() {
		t.Error("AnyUint should return symbolic value")
	}

	// AnyUint8
	uint8Val := AnyUint8(ctx, "u8")
	if uint8Val.IsConcrete() {
		t.Error("AnyUint8 should return symbolic value")
	}

	// AnyUint16
	uint16Val := AnyUint16(ctx, "u16")
	if uint16Val.IsConcrete() {
		t.Error("AnyUint16 should return symbolic value")
	}

	// AnyUint64
	uint64Val := AnyUint64(ctx, "u64")
	if uint64Val.IsConcrete() {
		t.Error("AnyUint64 should return symbolic value")
	}

	// AnyUintptr
	uintptrVal := AnyUintptr(ctx, "uptr")
	if uintptrVal.IsConcrete() {
		t.Error("AnyUintptr should return symbolic value")
	}
}

// TestIntStringMethods tests String method on Int types
func TestIntStringMethods(t *testing.T) {
	ctx := z3.NewContext(nil)

	// Int
	intVal := AnyInt(ctx, "i")
	if intVal.String() == "" {
		t.Error("Int.String() should not be empty")
	}
	concreteInt := Int{C: 42}
	if concreteInt.String() != "42" {
		t.Errorf("expected '42', got %s", concreteInt.String())
	}

	// Int8
	int8Val := AnyInt8(ctx, "i8")
	if int8Val.String() == "" {
		t.Error("Int8.String() should not be empty")
	}
	concreteInt8 := Int8{C: 42}
	if concreteInt8.String() != "42" {
		t.Errorf("expected '42', got %s", concreteInt8.String())
	}

	// Int16
	int16Val := AnyInt16(ctx, "i16")
	if int16Val.String() == "" {
		t.Error("Int16.String() should not be empty")
	}
	concreteInt16 := Int16{C: 42}
	if concreteInt16.String() != "42" {
		t.Errorf("expected '42', got %s", concreteInt16.String())
	}

	// Int64
	int64Val := AnyInt64(ctx, "i64")
	if int64Val.String() == "" {
		t.Error("Int64.String() should not be empty")
	}
	concreteInt64 := Int64{C: 42}
	if concreteInt64.String() != "42" {
		t.Errorf("expected '42', got %s", concreteInt64.String())
	}
}

// TestUintStringMethods tests String method on Uint types
func TestUintStringMethods(t *testing.T) {
	ctx := z3.NewContext(nil)

	// Uint
	uintVal := AnyUint(ctx, "u")
	if uintVal.String() == "" {
		t.Error("Uint.String() should not be empty")
	}
	concreteUint := Uint{C: 42}
	if concreteUint.String() != "42" {
		t.Errorf("expected '42', got %s", concreteUint.String())
	}

	// Uint8
	uint8Val := AnyUint8(ctx, "u8")
	if uint8Val.String() == "" {
		t.Error("Uint8.String() should not be empty")
	}
	concreteUint8 := Uint8{C: 42}
	if concreteUint8.String() != "42" {
		t.Errorf("expected '42', got %s", concreteUint8.String())
	}

	// Uint16
	uint16Val := AnyUint16(ctx, "u16")
	if uint16Val.String() == "" {
		t.Error("Uint16.String() should not be empty")
	}
	concreteUint16 := Uint16{C: 42}
	if concreteUint16.String() != "42" {
		t.Errorf("expected '42', got %s", concreteUint16.String())
	}

	// Uint64
	uint64Val := AnyUint64(ctx, "u64")
	if uint64Val.String() == "" {
		t.Error("Uint64.String() should not be empty")
	}
	concreteUint64 := Uint64{C: 42}
	if concreteUint64.String() != "42" {
		t.Errorf("expected '42', got %s", concreteUint64.String())
	}

	// Uintptr
	uintptrVal := AnyUintptr(ctx, "uptr")
	if uintptrVal.String() == "" {
		t.Error("Uintptr.String() should not be empty")
	}
	concreteUintptr := Uintptr{C: 42}
	if concreteUintptr.String() != "42" {
		t.Errorf("expected '42', got %s", concreteUintptr.String())
	}
}

// TestEvalMethods tests Eval methods on various types
func TestEvalMethods(t *testing.T) {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	// Int
	intVal := AnyInt(ctx, "i")
	solver.Assert(intVal.S.Eq(ctx.FromInt(42, ctx.BVSort(64)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := intVal.Eval(m)
		if result != 42 {
			t.Errorf("Int.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Int8
	int8Val := AnyInt8(ctx, "i8")
	solver.Assert(int8Val.S.Eq(ctx.FromInt(42, ctx.BVSort(8)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := int8Val.Eval(m)
		if result != 42 {
			t.Errorf("Int8.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Int16
	int16Val := AnyInt16(ctx, "i16")
	solver.Assert(int16Val.S.Eq(ctx.FromInt(42, ctx.BVSort(16)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := int16Val.Eval(m)
		if result != 42 {
			t.Errorf("Int16.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Int64
	int64Val := AnyInt64(ctx, "i64")
	solver.Assert(int64Val.S.Eq(ctx.FromInt(42, ctx.BVSort(64)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := int64Val.Eval(m)
		if result != 42 {
			t.Errorf("Int64.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Uint
	uintVal := AnyUint(ctx, "u")
	solver.Assert(uintVal.S.Eq(ctx.FromInt(42, ctx.BVSort(64)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := uintVal.Eval(m)
		if result != 42 {
			t.Errorf("Uint.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Uint8
	uint8Val := AnyUint8(ctx, "u8")
	solver.Assert(uint8Val.S.Eq(ctx.FromInt(42, ctx.BVSort(8)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := uint8Val.Eval(m)
		if result != 42 {
			t.Errorf("Uint8.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Uint16
	uint16Val := AnyUint16(ctx, "u16")
	solver.Assert(uint16Val.S.Eq(ctx.FromInt(42, ctx.BVSort(16)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := uint16Val.Eval(m)
		if result != 42 {
			t.Errorf("Uint16.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Uint64
	uint64Val := AnyUint64(ctx, "u64")
	solver.Assert(uint64Val.S.Eq(ctx.FromInt(42, ctx.BVSort(64)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := uint64Val.Eval(m)
		if result != 42 {
			t.Errorf("Uint64.Eval expected 42, got %d", result)
		}
	}
	solver.Reset()

	// Uintptr
	uintptrVal := AnyUintptr(ctx, "uptr")
	solver.Assert(uintptrVal.S.Eq(ctx.FromInt(42, ctx.BVSort(64)).(z3.BV)))
	if sat, _ := solver.Check(); sat {
		m := solver.Model()
		result := uintptrVal.Eval(m)
		if result != 42 {
			t.Errorf("Uintptr.Eval expected 42, got %d", result)
		}
	}
}
