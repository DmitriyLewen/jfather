package jfather

import (
	"fmt"
	"strconv"
	"unicode/utf16"
)

var escapes = map[rune]string{
	'\\': "\\",
	'/':  "/",
	'"':  "\"",
	'n':  "\n",
	'r':  "\r",
	'b':  "\b",
	'f':  "\f",
	't':  "\t",
}

func (p *parser) parseString() (Node, error) {

	n := p.newNode(KindString)

	b, err := p.next()
	if err != nil {
		return nil, err
	}

	if b != '"' {
		return nil, p.makeError("expecting string delimiter")
	}

	var str string

	var inEscape bool
	var inHex bool
	var hex []rune

	for {
		c, err := p.next()
		if err != nil {
			return nil, err
		}
		if inHex {
			switch {
			case c >= 'a' && c <= 'f', c >= 'A' && c <= 'F', c >= '0' && c <= '9':
				hex = append(hex, c)
				switch len(hex) {
				case 4:
					inHex = false
					// If we can't convert hex with 4 characters - we expect it to be a surrogate character
					// If not - we'll return an error later
					if char, ok := p.tryUnicode(hex); ok {
						str += char
						hex = nil
					}
				case 8: // surrogate
					char, err := p.trySurrogate(hex)
					if err != nil {
						return nil, p.makeError("invalid unicode character '%s'", err)
					}
					str += char
					inHex = false
					hex = nil
				}
			default:
				return nil, p.makeError("invalid hexedecimal escape sequence '\\%s%c'", string(hex), c)
			}
		} else if inEscape {
			inEscape = false
			if c == 'u' {
				inHex = true
				continue
			}
			seq, ok := escapes[c]
			if !ok {
				return nil, p.makeError("invalid escape sequence '\\%c'", c)
			}
			str += seq
		} else {
			switch c {
			case '\\':
				inEscape = true
			case '"':
				if err := p.hexFinished(hex); err != nil {
					return nil, err
				}
				n.raw = str
				n.end = p.position
				return n, nil
			default:
				if err := p.hexFinished(hex); err != nil {
					return nil, err
				}
				if c < 0x20 || c > 0x10FFFF {
					return nil, p.makeError("invalid unescaped character '0x%X'", c)
				}
				str += string(c)
			}
		}

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
