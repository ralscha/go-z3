// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "runtime"

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// A Solver is a collection of predicates that can be checked for
// satisfiability.
//
// These predicates form a stack that can be manipulated with
// Push/Pop.
type Solver struct {
	*solverImpl
	noEq
}

type solverImpl struct {
	ctx *Context
	c   C.Z3_solver
}

// NewSolver returns a new, empty solver.
func NewSolver(ctx *Context) *Solver {
	var impl *solverImpl
	ctx.do(func() {
		impl = &solverImpl{
			ctx,
			C.Z3_mk_solver(ctx.c),
		}
	})
	ctx.do(func() {
		C.Z3_solver_inc_ref(ctx.c, impl.c)
	})
	runtime.SetFinalizer(impl, func(impl *solverImpl) {
		impl.ctx.do(func() {
			C.Z3_solver_dec_ref(impl.ctx.c, impl.c)
		})
	})
	return &Solver{impl, noEq{}}
}

// Assert adds val to the set of predicates that must be satisfied.
func (s *Solver) Assert(val Bool) {
	s.ctx.do(func() {
		C.Z3_solver_assert(s.ctx.c, s.c, val.c)
	})
	runtime.KeepAlive(s)
	runtime.KeepAlive(val)
}

// Push saves the current state of the Solver so it can be restored
// with Pop.
func (s *Solver) Push() {
	s.ctx.do(func() {
		C.Z3_solver_push(s.ctx.c, s.c)
	})
	runtime.KeepAlive(s)
}

// Pop removes assertions that were added since the matching Push.
func (s *Solver) Pop() {
	s.ctx.do(func() {
		C.Z3_solver_pop(s.ctx.c, s.c, 1)
	})
	runtime.KeepAlive(s)
}

// Reset removes all assertions from the Solver and resets its stack.
func (s *Solver) Reset() {
	s.ctx.do(func() {
		C.Z3_solver_reset(s.ctx.c, s.c)
	})
	runtime.KeepAlive(s)
}

// ErrSatUnknown is produced when Z3 cannot determine satisfiability.
type ErrSatUnknown struct {
	// Reason gives a brief description of why Z3 could not
	// determine satisfiability.
	Reason string
}

// Error returns the reason Z3 could not determine satisfiability.
func (e *ErrSatUnknown) Error() string {
	return e.Reason
}

// Check determines whether the predicates in Solver s are satisfiable
// or unsatisfiable. If Z3 is unable to determine satisfiability, it
// returns an *ErrSatUnknown error.
func (s *Solver) Check() (sat bool, err error) {
	var res C.Z3_lbool
	s.ctx.do(func() {
		res = C.Z3_solver_check(s.ctx.c, s.c)
	})
	if res == C.Z3_L_UNDEF {
		// Get the reason.
		s.ctx.do(func() {
			cerr := C.Z3_solver_get_reason_unknown(s.ctx.c, s.c)
			err = &ErrSatUnknown{C.GoString(cerr)}
		})
	}
	runtime.KeepAlive(s)
	return res == C.Z3_L_TRUE, err
}

// Model returns the model for the last Check. Model panics if Check
// has not been called or the last Check did not return true.
func (s *Solver) Model() *Model {
	var model *Model
	s.ctx.do(func() {
		model = wrapModel(s.ctx, C.Z3_solver_get_model(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return model
}

// String returns a string representation of s.
func (s *Solver) String() string {
	var res string
	s.ctx.do(func() {
		res = C.GoString(C.Z3_solver_to_string(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return res
}

// NumScopes returns the number of backtracking points (Push calls
// without matching Pop calls).
func (s *Solver) NumScopes() uint {
	var res C.uint
	s.ctx.do(func() {
		res = C.Z3_solver_get_num_scopes(s.ctx.c, s.c)
	})
	runtime.KeepAlive(s)
	return uint(res)
}

// NumAssertions returns the number of assertions in the solver.
func (s *Solver) NumAssertions() uint {
	var res uint
	s.ctx.do(func() {
		vec := C.Z3_solver_get_assertions(s.ctx.c, s.c)
		res = uint(C.Z3_ast_vector_size(s.ctx.c, vec))
	})
	runtime.KeepAlive(s)
	return res
}

// Assertions returns the assertions in the solver.
func (s *Solver) Assertions() []Bool {
	var asts []C.Z3_ast
	s.ctx.do(func() {
		vec := C.Z3_solver_get_assertions(s.ctx.c, s.c)
		C.Z3_ast_vector_inc_ref(s.ctx.c, vec)
		defer C.Z3_ast_vector_dec_ref(s.ctx.c, vec)
		size := int(C.Z3_ast_vector_size(s.ctx.c, vec))
		asts = make([]C.Z3_ast, size)
		for i := 0; i < size; i++ {
			asts[i] = C.Z3_ast_vector_get(s.ctx.c, vec, C.uint(i))
		}
	})
	result := make([]Bool, len(asts))
	for i, ast := range asts {
		a := ast // capture for closure
		result[i] = Bool(wrapValue(s.ctx, func() C.Z3_ast { return a }))
	}
	runtime.KeepAlive(s)
	return result
}

// CheckAssumptions determines whether the predicates in Solver s
// together with the given assumptions are satisfiable or unsatisfiable.
// If Z3 is unable to determine satisfiability, it returns an *ErrSatUnknown error.
func (s *Solver) CheckAssumptions(assumptions ...Bool) (sat bool, err error) {
	cargs := make([]C.Z3_ast, len(assumptions))
	for i, arg := range assumptions {
		cargs[i] = arg.c
	}
	var res C.Z3_lbool
	s.ctx.do(func() {
		var cap *C.Z3_ast
		if len(cargs) > 0 {
			cap = &cargs[0]
		}
		res = C.Z3_solver_check_assumptions(s.ctx.c, s.c, C.uint(len(cargs)), cap)
	})
	if res == C.Z3_L_UNDEF {
		// Get the reason.
		s.ctx.do(func() {
			cerr := C.Z3_solver_get_reason_unknown(s.ctx.c, s.c)
			err = &ErrSatUnknown{C.GoString(cerr)}
		})
	}
	runtime.KeepAlive(s)
	runtime.KeepAlive(&cargs[0])
	return res == C.Z3_L_TRUE, err
}

// UnsatCore returns the subset of assumptions that were used in the
// unsatisfiability proof after a CheckAssumptions call that returned false.
func (s *Solver) UnsatCore() []Bool {
	var asts []C.Z3_ast
	s.ctx.do(func() {
		vec := C.Z3_solver_get_unsat_core(s.ctx.c, s.c)
		C.Z3_ast_vector_inc_ref(s.ctx.c, vec)
		defer C.Z3_ast_vector_dec_ref(s.ctx.c, vec)
		size := int(C.Z3_ast_vector_size(s.ctx.c, vec))
		asts = make([]C.Z3_ast, size)
		for i := 0; i < size; i++ {
			asts[i] = C.Z3_ast_vector_get(s.ctx.c, vec, C.uint(i))
		}
	})
	result := make([]Bool, len(asts))
	for i, ast := range asts {
		a := ast // capture for closure
		result[i] = Bool(wrapValue(s.ctx, func() C.Z3_ast { return a }))
	}
	runtime.KeepAlive(s)
	return result
}
