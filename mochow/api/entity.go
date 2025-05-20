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
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

type PartitionParams struct {
	PartitionType PartitionType `json:"partitionType,omitempty"`
	PartitionNum  uint32        `json:"partitionNum"`
}

type FieldSchema struct {
	FieldName     string      `json:"fieldName"`
	FieldType     FieldType   `json:"fieldType"`
	PrimaryKey    bool        `json:"primaryKey"`
	PartitionKey  bool        `json:"partitionKey"`
	AutoIncrement bool        `json:"autoIncrement"`
	NotNull       bool        `json:"notNull"`
	Dimension     uint32      `json:"dimension"`
	ElementType   ElementType `json:"elementType"`
	MaxCapacity   uint32      `json:"maxCapacity"`
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

	if len(f.ElementType) > 0 {
		fields["elementType"] = f.ElementType
	}
	if f.MaxCapacity != 0 {
		fields["maxCapacity"] = f.MaxCapacity
	}

	field, err := sonic.Marshal(fields)
	if err != nil {
		return nil, err
	}
	return field, nil
}

type IndexParams interface{}
type VectorIndexParams map[string]interface{}
type InvertedIndexParams struct {
	params map[string]interface{}
}

func NewInvertedIndexParams() *InvertedIndexParams {
	return &InvertedIndexParams{
		params: make(map[string]interface{}),
	}
}

func (p *InvertedIndexParams) Analyzer(analyzer InvertedIndexAnalyzer) *InvertedIndexParams {
	p.params["analyzer"] = analyzer
	return p
}

func (p *InvertedIndexParams) ParseMode(parseMode InvertedIndexParseMode) *InvertedIndexParams {
	p.params["parseMode"] = parseMode
	return p
}

func (p *InvertedIndexParams) AnalyzerCaseSensitive(caseSensitive bool) *InvertedIndexParams {
	p.params["analyzerCaseSensitive"] = caseSensitive
	return p
}

func (p *InvertedIndexParams) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(p.params)
}

type AutoBuildParams map[string]interface{}

type FilteringIndexField struct {
	Field              string             `json:"field,omitempty"`
	IndexStructureType IndexStructureType `json:"indexStructureType,omitempty"`
}

func (f *FilteringIndexField) FromMapInterface(i interface{}) error {
	var ok bool
	switch v := i.(type) {
	case map[string]interface{}:
		if field, exist := v["field"]; exist {
			if f.Field, ok = field.(string); !ok {
				return fmt.Errorf("field should be string")
			}
		}
		if indexStructureType, exist := v["indexStructureType"]; exist {
			var structureTypeStr string
			if structureTypeStr, ok = indexStructureType.(string); !ok {
				return fmt.Errorf("invalid indexStructureType enum value")
			}
			f.IndexStructureType = IndexStructureType(structureTypeStr)
		}
	}
	return nil
}

type IndexSchema struct {
	IndexName                    string
	IndexType                    IndexType
	MetricType                   MetricType
	Params                       IndexParams
	Field                        string
	InvertedIndexFields          []string                      // for inverted index
	InvertedIndexFieldAttributes []InvertedIndexFieldAttribute // for inverted index
	FilterIndexFields            []FilteringIndexField         // for filtering index
	State                        IndexState
	AutoBuild                    bool
	AutoBuildPolicy              AutoBuildParams
}

func (index IndexSchema) MarshalJSON() ([]byte, error) {
	params := make(map[string]interface{})
	// index name
	if len(index.IndexName) > 0 {
		params["indexName"] = index.IndexName
	}
	// index type
	params["indexType"] = index.IndexType
	// metric type
	if len(index.MetricType) > 0 {
		params["metricType"] = index.MetricType
	}
	// params
	if index.Params != nil {
		params["params"] = index.Params
	}
	// auto build
	if index.AutoBuild {
		params["autoBuild"] = index.AutoBuild
		params["autoBuildPolicy"] = index.AutoBuildPolicy
	}

	params["state"] = index.State

	// vector index and secondary index field
	if len(index.Field) > 0 {
		params["field"] = index.Field
	}

	// handle conflicted index fields
	switch index.IndexType {
	case InvertedIndex:
		// inverted index fields
		if len(index.InvertedIndexFields) > 0 {
			params["fields"] = index.InvertedIndexFields
		}
		// inverted index field attributes
		if len(index.InvertedIndexFieldAttributes) > 0 {
			params["fieldsIndexAttributes"] = index.InvertedIndexFieldAttributes
		}
	case FilteringIndex:
		// filtering index fields
		if len(index.FilterIndexFields) > 0 {
			params["fields"] = index.FilterIndexFields
		}
	}
	return sonic.Marshal(&params)
}

