package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TransformRule 转换规则
type TransformRule struct {
	// 匹配条件
	Service    string            // 服务名
	Path       string            // 路径
	Method     string            // HTTP方法
	Headers    map[string]string // 请求头匹配
	QueryParam map[string]string // 查询参数匹配

	// 转换操作
	RequestTransform  RequestTransformFunc  // 请求转换函数
	ResponseTransform ResponseTransformFunc // 响应转换函数
}

// RequestTransformFunc 请求转换函数
type RequestTransformFunc func(*http.Request) error

// ResponseTransformFunc 响应转换函数
type ResponseTransformFunc func(*http.Response) error

// Transformer 请求/响应转换器
type Transformer struct {
	rules []*TransformRule
}

// NewTransformer 创建新的转换器
func NewTransformer() *Transformer {
	return &Transformer{
		rules: make([]*TransformRule, 0),
	}
}

// AddRule 添加转换规则
func (t *Transformer) AddRule(rule *TransformRule) {
	t.rules = append(t.rules, rule)
}

// TransformRequest 转换请求
func (t *Transformer) TransformRequest(req *http.Request, service string) error {
	// 查找匹配的规则
	for _, rule := range t.rules {
		if rule.matches(req, service) && rule.RequestTransform != nil {
			if err := rule.RequestTransform(req); err != nil {
				return err
			}
		}
	}
	return nil
}

// TransformResponse 转换响应
func (t *Transformer) TransformResponse(resp *http.Response, service string) error {
	// 查找匹配的规则
	for _, rule := range t.rules {
		if rule.matches(resp.Request, service) && rule.ResponseTransform != nil {
			if err := rule.ResponseTransform(resp); err != nil {
				return err
			}
		}
	}
	return nil
}

// matches 检查是否匹配规则
func (r *TransformRule) matches(req *http.Request, service string) bool {
	// 检查服务名
	if r.Service != "" && r.Service != service {
		return false
	}

	// 检查路径
	if r.Path != "" && r.Path != req.URL.Path {
		return false
	}

	// 检查方法
	if r.Method != "" && r.Method != req.Method {
		return false
	}

	// 检查请求头
	for k, v := range r.Headers {
		if req.Header.Get(k) != v {
			return false
		}
	}

	// 检查查询参数
	for k, v := range r.QueryParam {
		if req.URL.Query().Get(k) != v {
			return false
		}
	}

	return true
}

// Common transform functions

// JSONKeyRename 重命名JSON字段
func JSONKeyRename(oldKey, newKey string) ResponseTransformFunc {
	return func(resp *http.Response) error {
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		// 解析JSON
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}

		// 重命名字段
		if val, ok := data[oldKey]; ok {
			data[newKey] = val
			delete(data, oldKey)
		}

		// 重新编码
		newBody, err := json.Marshal(data)
		if err != nil {
			return err
		}

		// 更新响应
		resp.Body = io.NopCloser(bytes.NewReader(newBody))
		resp.ContentLength = int64(len(newBody))
		resp.Header.Set("Content-Length", fmt.Sprint(len(newBody)))

		return nil
	}
}

// AddRequestHeader 添加请求头
func AddRequestHeader(key, value string) RequestTransformFunc {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

// AddResponseHeader 添加响应头
func AddResponseHeader(key, value string) ResponseTransformFunc {
	return func(resp *http.Response) error {
		resp.Header.Set(key, value)
		return nil
	}
}
