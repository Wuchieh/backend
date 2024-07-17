package pkg

import (
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	SUCCESS = "success"
)

// CreateResponse 創建 Response 結構
//
//	code 狀態碼
//	message 訊息
//	data 資料[可選]
func CreateResponse(code int, message string, data ...any) (int, *Response) {
	return code, &Response{
		Code:    code,
		Message: message,
		Data: func() any {
			if len(data) == 1 {
				return data[0]
			} else if len(data) > 1 {
				return data
			} else {
				return nil
			}
		}(),
	}
}

func CreateSuccessResponse(data ...interface{}) (int, *Response) {
	return CreateResponse(http.StatusOK, SUCCESS, data...)
}

func CreateSuccessResponseObj(data ...interface{}) *Response {
	_, obj := CreateSuccessResponse(data...)
	return obj
}
