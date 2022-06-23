package errx

import "encoding/json"

const LogicErr = -1
const DbExecuteErr = -2

// Error .
type Error struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// OutputParseString .
func OutputParseString(code int, msg string, data interface{}) string {
	e := Error{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	str, _ := json.Marshal(e)
	return string(str)
}

// ResponseError api err response
func ResponseError(code int, err error) string {
	return OutputParseString(code, err.Error(), struct{}{})
}

// SuccessWithoutData api success without data
func SuccessWithoutData() string {
	return OutputParseString(0, "", struct{}{})
}

// SuccessWithData api success with data
func SuccessWithData(d interface{}) string {
	return OutputParseString(0, "", d)
}
