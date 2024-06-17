package do

type ShortCodeError struct {
	Err  error
	Code int
}

func (e *ShortCodeError) IsBizErr() bool {
	return 500 != e.Code
}

func (e *ShortCodeError) IsNotFindErr() bool {
	return 404 == e.Code
}

func (e *ShortCodeError) IsSysErr() bool {
	return 500 == e.Code
}

func (e *ShortCodeError) GetMsg() string {
	if e.Err == nil {
		return "未知的错误信息"
	}
	return e.Err.Error()
}
