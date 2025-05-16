# 百度向量数据库 Mochow GO SDK

针对百度智能云向量数据库，我们推出了一套 GO SDK（下称Mochow SDK），方便用户通过代码调用百度向量数据库。

## 如何安装

直接从github下载，使用`go get`工具从github进行下载：
```
go get github.com/baidu/mochow-sdk-go/v2
```
目前GO SDK可以在go1.17及以上环境下运行。

## 快速使用

在使用Mochow SDK 之前，用户需要在百度智能云上创建向量数据库，以获得 API Key。API Key 是用户在调用Mochow SDK 时所需要的凭证。具体获取流程参见平台的[向量数据库使用说明文档](https://cloud.baidu.com/)。

获取到 API Key 后，用户还需要传递它们来初始化Mochow SDK。 可以通过如下方式初始化Mochow SDK：

```go
package main

import (
	"fmt"

	"github.com/baidu/mochow-sdk-go/v2/client"
	"github.com/baidu/mochow-sdk-go/v2/mochow"
)

func main() {
	clientConfig := &mochow.ClientConfiguration{
		Account:  "root",
		APIKey:   "your_api_key",
		Endpoint: "you_endpoint",   // example: http://127.0.0.1:8511
	}

	// 创建Mochow服务的Client
	client, err := mochow.NewClientWithConfig(clientConfig)
	if err != nil {
		fmt.Println("create client failed")
		// handle error
	}

	// 创建database
	err = client.CreateDatabase("your_database_name")
	if err != nil {
		fmt.Println("create database failed:", err)
		// handle error
	}
	fmt.Println("create database success")
}
```

## 功能

目前Mochow SDK 支持用户使用如下功能:

+ Databse 操作
+ Table 操作
+ Alias 操作
+ Index 操作
+ Row 操作

## License

Apache-2.0