func (index *IndexSchema) UnmarshalJSON(data []byte) error {
	ds := decoder.NewStreamDecoder(bytes.NewReader(data))
	ds.UseNumber()
	params := make(map[string]interface{})
	err := ds.Decode(&params)
	if err != nil {
		return err
	}

	// index name
	var ok bool
	if indexName, exist := params["indexName"]; exist {
		if index.IndexName, ok = indexName.(string); !ok {
			return fmt.Errorf("indexName should be string")
		}
	}
	// index type
	if indexType, exist := params["indexType"]; exist {
		var indexTypeStr string
		if indexTypeStr, ok = indexType.(string); !ok {
			return fmt.Errorf("indexType should be string")
		}
		index.IndexType = IndexType(indexTypeStr)
	}
	// metric type
	if metricType, exist := params["metricType"]; exist {
		var metricTypeStr string
		if metricTypeStr, ok = metricType.(string); !ok {
			return fmt.Errorf("metricType should be string")
		}
		index.MetricType = MetricType(metricTypeStr)
	}
	// index params
	if indexParams, exist := params["params"]; exist {
		index.Params = indexParams
	}
	// field
	if indexField, exist := params["field"]; exist {
		if index.Field, ok = indexField.(string); !ok {
			return fmt.Errorf("field should be string")
		}
	}
	// auto build
	if autoBuild, exist := params["autoBuild"]; exist {
		if index.AutoBuild, ok = autoBuild.(bool); !ok {
			return fmt.Errorf("autoBuild should be bool")
		}
	}
	if autoBuildPolicy, exist := params["autoBuildPolicy"]; exist {
		if index.AutoBuildPolicy, ok = autoBuildPolicy.(map[string]interface{}); !ok {
			return fmt.Errorf("invalid autoBuildPolicy")
		}
	}
	// index state
	if state, exist := params["state"]; exist {
		index.State = IndexState(state.(string))
	}
	switch index.IndexType {
	case InvertedIndex:
		// inverted index fields
		if fields, ok := params["fields"]; ok {
			switch v := fields.(type) {
			case []interface{}:
				index.InvertedIndexFields = make([]string, len(v))
				for i := 0; i < len(v); i++ {
					index.InvertedIndexFields[i] = v[i].(string)
				}
			}
		}
		if attributes, ok := params["fieldsIndexAttributes"]; ok {
			switch v := attributes.(type) {
			case []interface{}:
				index.InvertedIndexFieldAttributes = make([]InvertedIndexFieldAttribute, len(v))
				for i := 0; i < len(v); i++ {
					index.InvertedIndexFieldAttributes[i] = InvertedIndexFieldAttribute(v[i].(string))
				}
			}
		}
	case FilteringIndex:
		// fields for filtering index
		if fields, ok := params["fields"]; ok {
			switch t := fields.(type) {
			case []interface{}:
				index.FilterIndexFields = make([]FilteringIndexField, len(t))
				for i := 0; i < len(t); i++ {
					var filteringIndexField FilteringIndexField
					if err := filteringIndexField.FromMapInterface(t[i]); err != nil {
						return err
					}
					index.FilterIndexFields[i] = filteringIndexField
				}
			}
		}
	}
	return nil
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
	Params       *SearchParams `json:"params,omitempty"`
	Filter       string        `json:"filter,omitempty"`
}

type BatchANNSearchParams struct {
	VectorField  string        `json:"vectorField,omitempty"`
	VectorFloats [][]float32   `json:"vectorFloats,omitempty"`
	Params       *SearchParams `json:"params,omitempty"`
	Filter       string        `json:"filter,omitempty"`
}

type searchRequest interface {
	requestType() string
	toDict() map[string]interface{}
	isBatch() bool
}

/*
Optional configurable params for vector search.

For each index algorithm, the params that could be set are:

IndexType: HNSW
Params: ef, pruning

IndexType: HNSWPQ
Params: ef, pruning

IndexType: PUCK
Params: searchCoarseCount

IndexType: FLAT
Params: None
*/
type VectorSearchConfig struct {
	params map[string]interface{}
}

