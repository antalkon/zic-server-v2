package req

type UpdateGeneralSettingsReq struct {
	Language            *string `json:"language"`
	Timezone            *string `json:"timezone"`
	ServerName          *string `json:"server_name"`
	ServerType          *string `json:"server_type"`
	ServerURL           *string `json:"server_url"`
	ServerAddress       *string `json:"server_address"`
	ServerPhone         *string `json:"server_phone"`
	ServerEmail         *string `json:"server_email"`
	ServerContactPerson *string `json:"server_contact_person"`
}

type UpdateTelegramSettingsReq struct {
	Token           *string  `json:"token"`
	Timezone        *string  `json:"timezone"`
	AdminIDs        *[]int64 `json:"admin_ids"`
	TeacherIDs      *[]int64 `json:"teacher_ids"`
	MessageStart    *string  `json:"message_start"`
	MessageHelp     *string  `json:"message_help"`
	MessageSettings *string  `json:"message_settings"`
}

type UpdateApiSettingsReq struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

type UpdateLicenseSettingsReq struct {
	Token string `json:"token"`
}
