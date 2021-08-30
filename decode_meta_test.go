package jfather

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestParent struct {
	Child *TestChild `json:"child"`
}

type TestChild struct {
	Name   string
	Line   int
	Column int
}

func (t *TestChild) UnmarshalJSONWithMetadata(node Node) error {
	t.Line = node.Range().Start.Line
	t.Column = node.Range().Start.Column
	return node.Decode(&t.Name)
}

func Test_DecodeWithMetadata(t *testing.T) {
	example := []byte(`
{
	"child": "secret"
}
`)
	var parent TestParent
	require.NoError(t, Unmarshal(example, &parent))
	assert.Equal(t, 3, parent.Child.Line)
	assert.Equal(t, 11, parent.Child.Column)
	assert.Equal(t, "secret", parent.Child.Name)
}
