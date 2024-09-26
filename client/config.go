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

// config.go - define the client configuration for BCE

package client

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/baidu/mochow-sdk-go/auth"
)

// Constants and default values for the package bce
const (
	SdkVersion                      = "0.0.1"
	DefaultProtocol                 = "http"
	DefaultRegion                   = "bj"
	DefaultContentType              = "application/json;charset=utf-8"
	DefaultConnectionTimeoutInMills = 10 * 1000
	DefaultRequestTimeoutInMills    = 60 * 1000
	DefaultWarnLogTimeoutInMills    = 5 * 1000
)

var (
	DefaultUserAgent   string
	DefaultRetryPolicy = NewBackOffRetryPolicy(3, 20000, 300)
)

func init() {
	DefaultUserAgent = "mochow-sdk-go"
	DefaultUserAgent += "/" + SdkVersion
	DefaultUserAgent += "/" + runtime.Version()
	DefaultUserAgent += "/" + runtime.GOOS
	DefaultUserAgent += "/" + runtime.GOARCH
}

// BceClientConfiguration defines the config components structure.
type BceClientConfiguration struct {
	Endpoint                  string
	ProxyURL                  string
	Region                    string
	UserAgent                 string
	Credentials               *auth.BceCredentials
	SignOption                *auth.SignOptions
	Retry                     RetryPolicy
	ConnectionTimeoutInMillis int
	RequestTimeoutInMillis    int
	// CnameEnabled should be true when use custom domain as endpoint to visit bos resource
	CnameEnabled     bool
	BackupEndpoint   string
	RedirectDisabled bool
}

func (c *BceClientConfiguration) String() string {
	return fmt.Sprintf(`BceClientConfiguration [
        Endpoint=%s;
        ProxyURL=%s;
        Region=%s;
        UserAgent=%s;
        Credentials=%v;
        SignOption=%v;
        RetryPolicy=%v;
        ConnectionTimeoutInMillis=%v;
		RedirectDisabled=%v
    ]`, c.Endpoint, c.ProxyURL, c.Region, c.UserAgent, c.Credentials,
		c.SignOption, reflect.TypeOf(c.Retry).Name(), c.ConnectionTimeoutInMillis, c.RedirectDisabled)
}
