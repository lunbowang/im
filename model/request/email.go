package request

type ParamSendEmail struct {
	Email string `json:"email" binding:"required,lte=50"`
}
