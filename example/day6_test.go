package main

import (
	"geew"
	"net/http"
	"testing"
)

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

	r.Run(":9999")

}
