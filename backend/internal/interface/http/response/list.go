package response

// ListResponse 是通用的列表响应结构体。
type ListResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    ListData `json:"data"`
}

// ListData 是列表响应中的数据部分。
type ListData struct {
	Items      any `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// SuccessListResponse 构建一个成功的列表响应对象。
func SuccessListResponse(items any, page, pageSize, total int) ListResponse {
	totalPages := 0
	if pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	return ListResponse{
		Code:    200,
		Message: "success",
		Data: ListData{
			Items:      items,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// ErrorListResponse 构建一个错误的列表响应对象。
func ErrorListResponse(code int, message string) ListResponse {
	return ListResponse{
		Code:    code,
		Message: message,
		Data:    ListData{},
	}
}
