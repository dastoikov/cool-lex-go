// Copyright 2021-2024 The Cool-lex-Go Contributors, see the CONTRIBUTORS file.
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
	"math"
	"testing"
)

func testOverflow(op func(uint, uint) (uint, error), a, b uint, overflow bool, expect ...uint) error {
	r, err := op(a, b)
	if (err != nil) != overflow {
		return fmt.Errorf("overflow: expected %v, got (%v)", overflow, err)
	}
	if !overflow && r != expect[0] {
		return fmt.Errorf("overflow: expected %v, got %v, for a %d and b %d", expect[0], r, a, b)
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
