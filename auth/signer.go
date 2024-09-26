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

// signer.go - implement the specific sign algorithm of Mochow V1 protocol

package auth

import (
	"fmt"

	"github.com/baidu/mochow-sdk-go/http"
	"github.com/baidu/mochow-sdk-go/util/log"
)

// Signer abstracts the entity that implements the `Sign` method
type Signer interface {
	// Sign the given Request with the Credentials and SignOptions
	Sign(*http.Request, *BceCredentials, *SignOptions)
}

// SignOptions defines the data structure used by Signer
type SignOptions struct {
	HeadersToSign map[string]struct{}
	Timestamp     int64
	ExpireSeconds int
}

func (opt *SignOptions) String() string {
	return fmt.Sprintf(`SignOptions [
        HeadersToSign=%s;
        Timestamp=%d;
        ExpireSeconds=%d
    ]`, opt.HeadersToSign, opt.Timestamp, opt.ExpireSeconds)
}

// BceV1Signer implements the v1 sign algorithm
type BceV1Signer struct{}

// Sign - generate the authorization string from the BceCredentials and SignOptions
//
// PARAMS:
//   - req: *http.Request for this sign
//   - cred: *BceCredentials to access the serice
//   - opt: *SignOptions for this sign algorithm
func (b *BceV1Signer) Sign(req *http.Request, cred *BceCredentials, opt *SignOptions) {
	if req == nil {
		log.Fatal("request should not be null for sign")
		return
	}
	if cred == nil {
		log.Fatal("credentials should not be null for sign")
		return
	}

	// Generate auth string and add to the reqeust header
	account := cred.Account
	apiKey := cred.APIKey
	authStr := fmt.Sprintf("Bearer account=%s&api_key=%s", account, apiKey)
	log.Info("Authorization=" + authStr)

	req.SetHeader(http.Authorization, authStr)
}
