package hexerror

import (
	"encoding/json"
	"strings"
)

var _source string

func SetSource(source string) {
	_source = source
}


type HexError interface {
	error
	Code() string
	Description() string
	Source() string
	Detail() interface{}
}

type HexErrorImpl struct {
	HexCode        string      `json:"code"`
	HexDescription string      `json:"description"`
	HexSource      string      `json:"source"`
	HexDetail      interface{} `json:"detail,omitempty"`
}

func (e *HexErrorImpl) Code() string {
	return e.HexCode
}

func (e *HexErrorImpl) Source() string {
	return e.HexSource
}

func (e *HexErrorImpl) Description() string {
	return e.HexDescription
}

func (e *HexErrorImpl) Detail() interface{} {
	return e.HexDetail
}

func (e *HexErrorImpl) Error() string {
	data, _ := json.Marshal(e)
	return  string(data)
}

func NewHexError(code string, description string, detail ...interface{}) HexError {
	var _detail interface{}
	if len(detail) != 0 {
		_detail = detail[0]
	}
	return &HexErrorImpl{
		HexCode:        code,
		HexDescription: description,
		HexSource:      _source,
		HexDetail:      _detail,
	}
}

func NotFound(description string, detail ...interface{}) HexError {
	return NewHexError(notFoundCode, description, detail...)
}

func Unauthorized(description string, detail ...interface{}) HexError {
	return NewHexError(unauthorizedCode, description, detail...)
}

func InvalidAction(description string, detail ...interface{}) HexError {
	return NewHexError(invalidActionCode, description, detail...)
}

func DBError(description string, detail ...interface{}) HexError {
	return NewHexError(notFoundCode, description, detail...)
}

func InvalidData(description string, detail ...interface{}) HexError {
	return NewHexError(invalidDataCode, description, detail...)
}

func WrongPassword(description string, detail ...interface{}) HexError {
	return NewHexError(wrongPasswordCode, description, detail...)
}

func PermissionDenied(description string, detail ...interface{}) HexError {
	return NewHexError(permissionDeniedCode, description, detail...)
}

func WrongPartner(description string, detail ...interface{}) HexError {
	return NewHexError(wrongPartnerCode, description, detail...)
}

func Unknown(description string, detail ...interface{}) HexError {
	return NewHexError(unknownCode, description, detail...)
}

func HexErrorFromString(errStr string) HexError {
	if strings.Contains(errStr, "{") {
		index := strings.Index(errStr, "{")
		// TODO: 检查使用这种方法处理unicode是否有问题
		errStr = errStr[index:]
	}
	var hexerror HexErrorImpl
	err := json.Unmarshal([]byte(errStr), &hexerror)
	if err != nil {
		return Unknown("umarshal error fail: " + errStr)
	}
	return &hexerror
}

func getStringBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}