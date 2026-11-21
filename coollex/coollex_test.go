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
	"testing"
)

// coollexAlgorithm captures the common combination-generation interface of
// the various Cool-lex algorithms, in order to facilitate uniform testing.
type coollexAlgorithm interface {
	// Combinations returns an iterator over the generated combinations.
	// The name is chosen deliberately, to match actual implementations.
	Combinations() Combinations
}

// coollexConstructor captures the common constructor interface of the various
// Cool-lex algorithms, in order to facilitate uniform testing.
type coollexConstructor func(n, k uint) (coollexAlgorithm, error)

// construct adapts a Cool-lex algorithm constructor function so that it
// appears to be returning a `coollexAlgorithm`.
func construct[T any, PT interface {
	*T
	coollexAlgorithm
}](f func(uint, uint) (T, error)) coollexConstructor {
	return func(n, k uint) (coollexAlgorithm, error) {
		w, err := f(n, k)
		var ptr PT = &w
		return ptr, err
	}
}

const (
	// `n` and `k` for benchmarking the different Cool-lex algorithms
	benchAlgorithmN = 31
	benchAlgorithmK = 3
)

// verifyCombsEmpty verifies that the algorithm does not yield combinations
// for `n` and `k`, returning an error if any are yielded.
func (c coollexConstructor) verifyCombsEmpty(n, k uint) error {
	alg, _ := c(n, k)
	for combination := range alg.Combinations() {
		for element := range combination {
			return fmt.Errorf("element %v found for n %v and k %v", element, n, k)
		}
		return fmt.Errorf("combinations found for n %v and k %v", n, k)
	}
	return nil
}

// verifyCombs verifies that the algorithm yields combinations for `n` and `k`
// such that:
//   - the number of combinations yielded is correct
//   - the number of elements yielded in each combination is correct
//   - the number of combinations where each element occurred is correct
//
// It does not test combinations are yielded in Cool-lex order.
// It returns an error if `generator` fails to yield combinations.
func (c coollexConstructor) verifyCombs(n, k uint) error {
	// how many times (value at index) an element (index in slice) appeared in a combination
	hits := make([]uint, n)

	alg, err := c(n, k)
	if err != nil {
		return err
	}

	// total number of combinations yielded by the algorithm
	numComb := uint(0)
	for combination := range alg.Combinations() {
		numComb++
		// number of elements in this combination
		numElem := uint(0)
		for element := range combination {
			numElem++
			hits[element]++
		}
		if k != numElem {
			return fmt.Errorf("number of elements in a combination: expected %d, got %d, for n %d", k, numElem, n)
		}
	}
	expectedNumComb, err := simplemath.NumComb(n, k)
	if err != nil {
		return err
	}
	if numComb != expectedNumComb {
		return fmt.Errorf("number of combinations: expected %d, got %d, for n %d and k %d", expectedNumComb, numComb, n, k)
	}

	// number of combinations where each element occurred
	occur, err := simplemath.NumComb(n-1, k-1)
	if err != nil {
		return err
	}
	for element, hit := range hits {
		if occur != hit {
			return fmt.Errorf("number of combinations where each element appears: expected %d, got %d, for element %d, n %d, and k %d ", occur, hit, element, n, k)
		}
	}
	return nil
}
func testCoollex(t *testing.T, c coollexConstructor) {
	testCases := []struct {
		n, k uint
		test func(n, k uint) error
	}{
		{15, 8, c.verifyCombs},
		{15, 7, c.verifyCombs},
		{9, 9, c.verifyCombs},
		{1, 1, c.verifyCombs},
		{9, 0, c.verifyCombsEmpty},
		{2, 0, c.verifyCombsEmpty},
	}
	for _, tc := range testCases {
		err := tc.test(tc.n, tc.k)
		if err != nil {
			t.Fatal(err)
		}
	}
}
