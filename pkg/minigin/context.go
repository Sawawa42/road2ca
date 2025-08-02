package minigin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Context struct {
	Writer   *ResponseWriter
	Request  *http.Request
	handlers []HandlerFunc
	index    int
}

type ResponseWriter struct {
	Writer http.ResponseWriter
	status int
	size   int
}

// Header ラップされたhttp.ResponseWriterのHeaderメソッド
func (w *ResponseWriter) Header() http.Header {
	return w.Writer.Header()
}

// WriteHeader ラップされたhttp.ResponseWriterのWriteHeaderメソッド
func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.status == 0 {
		w.status = statusCode
		w.Writer.WriteHeader(statusCode)
	}
}

// Write ラップされたhttp.ResponseWriterのWriteメソッド
func (w *ResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.Writer.Write(data)
	w.size += size
	return size, err
}

// Status レスポンスのステータスコードを取得
func (w *ResponseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// Next はハンドラチェーンの次のハンドラを呼び出す
func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

// H はJSONレスポンスで使用されるヘルパー型
type H map[string]any

// JSON はレスポンスとしてJSONを返す
func (c *Context) JSON(code int, obj any) {
	json, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Failed to json.Marshal: %v", err)
		c.JSON(http.StatusInternalServerError, H{"error": "Internal server error"})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	c.Writer.Write(json)
}

// QueryInt はクエリパラメータを整数として取得する
func (c *Context) QueryInt(key string) (int, error) {
	values := c.Request.URL.Query()
	if value, ok := values[key]; ok && len(value) > 0 {
		return strconv.Atoi(value[0])
	}
	return 0, fmt.Errorf("query parameter %s not found", key)
}
