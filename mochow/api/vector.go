package api

import "encoding/base64"

type Vector interface {
	name() string
	representation() interface{}
}

type FloatVector []float32

func (v FloatVector) name() string {
	return "vectorFloats"
}

func (v FloatVector) representation() interface{} {
	return v
}

type BinaryVector []byte

func (b BinaryVector) name() string {
	return "vector"
}

func (b BinaryVector) representation() interface{} {
	return base64.StdEncoding.EncodeToString(b)
}

type SparseFloatVector map[string]float32

func (s SparseFloatVector) name() string {
	return "vector"
}

func (s SparseFloatVector) representation() interface{} {
	return s
}
