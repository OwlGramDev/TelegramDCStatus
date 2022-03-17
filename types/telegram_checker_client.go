package types

type TelegramCheckerClient struct {
	Client       *ClientContext
	FilesDC      []TelegramDCInfo
	StatusDC     []TelegramDCStatus
	LastRefresh  int64
	IsRefreshing bool
}
