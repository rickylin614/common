package cerrs

type ErrCode int

// ErrorCode
const (
	SUC            ErrCode = iota // 0 成功
	ERR_TOKEN_MISS                // 1 TOKEN遺失
	ERR_ACCOUNT                   // 2 帳號錯誤
	ERR_PWD                       // 3 密碼錯誤
)

func (e ErrCode) Int() int {
	return int(e)
}
