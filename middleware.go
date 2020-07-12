package geew

import (
	"fmt"
	"log"
	"runtime"
	"strings"
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

// recovery panic
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
			}
		}()
		c.Next()
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	// skip first 3 caller
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
