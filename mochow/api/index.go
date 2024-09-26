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

// index.go - the index APIs definition supported by the Mochow service

package api

import (
	"github.com/bytedance/sonic"

	"github.com/baidu/mochow-sdk-go/client"
	"github.com/baidu/mochow-sdk-go/http"
)

func CreateIndex(cli client.Client, args *CreateIndexArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getIndexURI())
	req.SetMethod(http.Post)
	req.SetParam("create", "")

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

func DescIndex(cli client.Client, args *DescIndexArgs) (*DescIndexResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getIndexURI())
	req.SetMethod(http.Post)
	req.SetParam("desc", "")

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
	result := &DescIndexResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func ModifyIndex(cli client.Client, args *ModifyIndexArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getIndexURI())
	req.SetMethod(http.Post)
	req.SetParam("modify", "")

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

func DropIndex(cli client.Client, database, table, indexName string) error {
	req := &client.BceRequest{}
	req.SetURI(getIndexURI())
	req.SetMethod(http.Delete)
	req.SetParam("database", database)
	req.SetParam("table", table)
	req.SetParam("indexName", indexName)

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

func RebuildIndex(cli client.Client, args *RebuildIndexArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getIndexURI())
	req.SetMethod(http.Post)
	req.SetParam("rebuild", "")

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
