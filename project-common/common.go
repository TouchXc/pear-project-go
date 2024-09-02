package common

// 基础序列化器  即接口返回的json结构体

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (rsp *Response) Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "Success",
		Data: data,
	}
}
func (rsp *Response) Failed(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
