package types

import (
	"TelegramServerChecker/client"
)

type TelegramCheckerClient struct {
	Client       *client.Context
	FilesDC      []TelegramDCInfo
	StatusDC     []TelegramDCStatus
	LastRefresh  int64
	IsRefreshing bool
}
