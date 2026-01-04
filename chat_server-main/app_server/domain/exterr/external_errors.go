package exterr

type ExternalError struct {
	Code  int32
	Msg   string
	CnMsg string
}

func (e ExternalError) Error() string {
	return e.Msg
}

func NewExternalError(code int32, msg, cnMsg string) ExternalError {
	return ExternalError{
		Code:  code,
		Msg:   msg,
		CnMsg: cnMsg,
	}
}

var (
	ErrInternal = NewExternalError(10000, "internal error", "内部错误")

	ErrParamMissing = NewExternalError(20001, "param missing", "必填参数缺失")
	ErrParamInvalid = NewExternalError(20002, "param invalid", "必填参数无效")
)
