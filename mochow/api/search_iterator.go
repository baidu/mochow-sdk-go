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

// search_iterator.go - implementation of search iterator for vector search

package api

import (
	"fmt"

	"github.com/baidu/mochow-sdk-go/v2/client"
)

// SearchIterator provides efficient pagination through large search result sets
type SearchIterator struct {
	table           string
	database        string
	client          client.Client
	request         vectorSearchRequest
	batchSize       uint32
	totalSize       uint32
	partitionKey    map[string]interface{}
	projections     []string
	readConsistency ReadConsistency
	iteratedIds     string
	returnedCount   uint32
	config          map[string]interface{}
}

// SearchIteratorOptions contains all the parameters needed to create a SearchIterator
type SearchIteratorOptions struct {
	Client          client.Client
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

// NewSearchIterator creates a new SearchIterator with the given options
func NewSearchIterator(opts *SearchIteratorOptions) (*SearchIterator, error) {
	// Validate request type
	switch opts.Request.(type) {
	case *VectorTopkSearchRequest, *MultivectorSearchRequest:
	default:
		return nil, fmt.Errorf("SearchIterator only supports VectorTopkSearchRequest and MultivectorSearchRequest")
	}

	// Validate sizes
	if opts.TotalSize < opts.BatchSize {
		return nil, fmt.Errorf("'totalSize' should not be less than 'batchSize'")
	}

	// Validate batch size and request limit
	switch r := opts.Request.(type) {
	case *VectorTopkSearchRequest:
		if r.limit != opts.BatchSize {
			return nil, fmt.Errorf("'request.limit' should be equal to 'batchSize'")
		}
	case *MultivectorSearchRequest:
		if r.limit != opts.BatchSize && r.isMarked("limit") {
			return nil, fmt.Errorf("'request.limit' should be equal to 'batchSize'")
		}
	}

	return &SearchIterator{
		table:           opts.Table,
		database:        opts.Database,
		client:          opts.Client,
		request:         opts.Request,
		batchSize:       opts.BatchSize,
		totalSize:       opts.TotalSize,
		partitionKey:    opts.PartitionKey,
		projections:     opts.Projections,
		readConsistency: opts.ReadConsistency,
		iteratedIds:     "",
		returnedCount:   0,
		config:          opts.Config,
	}, nil
}

// Next returns the next batch of search results
// Returns nil when the iterator is finished
func (si *SearchIterator) Next() ([]RowResult, error) {
	if si.returnedCount >= si.totalSize {
		return nil, nil
	}
	// 设置 iteratedIds
	si.request.SetIteratedIds(si.iteratedIds)
	// 调用 search
	result, err := search(si.client, si.database, si.table, si.request)
	if err != nil {
		return nil, err
	}
	if result.Rows == nil || result.Rows.IteratedIds == "" {
		return nil, fmt.Errorf("search iterator is not supported by the server")
	}
	si.iteratedIds = result.Rows.IteratedIds
	rowCount := uint32(len(result.Rows.Rows))
	if rowCount == 0 {
		return nil, nil
	}
	limit := min(si.totalSize-si.returnedCount, rowCount)
	si.returnedCount += limit
	if limit < rowCount {
		return result.Rows.Rows[:limit], nil
	}
	return result.Rows.Rows, nil
}

// Close cleans up any resources used by the iterator
func (si *SearchIterator) Close() {
	// No resources to clean up in this implementation
}

// min returns the minimum of two uint32 values
func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
