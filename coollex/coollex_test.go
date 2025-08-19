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

const (
	// `n` and `k` for benchmarking the different Cool-lex algorithms
	benchAlgorithmN = 31
	benchAlgorithmK = 3
)

// verifyNoCombs verifies that `generator` does not yield combinations for `n` and `k`.
// It returns an error if any are yielded.
func verifyNoCombs(n, k uint, generator func(n, k uint) (coollexAlgorithm, error)) error {
	alg, _ := generator(n, k)
	for combination := range alg.Combinations() {
		for element := range combination {
			return fmt.Errorf("element %v found for n %v and k %v", element, n, k)
		}
		return fmt.Errorf("combinations found for n %v and k %v", n, k)
	}
	return nil
}

// verifyCombs verifies that combinations yielded by `generator` for `n` and `k` have the following properties:
//   - the number of combinations generated is correct;
//   - the number of elements in each combination is correct;
//   - the number of combinations where each element occurs is correct.
//
// It does not test combinations are yielded in Cool-lex order.
// It returns an error if `generator` fails to yield combinations.
func verifyCombs(n, k uint, generator func(n, k uint) (coollexAlgorithm, error)) error {
	// array index denotes an element
	// value at given index denotes how many times this element appeared in a combination
	hits := make([]uint, n)

	alg, err := generator(n, k)
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
func testCoollex(t *testing.T, generator func(n, k uint) (coollexAlgorithm, error)) {
	testCases := []struct {
		n, k uint
		test func(n, k uint, generator func(n, k uint) (coollexAlgorithm, error)) error
	}{
		{15, 8, verifyCombs},
		{15, 7, verifyCombs},
		{9, 9, verifyCombs},
		{1, 1, verifyCombs},
		{9, 0, verifyNoCombs},
		{2, 0, verifyNoCombs},
	}
	for _, tc := range testCases {
		err := tc.test(tc.n, tc.k, generator)
		if err != nil {
			t.Fatal(err)
		}
	}
}
