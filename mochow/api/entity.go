/*
 * Copyright 2024 Baidu, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions
 * and limitations under the License.
 */

// entity.go - definitions of entity in Mochow service

package api

import (
	"bytes"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

type PartitionParams struct {
	PartitionType PartitionType `json:"partitionType,omitempty"`
	PartitionNum  uint32        `json:"partitionNum"`
}

type FieldSchema struct {
	FieldName     string    `json:"fieldName"`
	FieldType     FieldType `json:"fieldType"`
	PrimaryKey    bool      `json:"primaryKey"`
	PartitionKey  bool      `json:"partitionKey"`
	AutoIncrement bool      `json:"autoIncrement"`
	NotNull       bool      `json:"notNull"`
	Dimension     uint32    `json:"dimension"`
}

func (f *FieldSchema) MarshalJSON() ([]byte, error) {
	fields := make(map[string]interface{})
	if len(f.FieldName) > 0 {
		fields["fieldName"] = f.FieldName
	}
	if len(f.FieldType) > 0 {
		fields["fieldType"] = f.FieldType
	}
	if f.Dimension > 0 {
		fields["dimension"] = f.Dimension
	}
	fields["primaryKey"] = f.PrimaryKey
	fields["partitionKey"] = f.PartitionKey
	fields["autoIncrement"] = f.AutoIncrement
	fields["notNull"] = f.NotNull
	field, err := sonic.Marshal(fields)
	if err != nil {
		return nil, err
	}
	return field, nil
}

type VectorIndexParams map[string]interface{}

type AutoBuildParams map[string]interface{}

type IndexSchema struct {
	IndexName       string            `json:"indexName,omitempty"`
	IndexType       IndexType         `json:"indexType,omitempty"`
	MetricType      MetricType        `json:"metricType,omitempty"`
	Params          VectorIndexParams `json:"params,omitempty"`
	Field           string            `json:"field,omitempty"`
	State           IndexState        `json:"state,omitempty"`
	AutoBuild       bool              `json:"autoBuild,omitempty"`
	AutoBuildPolicy AutoBuildParams   `json:"autoBuildPolicy,omitempty"`
}

type TableSchema struct {
	Fields  []FieldSchema `json:"fields,omitempty"`
	Indexes []IndexSchema `json:"indexes,omitempty"`
}

type TableDescription struct {
	Database           string           `json:"database"`
	Table              string           `json:"table"`
	CreateTime         string           `json:"createTime"`
	Description        string           `json:"description"`
	Replication        uint32           `json:"replication"`
	Partition          *PartitionParams `json:"partition,omitempty"`
	EnableDynamicField bool             `json:"enableDynamicField"`
	State              TableState       `json:"state"`
	Aliases            []string         `json:"aliases,omitempty"`
	Schema             *TableSchema     `json:"schema,omitempty"`
}

type Row struct {
	Fields map[string]interface{} `json:"-"`
}

func (d *Row) MarshalJSON() ([]byte, error) {
	field, err := sonic.Marshal(d.Fields)
	if err != nil {
		return nil, err
	}
	return field, nil
}

func (d *Row) UnmarshalJSON(data []byte) error {
	ds := decoder.NewStreamDecoder(bytes.NewReader(data))
	ds.UseNumber()
	err := ds.Decode(&d.Fields)
	if err != nil {
		return err
	}
	return nil
}

type SearchParams struct {
	Params map[string]interface{}
}

func NewSearchParams() *SearchParams {
	return &SearchParams{
		Params: make(map[string]interface{}),
	}
}

func (h *SearchParams) AddEf(ef uint32) {
	h.Params["ef"] = ef
}

func (h *SearchParams) AddDistanceNear(distanceNear float64) {
	h.Params["distanceNear"] = distanceNear
}

func (h *SearchParams) AddDistanceFar(distanceFar float64) {
	h.Params["distanceFar"] = distanceFar
}

func (h *SearchParams) AddLimit(limit uint32) {
	h.Params["limit"] = limit
}

func (h *SearchParams) AddPruning(pruning bool) {
	h.Params["pruning"] = pruning
}

func (h *SearchParams) AddSearchCoarseCount(searchCoarseCount uint32) {
	h.Params["searchCoarseCount"] = searchCoarseCount
}

func (h *SearchParams) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(h.Params)
}

type ANNSearchParams struct {
	VectorField  string        `json:"vectorField,omitempty"`
	VectorFloats []float32     `json:"vectorFloats,omitempty"`
	Params       *SearchParams `json:"params,omitempty'"`
	Filter       string        `json:"filter,omitempty"`
}

type BatchANNSearchParams struct {
	VectorField  string        `json:"vectorField,omitempty"`
	VectorFloats [][]float32   `json:"vectorFloats,omitempty"`
	Params       *SearchParams `json:"params,omitempty'"`
	Filter       string        `json:"filter,omitempty"`
}

type AutoBuildPolicy interface {
	Params() map[string]interface{}
	AddTiming(timing string)
	AddPeriod(period uint64)
	AddRowCountIncrement(increment uint64)
	AddRowCountIncrementRatio(ratio float64)
}

type baseAutoBuildPolicy struct {
	params map[string]interface{}
}

func newBaseAutoBuildPolicy(policyType AutoBuildPolicyType) baseAutoBuildPolicy {
	return baseAutoBuildPolicy{
		params: map[string]interface{}{
			"policyType": policyType,
		},
	}
}

func (bp *baseAutoBuildPolicy) Params() map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range bp.params {
		params[k] = v
	}
	return params
}

func (bp *baseAutoBuildPolicy) AddTiming(timing string) {
	bp.params["timing"] = timing
}

func (bp *baseAutoBuildPolicy) AddPeriod(period uint64) {
	bp.params["periodInSecond"] = period
}

func (bp *baseAutoBuildPolicy) AddRowCountIncrement(increment uint64) {
	bp.params["rowCountIncrement"] = increment
}

func (bp *baseAutoBuildPolicy) AddRowCountIncrementRatio(ratio float64) {
	bp.params["rowCountIncrementRatio"] = ratio
}

type AutoBuildTimingPolicy struct {
	baseAutoBuildPolicy
}

func NewAutoBuildTimingPolicy() *AutoBuildTimingPolicy {
	return &AutoBuildTimingPolicy{
		baseAutoBuildPolicy: newBaseAutoBuildPolicy(AutoBuildPolicyTiming),
	}
}

type AutoBuildPeriodicalPolicy struct {
	baseAutoBuildPolicy
}

func NewAutoBuildPeriodicalPolicy() *AutoBuildPeriodicalPolicy {
	return &AutoBuildPeriodicalPolicy{
		baseAutoBuildPolicy: newBaseAutoBuildPolicy(AutoBuildPolicyPeriodical),
	}
}

type AutoBuildIncrementPolicy struct {
	baseAutoBuildPolicy
}

func NewAutoBuildIncrementPolicy() *AutoBuildIncrementPolicy {
	return &AutoBuildIncrementPolicy{
		baseAutoBuildPolicy: newBaseAutoBuildPolicy(AutoBuildPolicyIncrement),
	}
}
