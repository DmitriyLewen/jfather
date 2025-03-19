package jfather

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf16"
)

func (p *parser) parseString() (Node, error) {

	n := p.newNode(KindString)

	b, err := p.next()
	if err != nil {
		return nil, err
	}

	if b != '"' {
		return nil, p.makeError("expecting string delimiter")
	}

	// Start with opening quote
	buffer := []byte{'"'}
	for {
		c, err := p.next()
		if err != nil {
			return nil, err
		}
		if c == '"' {
			buffer = append(buffer, '"')
			var result string
			if err := json.Unmarshal(buffer, &result); err != nil {
				return nil, p.makeError("invalid JSON string: %v", err)
			}
			n.raw = result
			n.end = p.position
			return n, nil
		}

		buffer = append(buffer, byte(c))
	}
}

func (p *parser) tryUnicode(hex []rune) (string, bool) {
	char, err := strconv.Unquote(fmt.Sprintf("'\\u%s'", string(hex)))
	if err != nil {
		return "", false
	}
	return char, true
}

func (p *parser) trySurrogate(hex []rune) (string, error) {
	var surrogatePair []uint16

	var hexPairs [][]rune
	for _, hexPair := range append(hexPairs, hex[:4], hex[4:]) {
		parsedPair, err := strconv.ParseUint(string(hexPair), 16, 16)
		if err != nil {
			return "", p.makeError("invalid unicode character '%s': %s", string(hexPair), err)
		}
		surrogatePair = append(surrogatePair, uint16(parsedPair))

	}
	decoded := utf16.Decode(surrogatePair)
	return string(decoded), nil
}

func (p *parser) hexFinished(hex []rune) error {
	if len(hex) > 0 {
		return p.makeError("invalid hex character: %s", "\\u"+string(hex))
	}
	return nil
}
