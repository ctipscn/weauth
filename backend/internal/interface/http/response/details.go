package response

// FieldError 表示单个字段的校验错误详情。
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// DetailsResponse 是通用的带详情响应结构体，常用于参数校验失败等场景。
type DetailsResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data"`
	Details []FieldError `json:"details"`
}

// SuccessDetailsResponse 构建一个成功的带详情响应对象。
func SuccessDetailsResponse(data any) DetailsResponse {
	return DetailsResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Details: []FieldError{},
	}
}

// ErrorDetailsResponse 构建一个错误的带详情响应对象，携带字段级错误列表。
func ErrorDetailsResponse(code int, message string, details []FieldError) DetailsResponse {
	if details == nil {
		details = []FieldError{}
	}
	return DetailsResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		Details: details,
	}
}
