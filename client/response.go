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

// response.go - defines the common Mochow services response

package client

import (
	"io"
	"time"

	"github.com/bytedance/sonic/decoder"

	"github.com/baidu/mochow-sdk-go/http"
)

// BceResponse defines the response structure for receiving BCE services response.
type BceResponse struct {
	statusCode   int
	statusText   string
	requestID    string
	debugID      string
	response     *http.Response
	serviceError *BceServiceError
}

func (r *BceResponse) IsFail() bool {
	return r.response.StatusCode() >= 400
}

func (r *BceResponse) StatusCode() int {
	return r.statusCode
}

func (r *BceResponse) StatusText() string {
	return r.statusText
}

func (r *BceResponse) RequestID() string {
	return r.requestID
}

func (r *BceResponse) DebugID() string {
	return r.debugID
}

func (r *BceResponse) Header(key string) string {
	return r.response.GetHeader(key)
}

func (r *BceResponse) Headers() map[string]string {
	return r.response.GetHeaders()
}

func (r *BceResponse) Body() io.ReadCloser {
	return r.response.Body()
}

func (r *BceResponse) SetHTTPResponse(response *http.Response) {
	r.response = response
}

func (r *BceResponse) ElapsedTime() time.Duration {
	return r.response.ElapsedTime()
}

func (r *BceResponse) ServiceError() *BceServiceError {
	return r.serviceError
}

func (r *BceResponse) ParseResponse() {
	r.statusCode = r.response.StatusCode()
	r.statusText = r.response.StatusText()
	r.requestID = r.response.GetHeader(http.RequestID)
	if r.IsFail() {
		r.serviceError = NewBceServiceError(-1, r.statusText, r.requestID, r.statusCode)

		// First try to read the error `Code' and `Message' from body
		rawBody, _ := io.ReadAll(r.Body())
		defer r.Body().Close()
		if len(rawBody) != 0 {
			jsonDecoder := decoder.NewDecoder(string(rawBody))
			if err := jsonDecoder.Decode(r.serviceError); err != nil {
				r.serviceError = NewBceServiceError(
					-1,
					"Service json error message decode failed",
					r.requestID,
					r.statusCode)
			}
			return
		}
	}
}

func (r *BceResponse) ParseJSONBody(result interface{}) error {
	defer r.Body().Close()
	jsonDecoder := decoder.NewStreamDecoder(r.Body())
	return jsonDecoder.Decode(result)
}
