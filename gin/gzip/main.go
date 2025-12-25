package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/ping", func(c *gin.Context) {
		fmt.Println(c.GetHeader("Accept-Encoding"))
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	go func() {
		time.Sleep(time.Second)
		http.Get("http://localhost:8080/ping")
	}()

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
