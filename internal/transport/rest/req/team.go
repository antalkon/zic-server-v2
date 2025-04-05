package req

type UpdateRoleReq struct {
	Name *string `json:"name,omitempty"` // указатель = поле может отсутствовать
	Desc *string `json:"desc,omitempty"`
}
