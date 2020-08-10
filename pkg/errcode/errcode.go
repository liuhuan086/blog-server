package errcode

import (
	"fmt"
	"net/http"
)

// 在编写错误处理公共方法过程中，声明了error结构体，用于表示错误的响应结果。
// 然后把codes作为全局错误码的储存载体，以便查看当前的注册情况
// 最后在调用NewError创建新的Error时，进行重排校验
type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"detail"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d，错误信息: %s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) MsgF(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	e.details = []string{}
	for _, d := range details {
		e.details = append(e.details, d)
	}
	return e
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		// 使用http中的statusCode方法，针对一些特定错误码进行状态骂转换
		// 因为不同的内部错误码在HTTP状态码中表示不同的含义，所以将其区分开
		// 以便客户端及监控或报警等系统的识别和监听。
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
