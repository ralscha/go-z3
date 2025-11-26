// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ops

import (
	"go/token"
	"testing"
)

func TestIntBits(t *testing.T) {
	bits := intBits()
	// Should be 32 or 64 depending on platform
	if bits != 32 && bits != 64 {
		t.Errorf("intBits() = %d, expected 32 or 64", bits)
	}
}

func TestPtrBits(t *testing.T) {
	bits := ptrBits()
	// Should be 32 or 64 depending on platform
	if bits != 32 && bits != 64 {
		t.Errorf("ptrBits() = %d, expected 32 or 64", bits)
	}
}

func TestTypes(t *testing.T) {
	// Verify Types slice is correctly populated
	if len(Types) == 0 {
		t.Error("Types slice is empty")
	}

	// Check that Bool type exists
	foundBool := false
	for _, typ := range Types {
		if typ.StName == "Bool" {
			foundBool = true
			if typ.ConType != "bool" {
				t.Errorf("Bool.ConType = %s, expected bool", typ.ConType)
			}
			if typ.SymType != "Bool" {
				t.Errorf("Bool.SymType = %s, expected Bool", typ.SymType)
			}
			if typ.Flags&IsBool == 0 {
				t.Error("Bool should have IsBool flag")
			}
		}
	}
	if !foundBool {
		t.Error("Bool type not found in Types")
	}

	// Check Int32 type
	for _, typ := range Types {
		if typ.StName == "Int32" {
			if typ.Bits != 32 {
				t.Errorf("Int32.Bits = %d, expected 32", typ.Bits)
			}
			if typ.Flags&IsInteger == 0 {
				t.Error("Int32 should have IsInteger flag")
			}
			if typ.Flags&IsUnsigned != 0 {
				t.Error("Int32 should not have IsUnsigned flag")
			}
		}
	}

	// Check Uint32 type
	for _, typ := range Types {
		if typ.StName == "Uint32" {
			if typ.Bits != 32 {
				t.Errorf("Uint32.Bits = %d, expected 32", typ.Bits)
			}
			if typ.Flags&IsInteger == 0 {
				t.Error("Uint32 should have IsInteger flag")
			}
			if typ.Flags&IsUnsigned == 0 {
				t.Error("Uint32 should have IsUnsigned flag")
			}
		}
	}

	// Check Integer (big.Int) type
	for _, typ := range Types {
		if typ.StName == "Integer" {
			if typ.ConType != "*big.Int" {
				t.Errorf("Integer.ConType = %s, expected *big.Int", typ.ConType)
			}
			if typ.Flags&IsBigInt == 0 {
				t.Error("Integer should have IsBigInt flag")
			}
		}
	}

	// Check Real (big.Rat) type
	for _, typ := range Types {
		if typ.StName == "Real" {
			if typ.ConType != "*big.Rat" {
				t.Errorf("Real.ConType = %s, expected *big.Rat", typ.ConType)
			}
			if typ.Flags&IsBigRat == 0 {
				t.Error("Real should have IsBigRat flag")
			}
		}
	}
}

func TestBinOps(t *testing.T) {
	// Verify BinOps slice is correctly populated
	if len(BinOps) == 0 {
		t.Error("BinOps slice is empty")
	}

	// Check Add operation
	foundAdd := false
	for _, op := range BinOps {
		if op.Op == "+" && op.Tok == token.ADD {
			foundAdd = true
			if op.Method != "Add" {
				t.Errorf("Add method = %s, expected Add", op.Method)
			}
		}
	}
	if !foundAdd {
		t.Error("Add operation not found in BinOps")
	}

	// Check comparison operations have OpCompare flag
	for _, op := range BinOps {
		switch op.Op {
		case "==", "!=", "<", "<=", ">", ">=":
			if op.Flags&OpCompare == 0 {
				t.Errorf("Comparison op %s should have OpCompare flag", op.Op)
			}
		}
	}

	// Check shift operations have OpShift flag
	for _, op := range BinOps {
		if op.Op == "<<" || op.Op == ">>" {
			if op.Flags&OpShift == 0 {
				t.Errorf("Shift op %s should have OpShift flag", op.Op)
			}
		}
	}
}

func TestUnOps(t *testing.T) {
	// Verify UnOps slice is correctly populated
	if len(UnOps) == 0 {
		t.Error("UnOps slice is empty")
	}

	// Check Neg operation
	foundNeg := false
	for _, op := range UnOps {
		if op.Op == "-" && op.Tok == token.SUB {
			foundNeg = true
			if op.Method != "Neg" {
				t.Errorf("Neg method = %s, expected Neg", op.Method)
			}
		}
	}
	if !foundNeg {
		t.Error("Neg operation not found in UnOps")
	}

	// Check unary + has OpPos flag
	for _, op := range UnOps {
		if op.Op == "+" && op.Tok == token.ADD {
			if op.Flags&OpPos == 0 {
				t.Error("Unary + should have OpPos flag")
			}
		}
	}
}

func TestFlags(t *testing.T) {
	// Test Comparable flag
	if Comparable&IsBool == 0 {
		t.Error("Comparable should include IsBool")
	}
	if Comparable&IsInteger == 0 {
		t.Error("Comparable should include IsInteger")
	}
	if Comparable&IsFloat == 0 {
		t.Error("Comparable should include IsFloat")
	}
	if Comparable&IsBigInt == 0 {
		t.Error("Comparable should include IsBigInt")
	}
	if Comparable&IsBigRat == 0 {
		t.Error("Comparable should include IsBigRat")
	}

	// Test Ordered flag
	if Ordered&IsBool != 0 {
		t.Error("Ordered should not include IsBool")
	}
	if Ordered&IsInteger == 0 {
		t.Error("Ordered should include IsInteger")
	}
	if Ordered&IsFloat == 0 {
		t.Error("Ordered should include IsFloat")
	}
}
