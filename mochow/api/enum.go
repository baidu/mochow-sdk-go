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
	HNSW    IndexType = "HNSW"
	FLAT    IndexType = "FLAT"
	PUCK    IndexType = "PUCK"
	HNSWPQ  IndexType = "HNSWPQ"
	SPARSE  IndexType = "SPARSE_OPTIMIZED_FLAT"
	HNSWSQ  IndexType = "HNSWSQ"
	DISKANN IndexType = "DISKANN"
	IVF     IndexType = "IVF"
	IVFPQ   IndexType = "IVFPQ"
	IVFSQ   IndexType = "IVFSQ"

	// scalar index type
	SecondaryIndex                  IndexType = "SECONDARY"
	PersistentBitmapIndex           IndexType = "PERSISTENT_BITMAP"
	PersistentAggregatedBitmapIndex IndexType = "PERSISTENT_AGGREGATED_BITMAP"

	// filtering index type
	FilteringIndex IndexType = "FILTERING"

	// inverted index type
	InvertedIndex IndexType = "INVERTED"
)

type InvertedIndexAnalyzer string

const (
	EnglishAnalyzer InvertedIndexAnalyzer = "ENGLISH_ANALYZER"
	ChineseAnalyzer InvertedIndexAnalyzer = "CHINESE_ANALYZER"
	DefaultAnalyzer InvertedIndexAnalyzer = "DEFAULT_ANALYZER"
)

type InvertedIndexParseMode string

const (
	CoarseMode InvertedIndexParseMode = "COARSE_MODE"
	FineMode   InvertedIndexParseMode = "FINE_MODE"
)

type InvertedIndexFieldAttribute string

const (
	NotAnalyzed InvertedIndexFieldAttribute = "ATTRIBUTE_NOT_ANALYZED"
	Analyzed    InvertedIndexFieldAttribute = "ATTRIBUTE_ANALYZED"
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
	FieldTypeJSON        FieldType = "JSON"
	FieldTypeMap         FieldType = "MAP"

	// vector field type
	FieldTypeFloatVector  FieldType = "FLOAT_VECTOR"
	FieldTypeBinaryVector FieldType = "BINARY_VECTOR"
	FieldTypeSparseVector FieldType = "SPARSE_FLOAT_VECTOR"

	// array field type
	FieldTypeArray FieldType = "ARRAY"
)

type ElementType string

const (
	ElementTypeBool        ElementType = "BOOL"
	ElementTypeInt8        ElementType = "INT8"
	ElementTypeUint8       ElementType = "UINT8"
	ElementTypeInt16       ElementType = "INT16"
	ElementTypeUint16      ElementType = "UINT16"
	ElementTypeInt32       ElementType = "INT32"
	ElementTypeUint32      ElementType = "UINT32"
	ElementTypeInt64       ElementType = "INT64"
	ElementTypeUint64      ElementType = "UINT64"
	ElementTypeFloat       ElementType = "FLOAT"
	ElementTypeDouble      ElementType = "DOUBLE"
	ElementTypeDate        ElementType = "DATE"
	ElementTypeDatetime    ElementType = "DATETIME"
	ElementTypeTimestamp   ElementType = "TIMESTAMP"
	ElementTypeString      ElementType = "STRING"
	ElementTypeBinary      ElementType = "BINARY"
	ElementTypeUUID        ElementType = "UUID"
	ElementTypeText        ElementType = "TEXT"
	ElementTypeTextGBK     ElementType = "TEXT_GBK"
	ElementTypeTextGB18030 ElementType = "TEXT_GB18030"
	ElementTypeFloatVector ElementType = "FLOAT_VECTOR"
)

type MapKeyType string

