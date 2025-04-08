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
	benchElementsLowDensity32  int32 = 0b0110_0000_0000_0000_0000_0000_0000_0000
	benchElementsMidDensity32  int32 = 0b0111_1000_1111_0000_1111_0000_1111_0000
	benchElementsHighDensity32 int32 = 0b0001_1111_1111_1111_1111_1111_1111_1111
)

func BenchmarkComputerWord32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		w, _ := NewComputerWord32(benchAlgorithmN, benchAlgorithmK)
		for w.hasNext() {
			w.next()
		}
	}
}
func BenchmarkElementsLowDensity32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements32(benchElementsLowDensity32) {
			_ = e
		}
	}
}
func BenchmarkElementsMidDensity32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements32(benchElementsMidDensity32) {
			_ = e
		}
	}
}
func BenchmarkElementsHighDensity32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements32(benchElementsHighDensity32) {
			_ = e
		}
	}
}
func BenchmarkElements32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for e := range elements32(benchElementsLowDensity32) {
			_ = e
		}
		for e := range elements32(benchElementsMidDensity32) {
			_ = e
		}
		for e := range elements32(benchElementsHighDensity32) {
			_ = e
		}
	}
}
func TestElements32(t *testing.T) {
	coollex, _ := NewComputerWord32(31, 2)
	for combination := range coollex.Words() {
		e := elements32(combination)
		if word := toInt32(e); word != combination {
			t.Fatalf("expect %v, actual %v", combination, word)
		}
	}
	coollex, _ = NewComputerWord32(31, 29)
	for combination := range coollex.Words() {
		e := elements32(combination)
		if word := toInt32(e); word != combination {
			t.Fatalf("expect %v, actual %v", combination, word)
		}
	}
}

func TestComputerWord32(t *testing.T) {
	testCoollex(t, func(n, k uint) (coollexAlgorithm, error) {
		w, err := NewComputerWord32(n, k)
		return &w, err
	})

	if _, err := NewComputerWord32(32, 1); err == nil {
		t.Fatalf("error is expected for n>=32")
	}
}

func toInt32(elements Elements) int32 {
	result := int32(0)
	for element := range elements {
		result |= 1 << element
	}
	return result
}
