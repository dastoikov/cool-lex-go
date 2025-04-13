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
	"github.com/dastoikov/cool-lex-go/v2/simplemath"
	"iter"
	"math"
	"math/bits"
)

// ComputerWord32 implements the register-based (computer words) algorithm from the paper,
// see 3.3. Implementation in Computer Words, page 10.
// The implementation here is based on 32-bit registers, allowing for `n<=31`.
type ComputerWord32 struct {
	r2, r3 int32 // names as in the paper; r2 is mask, r3 stores the combination
}

// hasNext reports whether more combinations are available
func (word *ComputerWord32) hasNext() bool {
	return (word.r3 & word.r2) == 0
}

// next advances to the next combination in cool-lex order
func (word *ComputerWord32) next() {
	r3 := word.r3

	r0 := r3 & (r3 + 1)
	r1 := r0 ^ (r0 - 1)
	r0 = r1 + 1
	r1 = r1 & r3
	r0 = simplemath.DozB32(r0&r3, 1) // todo(das): allow external doz func?

	word.r3 = r3 + r1 - r0
}

func elements32(v int32) Elements {
	return func(yield func(uint) bool) {
		for r := uint32(v); r != 0; {
			ntz := bits.TrailingZeros32(r)
			if !yield(uint(ntz)) {
				break
			}
			// clear the rightmost 1-bit
			// see Hacker's Delight by Henry S. Warren, Jr., 2-1 Manipulating the Rightmost Bits, page 11.
			r &= r - 1
		}
	}
}

// newComputerWord32 initializes the algorithm for the specified number of 0-bits (s) and number of 1-bits (t).
// Precondition: `t>0`.
func newComputerWord32(s, t uint) ComputerWord32 {
	var r2 int32 = 1 << (s + t)
	var r3 int32 = (1 << t) - 1
	return ComputerWord32{r2, r3}
}

// Elements returns an iterator over the elements selected for the current combination.
func (word *ComputerWord32) Elements() Elements {
	return elements32(word.r3)
}

// Combinations returns an iterator over the generated combinations.
func (word *ComputerWord32) Combinations() Combinations {
	return func(yield func(Elements) bool) {
		for word.hasNext() && yield(word.Elements()) {
			word.next()
		}
	}
}

// Words returns an iterator over the generated combinations as follows:
// * a combination is represented by an `int32` value
// * in a combination, bits that are set represent the elements selected for the combination
// * the `n` least-significant bits store the combination, the other bits are cleared
func (word *ComputerWord32) Words() iter.Seq[int32] {
	return func(yield func(int32) bool) {
		for word.hasNext() && yield(word.r3) {
			word.next()
		}
	}
}

// NewComputerWord32 returns a combinations generator that yields combinations in Cool-lex order, working internally
// with 32-bit "registers".
//
// n: number of elements to combine; n>=k and n<32 must hold.
//
// k: number of elements in each combination.
//
// It is an error to pass arguments such that n < k.
// It is an error to pass arguments such that n >=32.
func NewComputerWord32(n, k uint) (ComputerWord32, error) {
	if n < k {
		return ComputerWord32{}, fmt.Errorf("n (%d) less than k (%d)", n, k)
	}
	if n >= 32 {
		return ComputerWord32{}, fmt.Errorf("n (%d) greater than 31, consider using LinkedList", n)
	}
	if k == 0 {
		return ComputerWord32{math.MinInt32, math.MinInt32}, nil // anything such that r2&r3 != 0
	}
	return newComputerWord32(n-k, k), nil
}