const (
	MapKeyTypeBool        MapKeyType = "BOOL"
	MapKeyTypeInt8        MapKeyType = "INT8"
	MapKeyTypeUint8       MapKeyType = "UINT8"
	MapKeyTypeInt16       MapKeyType = "INT16"
	MapKeyTypeUint16      MapKeyType = "UINT16"
	MapKeyTypeInt32       MapKeyType = "INT32"
	MapKeyTypeUint32      MapKeyType = "UINT32"
	MapKeyTypeInt64       MapKeyType = "INT64"
	MapKeyTypeUint64      MapKeyType = "UINT64"
	MapKeyTypeFloat       MapKeyType = "FLOAT"
	MapKeyTypeDouble      MapKeyType = "DOUBLE"
	MapKeyTypeDate        MapKeyType = "DATE"
	MapKeyTypeDatetime    MapKeyType = "DATETIME"
	MapKeyTypeTimestamp   MapKeyType = "TIMESTAMP"
	MapKeyTypeString      MapKeyType = "STRING"
	MapKeyTypeBinary      MapKeyType = "BINARY"
	MapKeyTypeUUID        MapKeyType = "UUID"
	MapKeyTypeText        MapKeyType = "TEXT"
	MapKeyTypeTextGBK     MapKeyType = "TEXT_GBK"
	MapKeyTypeTextGB18030 MapKeyType = "TEXT_GB18030"
)

type MapValueType string

const (
	MapValueTypeBool        MapValueType = "BOOL"
	MapValueTypeInt8        MapValueType = "INT8"
	MapValueTypeUint8       MapValueType = "UINT8"
	MapValueTypeInt16       MapValueType = "INT16"
	MapValueTypeUint16      MapValueType = "UINT16"
	MapValueTypeInt32       MapValueType = "INT32"
	MapValueTypeUint32      MapValueType = "UINT32"
	MapValueTypeInt64       MapValueType = "INT64"
	MapValueTypeUint64      MapValueType = "UINT64"
	MapValueTypeFloat       MapValueType = "FLOAT"
	MapValueTypeDouble      MapValueType = "DOUBLE"
	MapValueTypeDate        MapValueType = "DATE"
	MapValueTypeDatetime    MapValueType = "DATETIME"
	MapValueTypeTimestamp   MapValueType = "TIMESTAMP"
	MapValueTypeString      MapValueType = "STRING"
	MapValueTypeBinary      MapValueType = "BINARY"
	MapValueTypeUUID        MapValueType = "UUID"
	MapValueTypeText        MapValueType = "TEXT"
	MapValueTypeTextGBK     MapValueType = "TEXT_GBK"
	MapValueTypeTextGB18030 MapValueType = "TEXT_GB18030"
	MapValueTypeFloatVector MapValueType = "FLOAT_VECTOR"
	MapValueTypeArray       MapValueType = "ARRAY"
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
	ReadConsistencyEventual ReadConsistency = "EVENTUAL"
	ReadConsistencyStrong   ReadConsistency = "STRONG"
)

type TableState string

const (
	TableStateCreating TableState = "CREATING"
	TableStateNormal   TableState = "NORMAL"
	TableStateDeleting TableState = "DELETING"
)

type IndexState string

const (
	IndexStateInvalid  IndexState = "INVALID"
	IndexStateBuilding IndexState = "BUILDING"
	IndexStateNormal   IndexState = "NORMAL"
)

type IndexStructureType string

const (
	IndexStructureTypeDefault IndexStructureType = "DEFAULT"
	IndexStructureTypeBitmap  IndexStructureType = "BITMAP"
)

type ServerErrCode int32

const (
	OK                         ServerErrCode = 0
	InternalError              ServerErrCode = 1
	InvalidParameter           ServerErrCode = 2
	Timeout                    ServerErrCode = 3
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
	TableMemoryExceedLimit     ServerErrCode = 75
	FieldNotExist              ServerErrCode = 80
	FieldAlreadyExist          ServerErrCode = 81
	VectorFieldNotExist        ServerErrCode = 82
	InvalidIndexSchema         ServerErrCode = 90
	IndexNotExist              ServerErrCode = 91
	IndexAlreadyExist          ServerErrCode = 92
	IndexDuplicated            ServerErrCode = 93
	InvalidIndexState          ServerErrCode = 94
	PrimaryKeyDuplicated       ServerErrCode = 100
	RowKeyNotFound             ServerErrCode = 101
)

type FilterMode string

const (
	FilterModeAuto FilterMode = "AUTO"
	FilterModePost FilterMode = "POST"
)
