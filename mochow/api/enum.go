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

// enum.go - definitions of enumeration in Mochow service

package api

type MetricType string

const (
	L2     MetricType = "L2"
	IP     MetricType = "IP"
	COSINE MetricType = "COSINE"
)

type IndexType string

const (
	// vector index type
	HNSW   IndexType = "HNSW"
	FLAT   IndexType = "FLAT"
	PUCK   IndexType = "PUCK"
	HNSWPQ IndexType = "HNSWPQ"

	// scalar index type
	SecondaryIndex IndexType = "SECONDARY"
)

type FieldType string

const (
	// scalar field type
	FieldTypeBool        FieldType = "BOOL"
	FieldTypeInt8        FieldType = "INT8"
	FieldTypeUint8       FieldType = "UINT8"
	FieldTypeInt16       FieldType = "INT16"
	FieldTypeUint16      FieldType = "UINT16"
	FieldTypeInt32       FieldType = "INT32"
	FieldTypeUint32      FieldType = "UINT32"
	FieldTypeInt64       FieldType = "INT64"
	FieldTypeUint64      FieldType = "UINT64"
	FieldTypeFloat       FieldType = "FLOAT"
	FieldTypeDouble      FieldType = "DOUBLE"
	FieldTypeDate        FieldType = "DATE"
	FieldTypeDatetime    FieldType = "DATETIME"
	FieldTypeTimestamp   FieldType = "TIMESTAMP"
	FieldTypeString      FieldType = "STRING"
	FieldTypeBinary      FieldType = "BINARY"
	FieldTypeUUID        FieldType = "UUID"
	FieldTypeText        FieldType = "TEXT"
	FieldTypeTextGBK     FieldType = "TEXT_GBK"
	FieldTypeTextGB18030 FieldType = "TEXT_GB18030"

	// vector field type
	FieldTypeFloatVector FieldType = "FLOAT_VECTOR"
)

type AutoBuildPolicyType string

const (
	AutoBuildPolicyTiming     AutoBuildPolicyType = "TIMING"
	AutoBuildPolicyPeriodical AutoBuildPolicyType = "PERIODICAL"
	AutoBuildPolicyIncrement  AutoBuildPolicyType = "ROW_COUNT_INCREMENT"
)

type PartitionType string

const (
	HASH PartitionType = "HASH"
)

type ReadConsistency string

const (
	EVENTUAL ReadConsistency = "EVENTUAL"
	STRONG   ReadConsistency = "STRONG"
)

type TableState string

const (
	TableStateCreating TableState = "CREATING"
	TableStateNormal   TableState = "NORMAL"
	TableStateDeleting TableState = "DELETING"
)

type IndexState string

const (
	IndexStateBuilding IndexState = "BUILDING"
	IndexStateNormal   IndexState = "NORMAL"
)

type ServerErrCode int32

const (
	OK                         ServerErrCode = 0
	InternalError              ServerErrCode = 1
	InvalidParameter           ServerErrCode = 2
	InvalidHTTPURL             ServerErrCode = 10
	InvalidHTTPHeader          ServerErrCode = 11
	InvalidHTTPBody            ServerErrCode = 12
	MissSSLCertificates        ServerErrCode = 13
	UserNotExist               ServerErrCode = 20
	UserAlreadyExist           ServerErrCode = 21
	RoleNotExist               ServerErrCode = 22
	RoleAlreadyExist           ServerErrCode = 23
	AuthenticationFailed       ServerErrCode = 24
	PermissionDenied           ServerErrCode = 25
	DBNotExist                 ServerErrCode = 50
	DBAlreadyExist             ServerErrCode = 51
	DBTooManyTables            ServerErrCode = 52
	DBNotEmpty                 ServerErrCode = 53
	InvalidTableSchema         ServerErrCode = 60
	InvalidPartitionParameters ServerErrCode = 61
	TableTooManyFields         ServerErrCode = 62
	TableTooManyFamilies       ServerErrCode = 63
	TableTooManyPrimaryKeys    ServerErrCode = 64
	TableTooManyPartitionKeys  ServerErrCode = 65
	TableTooManyVectorFields   ServerErrCode = 66
	TableTooManyIndexes        ServerErrCode = 67
	DynamicSchemaError         ServerErrCode = 68
	TableNotExist              ServerErrCode = 69
	TableAlreadyExist          ServerErrCode = 70
	InvalidTableState          ServerErrCode = 71
	TableNotReady              ServerErrCode = 72
	AliasNotExist              ServerErrCode = 73
	AliasAlreadyExist          ServerErrCode = 74
	FieldNotExist              ServerErrCode = 80
	FieldAlreadyExist          ServerErrCode = 81
	VectorFieldNotExist        ServerErrCode = 82
	InvalidIndexSchema         ServerErrCode = 90
	IndexNotExist              ServerErrCode = 91
	IndexAlreadyExist          ServerErrCode = 92
	IndexDuplicated            ServerErrCode = 93
	InvalidIndexState          ServerErrCode = 94
	PrimaryKeyDuplicated       ServerErrCode = 100
)