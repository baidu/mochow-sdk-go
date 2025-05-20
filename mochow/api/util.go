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

// util.go - define the utilities for api package of Mochow service

package api

const (
	URIPrefixV1 = "/v1"

	RequestDatabaseURI = "/database"
	RequestTableURI    = "/table"
	RequestIndexURI    = "/index"
	RequestRowURI      = "/row"
	RequestRoleURI     = "/role"
	RequestUserURI     = "/user"
)

func getDatabaseURI() string {
	return URIPrefixV1 + RequestDatabaseURI
}

func getTableURI() string {
	return URIPrefixV1 + RequestTableURI
}

func getIndexURI() string {
	return URIPrefixV1 + RequestIndexURI
}

func getRowURI() string {
	return URIPrefixV1 + RequestRowURI
}

func getRoleURI() string {
	return URIPrefixV1 + RequestRoleURI
}

func getUserURI() string {
	return URIPrefixV1 + RequestUserURI
}
