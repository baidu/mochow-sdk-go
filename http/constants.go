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

// constants.go - defines constants of the Mochow http package including headers and methods

package http

// Constants of the supported HTTP methods
const (
	Get     = "GET"
	Put     = "PUT"
	Post    = "POST"
	Delete  = "DELETE"
	Head    = "HEAD"
	Options = "OPTIONS"
	Patch   = "PATCH"
)

// Constants of the HTTP headers
const (
	// Standard HTTP Headers
	Authorization    = "Authorization"
	ContentEncoding  = "Content-Encoding"
	ContentLanguage  = "Content-Language"
	ContentLength    = "Content-Length"
	ContentMD5       = "Content-Md5"
	ContentRange     = "Content-Range"
	ContentType      = "Content-Type"
	Date             = "Date"
	Expires          = "Expires"
	Host             = "Host"
	LastModified     = "Last-Modified"
	Location         = "Location"
	Server           = "Server"
	TransferEncoding = "Transfer-Encoding"
	UserAgent        = "User-Agent"

	// BCE Common HTTP Headers
	RequestID        = "Request-ID"
	RequestTimeoutMS = "Request-Timeout-MS"
)
