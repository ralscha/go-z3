// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"runtime"
	"unsafe"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// An Optimize is a collection of predicates that can be checked for
// satisfiability while optimizing objectives (minimizing or maximizing).
//
// These predicates form a stack that can be manipulated with
// Push/Pop.
type Optimize struct {
	*optimizeImpl
	noEq
}

type optimizeImpl struct {
	ctx *Context
	c   C.Z3_optimize
}

// NewOptimize returns a new, empty optimization context.
func NewOptimize(ctx *Context) *Optimize {
	var impl *optimizeImpl
	ctx.do(func() {
		impl = &optimizeImpl{
			ctx,
			C.Z3_mk_optimize(ctx.c),
		}
	})
	ctx.do(func() {
		C.Z3_optimize_inc_ref(ctx.c, impl.c)
	})
	runtime.SetFinalizer(impl, func(impl *optimizeImpl) {
		impl.ctx.do(func() {
			C.Z3_optimize_dec_ref(impl.ctx.c, impl.c)
		})
	})
	return &Optimize{impl, noEq{}}
}

// Assert adds val as a hard constraint to the optimization context.
func (o *Optimize) Assert(val Bool) {
	o.ctx.do(func() {
		C.Z3_optimize_assert(o.ctx.c, o.c, val.c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
}

// AssertAndTrack adds val as a hard constraint to the optimization context
// and associates it with the Boolean constant track for unsat core extraction.
func (o *Optimize) AssertAndTrack(val, track Bool) {
	o.ctx.do(func() {
		C.Z3_optimize_assert_and_track(o.ctx.c, o.c, val.c, track.c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	runtime.KeepAlive(track)
}

// AssertSoft adds val as a soft constraint with the given weight and optional id.
// Weight represents the penalty for violating the constraint (negative weights
// become rewards). ID provides a mechanism to group soft constraints.
// Returns the index of the soft constraint.
func (o *Optimize) AssertSoft(val Bool, weight string, id string) uint {
	sym := o.ctx.symbol(id)
	cweight := C.CString(weight)
	defer C.free(unsafe.Pointer(cweight))
	var handle C.uint
	o.ctx.do(func() {
		handle = C.Z3_optimize_assert_soft(o.ctx.c, o.c, val.c, cweight, sym)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	return uint(handle)
}

// Push saves the current state of the Optimize so it can be restored
// with Pop.
func (o *Optimize) Push() {
	o.ctx.do(func() {
		C.Z3_optimize_push(o.ctx.c, o.c)
	})
	runtime.KeepAlive(o)
}

// Pop removes all assertions added since the matching Push.
func (o *Optimize) Pop() {
	o.ctx.do(func() {
		C.Z3_optimize_pop(o.ctx.c, o.c)
	})
	runtime.KeepAlive(o)
}

// Objective is a handle to an optimization objective that can be used to
// retrieve the upper/lower bounds of an optimization solution.
type Objective struct {
	opt    *Optimize
	handle C.uint
}

// Maximize adds a maximization objective for the given value.
// Returns an Objective handle that can be used to retrieve bounds.
func (o *Optimize) Maximize(val Value) *Objective {
	var handle C.uint
	o.ctx.do(func() {
		handle = C.Z3_optimize_maximize(o.ctx.c, o.c, val.impl().c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	return &Objective{o, handle}
}

// Minimize adds a minimization objective for the given value.
// Returns an Objective handle that can be used to retrieve bounds.
func (o *Optimize) Minimize(val Value) *Objective {
	var handle C.uint
	o.ctx.do(func() {
		handle = C.Z3_optimize_minimize(o.ctx.c, o.c, val.impl().c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	return &Objective{o, handle}
}

// Lower returns the lower bound of the objective after a successful Check.
func (obj *Objective) Lower() Value {
	var ast AST
	obj.opt.ctx.do(func() {
		cast := C.Z3_optimize_get_lower(obj.opt.ctx.c, obj.opt.c, obj.handle)
		ast = wrapAST(obj.opt.ctx, cast)
	})
	runtime.KeepAlive(obj)
	return ast.AsValue()
}

// Upper returns the upper bound of the objective after a successful Check.
func (obj *Objective) Upper() Value {
	var ast AST
	obj.opt.ctx.do(func() {
		cast := C.Z3_optimize_get_upper(obj.opt.ctx.c, obj.opt.c, obj.handle)
		ast = wrapAST(obj.opt.ctx, cast)
	})
	runtime.KeepAlive(obj)
	return ast.AsValue()
}

// Check determines whether the predicates in the Optimize context are
// satisfiable and produces optimal values. If Z3 is unable to determine
// satisfiability, it returns an *ErrSatUnknown error.
func (o *Optimize) Check() (sat bool, err error) {
	var res C.Z3_lbool
	o.ctx.do(func() {
		res = C.Z3_optimize_check(o.ctx.c, o.c, 0, nil)
	})
	if res == C.Z3_L_UNDEF {
		// Get the reason.
		o.ctx.do(func() {
			cerr := C.Z3_optimize_get_reason_unknown(o.ctx.c, o.c)
			err = &ErrSatUnknown{C.GoString(cerr)}
		})
	}
	runtime.KeepAlive(o)
	return res == C.Z3_L_TRUE, err
}

// CheckAssumptions determines whether the predicates in the Optimize context
// together with the given assumptions are satisfiable and produces optimal values.
// If Z3 is unable to determine satisfiability, it returns an *ErrSatUnknown error.
func (o *Optimize) CheckAssumptions(assumptions ...Bool) (sat bool, err error) {
	cargs := make([]C.Z3_ast, len(assumptions))
	for i, arg := range assumptions {
		cargs[i] = arg.c
	}
	var res C.Z3_lbool
	o.ctx.do(func() {
		var cap *C.Z3_ast
		if len(cargs) > 0 {
			cap = &cargs[0]
		}
		res = C.Z3_optimize_check(o.ctx.c, o.c, C.uint(len(cargs)), cap)
	})
	if res == C.Z3_L_UNDEF {
		// Get the reason.
		o.ctx.do(func() {
			cerr := C.Z3_optimize_get_reason_unknown(o.ctx.c, o.c)
			err = &ErrSatUnknown{C.GoString(cerr)}
		})
	}
	runtime.KeepAlive(o)
	if len(cargs) > 0 {
		runtime.KeepAlive(&cargs[0])
	}
	return res == C.Z3_L_TRUE, err
}

// Model returns the model for the last Check. Model panics if Check
// has not been called or the last Check did not return true.
func (o *Optimize) Model() *Model {
	var model *Model
	o.ctx.do(func() {
		model = wrapModel(o.ctx, C.Z3_optimize_get_model(o.ctx.c, o.c))
	})
	runtime.KeepAlive(o)
	return model
}

// UnsatCore returns the subset of assumptions that were used in the
// unsatisfiability proof after a CheckAssumptions call that returned false.
func (o *Optimize) UnsatCore() []Bool {
	var asts []C.Z3_ast
	o.ctx.do(func() {
		vec := C.Z3_optimize_get_unsat_core(o.ctx.c, o.c)
		C.Z3_ast_vector_inc_ref(o.ctx.c, vec)
		defer C.Z3_ast_vector_dec_ref(o.ctx.c, vec)
		size := int(C.Z3_ast_vector_size(o.ctx.c, vec))
		asts = make([]C.Z3_ast, size)
		for i := 0; i < size; i++ {
			asts[i] = C.Z3_ast_vector_get(o.ctx.c, vec, C.uint(i))
		}
	})
	result := make([]Bool, len(asts))
	for i, ast := range asts {
		a := ast // capture for closure
		result[i] = Bool(wrapValue(o.ctx, func() C.Z3_ast { return a }))
	}
	runtime.KeepAlive(o)
	return result
}

// String returns a string representation of o.
func (o *Optimize) String() string {
	var res string
	o.ctx.do(func() {
		res = C.GoString(C.Z3_optimize_to_string(o.ctx.c, o.c))
	})
	runtime.KeepAlive(o)
	return res
}

// SetParams sets parameters on the optimization context.
func (o *Optimize) SetParams(config *Config) {
	cparams := config.toC(o.ctx)
	o.ctx.do(func() {
		C.Z3_optimize_set_params(o.ctx.c, o.c, cparams)
	})
	o.ctx.do(func() {
		C.Z3_params_dec_ref(o.ctx.c, cparams)
	})
	runtime.KeepAlive(o)
}

// Assertions returns the assertions in the optimization context.
func (o *Optimize) Assertions() []Bool {
	var asts []C.Z3_ast
	o.ctx.do(func() {
		vec := C.Z3_optimize_get_assertions(o.ctx.c, o.c)
		C.Z3_ast_vector_inc_ref(o.ctx.c, vec)
		defer C.Z3_ast_vector_dec_ref(o.ctx.c, vec)
		size := int(C.Z3_ast_vector_size(o.ctx.c, vec))
		asts = make([]C.Z3_ast, size)
		for i := 0; i < size; i++ {
			asts[i] = C.Z3_ast_vector_get(o.ctx.c, vec, C.uint(i))
		}
	})
	result := make([]Bool, len(asts))
	for i, ast := range asts {
		a := ast // capture for closure
		result[i] = Bool(wrapValue(o.ctx, func() C.Z3_ast { return a }))
	}
	runtime.KeepAlive(o)
	return result
}

// FromString parses an SMT-LIB2 string with assertions, soft constraints
// and optimization objectives and adds them to the optimization context.
func (o *Optimize) FromString(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	o.ctx.do(func() {
		C.Z3_optimize_from_string(o.ctx.c, o.c, cs)
	})
	runtime.KeepAlive(o)
}

// FromFile parses an SMT-LIB2 file with assertions, soft constraints
// and optimization objectives and adds them to the optimization context.
func (o *Optimize) FromFile(path string) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	o.ctx.do(func() {
		C.Z3_optimize_from_file(o.ctx.c, o.c, cpath)
	})
	runtime.KeepAlive(o)
}

// Help returns a string describing the parameters accepted by the optimizer.
func (o *Optimize) Help() string {
	var res string
	o.ctx.do(func() {
		res = C.GoString(C.Z3_optimize_get_help(o.ctx.c, o.c))
	})
	runtime.KeepAlive(o)
	return res
}
