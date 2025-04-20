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

/*
Package coollex provides implementations of the different Cool-lex algorithms presented in
"The Coolest Way to Generate Combinations" paper by Frank Ruskey and Aaron Williams.

Hats off.

Refer to http://webhome.cs.uvic.ca/~ruskey/Publications/Coollex/CoolComb.html:

	Section 3.2. Iterative Algorithms.
	Section 3.3. Implementation in Computer Words.
*/
package coollex

import "iter"

// Elements is an iterator over the elements of a combination. The iterator yields exactly
// k elements. Every element is in the range [0, n).
type Elements = iter.Seq[uint]

// Combinations is an iterator over combinations.
type Combinations = iter.Seq[Elements]

// internal type to facilitate generic testing
type coollexAlgorithm interface {
	Combinations() Combinations
}
