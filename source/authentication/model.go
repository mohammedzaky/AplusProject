package authentication

//LoginAuthentication check if username and password are vaild or not
type LoginAuthentication struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required,gte=6,lte=13"`
}
