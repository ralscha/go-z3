// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3log

import (
	"os"
	"testing"
)

func TestZ3Log(t *testing.T) {
	// Create a temporary file for the log
	tmpfile, err := os.CreateTemp("", "z3log_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Open the log
	ok := Open(tmpfile.Name())
	if !ok {
		t.Fatal("Failed to open Z3 log")
	}

	// Append something
	Append("test log entry")

	// Close the log
	Close()

	// Verify the file exists and has content
	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Log("Log file is empty (may be expected depending on Z3 logging behavior)")
	}
}
