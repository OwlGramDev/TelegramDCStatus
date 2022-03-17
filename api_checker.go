package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"TelegramServerChecker/consts"
	"TelegramServerChecker/telegramInfo"
	"TelegramServerChecker/types"
	tdLib "github.com/Arman92/go-tdlib"
	"github.com/go-ping/ping"
	"github.com/valyala/fasthttp"
)

func TelegramServerChecker() *types.TelegramCheckerClient {
	instance := Client()
	instance.Login()
	var listDCInfo []types.TelegramDCInfo
	var listStatus []types.TelegramDCStatus
	messageList := instance.GetMessageList()
	for i := 0; i < len(messageList); i++ {
		message := messageList[i].Content
		if message.GetMessageContentEnum() == "messageAnimation" {
			file := message.(*tdLib.MessageAnimation)
			dcIDTmp := strings.ReplaceAll(file.Animation.FileName, "st-", "")
			dcIDTmp = strings.ReplaceAll(dcIDTmp, ".gif.mp4", "")
			dcID, _ := strconv.Atoi(dcIDTmp)
			listDCInfo = append(listDCInfo, types.TelegramDCInfo{
				ID:     int8(dcID),
				FileID: file.Animation.Animation.ID,
			})
			listStatus = append(listStatus, types.TelegramDCStatus{
				Id:     int8(dcID),
				Status: -1,
			})
		}
	}
	fmt.Println("\nStarted Telegram DC Checker!")
	return &types.TelegramCheckerClient{
		Client:       instance,
		FilesDC:      listDCInfo,
		StatusDC:     listStatus,
		IsRefreshing: true,
	}
}

func (tg *types.TelegramCheckerClient) readBackup() {
	r, err := os.ReadFile(consts.backupFolder)
	if err == nil {
		var recovery []types.TelegramDCStatus
		_ = json.Unmarshal(r, &recovery)
		tg.StatusDC = recovery
	}
}

func (tg *types.TelegramCheckerClient) doBackup() {
	r, _ := json.Marshal(tg.StatusDC)
	_ = os.WriteFile(consts.backupFolder, r, 0644)
}

func (tg *types.TelegramCheckerClient) runDownloadWithTimeout(fileId int32) int8 {
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

func (tg *types.TelegramCheckerClient) pingWithTimeout(address string) int64 {
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

func (tg *types.TelegramCheckerClient) Run() {
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

				_ = os.Remove(fmt.Sprintf("%s/td_files/animations/st-%d.gif.mp4", consts.tdSessionFiles, tg.FilesDC[i].ID))
			} else {
				listStatus = append(listStatus, tg.StatusDC[i])
			}
		}
		tg.StatusDC = listStatus
		tg.IsRefreshing = false
		t = time.Now()
		tg.LastRefresh = t.Unix()
		tg.doBackup()
		SendData(listStatus)
		time.Sleep(time.Second * time.Duration(60-t.Second()))
	}
}

func SendData(data []types.TelegramDCStatus) string {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(consts.apiEndpoint + "sendData")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	marshal, _ := json.Marshal(data)
	req.SetBody(marshal)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return ""
	}

	statusCode := resp.StatusCode()
	if statusCode != fasthttp.StatusOK {
		if statusCode == fasthttp.StatusTooManyRequests {
			log.Println("You have been banned from api.owlgram.org because you are flooding, wait and retry after 30 minutes")
		}
		return ""
	}

	return string(resp.Body())
}
