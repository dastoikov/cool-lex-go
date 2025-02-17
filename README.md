# cool-lex-go

The Cool-lex order and algorithms were invented and authored by Frank Ruskey and Aaron Williams (<http://webhome.cs.uvic.ca/~ruskey/Publications/Coollex/CoolComb.html>).
You may need to obtain permission to use Cool-lex, as governed by applicable laws and academic practices.

The code in this repository is authored by the cool-lex-go [contributors](CONTRIBUTORS), and is licensed under Apache License, version 2.0 license.

## Example
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