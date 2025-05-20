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

// row.go - the row APIs definition supported by the Mochow service

package api

import (
	"github.com/bytedance/sonic"

	"github.com/baidu/mochow-sdk-go/v2/client"
	"github.com/baidu/mochow-sdk-go/v2/http"
)

func InsertRow(cli client.Client, args *InsertRowArgs) (*InsertRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("insert", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &InsertRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func UpsertRow(cli client.Client, args *UpsertRowArg) (*UpsertRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("upsert", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &UpsertRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteRow(cli client.Client, args *DeleteRowArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("delete", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

func QueryRow(cli client.Client, args *QueryRowArgs) (*QueryRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("query", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &QueryRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func BatchQueryRow(cli client.Client, args *BatchQueryRowArgs) (*BatchQueryRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("batchQuery", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &BatchQueryRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func VectorSearch(cli client.Client, args *VectorSearchArgs) (*SearchResult, error) {
	return search(cli, args.Database, args.Table, args.Request)
}

func BM25Search(cli client.Client, args *BM25SearchArgs) (*SearchResult, error) {
	return search(cli, args.Database, args.Table, args.Request)
}

func HybridSearch(cli client.Client, args *HybridSearchArgs) (*SearchResult, error) {
	return search(cli, args.Database, args.Table, args.Request)
}

func MultiVectorSearch(cli client.Client, args *MultivectorSearchArgs) (*SearchResult, error) {
	return search(cli, args.Database, args.Table, args.Request)
}

func search(cli client.Client, database string, table string, request searchRequest) (*SearchResult, error) {
	args := request.toDict()
	args["database"] = database
	args["table"] = table

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}

	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam(request.requestType(), "")
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &SearchResult{IsBatch: request.isBatch()}
	if result.IsBatch {
		result.BatchRows = &BatchSearchRowResult{}
		if err := resp.ParseJSONBody(result.BatchRows); err != nil {
			return nil, err
		}
	} else {
		result.Rows = &SearchRowResult{}
		if err := resp.ParseJSONBody(result.Rows); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Deprecated: you should use VectorSearch with VectorTopkSearchRequest or VectorRangeSearchRequest instead.
func SearchRow(cli client.Client, args *SearchRowArgs) (*SearchRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("search", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &SearchRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateRow(cli client.Client, args *UpdateRowArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("update", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

func SelectRow(cli client.Client, args *SelectRowArgs) (*SelectRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("select", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &SelectRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Deprecated: you should use VectorSearch with VectorBatchSearchRequest instead.
func BatchSearchRow(cli client.Client, args *BatchSearchRowArgs) (*BatchSearchRowResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getRowURI())
	req.SetMethod(http.Post)
	req.SetParam("batchSearch", "")

	jsonBytes, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}
	body, err := client.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	resp := &client.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	result := &BatchSearchRowResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}
