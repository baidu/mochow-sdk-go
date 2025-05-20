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

// client.go - definiton the BceClientConfiguration and BceClient structure

// Package client implements the infrastructure to access BCE services.
//
//   - BceClient:
//     It is the general client of BCE to access all services. It builds http request to access the
//     services based on the given client configuration.
//
//   - BceClientConfiguration:
//     The client configuration data structure which contains endpoint, region, credentials, retry
//     policy, sign options and so on. It supports most of the default value and user can also
//     access or change the default with its public fields' name.
//
//   - Error types:
//     The error types when making request or receiving response to the BCE services contains two
//     types: the BceClientError when making request to BCE services and the BceServiceError when
//     recieving response from them.
//
//   - BceRequest:
//     The request instance stands for an request to access the BCE services.
//
//   - BceResponse:
//     The response instance stands for an response from the BCE services.

package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/baidu/mochow-sdk-go/v2/auth"
	"github.com/baidu/mochow-sdk-go/v2/http"
	"github.com/baidu/mochow-sdk-go/v2/util"
	"github.com/baidu/mochow-sdk-go/v2/util/log"
)

// Client is the general interface which can perform sending request. Different service
// will define its own client in case of specific extension.
type Client interface {
	SendRequest(*BceRequest, *BceResponse) error
	SendRequestFromBytes(*BceRequest, *BceResponse, []byte) error
	GetBceClientConfig() *BceClientConfiguration
}

// BceClient defines the general client to access the BCE services.
type BceClient struct {
	Config *BceClientConfiguration
	Signer auth.Signer // the sign algorithm
}

// BuildHttpRequest - the helper method for the client to build http request
//
// PARAMS:
//   - request: the input request object to be built
func (c *BceClient) buildHTTPRequest(request *BceRequest) {
	// Construct the http request instance for the special fields
	request.BuildHTTPRequest()

	// Set the client specific configurations
	if request.Endpoint() == "" {
		request.SetEndpoint(c.Config.Endpoint)
	}
	if request.Protocol() == "" {
		request.SetProtocol(DefaultProtocol)
	}
	if len(c.Config.ProxyURL) != 0 {
		request.SetProxyURL(c.Config.ProxyURL)
	}
	request.SetTimeout(c.Config.RequestTimeoutInMillis / 1000)

	// Set the BCE request headers
	request.SetHeader(http.Host, request.Host())
	request.SetHeader(http.UserAgent, c.Config.UserAgent)
	request.SetHeader(http.Date, util.FormatISO8601Date(util.NowUTCSeconds()))
	request.SetHeader(http.RequestTimeoutMS, strconv.Itoa(c.Config.RequestTimeoutInMillis))

	//set default content-type if null
	if request.Header(http.ContentType) == "" {
		request.SetHeader(http.ContentType, DefaultContentType)
	}

	// Generate the auth string if needed
	if c.Config.Credentials != nil {
		c.Signer.Sign(&request.Request, c.Config.Credentials, c.Config.SignOption)
	}
}

// SendRequest - the client performs sending the http request with retry policy and receive the
// response from the BCE services.
//
// PARAMS:
//   - req: the request object to be sent to the BCE service
//   - resp: the response object to receive the content from BCE service
//
// RETURNS:
//   - error: nil if ok otherwise the specific error
func (c *BceClient) SendRequest(req *BceRequest, resp *BceResponse) error {
	// Return client error if it is not nil
	if req.ClientError() != nil {
		return req.ClientError()
	}

	// Build the http request and prepare to send
	c.buildHTTPRequest(req)
	log.Infof("send http request: %v", req)

	// Send request with the given retry policy
	retries := 0
	if req.Body() != nil {
		defer req.Body().Close() // Manually close the ReadCloser body for retry
	}
	for {
		// The request body should be temporarily saved if retry to send the http request
		var retryBuf bytes.Buffer
		var teeReader io.Reader
		if c.Config.Retry.ShouldRetry(nil, 0) && req.Body() != nil {
			teeReader = io.TeeReader(req.Body(), &retryBuf)
			req.Request.SetBody(ioutil.NopCloser(teeReader))
		}
		httpResp, err := http.Execute(&req.Request)

		if err != nil {
			if c.Config.Retry.ShouldRetry(err, retries) {
				delayInMills := c.Config.Retry.GetDelayBeforeNextRetryInMillis(err, retries)
				time.Sleep(delayInMills)
			} else {
				return &BceClientError{
					fmt.Sprintf("execute http request failed! Retried %d times, error: %v",
						retries, err)}
			}
			retries++
			log.Warnf("send request failed: %v, retry for %d time(s)", err, retries)
			if req.Body() != nil {
				_, _ = io.ReadAll(teeReader)
				req.Request.SetBody(ioutil.NopCloser(&retryBuf))
			}
			continue
		}
		resp.SetHTTPResponse(httpResp)
		resp.ParseResponse()

		log.Infof("receive http response: status: %s, debugId: %s, requestId: %s, elapsed: %v",
			resp.StatusText(), resp.DebugID(), resp.RequestID(), resp.ElapsedTime())

		if resp.ElapsedTime().Milliseconds() > DefaultWarnLogTimeoutInMills {
			log.Warnf("request time more than 5 second, debugId: %s, requestId: %s, elapsed: %v",
				resp.DebugID(), resp.RequestID(), resp.ElapsedTime())
		}
		for k, v := range resp.Headers() {
			log.Debugf("%s=%s", k, v)
		}
		if resp.IsFail() {
			err := resp.ServiceError()
			if c.Config.Retry.ShouldRetry(err, retries) {
				delayInMills := c.Config.Retry.GetDelayBeforeNextRetryInMillis(err, retries)
				time.Sleep(delayInMills)
			} else {
				return err
			}
			retries++
			log.Warnf("send request failed, retry for %d time(s)", retries)
			if req.Body() != nil {
				_, _ = io.ReadAll(teeReader)
				req.Request.SetBody(ioutil.NopCloser(&retryBuf))
			}
			continue
		}
		return nil
	}
}

