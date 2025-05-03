package res

type RoomStatusRes struct {
	Status    string `json:"status"`
	PcOnline  int    `json:"pc_online"`
	PcOffline int    `json:"pc_offline"`
	PcTotal   int    `json:"pc_total"`
	Blocked   int    `json:"blocked"`
	Free      int    `json:"free"`
}
