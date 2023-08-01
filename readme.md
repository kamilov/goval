Go library that can populate variable value from a string with data type conversion
```go
package main

import (
	"fmt"
	"github.com/kamilov/goval"
)

func main() {
	var b bool
	var s string
	var m map[string]uint8

	_ = goval.Val(&b, "True")
	_ = goval.Val(&s, "String")
	_ = goval.Val(&m, `{"a": 1, "b": 2}`)

	fmt.Printf("%V\n", b)
	fmt.Printf("%V\n", s)
	fmt.Printf("%V\n", m)
}
```