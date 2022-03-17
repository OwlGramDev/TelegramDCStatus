package client_telegram

import (
	"TelegramServerChecker/client_raw"
	"TelegramServerChecker/types"
)

type Client struct {
	Client       *client_raw.Context
	FilesDC      []types.TelegramDCInfo
	StatusDC     []types.TelegramDCStatus
	LastRefresh  int64
	IsRefreshing bool
}
