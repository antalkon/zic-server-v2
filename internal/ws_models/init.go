package wsmodels

type InitPayload struct {
	ComputerID    string `json:"computer_id"`
	OS            string `json:"os"`
	ClientVersion string `json:"client_version"`
	PublicIP      string `json:"public_ip"`
	JwtToken      string `json:"jwt"`
	LocalIP       string `json:"local_ip"`
}
