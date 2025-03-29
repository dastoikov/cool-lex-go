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
	"slices"
	"testing"
)

func BenchmarkLinkedList(b *testing.B) {
	for n := 0; n < b.N; n++ {
		list, _ := NewLinkedList(benchAlgorithmN, benchAlgorithmK)
		_ = list
		for list.hasNext() {
			list.next()
		}
	}
}

// Number of nodes for benchmarking newLinkedList.
// That is, the number of elements to choose from, or n.
const (
	HIGH   = uint(360_000_000) //
	MEDIUM = uint(780_000)
	LOW    = uint(90)
)

func BenchmarkNewLinkedListLow(b *testing.B) {
	lo := LOW / 3
	for i := 0; i < b.N; i++ {
		_ = newLinkedList(lo*2, lo)
	}
}

func BenchmarkNewLinkedListMedium(b *testing.B) {
	md := MEDIUM / 3
	for i := 0; i < b.N; i++ {
		_ = newLinkedList(md*2, md)
	}
}

func BenchmarkNewLinkedListHigh(b *testing.B) {
	up := HIGH / 3
	for i := 0; i < b.N; i++ {
		_ = newLinkedList(up*2, up)
	}
}

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
	testCoollex(t, func(n, k uint) (coollexAlgorithm, error) {
		list, err := NewLinkedList(n, k)
		return &list, err
	})
}
