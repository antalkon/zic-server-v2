package req

type SignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}
