package example

import (
	"geew"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	r := geew.New()
	r.GET("/", func(c *geew.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
	})

	r.GET("/hello", func(c *geew.Context) {
		c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
	})
	// Resutful api
	r.GET("/hello/:name", func(c *geew.Context) {
		c.String(http.StatusOK, "Hello %s, you`re at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *geew.Context) {
		c.JSON(http.StatusOK, geew.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
