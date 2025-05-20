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

import "github.com/bytedance/sonic"

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

type QueryKey struct {
	PrimaryKey   map[string]interface{} `json:"primaryKey,omitempty"`
	PartitionKey map[string]interface{} `json:"partitionKey,omitempty"`
}

type BatchQueryRowArgs struct {
	Database        string          `json:"database"`
	Table           string          `json:"table"`
	Keys            []QueryKey      `json:"keys,omitempty"`
	Projections     []string        `json:"projections,omitempty"`
	RetrieveVector  bool            `json:"retrieveVector,omitempty"`
	ReadConsistency ReadConsistency `json:"readConsistency,omitempty"`
}

type BatchQueryRowResult struct {
	Row []Row `json:"rows,omitempty"`
}

// Deprecated
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
	Score    float64 `json:"score"`
}

type SearchRowResult struct {
	SearchVectorFloats []float32   `json:"searchVectorFloats,omitempty"`
	Rows               []RowResult `json:"rows,omitempty"`
	IteratedIds        string      `json:"iteratedIds,omitempty"`
}

// vector topk search, range search and batch search
type VectorSearchArgs struct {
	Database string
	Table    string
	Request  vectorSearchRequest
}

// BM25 search
type BM25SearchArgs struct {
	Database string
	Table    string
	Request  bm25SearchRequest
}

// hybrid search (vector + BM25)
type HybridSearchArgs struct {
	Database string
	Table    string
	Request  hybridSearchRequest
}

// multi vector search (vector + BM25)
type MultivectorSearchArgs struct {
	Database string
	Table    string
	Request  vectorSearchRequest
}

// search iterator
type SearchIteratorArgs struct {
	Database        string
	Table           string
	Request         vectorSearchRequest
	BatchSize       uint32
	TotalSize       uint32
	PartitionKey    map[string]interface{}
	Projections     []string
	ReadConsistency ReadConsistency
	Config          map[string]interface{}
}

type SearchResult struct {
	IsBatch   bool
	Rows      *SearchRowResult      // for single search
	BatchRows *BatchSearchRowResult // for batch search
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

///////////////////////////////  RBAC  ///////////////////////////////

type PrivilegeTuple struct {
	Database   string   `json:"database,omitempty"`
	Table      string   `json:"table,omitempty"`
	Privileges []string `json:"privileges,omitempty"`
}

func (p *PrivilegeTuple) UnmarshalJSON(data []byte) error {

	tmp := struct {
		Database   string   `json:"database,omitempty"`
		Table      string   `json:"table,omitempty"`
		Privileges []string `json:"privileges,omitempty"`
		Privilege  []string `json:"privilege,omitempty"`
	}{}
	if err := sonic.Unmarshal(data, &tmp); err != nil {
		return err
	}
	p.Database = tmp.Database
	p.Table = tmp.Table
	p.Privileges = tmp.Privileges
	if len(p.Privileges) == 0 && len(tmp.Privilege) > 0 {
		p.Privileges = tmp.Privilege
	}
	return nil
}

///////////////////////////////  role  ///////////////////////////////

type CreateRoleArgs struct {
	Role string `json:"role,omitempty"`
}

type DropRoleArgs struct {
	Role string `json:"role,omitempty"`
}

type GrantRolePrivilegesArgs struct {
	Role            string           `json:"role,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type RevokeRolePrivilegesArgs struct {
	Role            string           `json:"role,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type ShowRolePrivilegesArgs struct {
	Role string `json:"role,omitempty"`
}

type ShowRolePrivilegesResult struct {
	Users           []string         `json:"users,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type SelectRoleArgs struct {
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type SelectRoleResult struct {
	Roles []string `json:"Roles"`
}

///////////////////////////////  user  ///////////////////////////////

type CreateUserArgs struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type DropUserArgs struct {
	Username string `json:"username,omitempty"`
}

type ChangeUserPasswordArgs struct {
	Username    string `json:"username,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

type GrantUserRolesArgs struct {
	Username string   `json:"username,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

type RevokeUserRolesArgs struct {
	Username string   `json:"username,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

type GrantUserPrivilegesArgs struct {
	Username        string           `json:"username,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type RevokeUserPrivilegesArgs struct {
	Username        string           `json:"username,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type ShowUserPrivilegesArgs struct {
	Username string `json:"username,omitempty"`
}

type ShowUserPrivilegesResult struct {
	Roles []struct {
		Role            string           `json:"role,omitempty"`
		PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
	} `json:"roles,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type SelectUserArgs struct {
	Roles           []string         `json:"roles,omitempty"`
	PrivilegeTuples []PrivilegeTuple `json:"privilegeTuples,omitempty"`
}

type SelectUserResult struct {
	Users []string `json:"users,omitempty"`
}
