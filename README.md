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
	"github.com/dastoikov/cool-lex-go/v2/coollex"
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
	"github.com/dastoikov/cool-lex-go/v2/coollex"
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
	"github.com/dastoikov/cool-lex-go/v2/coollex"
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

Note: the following code is illustrative.

popcnt_asm_amd64.s:
```asm
#include "textflag.h"

// func popCount(x uint64) uint64
TEXT ·popCount(SB),$0-16
	MOVQ    x+0(FP), CX
	POPCNTQ CX, CX
	MOVQ    CX, ret+8(FP)
	RET
```

popcnt_asm_amd64.go:
```go
//go:nosplit
//go:noinline
func popCount(x uint64) uint64
```

```go
package main

import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/v2/coollex"
)

func main() {
	elements := func(combination int64) coollex.Elements {
		return func(yield func(uint) bool) {
			for r := uint64(combination); r != 0; {
				ntz := popCount(^r & (r - 1)) // 101000 -> popCount(000111)
				if !yield(uint(ntz)) {
					break
				}
				r &= r - 1 // clear the rightmost 1-bit
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

**ComputerWord, 32-bit**

Its usage is analogous to the 64-bit algorithm, with `NewComputerWord32` as
its constructor (compared to `NewComputerWord64`).

```go
package main

import (
	"fmt"
	"github.com/dastoikov/cool-lex-go/v2/coollex"
)

func main() {
	// no error for n=3, k=2
	generator, _ := coollex.NewComputerWord32(3, 2)
	for combination := range generator.Words() {
		fmt.Printf("%03b\n", combination)
	}
	// prints:
	// 011
	// 110
	// 101

	generator, _ = coollex.NewComputerWord32(3, 2)
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

A tailored function may be used to retrieve the set bits from an `int32`
combination.

The following code highlights the `Words()` function. It makes combinations
available as `int32` values. The combination is stored in the `n` LSB bits,
with `1`-bits indicating selected elements. The `32-n` MSB bits are `0`.

The assembly code is illustrative. It may differ depending on the target
architecture.

```asm
// func popCount(x uint32) uint32
TEXT ·popCount32(SB),$0-12
	MOVL    x+0(FP), CX
	POPCNTL CX, CX
	MOVL    CX, ret+8(FP)
	RET
```

```go
//go:nosplit
//go:noinline
func popCount32(x uint32) uint32
func main() {
	elements := func(combination int32) coollex.Elements {
		return func(yield func(uint) bool) {
			for r := uint32(combination); r != 0; {
				ntz := popCount32(^r & (r - 1)) // 101000 -> popCount(000111)
				if !yield(uint(ntz)) {
					break
				}
				r &= r - 1 // clear the rightmost 1-bit
			}
		}
	}

	generator, _ := coollex.NewComputerWord32(3, 2)
	for combination := range generator.Words() {
		for element := range elements(combination) {
			fmt.Print(element)
		}
		fmt.Println()
	}
}
```
