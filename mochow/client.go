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

// client.go - define the client for Mochow service

package mochow

import (
	"errors"

	"github.com/baidu/mochow-sdk-go/auth"
	"github.com/baidu/mochow-sdk-go/client"
	"github.com/baidu/mochow-sdk-go/mochow/api"
)

type Client struct {
	*client.BceClient
}

type ClientConfiguration struct {
	Account             string
	APIKey              string
	Endpoint            string
	RedirectDisabled    bool
	ConnectionTimeoutMS int
	RequestTimeoutMS    int
	MaxRetry            int
}

// NewClient make the Mochow service client with default configuration.
// Use `cli.Config.xxx` to access the config or change it to non-default value.
func NewClient(account, apiKey, endpoint string) (*Client, error) {
	return NewClientWithConfig(&ClientConfiguration{
		Account:          account,
		APIKey:           apiKey,
		Endpoint:         endpoint,
		RedirectDisabled: false,
	})
}

func NewClientWithConfig(config *ClientConfiguration) (*Client, error) {
	var credentials *auth.BceCredentials
	var err error

	// Init credentials with account and apikey
	account, apiKey, endpoint := config.Account, config.APIKey, config.Endpoint
	if len(account) == 0 || len(apiKey) == 0 || len(endpoint) == 0 {
		return nil, errors.New("account, apiKey and endpoint missing for creating mochow client")
	}
	credentials, err = auth.NewBceCredentials(account, apiKey)
	if err != nil {
		return nil, err
	}

	defaultConf := &client.BceClientConfiguration{
		Endpoint:                  endpoint,
		Region:                    client.DefaultRegion,
		UserAgent:                 client.DefaultUserAgent,
		Credentials:               credentials,
		SignOption:                nil,
		Retry:                     client.DefaultRetryPolicy,
		ConnectionTimeoutInMillis: client.DefaultConnectionTimeoutInMills,
		RequestTimeoutInMillis:    client.DefaultRequestTimeoutInMills,
		RedirectDisabled:          config.RedirectDisabled}

	// Check timeout options
	if config.ConnectionTimeoutMS < 0 || config.RequestTimeoutMS < 0 {
		return nil, errors.New("connection and request timeout is negative")
	}
	if config.ConnectionTimeoutMS > 0 {
		defaultConf.ConnectionTimeoutInMillis = config.ConnectionTimeoutMS
	}
	if config.RequestTimeoutMS > 0 {
		defaultConf.RequestTimeoutInMillis = config.RequestTimeoutMS
	}
	if defaultConf.RequestTimeoutInMillis <= defaultConf.ConnectionTimeoutInMillis {
		return nil, errors.New("request timeout should greater than connection timeout")
	}

	// Check max retry option
	if config.MaxRetry < 0 {
		// negative max retry means no retry
		defaultConf.Retry = client.NewNoRetryPolicy()
	} else if config.MaxRetry > 0 {
		defaultConf.Retry = client.NewBackOffRetryPolicy(config.MaxRetry, 20000, 300)
	}

	v1Signer := &auth.BceV1Signer{}
	client := &Client{client.NewBceClient(defaultConf, v1Signer)}
	return client, nil
}

/********************* Database interfaces *********************/
func (c *Client) CreateDatabase(database string) error {
	args := &api.CreateDatabaseArgs{Database: database}
	return api.CreateDatabase(c, args)
}

func (c *Client) DropDatabase(database string) error {
	return api.DropDatabase(c, database)
}

func (c *Client) ListDatabase() (*api.ListDatabaseResult, error) {
	return api.ListDatabase(c)
}

func (c *Client) HasDatabase(database string) (bool, error) {
	listDatabaseResult, err := c.ListDatabase()
	if err != nil {
		return false, err
	}
	for _, db := range listDatabaseResult.Databases {
		if db == database {
			return true, nil
		}
	}
	return false, nil
}

/********************* Table interfaces *********************/
func (c *Client) CreateTable(args *api.CreateTableArgs) error {
	return api.CreateTable(c, args)
}

func (c *Client) DropTable(database, table string) error {
	return api.DropTable(c, database, table)
}

func (c *Client) ListTable(database string) (*api.ListTableResult, error) {
	args := &api.ListTableArgs{Database: database}
	return api.ListTable(c, args)
}

func (c *Client) HasTable(database, table string) (bool, error) {
	listTableResult, err := c.ListTable(database)
	if err != nil {
		return false, err
	}
	for _, tableName := range listTableResult.Tables {
		if tableName == table {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) DescTable(database, table string) (*api.DescTableResult, error) {
	args := &api.DescTableArgs{Database: database, Table: table}
	return api.DescTable(c, args)
}

func (c *Client) AddField(args *api.AddFieldArgs) error {
	return api.AddField(c, args)
}

func (c *Client) AliasTable(database, table, alias string) error {
	args := &api.AliasTableArgs{Database: database, Table: table, Alias: alias}
	return api.AliasTable(c, args)
}

func (c *Client) UnaliasTable(database, table, alias string) error {
	args := &api.UnaliasTableArgs{Database: database, Table: table, Alias: alias}
	return api.UnaliasTable(c, args)
}

func (c *Client) ShowTableStats(database, table string) (*api.ShowTableStatsResult, error) {
	args := &api.ShowTableStatsArgs{Database: database, Table: table}
	return api.ShowTableStats(c, args)
}

func (c *Client) CreateIndex(args *api.CreateIndexArgs) error {
	return api.CreateIndex(c, args)
}

func (c *Client) DescIndex(database, table, indexName string) (*api.DescIndexResult, error) {
	args := &api.DescIndexArgs{Database: database, Table: table, IndexName: indexName}
	return api.DescIndex(c, args)
}

func (c *Client) ModifyIndex(args *api.ModifyIndexArgs) error {
	return api.ModifyIndex(c, args)
}

func (c *Client) DropIndex(database, table, indexName string) error {
	return api.DropIndex(c, database, table, indexName)
}

func (c *Client) RebuildIndex(database, table, indexName string) error {
	args := &api.RebuildIndexArgs{Database: database, Table: table, IndexName: indexName}
	return api.RebuildIndex(c, args)
}

func (c *Client) InsertRow(args *api.InsertRowArgs) (*api.InsertRowResult, error) {
	return api.InsertRow(c, args)
}

func (c *Client) UpsertRow(args *api.UpsertRowArg) (*api.UpsertRowResult, error) {
	return api.UpsertRow(c, args)
}

func (c *Client) DeleteRow(args *api.DeleteRowArgs) error {
	return api.DeleteRow(c, args)
}

func (c *Client) QueryRow(args *api.QueryRowArgs) (*api.QueryRowResult, error) {
	return api.QueryRow(c, args)
}

func (c *Client) SearchRow(args *api.SearchRowArgs) (*api.SearchRowResult, error) {
	return api.SearchRow(c, args)
}

func (c *Client) UpdateRow(args *api.UpdateRowArgs) error {
	return api.UpdateRow(c, args)
}

func (c *Client) SelectRow(args *api.SelectRowArgs) (*api.SelectRowResult, error) {
	return api.SelectRow(c, args)
}

func (c *Client) BatchSearchRow(args *api.BatchSearchRowArgs) (*api.BatchSearchRowResult, error) {
	return api.BatchSearchRow(c, args)
}
