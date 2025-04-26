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
	"math/big"
)

// ComputerWordBig implements the register-based (computer words) algorithm presented in the paper,
// see 3.3. Implementation in Computer Words, page 10.
// The implementation here is based on `big.Int`, allowing for arbitrary `n` and maintaining
// internally a state of just four `big.Int` numbers.
type ComputerWordBig struct {
	// names as in the paper
	r2, r3 *big.Int // r2 is a mask, r3 stores the combination
	r0, r1 *big.Int // auxiliaries, kept in the struct to avoid memory allocations
}

var bigOne = big.NewInt(1)

// hasNext reports whether more combinations are available
func (word *ComputerWordBig) hasNext() bool {
	word.r0.And(word.r3, word.r2)
	return len(word.r0.Bits()) == 0
}

// next advances to the next combination in cool-lex order
func (word *ComputerWordBig) next() {

	// 1: r0 := r3 & (r3 + 1)
	// 2: r1 := r0 ^ (r0 - 1)
	// 3: r0 = r1 + 1
	// 4: r1 = r1 & r3
	// 5: r0 = DOZ(r0&r3, 1)
	// 6: r3 = r3 + r1 - r0

	r3 := word.r3
	r0 := word.r0.Set(r3)
	r0.Add(r0, bigOne).And(r0, r3)

	r1 := word.r1.Set(r0)
	r1.Sub(r1, bigOne).Xor(r1, r0)

	r0.Add(r1, bigOne)
	r1.And(r1, r3)

	r0.And(r0, r3)
	if r0.Cmp(bigOne) >= 0 {
		r0.Sub(r0, bigOne)
	}

	word.r3.Add(r3, r1).Sub(r3, r0)
}

// newComputerWordBig initializes the algorithm for the specified number of 0-bits (s) and number of 1-bits (t).
// Precondition: `t>0`.
func newComputerWordBig(s, t uint) ComputerWordBig {
	aux := make([]big.Int, 4)

	r2 := aux[0].Lsh(bigOne, s+t)
	r3 := aux[1].Lsh(bigOne, t)
	r3.Sub(r3, bigOne)
	return ComputerWordBig{
		r2: r2,
		r3: r3,
		r0: &aux[2],
		r1: &aux[3],
	}
}

// newComputerWordBigEnded initializes the algorithm such that it does not yield any combinations.
// In other words alg.hasNext() returns false.
func newComputerWordBigEnded() ComputerWordBig {
	return ComputerWordBig{
		// anything such that r2&r3 != 0
		r2: bigOne,
		r3: bigOne,
		r0: new(big.Int), // r0 is used by hasNext
	}
}

// Elements returns an iterator over the elements selected for the current combination.
func (word *ComputerWordBig) Elements() Elements {
	return func(yield func(uint) bool) {
		r0 := word.r0.Set(word.r3)
		r1 := word.r1
		for len(r0.Bits()) != 0 {
			ntz := r0.TrailingZeroBits()
			if !yield(ntz) {
				break
			}
			r1.Sub(r0, bigOne) // here and below: clear r0's rightmost set bit
			r0.And(r0, r1)
		}
	}
}

// Combinations returns an iterator over the generated combinations.
func (word *ComputerWordBig) Combinations() Combinations {
	return func(yield func(Elements) bool) {
		for word.hasNext() && yield(word.Elements()) {
			word.next()
		}
	}
}

// Words returns an iterator over the generated combinations as follows:
//   - a combination is represented by a `big.Int` value
//   - in a combination, bits that are set represent the elements selected for the combination
//   - the `n` least-significant bits store the combination, with `k` bits set
//
// Note: Words provides raw access to the internal state of the algorithm and should only be
// used for bit-reading.
func (word *ComputerWordBig) Words() iter.Seq[*big.Int] {
	return func(yield func(*big.Int) bool) {
		for word.hasNext() && yield(word.r3) {
			word.next()
		}
	}
}

// NewComputerWordBig returns a combinations generator that yields combinations in Cool-lex order,
// working internally with arbitrary-size "registers".
//
// n: number of elements to combine; n>=k.
//
// k: number of elements in each combination.
//
// It is an error to pass arguments such that n < k.
func NewComputerWordBig(n, k uint) (ComputerWordBig, error) {
	if n < k {
		return ComputerWordBig{}, fmt.Errorf("n (%d) less than k (%d)", n, k)
	}
	if k == 0 {
		return newComputerWordBigEnded(), nil
	}
	return newComputerWordBig(n-k, k), nil
}
