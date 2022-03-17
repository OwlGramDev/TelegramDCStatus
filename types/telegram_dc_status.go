package types

type TelegramDCStatus struct {
	Id       int8  `json:"dc_id"`
	Ping     int64 `json:"ping"`
	Status   int8  `json:"dc_status"`
	LastDown int64 `json:"last_down"`
	LastLag  int64 `json:"last_lag"`
}
