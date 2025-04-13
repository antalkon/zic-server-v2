package req

type UpdateComputerReq struct {
	Name          *string `json:"name,omitempty"`
	Location      *string `json:"location,omitempty"`
	Description   *string `json:"description,omitempty"`
	PublicIP      *string `json:"public_ip,omitempty"`
	LocalIP       *string `json:"local_ip,omitempty"`
	OS            *string `json:"os,omitempty"`
	ClientVersion *string `json:"client_version,omitempty"`
	Comment       *string `json:"comment,omitempty"`
}
