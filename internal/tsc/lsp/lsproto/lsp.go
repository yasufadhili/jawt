package lsproto

import (
	"encoding/json"
	"fmt"
)

type DocumentUri string // !!!

type URI string // !!!

type Method string

type Nullable[T any] struct {
	Value T
	Null  bool
}

func ToNullable[T any](v T) Nullable[T] {
	return Nullable[T]{Value: v}
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{Null: true}
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.Null {
		return []byte(`null`), nil
	}
	return json.Marshal(n.Value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	*n = Nullable[T]{}
	if string(data) == `null` {
		n.Null = true
		return nil
	}
	return json.Unmarshal(data, &n.Value)
}

func unmarshalPtrTo[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T: %w", (*T)(nil), err)
	}
	return &v, nil
}

func unmarshalAny(data []byte) (any, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal any: %w", err)
	}
	return v, nil
}

func unmarshalEmpty(data []byte) (any, error) {
	if len(data) != 0 {
		return nil, fmt.Errorf("expected empty, got: %s", string(data))
	}
	return nil, nil
}

func assertOnlyOne(message string, values ...bool) {
	count := 0
	for _, v := range values {
		if v {
			count++
		}
	}
	if count != 1 {
		panic(message)
	}
}
