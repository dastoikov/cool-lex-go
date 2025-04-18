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
	"testing"
)

const (
	benchElementsLowDensity64  int64 = 0b0110_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000_0000
	benchElementsMidDensity64  int64 = 0b0111_1000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000
	benchElementsHighDensity64 int64 = 0b0001_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111
)

func BenchmarkComputerWord64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		w, _ := NewComputerWord64(benchAlgorithmN, benchAlgorithmK)
		for w.hasNext() {
			w.next()
		}
	}
}
func BenchmarkElementsLowDensity64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements64(benchElementsLowDensity64) {
			_ = e
		}
	}
}
func BenchmarkElementsMidDensity64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements64(benchElementsMidDensity64) {
			_ = e
		}
	}
}
func BenchmarkElementsHighDensity64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements64(benchElementsHighDensity64) {
			_ = e
		}
	}
}
func BenchmarkElements64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements64(benchElementsLowDensity64) {
			_ = e
		}
		for e := range elements64(benchElementsMidDensity64) {
			_ = e
		}
		for e := range elements64(benchElementsHighDensity64) {
			_ = e
		}
	}
}
func TestElements64(t *testing.T) {
	coollex, _ := NewComputerWord64(63, 2)
	for combination := range coollex.Words() {
		e := elements64(combination)
		if word := toInt64(e); word != combination {
			t.Fatalf("expect %v, actual %v", combination, word)
		}
	}
	coollex, _ = NewComputerWord64(63, 61)
	for combination := range coollex.Words() {
		e := elements64(combination)
		if word := toInt64(e); word != combination {
			t.Fatalf("expect %v, actual %v", combination, word)
		}
	}
}

func TestComputerWord64(t *testing.T) {
	testCoollex(t, func(n, k uint) (coollexAlgorithm, error) {
		w, err := NewComputerWord64(n, k)
		return &w, err
	})

	if _, err := NewComputerWord64(64, 1); err == nil {
		t.Fatalf("error is expected for n>=64")
	}
}

func toInt64(elements Elements) int64 {
	result := int64(0)
	for element := range elements {
		result |= 1 << element
	}
	return result
}
