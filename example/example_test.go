package main

import (
	"geew"
	"net/http"
	"testing"
)

const port = ":65533"

// Example for Use
func TestUse(t *testing.T) {
	r := geew.New()

	r.GET("/", func(c *geew.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
	})

	// expect /hello?name=geew
	r.GET("/hello", func(c *geew.Context) {
		c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
	})

	// expect /hello/geew
	// restful api
	r.GET("/hello/:name", func(c *geew.Context) {
		c.String(http.StatusOK, "hello %s, you`re at %s\nn", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *geew.Context) {
		c.JSON(http.StatusOK, geew.H{
			"filepath": c.Param("filepath"),
		})
	})

	r.POST("/login", func(c *geew.Context) {
		c.JSON(http.StatusOK, geew.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(port)
}

// Example for router group
func TestGroup(t *testing.T) {
	r := geew.New()
	r.GET("/", func(c *geew.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *geew.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
		})

		// expect /hello?name=geew
		v1.GET("/hello", func(c *geew.Context) {
			c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		// expect /hello/geew
		// restful api
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

	r.Run(port)
}

// Testing for middleware(TimerOperation, Recovery)
func TestRecovery(t *testing.T) {
	r := geew.New()
	r.Use(geew.TimerOperation(), geew.Recovery())

	r.GET("/", func(c *geew.Context) {
		c.String(http.StatusOK, "Hello Geew")
	})

	r.GET("/panic", func(c *geew.Context) {
		ns := []string{"geew"}
		c.String(http.StatusOK, ns[2])
	})

	r.Run(port)
}
