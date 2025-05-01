package req

// step 0: ACTIVATE LICENSE

// Step 1
type FristStartForm struct {
	Name          string `json:"name" validate:"required"`
	Type          string `json:"type" validate:"required"`
	Url           string `json:"url" validate:"required"`
	Address       string `json:"address" validate:"required"`
	Phone         string `json:"phone" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	ContactPerson string `json:"contact_person" validate:"required"`
}

// Step 2: API (zpi.zic.zentas.ru)
type FristStartAPI struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

// Step 3: Change admin password\
type CgangeAdminPassword struct {
	NewPassword string `json:"new_password" validate:"required"`
}
