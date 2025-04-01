# cool-lex-go

The Cool-lex order and algorithms were invented and authored by Frank Ruskey and Aaron Williams (<http://webhome.cs.uvic.ca/~ruskey/Publications/Coollex/CoolComb.html>).
You may need to obtain permission to use Cool-lex, as governed by applicable laws and academic practices.

The code in this repository is authored by the cool-lex-go [contributors](CONTRIBUTORS), and is licensed under Apache License, version 2.0 license.

## Examples

**LinkedList**

```go
package main
import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/coollex"
)

func main() {
	// no error for n=3, k=2
	generator, _ := coollex.NewLinkedList(3, 2)
	for combination := range generator.Combinations() {
		for element := range combination {
			fmt.Print(element)
		}
		fmt.Println()
	}
	// prints:
	// 01
	// 12
	// 02
}
```

**ComputerWord, 64-bit**

```go
package main
import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/coollex"
)

func main() {
	// no error for n=3, k=2
	generator, _ := coollex.NewComputerWord64(3, 2)
	for combination := range generator.Combinations() {
		for element := range combination {
			fmt.Print(element)
		}
		fmt.Println()
	}
	// prints:
	// 01
	// 12
	// 02
}
```

Alternatively, if it is preferable to work with the combination stored in an `int64` value directly:

```go
package main

import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/coollex"
)

func main() {
	// no error for n=3, k=2
	generator, _ := coollex.NewComputerWord64(3, 2)
	for combination := range generator.Words() {
		fmt.Printf("%03b\n", combination)
	}
	// prints:
	// 011
	// 110
	// 101
}
```

Alternatively, if it is preferable to use a custom function for retrieving the set bits from an `int64` combination:

```go
package main

import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/coollex"
	"math/bits"
)

func main() {
	// If popcount is fast...
	// (Note: the code below is just for illustration purposes as the built-in elements() is generally faster.)
	elements := func(combination int64) coollex.Elements {
		return func(yield func(uint) bool) {
			for r := uint64(combination); r != 0; {
				ntz := bits.OnesCount64(^r & (r - 1)) // set trailing 0-bits, clear all others
				if !yield(uint(ntz)) {
					break
				}
				r &= r - 1 // turn off the right-most 1-bit
			}
		}
	}

	generator, _ := coollex.NewComputerWord64(3, 2)
	for combination := range generator.Words() {
		for element := range elements(combination) {
			fmt.Print(element)
		}
		fmt.Println()
	}
	// prints:
	// 01
	// 12
	// 02
}
```