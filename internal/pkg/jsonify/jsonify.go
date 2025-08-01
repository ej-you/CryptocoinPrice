// Package jsonify provides JSONify interface to marshal/unmarshal any structs.
package jsonify

import "encoding/json"

// Jsonify provides methods to (de)serialize any struct.
type Jsonify interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// Jsonify implementation.
type standartJsonify struct{}

// New returns new Jsonify.
func New() Jsonify {
	return &standartJsonify{}
}

// Marshal serializes JSON into bytes.
func (j standartJsonify) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal deserializes JSON from bytes.
func (j standartJsonify) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
