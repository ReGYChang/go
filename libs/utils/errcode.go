package utils

const (
	Success      = 200   // success
	ReqError     = 10001 // 連線參數錯誤
	ListFailed   = 10002 // 取得清單失敗
	ListNotFound = 10003 // 無清單列表
	CreateFailed = 10004 // 建立失敗
	UpdateFailed = 10005 // 更新失敗
	DelFailed    = 10006 // 刪除失敗
	InfoFailed   = 10007 // 取得資訊失敗  // 查單筆資料
	InfoNotFound = 10008 // 查無項目資訊
)

var errorDescription = map[int]string{
	Success:      "Success",
	ReqError:     "incorrect data params",
	ListFailed:   "ListFailed",
	ListNotFound: "ListNotFound",
	CreateFailed: "CreateFailed",
	UpdateFailed: "UpdateFailed",
	DelFailed:    "DelFailed",
	InfoFailed:   "InfoFailed",
	InfoNotFound: "InfoNotFound",
}

func ErrorText(code int) string {
	return errorDescription[code]
}
