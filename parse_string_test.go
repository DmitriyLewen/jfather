package jfather

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	example := []byte(`"hello"`)
	var output string
	err := Unmarshal(example, &output)
	require.NoError(t, err)
	assert.Equal(t, "hello", output)
}

func Test_StringWithUnicode(t *testing.T) {
	example := []byte(`"\u0440 and \u0441"`)
	var output string
	err := Unmarshal(example, &output)
	require.NoError(t, err)
	assert.Equal(t, "р and с", output)
}

func Test_StringWithSurrogateUnicode(t *testing.T) {
	example := []byte(`"\ud83d\udee0 and \ud83d\udee1"`)
	var output string
	err := Unmarshal(example, &output)
	require.NoError(t, err)
	assert.Equal(t, "🛠 and 🛡", output)
}

func Test_StringWithOnlyInvalidUnicode(t *testing.T) {
	example := []byte(`"\ud83d"`)
	var output string
	err := Unmarshal(example, &output)
	require.NoError(t, err)
	assert.Equal(t, "�", output)
}

func Test_StringWithInvalidUnicode(t *testing.T) {
	example := []byte(`"\ud83d something"`)
	var output string
	err := Unmarshal(example, &output)
	require.NoError(t, err)
	assert.Equal(t, "� something", output)
}

func Test_StringToUninitialisedPointer(t *testing.T) {
	example := []byte(`"hello"`)
	var str *string
	err := Unmarshal(example, str)
	require.Error(t, err)
	assert.Nil(t, str)
}

func Test_String_ToInterface(t *testing.T) {
	example := []byte(`"hello"`)
	var output interface{}
	err := Unmarshal(example, &output)
	require.NoError(t, err)
	assert.Equal(t, "hello", output)
}
