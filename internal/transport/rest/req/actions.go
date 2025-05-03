package req

type SendRebootReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
	Delay      int    `json:"delay"`
}

type SendShutdownReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
	Delay      int    `json:"delay"`
}

type SendBlockReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
}

type SendUnblockReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
}

type SendLockScreenReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
}

type SendUrlReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
	Url        string `json:"url" validate:"required,url"`
}

type SendMessageReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
	Message    string `json:"message" validate:"required"`
	Type       string `json:"type" validate:"required"`
}
