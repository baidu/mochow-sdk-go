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

// error.go - define the error types for Mochow

package client

import "strconv"

const (
	accessDenied          = "AccessDenied"
	inappropriateJSON     = "InappropriateJSON"
	internalError         = "InternalError"
	invalidHTTPAuthHeader = "InvalidHTTPAuthHeader"
	invalidHTTPRequest    = "InvalidHTTPRequest"
	invalidURI            = "InvalidURI"
	malformedJSON         = "MalformedJSON"
	preConditionFailed    = "PreconditionFailed"
	requestExpired        = "RequestExpired"
)

// BceError abstracts the error for BCE
type BceError interface {
	error
}

// BceClientError defines the error struct for the client when making request
type BceClientError struct{ Message string }

func (b *BceClientError) Error() string { return b.Message }

func NewBceClientError(msg string) *BceClientError { return &BceClientError{msg} }

// BceServiceError defines the error struct for the BCE service when receiving response
type BceServiceError struct {
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	RequestID  string
	StatusCode int
}

func (b *BceServiceError) Error() string {
	ret := "[Code: " + strconv.Itoa(b.Code)
	ret += "; Message: " + b.Message
	ret += "; RequestId: " + b.RequestID + "]"
	return ret
}

func NewBceServiceError(code int, msg, reqID string, status int) *BceServiceError {
	return &BceServiceError{code, msg, reqID, status}
}
