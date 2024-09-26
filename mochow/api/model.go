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

// model.go - definitions of the request arguments and results data structure model

package api

type CreateDatabaseArgs struct {
	Database string `json:"database"`
}

type ListDatabaseResult struct {
	Databases []string `json:"databases,omitempty"`
}

type CreateTableArgs struct {
	Database           string           `json:"database"`
	Table              string           `json:"table"`
	Description        string           `json:"description"`
	Replication        uint32           `json:"replication"`
	Partition          *PartitionParams `json:"partition,omitempty"`
	EnableDynamicField bool             `json:"enableDynamicField,omitempty"`
	Schema             *TableSchema     `json:"schema,omitempty"`
}

type ListTableArgs struct {
	Database string `json:"database"`
}

type ListTableResult struct {
	Tables []string `json:"tables,omitempty"`
}

type DescTableArgs struct {
	Database string `json:"database"`
	Table    string `json:"table"`
}

type DescTableResult struct {
	Table *TableDescription `json:"table"`
}

type AddFieldArgs struct {
	Database string       `json:"database"`
	Table    string       `json:"table"`
	Schema   *TableSchema `json:"schema,omitempty"`
}

type AliasTableArgs struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	Alias    string `json:"alias"`
}

type UnaliasTableArgs struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	Alias    string `json:"alias"`
}

type ShowTableStatsArgs struct {
	Database string `json:"database"`
	Table    string `json:"table"`
}

type ShowTableStatsResult struct {
	RowCount         uint64 `json:"rowCount"`
	MemorySizeInByte uint64 `json:"memorySizeInByte"`
	DiskSizeInByte   uint64 `json:"diskSizeInByte"`
}

type CreateIndexArgs struct {
	Database string        `json:"database"`
	Table    string        `json:"table"`
	Indexes  []IndexSchema `json:"indexes"`
}

type DescIndexArgs struct {
	Database  string `json:"database"`
	Table     string `json:"table"`
	IndexName string `json:"indexName"`
}

type DescIndexResult struct {
	Index IndexSchema `json:"index"`
}

type ModifyIndexArgs struct {
	Database string      `json:"database"`
	Table    string      `json:"table"`
	Index    IndexSchema `json:"index"`
}

type RebuildIndexArgs struct {
	Database  string `json:"database"`
	Table     string `json:"table"`
	IndexName string `json:"indexName"`
}

type InsertRowArgs struct {
	Database string `json:"database,omitempty"`
	Table    string `json:"table,omitempty"`
	Rows     []Row  `json:"rows,omitempty"`
}

type InsertRowResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type UpsertRowArg InsertRowArgs

type UpsertRowResult InsertRowResult

type DeleteRowArgs struct {
	Database     string                 `json:"database"`
	Table        string                 `json:"table"`
	PrimaryKey   map[string]interface{} `json:"primaryKey,omitempty"`
	PartitionKey map[string]interface{} `json:"partitionKey,omitempty"`
	Filter       string                 `json:"filter,omitempty"`
}

type QueryRowArgs struct {
	Database        string                 `json:"database"`
	Table           string                 `json:"table"`
	PrimaryKey      map[string]interface{} `json:"primaryKey,omitempty"`
	PartitionKey    map[string]interface{} `json:"partitionKey,omitempty"`
	Projections     []string               `json:"projections,omitempty"`
	RetrieveVector  bool                   `json:"retrieveVector,omitempty"`
	ReadConsistency ReadConsistency        `json:"readConsistency,omitempty"`
}

type QueryRowResult struct {
	Row Row `json:"row,omitempty"`
}

type SearchRowArgs struct {
	Database        string                 `json:"database"`
	Table           string                 `json:"table"`
	ANNS            *ANNSearchParams       `json:"anns,omitempty"`
	PartitionKey    map[string]interface{} `json:"partitionKey,omitempty"`
	RetrieveVector  bool                   `json:"retrieveVector,omitempty"`
	Projections     []string               `json:"projections,omitempty"`
	ReadConsistency ReadConsistency        `json:"readConsistency,omitempty"`
}

type RowResult struct {
	Row      Row     `json:"row"`
	Distance float64 `json:"distance"`
}

type SearchRowResult struct {
	SearchVectorFloats []float32   `json:"searchVectorFloats,omitempty"`
	Rows               []RowResult `json:"rows,omitempty"`
}

type UpdateRowArgs struct {
	Database     string                 `json:"database"`
	Table        string                 `json:"table"`
	PrimaryKey   map[string]interface{} `json:"primaryKey,omitempty"`
	PartitionKey map[string]interface{} `json:"partitionKey,omitempty"`
	Update       map[string]interface{} `json:"update,omitempty"`
}

type SelectRowArgs struct {
	Database        string                 `json:"database"`
	Table           string                 `json:"table"`
	Filter          string                 `json:"filter,omitempty"`
	Marker          map[string]interface{} `json:"marker,omitempty"`
	Limit           uint64                 `json:"limit"`
	Projections     []string               `json:"projections,omitempty"`
	ReadConsistency ReadConsistency        `json:"readConsistency,omitempty"`
}

type SelectRowResult struct {
	IsTruncated bool                   `json:"isTruncated"`
	NextMarker  map[string]interface{} `json:"nextMarker,omitempty"`
	Rows        []Row                  `json:"rows,omitempty"`
}

type BatchSearchRowArgs struct {
	Database        string                 `json:"database"`
	Table           string                 `json:"table"`
	ANNS            *BatchANNSearchParams  `json:"anns,omitempty"`
	PartitionKey    map[string]interface{} `json:"partitionKey,omitempty"`
	RetrieveVector  bool                   `json:"retrieveVector,omitempty"`
	ReadConsistency ReadConsistency        `json:"readConsistency,omitempty"`
}

type BatchSearchRowResult struct {
	Results []SearchRowResult `json:"results,omitempty"`
}