func (h VectorSearchConfig) New() *VectorSearchConfig {
	return &VectorSearchConfig{
		params: make(map[string]interface{}),
	}
}

func (h *VectorSearchConfig) Ef(ef uint32) *VectorSearchConfig {
	h.params["ef"] = ef
	return h
}

func (h *VectorSearchConfig) Pruning(pruning bool) *VectorSearchConfig {
	h.params["pruning"] = pruning
	return h
}

func (h *VectorSearchConfig) SearchCoarseCount(searchCoarseCount uint32) *VectorSearchConfig {
	h.params["searchCoarseCount"] = searchCoarseCount
	return h
}

func (h *VectorSearchConfig) FilterMode(filterMode FilterMode) *VectorSearchConfig {
	h.params["filterMode"] = filterMode
	return h
}

func (h *VectorSearchConfig) PostFilterAmplicationFactor(amplicationFactor float32) *VectorSearchConfig {
	h.params["postFilterAmplificationFactor"] = amplicationFactor
	return h
}

type vectorSearchRequest interface {
	searchRequest
	vectorSearchRequestDummyInterface()
	SetIteratedIds(string)
}

type request struct {
	set map[string]bool
}

func (r request) mark(key string) {
	r.set[key] = true
}

func (r request) isMarked(key string) bool {
	_, ok := r.set[key]
	return ok
}

type searchCommonFields struct {
	request
	partitionKey    map[string]interface{}
	projections     []string
	readConsistency ReadConsistency
	limit           uint32
	filter          string
}

func searchCommonFieldsToMap(r *searchCommonFields) map[string]interface{} {
	fields := make(map[string]interface{})
	if r.isMarked("partitionKey") {
		fields["partitionKey"] = r.partitionKey
	}
	if r.isMarked("projections") {
		fields["projections"] = r.projections
	}
	if r.isMarked("readConsistency") {
		fields["readConsistency"] = r.readConsistency
	}
	if r.isMarked("filter") {
		fields["filter"] = r.filter
	}
	if r.isMarked("limit") {
		fields["limit"] = r.limit
	}
	return fields
}

type vectorSearchFields struct {
	searchCommonFields
	vectorField  string
	vector       Vector
	vectors      []Vector
	distanceNear float64
	distanceFar  float64
	config       *VectorSearchConfig
	options      *AdvancedOptions
}

func (r vectorSearchFields) fillSearchFields(fields *map[string]interface{}) {
	anns := make(map[string]interface{})
	if r.isMarked("vectorField") {
		anns["vectorField"] = r.vectorField
	}
	if r.isMarked("vector") {
		anns[r.vector.name()] = r.vector.representation()
	}
	if r.isMarked("vectors") && len(r.vectors) != 0 {
		vectors := make([]interface{}, 0, len(r.vectors))
		for _, vec := range r.vectors {
			vectors = append(vectors, vec.representation())
		}
		anns[r.vectors[0].name()] = vectors
	}
	if r.isMarked("filter") {
		anns["filter"] = r.filter
	}

	params := make(map[string]interface{})
	if r.isMarked("config") {
		for k, v := range r.config.params {
			params[k] = v
		}
	}
	if r.isMarked("distanceNear") {
		params["distanceNear"] = r.distanceNear
	}
	if r.isMarked("distanceFar") {
		params["distanceFar"] = r.distanceFar
	}
	if r.isMarked("limit") {
		params["limit"] = r.limit
	}
	if len(params) != 0 {
		anns["params"] = params
	}
	if len(anns) != 0 {
		(*fields)["anns"] = anns
	}
	if r.isMarked("options") {
		(*fields)["advancedOptions"] = r.options
	}

	for k, v := range searchCommonFieldsToMap(&r.searchCommonFields) {
		if k == "filter" || k == "limit" { // in "anns"
			continue
		}
		(*fields)[k] = v
	}
}

type fusionRankPolicy interface {
	fusionRankPolicyDummyInterface()
	Params() map[string]interface{}
}

type RRFRank struct {
	fusionRankPolicy
	k int64
}

func (rank RRFRank) New(k int64) *RRFRank {
	rank.k = k
	return &rank
}

