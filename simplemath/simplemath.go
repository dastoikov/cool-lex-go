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

/*
Package simplemath offers simple/naive implementations of math operations for the purposes of aiding in writing auto tests.
*/
package simplemath

import (
	"fmt"
	"math/bits"
)

// MulRange calculates the product of all integers within the inclusive interval [n1, n2] (or [n2, n1] if n2 < n1).
//
// It returns the product if no numeric overflow occurs.  Otherwise, it returns an error,
// and the partial product calculated before the overflow is stored in 'r'.
func MulRange(n2, n1 uint) (r uint, err error) {
	if n1 > n2 {
		n1, n2 = n2, n1
	}

	for r = n1; n1 < n2 && err == nil; r, err = Mul(r, n1) {
		n1++
	}
	return
}

// Add returns the sum of a and b if no numeric overflow occurs, or an error otherwise.
func Add(a, b uint) (uint, error) {
	sum, carry := bits.Add(a, b, 0)
	if carry == 1 {
		return 0, fmt.Errorf("numeric overflow occurred adding %d and %d", a, b)
	}
	return sum, nil
}

// Mul returns the product of a and b, or an error if numeric overflow occurs.
func Mul(a, b uint) (uint, error) {
	overflow, product := bits.Mul(a, b)
	if overflow != 0 {
		return 0, fmt.Errorf("numeric overflow occurred multiplying %d and %d", a, b)
	}
	return product, nil
}

// NumComb calculates the number of combinations for the specified k and n.
//
// n: number of elements to combine; n>=k must hold.
// k: number of elements in a combination.
//
// Error is reported if numeric overflow occurs.
func NumComb(n, k uint) (uint, error) {
	if k > n {
		return 0, fmt.Errorf("k (%d) > n (%d)", k, n)
	}
	if k == 0 {
		return 1, nil
	}

	cn, err := MulRange(n, n-k+1)
	if err != nil {
		return 0, err
	}
	ck, err := Factorial(k)
	if err != nil {
		return 0, err
	}
	return cn / ck, nil
}

// Factorial returns the factorial of n, or 1 for n=0.
// Error is reported if numeric overflow occurs.
func Factorial(n uint) (uint, error) {
	if n == 0 {
		return 1, nil
	}
	return MulRange(n, 1)
}

// BitEq64 implements the bitwise-equivalence operation, that is, NOT(XOR(a,b)).
//
//go:inline
func BitEq64(a, b int64) int64 {
	return ^(a ^ b)
}

// DozB64 implements the difference-or-zero function.
// DozB64(a,b) is a - b if a >= b, and is 0 if a < b.
//
// It is named "saturated subtraction" in The Coolest Way to Generate Combinations paper
// by Frank Ruskey and Aaron Williams, see 3.3. Implementation in Computer words, page 10.
//
// See also: `simplemath.Doz64`.
func DozB64(a, b int64) int64 {
	if a < b {
		return 0
	}
	return a - b
}

// Doz64 implements the difference-or-zero function, in a branchless fashion.
// See `simplemath.DozB64` for details.
func Doz64(a, b int64) int64 {
	// see Hacker's Delight by Henry S. Warren, Jr., 2-18 Doz, Max, Min, pages 37-38.
	d := a - b
	return d & (BitEq64(d, (a^b)&(a^d)) >> 63)
}
