/*
 * Copyright 2017 Baidu, Inc.
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

// credentials.go - the credentials data structure definition

// Package auth implements the authorization functionality for mochow api.
// It use the account and api key with the specific sign algorithm to generate the authorization string.
package auth

import "errors"

// BceCredentials define the data structure for authorization
type BceCredentials struct {
	Account string // account name to the service
	APIKey  string // api key to the service
}

func (b *BceCredentials) String() string {
	str := "account: " + b.Account + ", api_key: " + b.APIKey
	return str
}

func NewBceCredentials(account, apiKey string) (*BceCredentials, error) {
	if len(account) == 0 {
		return nil, errors.New("account should not be empty")
	}
	if len(apiKey) == 0 {
		return nil, errors.New("apiKey should not be empty")
	}

	return &BceCredentials{account, apiKey}, nil
}
