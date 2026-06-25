package dto

type JWTResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type LoginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
