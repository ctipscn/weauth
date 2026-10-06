package request

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Meta 请求元数据，包含请求 ID 和时间戳
type Meta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

// BaseRequest 前端请求约定格式
type BaseRequest struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

// BindJSON 将请求体按 BaseRequest 信封格式解析，提取 Data 到指定类型 T 并执行 gin 绑定校验
func BindJSON[T any](c *gin.Context) (Meta, T, error) {
	body, err := c.GetRawData()
	if err != nil {
		return Meta{}, *new(T), err
	}
	// 恢复请求体，不影响下游可能再次调用的逻辑
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var baseReq BaseRequest
	if err := json.Unmarshal(body, &baseReq); err != nil {
		return Meta{}, *new(T), err
	}

	dataBytes, err := json.Marshal(baseReq.Data)
	if err != nil {
		return Meta{}, *new(T), err
	}

	var req T
	if err := json.Unmarshal(dataBytes, &req); err != nil {
		return Meta{}, *new(T), err
	}

	// 执行 gin 结构体校验（binding tag）
	if err := binding.Validator.ValidateStruct(req); err != nil {
		return Meta{}, *new(T), err
	}

	return baseReq.Meta, req, nil
}
