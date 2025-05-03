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
