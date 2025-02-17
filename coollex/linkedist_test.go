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
	"slices"
	"testing"
)

func TestValueTrueNodes(t *testing.T) {
	listOf := func(values []bool) *node {
		if len(values) == 0 {
			return nil
		}
		nodes := make([]node, len(values))
		nodes[0].value = values[0]
		for i := 1; i < len(values); i++ {
			nodes[i].value = values[i]
			nodes[i-1].next = &nodes[i]
		}
		return &nodes[0]
	}
	test := func(values []bool, expect []uint) {
		actual := make([]uint, 0, len(expect))
		for elem := range listOf(values).valueTrueNodes() {
			actual = append(actual, elem)
		}
		slices.Sort(expect)
		if !slices.Equal(expect, actual) {
			t.Fatalf("Expected %v, got %v, for values %v", expect, actual, values)
		}
	}
	test([]bool{true}, []uint{0})
	test([]bool{false}, []uint{})
	test([]bool{false, false}, []uint{})
	test([]bool{false, true}, []uint{1})
	test([]bool{true, false}, []uint{0})
	test([]bool{false, true, true}, []uint{1, 2})
	test([]bool{true, false, true}, []uint{0, 2})
	test([]bool{true, true, false}, []uint{0, 1})
	test([]bool{true, true, true}, []uint{0, 1, 2})
	test([]bool{false, false, false}, []uint{})
}

func TestLinkedList(t *testing.T) {
	testNoCombs := func(n, k uint) {
		alg, _ := NewLinkedList(n, k)
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

		alg, err := NewLinkedList(n, k)
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
	test(49, 6)
	test(9, 9)
	test(1, 1)

	testNoCombs(2, 0)
	testNoCombs(9, 0)
}
