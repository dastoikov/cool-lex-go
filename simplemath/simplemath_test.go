// Copyright 2025 The Cool-lex-Go Contributors, see the CONTRIBUTORS file.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing permissions and
// limitations under the License.
package simplemath

import (
	"fmt"
	"iter"
	"math"
	"testing"
)

func testOverflow(op func(uint, uint) (uint, error), a, b uint, overflow bool, expect ...uint) error {
	r, err := op(a, b)
	if (err != nil) != overflow {
		return fmt.Errorf("overflow: expected %v, got (%v)", overflow, err)
	}
	if len(expect) > 0 && r != expect[0] {
		return fmt.Errorf("overflow: expected %v, got %v, for a=%d and b=%d and overflow=%t", expect[0], r, a, b, overflow)
	}
	return nil
}

func TestMul(t *testing.T) {
	test := func(a, b uint, overflow bool, expect ...uint) {
		err := testOverflow(Mul, a, b, overflow, expect...)
		if err != nil {
			t.Fatal(err)
		}
	}
	test(math.MaxUint, 1, false, math.MaxUint)
	test(math.MaxUint, 2, true)
	test(7, 7, false, 49)
}

func TestMulRange(t *testing.T) {
	test := func(a, b uint, overflow bool, expect ...uint) {
		err := testOverflow(MulRange, a, b, overflow, expect...)
		if err != nil {
			t.Fatal(err)
		}
	}
	test(math.MaxUint, math.MaxUint, false, math.MaxUint)
	test(math.MaxUint, math.MaxUint-1, true)
	test(3, 1, false, 6)
	test(1, 3, false, 6)
	test(0, 1, false, 0)
	test(0, 2, false, 0)
	test(3, 3, false, 3)
	test(1, 21, true, 2432902008176640000 /*20!*/)
}

func TestAdd(t *testing.T) {
	test := func(a, b uint, overflow bool, expect ...uint) {
		err := testOverflow(Add, a, b, overflow, expect...)
		if err != nil {
			t.Fatal(err)
		}
	}
	test(0, 0, false, 0)
	test(0, 1, false, 1)
	test(1, 3, false, 4)
	test(math.MaxUint, 0, false, math.MaxUint)
	test(math.MaxUint, 1, true)
	test(2, math.MaxUint, true)
}

func TestNumComb(t *testing.T) {
	test := func(n, k uint, expect uint) {
		actual, err := NumComb(n, k)
		if err != nil {
			t.Fatal(err)
		}
		if expect != actual {
			t.Fatalf("num comb: expected %d, got %d", expect, actual)
		}
	}
	test(2, 0, 1)
	test(3, 2, 3)
}

// --- DOZ ---

var (
	benchmarkDataDoz = map[int64]int64{
		-128:                128,
		math.MaxInt64 - 128: math.MaxInt64,
		math.MinInt64:       math.MinInt64 + 128,
	}
)

// rangeCartesianProduct applies biFunc to the Cartesian product of the specified range (inclusive), pairwise.
func rangeCartesianProduct(lo, hi int64, biFunc func(int64, int64) int64) {
	if lo > hi {
		lo, hi = hi, lo
	}
	for a := lo; ; a++ {
		for b := lo; ; b++ {
			_ = biFunc(a, b)
			if b == hi {
				break
			}
		}
		if a == hi {
			break
		}
	}
}

func BenchmarkDoz64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for lo, hi := range benchmarkDataDoz {
			rangeCartesianProduct(lo, hi, Doz64)
		}
	}
}

func BenchmarkDozB64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for lo, hi := range benchmarkDataDoz {
			rangeCartesianProduct(lo, hi, DozB64)
		}
	}
}

func TestDoz64(t *testing.T) {
	test := func(a, b int64) {
		t.Helper()
		expect := DozB64(a, b)
		actual := Doz64(a, b)
		if actual != expect {
			t.Fatalf("doz: for a=%v, b=%v expected %v, got %v (%b)", a, b, expect, actual, actual)
		}
	}

	// use only a handful of all numbers for sanity-testing the function
	for a := range allInt8() {
		for b := range allInt8() {
			test(int64(a), int64(b))
		}
	}

	test(math.MinInt64, math.MinInt64)
	test(math.MaxInt64, math.MaxInt64)
	test(math.MaxInt64, math.MinInt64)
	test(math.MinInt64, math.MaxInt64)
}

func TestDoz32(t *testing.T) {
	test := func(a, b int32) {
		t.Helper()
		expect := DozB32(a, b)
		actual := Doz32(a, b)
		if actual != expect {
			t.Fatalf("doz: for a=%v, b=%v expected %v, got %v (%b)", a, b, expect, actual, actual)
		}
	}

	// use only a handful of all numbers for sanity-testing the function
	for a := range allInt8() {
		for b := range allInt8() {
			test(int32(a), int32(b))
		}
	}

	test(math.MinInt32, math.MinInt32)
	test(math.MaxInt32, math.MaxInt32)
	test(math.MaxInt32, math.MinInt32)
	test(math.MinInt32, math.MaxInt32)
}

func allInt8() iter.Seq[int8] {
	return func(yield func(int8) bool) {
		for i := int8(math.MinInt8); ; i++ {
			if !yield(i) || i == math.MaxInt8 {
				return
			}
		}
	}
}
