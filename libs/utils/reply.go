package utils

type ErrorResponse struct {
	Error       interface{} `json:"error"`
	Description interface{} `json:"description"`
}

// 給 err 回傳錯誤訊息
func ReplyError(errCode int) ErrorResponse {
	res := ErrorResponse{
		Error:       errCode,
		Description: ErrorText(errCode),
	}
	return res
}

type replySuccess struct {
	Status interface{} `json:"status"`
}

// 只回傳成功狀態
func ReplySuccess() replySuccess {
	res := replySuccess{
		Status: Success,
	}
	return res
}

type RepError struct {
	Error interface{} `json:"error"`
}