func (rank *RRFRank) Params() map[string]interface{} {
	return map[string]interface{}{
		"strategy": "rrf",
		"params": map[string]interface{}{
			"k": rank.k,
		},
	}
}

func (rank *RRFRank) fusionRankPolicyDummyInterface() {}

type WeightedRank struct {
	fusionRankPolicy
	weights []float64
}

func (rank WeightedRank) New(weights []float64) *WeightedRank {
	rank.weights = weights
	return &rank
}

func (rank *WeightedRank) Params() map[string]interface{} {
	return map[string]interface{}{
		"strategy": "ws",
		"params": map[string]interface{}{
			"weights": rank.weights,
		},
	}
}

func (rank *WeightedRank) fusionRankPolicyDummyInterface() {}

/**** Vector Topk Search ****/

type VectorTopkSearchRequest struct {
	vectorSearchRequest // interface
	vectorSearchFields  // common fields
	iteratedIds         string
}

func (r VectorTopkSearchRequest) New(vectorField string, vector Vector, limit uint32) *VectorTopkSearchRequest {
	r.set = make(map[string]bool, 0)

	r.mark("vectorField")
	r.vectorField = vectorField

	r.mark("vector")
	r.vector = vector

	r.mark("limit")
	r.limit = limit
	return &r
}

func (r *VectorTopkSearchRequest) String() string {
	return fmt.Sprintf("VectorTopkSearchRequest:%v", r.toDict())
}

func (r *VectorTopkSearchRequest) PartitionKey(partitionKey map[string]interface{}) *VectorTopkSearchRequest {
	r.mark("partitionKey")
	r.partitionKey = partitionKey
	return r
}

func (r *VectorTopkSearchRequest) ReadConsistency(readConsistency ReadConsistency) *VectorTopkSearchRequest {
	r.mark("readConsistency")
	r.readConsistency = readConsistency
	return r
}

func (r *VectorTopkSearchRequest) Projections(projections []string) *VectorTopkSearchRequest {
	r.mark("projections")
	r.projections = projections
	return r
}

func (r *VectorTopkSearchRequest) Filter(filter string) *VectorTopkSearchRequest {
	r.mark("filter")
	r.filter = filter
	return r
}

func (r *VectorTopkSearchRequest) Config(config *VectorSearchConfig) *VectorTopkSearchRequest {
	r.mark("config")
	r.config = config
	return r
}

func (r *VectorTopkSearchRequest) AdvancedOptions(options *AdvancedOptions) *VectorTopkSearchRequest {
	r.mark("advancedOptions")
	r.options = options
	return r
}

func (r *VectorTopkSearchRequest) requestType() string {
	return "search"
}

func (r *VectorTopkSearchRequest) isBatch() bool {
	return false
}

func (r *VectorTopkSearchRequest) toDict() map[string]interface{} {
	fields := make(map[string]interface{})
	r.fillSearchFields(&fields)
	return fields
}

func (r *VectorTopkSearchRequest) vectorSearchRequestDummyInterface() {
}

func (r *VectorTopkSearchRequest) SetIteratedIds(ids string) {
	r.iteratedIds = ids
}

/**** Vector Range Search ****/
type DistanceRange struct {
	Min, Max float64
}

type VectorRangeSearchRequest struct {
	vectorSearchRequest // interface
	vectorSearchFields  // common fields
}

func (r VectorRangeSearchRequest) New(vectorField string, vector Vector, distanceRange DistanceRange) *VectorRangeSearchRequest {
	r.set = make(map[string]bool, 0)

	r.mark("vectorField")
	r.vectorField = vectorField

	r.mark("vector")
	r.vector = vector

	r.mark("distanceNear")
	r.distanceNear = distanceRange.Min

	r.mark("distanceFar")
	r.distanceFar = distanceRange.Max
	return &r
}

func (r *VectorRangeSearchRequest) String() string {
	return fmt.Sprintf("VectorRangeSearchRequest:%v", r.toDict())
}

func (r *VectorRangeSearchRequest) PartitionKey(partitionKey map[string]interface{}) *VectorRangeSearchRequest {
	r.mark("partitionKey")
	r.partitionKey = partitionKey
	return r
}

func (r *VectorRangeSearchRequest) ReadConsistency(readConsistency ReadConsistency) *VectorRangeSearchRequest {
	r.mark("readConsistency")
	r.readConsistency = readConsistency
	return r
}

