package req

type SendRebootReq struct {
	ComputerID string `json:"computer_id" validate:"required"`
	Delay      int    `json:"delay"`
}
