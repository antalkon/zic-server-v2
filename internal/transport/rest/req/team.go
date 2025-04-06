package req

type UpdateRoleReq struct {
	Name *string `json:"name,omitempty"` // указатель = поле может отсутствовать
	Desc *string `json:"desc,omitempty"`
}
type UpdateUserReq struct {
	Name    *string `json:"name,omitempty"` // указатель = поле может отсутствовать
	Surname *string `json:"surname,omitempty"`
	Email   *string `json:"email,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Role    *string `json:"role,omitempty"`
	Desc    *string `json:"desc,omitempty"`
}

type UpdatePasswordReq struct {
	NewPassword string `json:"new_password"`
}