func (r *VectorRangeSearchRequest) Projections(projections []string) *VectorRangeSearchRequest {
	r.mark("projections")
	r.projections = projections
	return r
}

func (r *VectorRangeSearchRequest) Limit(limit uint32) *VectorRangeSearchRequest {
	r.mark("limit")
	r.limit = limit
	return r
}

func (r *VectorRangeSearchRequest) Filter(filter string) *VectorRangeSearchRequest {
	r.mark("filter")
	r.filter = filter
	return r
}

func (r *VectorRangeSearchRequest) Config(config *VectorSearchConfig) *VectorRangeSearchRequest {
	r.mark("config")
	r.config = config
	return r
}

func (r *VectorRangeSearchRequest) requestType() string {
	return "search"
}

func (r *VectorRangeSearchRequest) isBatch() bool {
	return false
}

func (r *VectorRangeSearchRequest) toDict() map[string]interface{} {
	fields := make(map[string]interface{})
	r.fillSearchFields(&fields)
	return fields
}

func (r *VectorRangeSearchRequest) vectorSearchRequestDummyInterface() {
}

/**** Vector Batch Search ****/
type VectorBatchSearchRequest struct {
	vectorSearchRequest // interface
	vectorSearchFields  // common fields
}

func (r VectorBatchSearchRequest) New(vectorField string, vectors []Vector) *VectorBatchSearchRequest {
	r.set = make(map[string]bool, 0)

	r.mark("vectorField")
	r.vectorField = vectorField

	r.mark("vectors")
	r.vectors = vectors
	return &r
}

func (r *VectorBatchSearchRequest) String() string {
	return fmt.Sprintf("VectorBatchSearchRequest:%v", r.toDict())
}

func (r *VectorBatchSearchRequest) PartitionKey(partitionKey map[string]interface{}) *VectorBatchSearchRequest {
	r.mark("partitionKey")
	r.partitionKey = partitionKey
	return r
}

func (r *VectorBatchSearchRequest) ReadConsistency(readConsistency ReadConsistency) *VectorBatchSearchRequest {
	r.mark("readConsistency")
	r.readConsistency = readConsistency
	return r
}

func (r *VectorBatchSearchRequest) Projections(projections []string) *VectorBatchSearchRequest {
	r.mark("projections")
	r.projections = projections
	return r
}

func (r *VectorBatchSearchRequest) Limit(limit uint32) *VectorBatchSearchRequest {
	r.mark("limit")
	r.limit = limit
	return r
}

func (r *VectorBatchSearchRequest) DistanceRange(distanceRange DistanceRange) *VectorBatchSearchRequest {
	r.mark("distanceNear")
	r.distanceNear = distanceRange.Min

	r.mark("distanceFar")
	r.distanceFar = distanceRange.Max
	return r
}

func (r *VectorBatchSearchRequest) Filter(filter string) *VectorBatchSearchRequest {
	r.mark("filter")
	r.filter = filter
	return r
}

func (r *VectorBatchSearchRequest) Config(config *VectorSearchConfig) *VectorBatchSearchRequest {
	r.mark("config")
	r.config = config
	return r
}

func (r *VectorBatchSearchRequest) isBatch() bool {
	return true
}

func (r *VectorBatchSearchRequest) requestType() string {
	return "batchSearch"
}

func (r *VectorBatchSearchRequest) toDict() map[string]interface{} {
	fields := make(map[string]interface{})
	r.fillSearchFields(&fields)
	return fields
}

func (r *VectorBatchSearchRequest) vectorSearchRequestDummyInterface() {
}

/**** BM25 Search ****/
type bm25SearchRequest interface {
	searchRequest

	// Make sure user not pass e.g. 'VectorSearchRequest' to BM25Search api
	bm25SearchRequestDummyInterface()
}

type BM25SearchRequest struct {
	bm25SearchRequest  // interface
	searchCommonFields // common fields
	indexName          string
	searchText         string
}

func (r BM25SearchRequest) New(indexName string, searchText string) *BM25SearchRequest {
	r.set = make(map[string]bool, 0)

	r.indexName = indexName
	r.searchText = searchText
	return &r
}

func (r *BM25SearchRequest) String() string {
	return fmt.Sprintf("BM25SearchRequest:%v", r.toDict())
}

