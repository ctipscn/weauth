package response

type BaseResponse struct {
	Code int    `json:"code"` // 业务自定义错误码
	Data any    `json:"data"` // 返回的具体内容
	Msg  string `json:"msg"`  // 提示信息
}

// 成功响应
func SuccessBaseResponse(data any) BaseResponse {
	return BaseResponse{
		Code: 0,
		Data: data,
		Msg:  "success",
	}
}

// 失败响应
func ErrorBaseResponse(httpStatus int, code int, msg string) BaseResponse {
	return BaseResponse{
		Code: code,
		Data: nil,
		Msg:  msg,
	}
}
