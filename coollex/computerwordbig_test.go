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

func BenchmarkComputerWordBig(b *testing.B) {
	for n := 0; n < b.N; n++ {
		w, _ := NewComputerWordBig(benchAlgorithmN, benchAlgorithmK)
		for w.hasNext() {
			w.next()
		}
	}
}

func TestComputerWordBig(t *testing.T) {
	testCoollex(t, func(n, k uint) (coollexAlgorithm, error) {
		w, err := NewComputerWordBig(n, k)
		return &w, err
	})
}
