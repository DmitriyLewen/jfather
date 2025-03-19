package jfather

import (
	"encoding/json"
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
