package req

type UpdateUserDataReq struct {
	FirstName *string `json:"name"` // теперь опционально
	SurName   *string `json:"surname"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
}

type UpdateDataPasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
