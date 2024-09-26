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

// table.go - the table APIs definition supported by the Mochow service

package api

import (
	"github.com/bytedance/sonic"

	"github.com/baidu/mochow-sdk-go/client"
	"github.com/baidu/mochow-sdk-go/http"
)

func CreateTable(cli client.Client, args *CreateTableArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
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

func DropTable(cli client.Client, database, table string) error {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
	req.SetMethod(http.Delete)
	req.SetParam("database", database)
	req.SetParam("table", table)

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

func ListTable(cli client.Client, args *ListTableArgs) (*ListTableResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
	req.SetMethod(http.Post)
	req.SetParam("list", "")

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
	result := &ListTableResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func DescTable(cli client.Client, args *DescTableArgs) (*DescTableResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
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
	result := &DescTableResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}

func AddField(cli client.Client, args *AddFieldArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
	req.SetMethod(http.Post)
	req.SetParam("addField", "")

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

func AliasTable(cli client.Client, args *AliasTableArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
	req.SetMethod(http.Post)
	req.SetParam("alias", "")

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

func UnaliasTable(cli client.Client, args *UnaliasTableArgs) error {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
	req.SetMethod(http.Post)
	req.SetParam("unalias", "")

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

func ShowTableStats(cli client.Client, args *ShowTableStatsArgs) (*ShowTableStatsResult, error) {
	req := &client.BceRequest{}
	req.SetURI(getTableURI())
	req.SetMethod(http.Post)
	req.SetParam("stats", "")

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
	result := &ShowTableStatsResult{}
	if err := resp.ParseJSONBody(result); err != nil {
		return nil, err
	}
	return result, nil
}
