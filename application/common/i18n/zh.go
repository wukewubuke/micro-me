package i18n

const (
	ErrParam  string = "参数错误!!!"
	ErrServer string = "服务器繁忙,请稍后再试!!!"
	TimeLayOut string = "2006-01-02 15:04:05"
)

var ZhMessage = map[string]string{
	"LoginRequest.Username.required": "手机号不能为空",
	"LoginRequest.Password.required": "密码不能为空",



	"SendRequest.FromToken.required": "Token不能为空",
	"SendRequest.ToToken.required": "对方token不能为空",
	"SendRequest.Body.required": "发送内容不能为空",



}


