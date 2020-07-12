package example

import (
	"geew"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestMiddleware(t *testing.T) {
	r := geew.New()
	// global middleware
	r.Use(geew.TimerOperation())
	r.GET("/", func(c *geew.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(func(c *geew.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	})
	{
		v2.GET("/hello/:name", func(c *geew.Context) {
			c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
