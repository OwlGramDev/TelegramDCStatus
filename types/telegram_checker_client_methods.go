package types

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"TelegramServerChecker/api_client"
	"TelegramServerChecker/consts"
	"TelegramServerChecker/telegramInfo"
	"github.com/go-ping/ping"
)

func (tg *TelegramCheckerClient) readBackup() {
	r, err := os.ReadFile(consts.BackupFolder)
	if err == nil {
		var recovery []TelegramDCStatus
		_ = json.Unmarshal(r, &recovery)
		tg.StatusDC = recovery
	}
}

func (tg *TelegramCheckerClient) doBackup() {
	r, _ := json.Marshal(tg.StatusDC)
	_ = os.WriteFile(consts.BackupFolder, r, 0644)
}

func (tg *TelegramCheckerClient) runDownloadWithTimeout(fileId int32) int8 {
	waitChannel := make(chan int8, 1)
	start := time.Now()
	go func() {
		err := tg.Client.DownloadFile(fileId)
		if err != nil {
			waitChannel <- 0
		} else {
			waitChannel <- 1
		}
	}()
	select {
	case res := <-waitChannel:
		if time.Since(start).Seconds() >= 3 {
			return 2
		} else {
			return res
		}
	case <-time.After(10 * time.Second):
		tg.Client.CancelDownloadFile(fileId)
		return 0
	}
}

func (tg *TelegramCheckerClient) pingWithTimeout(address string) int64 {
	waitChannel := make(chan int64, 1)
	pingRequest, err := ping.NewPinger(address)
	if err != nil {
		log.Println(err)
		return 0
	}
	go func() {
		pingRequest.SetPrivileged(true)
		pingRequest.Count = 1
		_ = pingRequest.Run()
		waitChannel <- pingRequest.Statistics().AvgRtt.Milliseconds()
	}()
	select {
	case res := <-waitChannel:
		return res
	case <-time.After(2 * time.Second):
		pingRequest.Stop()
		return 2
	}
}

func (tg *TelegramCheckerClient) Run() {
	tg.readBackup()
	updateFloodWait := int64(60 * 10)
	for {
		tg.IsRefreshing = true
		var listStatus []TelegramDCStatus
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
				listStatus = append(listStatus, TelegramDCStatus{
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
