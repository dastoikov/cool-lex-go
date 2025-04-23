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

package coollex

import (
	"fmt"
	"iter"
	"math"
	"math/bits"
)

// ComputerWord64 implements the register-based (computer words) algorithm from the paper,
// see 3.3. Implementation in Computer Words, page 10.
// The implementation here is based on 64-bit registers, allowing for `n<=63`.
type ComputerWord64 struct {
	r2, r3 int64 // names as in the paper; r2 is mask, r3 stores the combination
}

// hasNext reports whether more combinations are available
func (word *ComputerWord64) hasNext() bool {
	return (word.r3 & word.r2) == 0
}

// next advances to the next combination in cool-lex order
func (word *ComputerWord64) next() {
	r3 := word.r3

	r0 := r3 & (r3 + 1)
	r1 := r0 ^ (r0 - 1)
	r0 = r1 + 1
	r1 = r1 & r3

	// equivalent to r0 = simplemath.DozB64(r0&r3, 1)
	// however, rephrasing it as follows results in a slightly better assembly
	if r0&r3 > 0 {
		r0 -= 1
	} else {
		r0 = 0
	}

	word.r3 = r3 + r1 - r0
}

func elements64(v int64) Elements {
	return func(yield func(uint) bool) {
		for r := uint64(v); r != 0; {
			ntz := bits.TrailingZeros64(r)
			if !yield(uint(ntz)) {
				break
			}
			// clear the rightmost 1-bit
			// see Hacker's Delight by Henry S. Warren, Jr., 2-1 Manipulating the Rightmost Bits, page 11.
			r &= r - 1
		}
	}
}

// newComputerWord64 initializes the algorithm for the specified number of 0-bits (s) and number of 1-bits (t).
// Precondition: `t>0`.
func newComputerWord64(s, t uint) ComputerWord64 {
	var r2 int64 = 1 << (s + t)
	var r3 int64 = (1 << t) - 1
	return ComputerWord64{r2, r3}
}

// Elements returns an iterator over the elements selected for the current combination.
func (word *ComputerWord64) Elements() Elements {
	return elements64(word.r3)
}

// Combinations returns an iterator over the generated combinations.
func (word *ComputerWord64) Combinations() Combinations {
	return func(yield func(Elements) bool) {
		for word.hasNext() && yield(word.Elements()) {
			word.next()
		}
	}
}

// Words returns an iterator over the generated combinations as follows:
// * a combination is represented by an `int64` value
// * in a combination, bits that are set represent the elements selected for the combination
// * the `n` least-significant bits store the combination, the other bits are cleared
func (word *ComputerWord64) Words() iter.Seq[int64] {
	return func(yield func(int64) bool) {
		for word.hasNext() && yield(word.r3) {
			word.next()
		}
	}
}

// NewComputerWord64 returns a combinations generator that yields combinations in Cool-lex order, working internally
// with 64-bit "registers".
//
// n: number of elements to combine; n>=k and n<64 must hold.
//
// k: number of elements in each combination.
//
// It is an error to pass arguments such that n < k.
// It is an error to pass arguments such that n >=64.
func NewComputerWord64(n, k uint) (ComputerWord64, error) {
	if n < k {
		return ComputerWord64{}, fmt.Errorf("n (%d) less than k (%d)", n, k)
	}
	if n >= 64 {
		return ComputerWord64{}, fmt.Errorf("n (%d) greater than 63, consider using LinkedList", n)
	}

	if k == 0 {
		return ComputerWord64{math.MinInt64, math.MinInt64}, nil // anything such that r2&r3 != 0
	}
	return newComputerWord64(n-k, k), nil
}
