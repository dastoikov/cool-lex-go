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
	"github.com/dastoikov/cool-lex-go/simplemath"
	"testing"
)

const (
	// `n` and `k` for benchmarking the different Cool-lex algorithms
	benchAlgorithmN = 33
	benchAlgorithmK = 3
)

func testCoollex(t *testing.T, generator func(n, k uint) (coollexAlgorithm, error)) {
	testNoCombs := func(n, k uint) {
		alg, _ := generator(n, k)
		for combination := range alg.Combinations() {
			for element := range combination {
				t.Fatalf("element %v found for n %v and k %v", element, n, k)
			}
			t.Fatalf("combinations found for n %v and k %v", n, k)
		}
	}
	test := func(n, k uint) {
		// array index denotes an element
		// value at given index denotes how many times this element appeared in a combination
		hits := make([]uint, n)

		alg, err := generator(n, k)
		if err != nil {
			t.Fatal(err)
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
				t.Fatalf("number of elements in a combination: expected %d, got %d, for n %d", k, numElem, n)
			}
		}
		expectedNumComb, err := simplemath.NumComb(n, k)
		if err != nil {
			t.Fatal(err)
		}
		if numComb != expectedNumComb {
			t.Fatalf("number of combinations: expected %d, got %d, for n %d", expectedNumComb, numComb, n)
		}
		occur, err := simplemath.NumComb(n-1, k-1)
		if err != nil {
			t.Fatal(err)
		}
		for element, hit := range hits {
			if occur != hit {
				t.Fatalf("number of combinations where each element appears: expected %d, got %d, for element %d, k %d, and n %d", occur, hit, element, k, n)
			}
		}
	}
	test(15, 6)
	test(15, 7)
	test(25, 13)
	test(9, 9)
	test(1, 1)

	testNoCombs(2, 0)
	testNoCombs(9, 0)
}
