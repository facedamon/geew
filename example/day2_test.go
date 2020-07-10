package example

import (
	"geew"
	"net/http"
	"testing"
)

func TestDay2(t *testing.T) {
	r := geew.New()
	r.GET("/", func(c *geew.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
	})
	r.GET("/hello", func(c *geew.Context) {
		c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *geew.Context) {
		c.JSON(http.StatusOK, geew.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
