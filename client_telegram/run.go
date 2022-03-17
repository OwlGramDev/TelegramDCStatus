package client_telegram

import (
	"fmt"
	"log"
	"os"
	"time"

	"TelegramServerChecker/api_client"
	"TelegramServerChecker/consts"
	"TelegramServerChecker/telegramInfo"
	"TelegramServerChecker/types"
)

func (tg *Client) Run() {
	tg.readBackup()
	updateFloodWait := int64(60 * 10)
	for {
		tg.IsRefreshing = true
		var listStatus []types.TelegramDCStatus
		t := time.Now()
		for i := 0; i < len(tg.FilesDC); i++ {
			canUpdate := false
			ipAddress, err := telegramInfo.GetIPFromDC(tg.StatusDC[i].Id)
			if err != nil {
				log.Println("[TelegramDC]", err)
				continue
			}
			pingResult := tg.pingWithTimeout(ipAddress.String())
			if (t.Unix() - tg.StatusDC[i].LastDown) >= updateFloodWait {
				canUpdate = true
			} else if (t.Unix() - tg.StatusDC[i].LastLag) >= updateFloodWait {
				canUpdate = true
			}
			if canUpdate {
				res := tg.runDownloadWithTimeout(tg.FilesDC[i].FileID)
				if pingResult == 2 && res == 1 {
					res = 1
				} else if pingResult == 2 && res == 2 {
					res = 0
				}
				if pingResult == 2 {
					pingResult = int64(2 * time.Second)
					if pingResult >= 1000 {
						res = 0
					} else if pingResult >= 300 {
						res = 2
					}
				}
				if res == 0 {
					tg.StatusDC[i].LastDown = t.Unix()
				} else if res == 2 {
					tg.StatusDC[i].LastLag = t.Unix()
				}
				listStatus = append(listStatus, types.TelegramDCStatus{
					tg.FilesDC[i].ID,
					pingResult,
					res,
					tg.StatusDC[i].LastDown,
					tg.StatusDC[i].LastLag,
				})

				_ = os.Remove(fmt.Sprintf("%s/td_files/animations/st-%d.gif.mp4", consts.TdSessionFiles, tg.FilesDC[i].ID))
			} else {
				listStatus = append(listStatus, tg.StatusDC[i])
			}
		}
		tg.StatusDC = listStatus
		tg.IsRefreshing = false
		t = time.Now()
		tg.LastRefresh = t.Unix()
		tg.doBackup()
		api_client.SendData(listStatus)
		time.Sleep(time.Second * time.Duration(60-t.Second()))
	}
}
