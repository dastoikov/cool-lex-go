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
Package coollex provides implementations of the different Cool-lex algorithms presented in the
_The Coolest Way to Generate Combinations_ paper by Frank Ruskey and Aaron Williams.

Hats off.

See http://webhome.cs.uvic.ca/~ruskey/Publications/Coollex/CoolComb.html, section **3.2. Iterative Algorithms.**
*/
package coollex

import (
	"fmt"
	"iter"
)

type node struct {
	value bool
	next  *node
}

// valueTrueNodes returns an iterator that yields the indices of nodes with a true value
// in the list of nodes reachable from node n.
//
// Example:
// (1: node value is true; 0: node value is false)
//
//	           list:  1101001
//	                  ^^ ^  ^
//	iterator yields:  01 3  6
func (n *node) valueTrueNodes() iter.Seq[uint] {
	return func(yield func(uint) bool) {
		var i uint
		for curr := n; curr != nil && (!curr.value || yield(i)); curr = curr.next {
			i++
		}
	}
}

func (n *node) createNext(value bool) *node {
	next := &node{
		value, nil,
	}
	n.next = next
	return next
}

// LinkedList implements the LinkedList algorithm from the paper.
type LinkedList struct {
	// b, x are named as found in the research paper
	// b - the head of the list; this is the node with the greatest "index"
	// x - the first node, tail-to-head, whose value is 1 and whose predecessor's value is 0
	b, x *node
}

// newLinkedList creates a new LinkedList with the specified number of 0-bits (s) and number of 1-bits (t; precondition: t>0).
func newLinkedList(s, t uint) LinkedList {
	// batch allocation of nodes still needs proper benchmarking before it can replace the code below

	b := &node{
		value: true,
		next:  nil,
	}
	x := b

	// initial state: ones to the left, zeros to the right
	for ; t > 1; t-- {
		x = x.createNext(true)
	}
	var last = x
	for ; s > 0; s-- {
		last = last.createNext(false)
	}
	return LinkedList{b, x}
}

// hasNext reports whether more combinations are available
func (list *LinkedList) hasNext() bool {
	return list.x.next != nil
}

// next advances the LinkedList to the next combination in cool-lex order
func (list *LinkedList) next() {
	y := list.x.next
	list.x.next = list.x.next.next
	y.next = list.b
	list.b = y

	if !list.b.value && list.b.next.value {
		list.x = list.b.next
	}
}

// ElementIterator is an iterator over the elements of a combination.
type ElementIterator = iter.Seq[uint]

// CombinationIterator is an iterator over combinations.
type CombinationIterator = iter.Seq[ElementIterator]

// Elements returns an iterator over the elements selected for the current combination.
func (list *LinkedList) Elements() ElementIterator {
	return list.b.valueTrueNodes()
}

// Combinations returns an iterator over the generated combinations.
func (list *LinkedList) Combinations() CombinationIterator {
	// k=0 -> b=nil (the list has no head)
	if list.b == nil {
		return func(yield func(ElementIterator) bool) {}
	}
	return func(yield func(ElementIterator) bool) {
		//the algorithm is initially positioned at the first combination
		for yield(list.Elements()) && list.hasNext() {
			list.next()
		}
	}
}

// NewLinkedList returns a combinations generator that yields combinations in Cool-lex order, implementing internally
// the LinkedList Cool-lex algorithm.
//
// n: number of elements to combine; n>=k must hold.
//
// k: number of elements in each combination.
//
// It is an error to pass arguments such that n < k.
func NewLinkedList(n, k uint) (LinkedList, error) {
	if n < k {
		return LinkedList{}, fmt.Errorf("n (%d) less than k (%d)", n, k)
	}
	if k == 0 {
		return LinkedList{}, nil
	}
	return newLinkedList(n-k, k), nil
}
