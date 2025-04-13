package req

type UpdateRoomRequest struct {
	Name        *string `json:"name,omitempty"`
	Number      *int    `json:"number,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageId     *string `json:"image_id,omitempty"`
}
