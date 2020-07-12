package geew

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware support
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	c.Req.ParseMultipartForm(MaxMultipartMemory)
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	c.Req.ParseMultipartForm(MaxMultipartMemory)
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	c.Req.ParseMultipartForm(MaxMultipartMemory)
	v, _ := c.Params[key]
	return v
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Req.MultipartForm == nil {
		if err := c.Req.ParseMultipartForm(MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Req.ParseMultipartForm(MaxMultipartMemory)
	return c.Req.MultipartForm, err
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
