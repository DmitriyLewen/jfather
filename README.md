# jfather

Parse JSON with line numbers and more!

This is a JSON parsing module that provides additional information during the unmarshalling process, such as line numbers, columns etc.

You can use jfather to unmarshal JSON just like the `encoding/json` package, except implementing your own unmarshalling functionality to gather metadata can be done by implementing the `jfather.Umarshaller` interface, which requires a method with the signature `UnmarshalJSONWithMetadata(node jfather.Node) error`. A full example is below.

## Full Example

```golang
package main

import (
	"fmt"

	"github.com/liamg/jfather"
)

type ExampleParent struct {
	Child *ExampleChild `json:"child"`
}

type ExampleChild struct {
	Name   string
	Line   int
	Column int
}

func (t *ExampleChild) UnmarshalJSONWithMetadata(node jfather.Node) error {
	t.Line = node.Range().Start.Line
	t.Column = node.Range().Start.Column
	return node.Decode(&t.Name)
}

func main() {
	input := []byte(`{
	"child": "secret"
}`)
	var parent ExampleParent
	if err := jfather.Unmarshal(input, &parent); err != nil {
		panic(err)
	}

	fmt.Printf("Child value is at line %d, column %d, and is set to '%s'\n",
		parent.Child.Line, parent.Child.Column, parent.Child.Name)

	// outputs:
	//  Child value is at line 2, column 12, and is set to 'secret'
}
```
