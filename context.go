package geew

import (
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"sync"
)

type H map[string]interface{}

// handlers length  boundary
const abortIndex int8 = math.MaxInt8 / 2

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
	handlers HandlersChain
	index    int8
	// mutex protect keys map
	mu sync.RWMutex
	// queryCache use url.ParseQuery cached the param query result from c.Req.URL.Query
	queryCache url.Values
	//formCache use url.ParseQuery cached PostForm contains the parsed form data from POST,
	//PUT body parameters.
	formCache url.Values
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		index:      -1,
		queryCache: nil,
		formCache:  nil,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < int8(s); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

// IsAbort returns true if the current context was aborted.
func (c *Context) IsAbort() bool {
	return c.index >= abortIndex
}

func (c *Context) AbortWithStatusJSON(code int, o interface{}) {
	c.Abort()
	c.JSON(code, o)
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		c.queryCache = c.Req.URL.Query()
	}
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// QueryArray returns a slice of strings for given query key
// the length od the slice depends on the number of params with the given key
func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		r := c.Req
		if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
			if err != http.ErrNotMultipart {
				L.Error("error on parse multipart form array: %v", err)
			}
		}
		c.formCache = r.PostForm
	}
}

// GetPostFormArray returns a slice of strings for a given form key
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.initFormCache()
	if values := c.formCache[key]; len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// PostFormArray returns a slice of strings for a given form key
// the length of slice depends on the number of params with the given key
func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

// PostForm returns the specified key from a POST urlencoded from or multipart form
func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
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

// nameOfFunction return func name
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// HandlerName returns the main handler`s name.
// For example if the handler is 'handlerGetUsers()'
// this function will return 'main.handlerGetUsers'.
func (c *Context) HandlerName() string {
	return nameOfFunction(c.handlers.Last())
}
