package example

import (
	"geew"
	"net/http"
	"testing"
)

func TestGroup(t *testing.T) {
	r := geew.New()
	r.GET("/index", func(c *geew.Context) {
		c.HTML(http.StatusOK, "<h1>old router supported</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *geew.Context) {
			c.HTML(http.StatusOK, "<h1>Hello group v1</h1>")
		})
		v1.GET("/hello", func(c *geew.Context) {
			c.String(http.StatusOK, "hello %s, you`re at %s\nn", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		// Restful api
		v2.GET("/hello/:name", func(c *geew.Context) {
			c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *geew.Context) {
			c.JSON(http.StatusOK, geew.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":9999")
}
