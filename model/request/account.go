package request

/*
定义账户相关的请求参数结构体
*/

type ParamCreateAccount struct {
	Name      string `json:"name" binding:"required,gte=1,lte=20"`       // 账户名（唯一）
	Gender    string `json:"gender" binding:"required,oneof= 男 女 未知"`    // 性别
	Signature string `json:"signature" binding:"required,gte=0,lte=100"` // 签名
}
