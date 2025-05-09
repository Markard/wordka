package login

type Response struct {
	Token string `json:"token"`
}

func NewResponse(tokenString string) *Response {
	return &Response{tokenString}
}
