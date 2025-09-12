package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	resource = "http://site.cdn.internal:8080/index.js"
	all      = "http://site.cdn.internal:8080/8BFE-656DC3564C05"
)

func main() {
	// 要解析的域名、端口和目标 IP
	resolvePort := "8080"
	resolveIP := "127.0.0.1"

	// 1. 创建一个自定义的 http.Transport
	// 这部分的逻辑和标准库的实现完全相同
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout: 3 * time.Second,
			}
			return dialer.DialContext(ctx, network, net.JoinHostPort(resolveIP, resolvePort))
		},
	}

	// 2. 创建一个 resty 客户端
	client := resty.New()

	// 3. 将自定义的 Transport 设置给 resty 客户端. [1]
	client.SetTransport(transport)

	resp, err := client.R().Execute("PURGE", all) // 使用自定义方法 "PURGE" 执行请求

	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	fmt.Println("状态码:", resp.Status())
}
