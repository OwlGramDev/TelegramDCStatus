package client_telegram

import (
	"encoding/json"
	"os"

	"TelegramServerChecker/consts"
	"TelegramServerChecker/types"
)

func (tg *Client) readBackup() {
	r, err := os.ReadFile(consts.BackupFolder)
	if err == nil {
		var recovery []types.TelegramDCStatus
		_ = json.Unmarshal(r, &recovery)
		tg.StatusDC = recovery
	}
}

func (tg *Client) doBackup() {
	r, _ := json.Marshal(tg.StatusDC)
	_ = os.WriteFile(consts.BackupFolder, r, 0644)
}
