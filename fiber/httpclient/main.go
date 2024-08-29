package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/gofiber/fiber/v3/client"
	"github.com/valyala/fasthttp"
)

// resty
/**
client := resty.New()
client.GetClient().Transport = &http.Transport{
	DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
		fmt.Println(network, addr)
		return net.Dial(network, "127.0.0.1:9090")
	},
}
*/

func main() {
	c := client.New().SetDial(func(addr string) (net.Conn, error) {
		fmt.Println("dial raw addr", addr)
		addr = "127.0.0.1:9091" // proxyman addr
		return fasthttp.Dial(addr)
	})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		resp, err := c.SetHeader("X-Client", "gogo").
			Get("http://baidu.com:123", client.Config{Header: map[string]string{"Host": "github.com"}}) // host 不会生效，会被 URL 中的地址覆盖掉
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp)
	}()

	go func() {
		defer wg.Done()
		resp, err := c.Get("http://google.com")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp)
	}()
	wg.Wait()

	// 自己创建 request 来修改 host
	req := client.AcquireRequest()
	defer client.ReleaseRequest(req)

	req.RawRequest.UseHostHeader = true // 没有这个就 Host 就会被重新覆盖
	req.SetHeader("Host", "baidu.com:110")
	req.SetClient(c)
	resp, err := req.Get("http://baidu.com:321")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()
	fmt.Println(resp.String())
}
