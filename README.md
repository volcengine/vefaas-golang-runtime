# veFaaS Golang 运行时

用于开发 veFaaS Golang 函数的 SDK 及代码样例，更多信息请参考 [Golang 函数开发指南](https://www.volcengine.com/docs/6662/106051)。

## 快速上手
```Golang
// main.go
package main

import (
	"context"
	"fmt"

	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/vefaas"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
)

func main() {
	// Start your vefaas function =D.
	vefaas.Start(handler)
}

// Define your handler function.
func handler(ctx context.Context, r *events.HTTPRequest) (*events.EventResponse, error) {
	fmt.Printf("request id: %s", vefaascontext.RequestIdFromContext(ctx))
	fmt.Printf("request headers: %v", r.Headers)

	return &events.EventResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte("Hello veFaaS!"),
	}, nil
}
```

## 构建 veFaaS Golang 函数
对于 veFaaS Golang 运行时函数，在部署之前，开发者需要将其在本地构建为可在 Linux 环境 amd64 架构下执行的、命名为 main 的二进制文件，并将其打包在 .zip 文件中。

### 对于 Linux/macOS 开发者
可通过如下 shell 指令对程序进行编译和打包：
```shell
# Build your program that's executable for Linux under architecture amd64.
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main
# Zip output.
zip main.zip main
```
### 对于 Windows 开发者

1. 通过快捷键 `Windows + R` 唤出命令窗口，输入 `wt` 并键入 `Enter` 打开命令行终端。通过以下命令执行编译：
```shell
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o main
```
2. 使用打包工具对编译好的二进制 `main` 进行打包，`main` 需要位于 .zip 文件根目录下。


## 发布 veFaaS Golang 函数
对于 Golang 这类编译型函数，veFaaS 暂不支持通过控制台在线编辑函数进而构建和发布，需要开发者在本地进行代码开发，按照 [构建 veFaaS Golang 函数](#构建-vefaas-golang-函数) 中的步骤对函数编译和打包，然后将打包好的 .zip 文件上传至 veFaaS 控制台，进行函数发布。


## 注意事项
- 编译后的二进制文件须为命名为 main 的、可在 Linux 环境 amd64 架构下执行的二进制文件
- 打包后的 main 须位于 .zip 文件的根目录下，不能嵌套在其它文件夹下
- 如果主程序的执行依赖其它本地文件，如 config.yaml，在打包 .zip 文件时，这些文件依赖要一同打包，代码中通过相对位置来进行引用
