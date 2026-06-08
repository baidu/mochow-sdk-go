package api

import (
	"encoding/base64"
	"sort"
	"strconv"

	"github.com/bytedance/sonic"
)

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
	return s.pairs()
}

func (s SparseFloatVector) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(s.pairs())
}

func (s SparseFloatVector) pairs() [][]interface{} {
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		left, leftErr := strconv.ParseInt(keys[i], 10, 64)
		right, rightErr := strconv.ParseInt(keys[j], 10, 64)
		if leftErr == nil && rightErr == nil {
			return left < right
		}
		if leftErr == nil {
			return true
		}
		if rightErr == nil {
			return false
		}
		return keys[i] < keys[j]
	})

	pairs := make([][]interface{}, 0, len(keys))
	for _, key := range keys {
		if idx, err := strconv.ParseInt(key, 10, 64); err == nil {
			pairs = append(pairs, []interface{}{idx, s[key]})
			continue
		}
		pairs = append(pairs, []interface{}{key, s[key]})
	}
	return pairs
}