// SendRequestFromBytes - the client performs sending the http request with retry policy and receive the
// response from the BCE services.
//
// PARAMS:
//   - req: the request object to be sent to the BCE service
//   - resp: the response object to receive the content from BCE service
//   - content: the content of body
//
// RETURNS:
//   - error: nil if ok otherwise the specific error
func (c *BceClient) SendRequestFromBytes(req *BceRequest, resp *BceResponse, content []byte) error {
	// Return client error if it is not nil
	if req.ClientError() != nil {
		return req.ClientError()
	}
	// Build the http request and prepare to send
	c.buildHTTPRequest(req)
	log.Infof("send http request: %v", req)
	// Send request with the given retry policy
	retries := 0
	for {
		// The request body should be temporarily saved if retry to send the http request
		buf := bytes.NewBuffer(content)
		req.Request.SetBody(ioutil.NopCloser(buf))
		defer req.Request.Body().Close() // Manually close the ReadCloser body for retry
		httpResp, err := http.Execute(&req.Request)
		if err != nil {
			if c.Config.Retry.ShouldRetry(err, retries) {
				delayInMills := c.Config.Retry.GetDelayBeforeNextRetryInMillis(err, retries)
				time.Sleep(delayInMills)
			} else {
				return &BceClientError{
					fmt.Sprintf("execute http request failed! Retried %d times, error: %v",
						retries, err)}
			}
			retries++
			log.Warnf("send request failed: %v, retry for %d time(s)", err, retries)
			continue
		}
		resp.SetHTTPResponse(httpResp)
		resp.ParseResponse()
		log.Infof("receive http response: status: %s, debugId: %s, requestId: %s, elapsed: %v",
			resp.StatusText(), resp.DebugID(), resp.RequestID(), resp.ElapsedTime())
		for k, v := range resp.Headers() {
			log.Debugf("%s=%s", k, v)
		}
		if resp.IsFail() {
			err := resp.ServiceError()
			if c.Config.Retry.ShouldRetry(err, retries) {
				delayInMills := c.Config.Retry.GetDelayBeforeNextRetryInMillis(err, retries)
				time.Sleep(delayInMills)
			} else {
				return err
			}
			retries++
			log.Warnf("send request failed, retry for %d time(s)", retries)
			continue
		}
		return nil
	}
}

func (c *BceClient) GetBceClientConfig() *BceClientConfiguration {
	return c.Config
}

func NewBceClient(conf *BceClientConfiguration, sign auth.Signer) *BceClient {
	clientConfig := http.ClientConfig{
		RedirectDisabled:         conf.RedirectDisabled,
		ConnectionTimeoutInMills: conf.ConnectionTimeoutInMillis,
	}
	http.InitClient(clientConfig)
	return &BceClient{conf, sign}
}

func NewBceClientWithAPIKey(account, apiKey, endPoint string) (*BceClient, error) {
	credentials, err := auth.NewBceCredentials(account, apiKey)
	if err != nil {
		return nil, err
	}
	defaultConf := &BceClientConfiguration{
		Endpoint:                  endPoint,
		Region:                    DefaultRegion,
		UserAgent:                 DefaultUserAgent,
		Credentials:               credentials,
		SignOption:                nil,
		Retry:                     DefaultRetryPolicy,
		ConnectionTimeoutInMillis: DefaultConnectionTimeoutInMills,
		RedirectDisabled:          false}
	v1Signer := &auth.BceV1Signer{}

	return NewBceClient(defaultConf, v1Signer), nil
}