func (r *BM25SearchRequest) PartitionKey(partitionKey map[string]interface{}) *BM25SearchRequest {
	r.mark("partitionKey")
	r.partitionKey = partitionKey
	return r
}

func (r *BM25SearchRequest) ReadConsistency(readConsistency ReadConsistency) *BM25SearchRequest {
	r.mark("readConsistency")
	r.readConsistency = readConsistency
	return r
}

func (r *BM25SearchRequest) Projections(projections []string) *BM25SearchRequest {
	r.mark("projections")
	r.projections = projections
	return r
}

func (r *BM25SearchRequest) Limit(limit uint32) *BM25SearchRequest {
	r.mark("limit")
	r.limit = limit
	return r
}

func (r *BM25SearchRequest) Filter(filter string) *BM25SearchRequest {
	r.mark("filter")
	r.filter = filter
	return r
}

func (r *BM25SearchRequest) toDict() map[string]interface{} {
	fields := make(map[string]interface{})
	for k, v := range searchCommonFieldsToMap(&r.searchCommonFields) {
		fields[k] = v
	}

	bm25Params := make(map[string]interface{})
	bm25Params["indexName"] = r.indexName
	bm25Params["searchText"] = r.searchText
	fields["BM25SearchParams"] = bm25Params

	return fields
}

func (r *BM25SearchRequest) requestType() string {
	return "search"
}

func (r *BM25SearchRequest) isBatch() bool {
	return false
}

func (r *BM25SearchRequest) bm25SearchRequestDummyInterface() {
}

/**** Hybrid Search ****/

type hybridSearchRequest interface {
	searchRequest

	// Make sure user not pass e.g. 'VectorSearchRequest' to HybridSearch api
	hybridSearchRequestDummyInterface()
}

type HybridSearchRequest struct {
	hybridSearchRequest // interface
	searchCommonFields  // common fields

	vectorRequest vectorSearchRequest
	bm25Request   bm25SearchRequest
	vectorWeight  float32
	bm25Weight    float32
}

/*
Note: 'limit' and 'filter' are global settings, and they will
apply to both vector search and BM25 search. Avoid setting them in
'bm25Request' or 'vectorRequest'.  Any settings in 'vectorRequest'
or 'bm25Request' for 'limit' or 'filter' will be overridden by the
general settings.
*/
func (r HybridSearchRequest) New(
	vectorRequest vectorSearchRequest,
	bm25Request bm25SearchRequest,
	vectorWeight float32,
	bm25Weight float32,
) *HybridSearchRequest {
	r.set = make(map[string]bool, 0)

	r.vectorRequest = vectorRequest
	r.bm25Request = bm25Request
	r.vectorWeight = vectorWeight
	r.bm25Weight = bm25Weight
	return &r
}

func (r *HybridSearchRequest) String() string {
	return fmt.Sprintf("HybridSearchRequest:%v", r.toDict())
}

func (r *HybridSearchRequest) PartitionKey(partitionKey map[string]interface{}) *HybridSearchRequest {
	r.mark("partitionKey")
	r.partitionKey = partitionKey
	return r
}

func (r *HybridSearchRequest) ReadConsistency(readConsistency ReadConsistency) *HybridSearchRequest {
	r.mark("readConsistency")
	r.readConsistency = readConsistency
	return r
}

func (r *HybridSearchRequest) Projections(projections []string) *HybridSearchRequest {
	r.mark("projections")
	r.projections = projections
	return r
}

func (r *HybridSearchRequest) Limit(limit uint32) *HybridSearchRequest {
	r.mark("limit")
	r.limit = limit
	return r
}

func (r *HybridSearchRequest) Filter(filter string) *HybridSearchRequest {
	r.mark("filter")
	r.filter = filter
	return r
}

func (r *HybridSearchRequest) toDict() map[string]interface{} {
	fields := make(map[string]interface{})

	for k, v := range r.bm25Request.toDict() {
		fields[k] = v
	}
	for k, v := range r.vectorRequest.toDict() {
		fields[k] = v
	}

	for k, v := range searchCommonFieldsToMap(&r.searchCommonFields) {
		fields[k] = v
	}

	_, ok := fields["anns"]
	if ok {
		fields["anns"].(map[string]interface{})["weight"] = r.vectorWeight
	}

	_, ok = fields["BM25SearchParams"]
	if ok {
		fields["BM25SearchParams"].(map[string]interface{})["weight"] = r.bm25Weight
	}

	return fields
}

