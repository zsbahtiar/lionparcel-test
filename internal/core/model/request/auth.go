package request

type Register struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required|min=8"`
	Email    string `json:"email" validate:"required|email"`
	IsAdmin  bool   `json:"is_admin"`
}

type Login struct {
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required|min=8"`
}
