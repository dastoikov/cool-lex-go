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

/*
Package simplemath offers simple/naive implementations of math operations for the purposes of aiding in writing auto tests.
*/
package simplemath

import (
	"fmt"
	"math/bits"
)

func MulRange(n2, n1 uint) (r uint, err error) {
	if n1 > n2 {
		n1, n2 = n2, n1
	}

	r = n1
	for n1 < n2 {
		n1, err = Add(1, n1)
		if err != nil {
			return
		}
		r, err = Mul(r, n1)
		if err != nil {
			return
		}
	}
	return
}

func Add(a, b uint) (uint, error) {
	sum, carry := bits.Add(a, b, 0)
	if carry == 1 {
		return 0, fmt.Errorf("numeric overflow occurred adding %d and %d", a, b)
	}
	return sum, nil
}

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