func (r *HybridSearchRequest) isBatch() bool {
	return false
}

func (r *HybridSearchRequest) requestType() string {
	return "search"
}

func (r *HybridSearchRequest) hybridSearchRequestDummyInterface() {
}

/**** MultiVector Search ****/
type MultivectorSearchRequest struct {
	vectorSearchRequest
	searchCommonFields // common fields
	ranking            fusionRankPolicy
	vectorRequests     []vectorSearchRequest
	iteratedIds        string
}

/*
Note: 'limit' and 'filter' are global settings. Avoid setting them in
'bm25Request' or 'vectorRequest'.  Any settings in 'vectorRequests'
for 'limit' or 'filter' will be overridden by the
general settings.
*/
func (r MultivectorSearchRequest) New() *MultivectorSearchRequest {
	r.set = make(map[string]bool, 0)
	return &r
}

func (r *MultivectorSearchRequest) AddSingleVectorSearchRequest(vectorRequest vectorSearchRequest) *MultivectorSearchRequest {
	r.vectorRequests = append(r.vectorRequests, vectorRequest)
	return r
}

func (r *MultivectorSearchRequest) String() string {
	return fmt.Sprintf("MultivectorSearchRequest:%v", r.toDict())
}

func (r *MultivectorSearchRequest) PartitionKey(partitionKey map[string]interface{}) *MultivectorSearchRequest {
	r.mark("partitionKey")
	r.partitionKey = partitionKey
	return r
}

func (r *MultivectorSearchRequest) ReadConsistency(readConsistency ReadConsistency) *MultivectorSearchRequest {
	r.mark("readConsistency")
	r.readConsistency = readConsistency
	return r
}

func (r *MultivectorSearchRequest) Projections(projections []string) *MultivectorSearchRequest {
	r.mark("projections")
	r.projections = projections
	return r
}

func (r *MultivectorSearchRequest) Limit(limit uint32) *MultivectorSearchRequest {
	r.mark("limit")
	r.limit = limit
	return r
}

func (r *MultivectorSearchRequest) Filter(filter string) *MultivectorSearchRequest {
	r.mark("filter")
	r.filter = filter
	return r
}

func (r *MultivectorSearchRequest) Rank(rank fusionRankPolicy) *MultivectorSearchRequest {
	r.ranking = rank
	return r
}

func (r *MultivectorSearchRequest) toDict() map[string]interface{} {
	fields := make(map[string]interface{})
	searchList := []interface{}{}
	for _, singleVectorRequest := range r.vectorRequests {
		singleVectorRequestParam := singleVectorRequest.toDict()
		searchList = append(searchList, singleVectorRequestParam["anns"])
	}
	fields["search"] = searchList
	for k, v := range searchCommonFieldsToMap(&r.searchCommonFields) {
		fields[k] = v
	}
	if r.ranking != nil {
		fields["ranking"] = r.ranking.Params()
	}
	return fields
}

func (r *MultivectorSearchRequest) isBatch() bool {
	return false
}

func (r *MultivectorSearchRequest) requestType() string {
	return "multiVectorSearch"
}

func (r *MultivectorSearchRequest) multiVectorSearchRequestDummyInterface() {
}

func (r *MultivectorSearchRequest) SetIteratedIds(ids string) {
	r.iteratedIds = ids
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

type AdvancedOptions struct {
	options map[string]interface{} `json:"-"`
}

func NewAdvancedOptions() *AdvancedOptions {
	return &AdvancedOptions{
		options: make(map[string]interface{}),
	}
}

func (a *AdvancedOptions) AcceptPartialSuccessOnMPP(acceptPartialSuccessOnMPP bool) *AdvancedOptions {
	a.options["acceptPartialSuccessOnMPP"] = acceptPartialSuccessOnMPP
	return a
}

func (a *AdvancedOptions) SuccessRateLowerBoundOnMPP(successRateLowerBoundOnMPP float32) *AdvancedOptions {
	a.options["successRateLowerBoundOnMPP"] = successRateLowerBoundOnMPP
	return a
}

func (a *AdvancedOptions) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(a.options)
}
