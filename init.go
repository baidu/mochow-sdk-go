// init.go - just import the sub packages

// Package sdk imports all sub packages to build all of them when calling `go install', `go build'
// or `go get' commands.
package sdk

import (
	_ "github.com/baidu/mochow-sdk-go/auth"       // register auth package
	_ "github.com/baidu/mochow-sdk-go/client"     // register client package
	_ "github.com/baidu/mochow-sdk-go/http"       // register http package
	_ "github.com/baidu/mochow-sdk-go/mochow"     // register mochow package
	_ "github.com/baidu/mochow-sdk-go/mochow/api" // register api package
	_ "github.com/baidu/mochow-sdk-go/util"       // register util package
	_ "github.com/baidu/mochow-sdk-go/util/log"   // register log package
)
