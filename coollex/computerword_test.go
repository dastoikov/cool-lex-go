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

func BenchmarkComputerWord(b *testing.B) {
	for n := 0; n < b.N; n++ {
		w, _ := NewComputerWord64(benchAlgorithmN, benchAlgorithmK)
		for w.hasNext() {
			w.next()
		}
	}
}

func BenchmarkElements(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := uint64(0b11111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_00_1)
		for e := range elements(int64(v)) {
			_ = e
		}
	}
}

func TestElements(t *testing.T) {
	v := uint64(0b11111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_0000_1111_00_1)
	expect := [34]uint64{
		0,
		3, 4, 5, 6,
		11, 12, 13, 14,
		19, 20, 21, 22,
		27, 28, 29, 30,
		35, 36, 37, 38,
		43, 44, 45, 46,
		51, 52, 53, 54,
		59, 60, 61, 62, 63,
	}
	actual := [34]uint64{}

	i := 0
	for e := range elements(int64(v)) {
		actual[i] = uint64(e)
		i++
	}
	if actual != expect {
		t.Fatalf("expect %v, actual %v", expect, actual)
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
