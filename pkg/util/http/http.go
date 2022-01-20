package http

import (
	"net"
	"net/http"
	"time"
)

// NewKeepAliveHttpClient 创建http client，设置长连接参数，适配http client keep alive
// 可以设置连接池的http client 类似java的ok http
func NewKeepAliveHttpClient() *http.Client {
	var netTransport = &http.Transport{
		DialTLSContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 90 * time.Second,
		}).DialContext,
		// 最大空闲连接数
		MaxIdleConns: 100,
		// 最大持久连接数，超过了会新建，默认为2
		MaxIdleConnsPerHost: 150,
		// tls 握手超时时间
		TLSHandshakeTimeout: 10 * time.Second,
		// 空闲连接超时时间
		IdleConnTimeout: 90 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   15 * time.Second,
		Transport: netTransport,
	}
	return netClient
}
