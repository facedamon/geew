package geew

import (
	"log"
	"time"
)

// calc the time of func
func TimerOperation() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v ", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
