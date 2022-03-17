package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tdLib "github.com/Arman92/go-tdlib"
	"github.com/go-ping/ping"
	"github.com/valyala/fasthttp"
)

func TelegramServerChecker() *TgCheckerClient {
	instance := Client()
	instance.Login()
	var listDCInfo []TelegramDCInfo
	var listStatus []TelegramDCStatus
	messageList := instance.GetMessageList()
	for i := 0; i < len(messageList); i++ {
		message := messageList[i].Content
		if message.GetMessageContentEnum() == "messageAnimation" {
			file := message.(*tdLib.MessageAnimation)
			dcIDTmp := strings.ReplaceAll(file.Animation.FileName, "st-", "")
			dcIDTmp = strings.ReplaceAll(dcIDTmp, ".gif.mp4", "")
			dcID, _ := strconv.Atoi(dcIDTmp)
			listDCInfo = append(listDCInfo, TelegramDCInfo{
				int8(dcID),
				file.Animation.Animation.ID,
			})
			listStatus = append(listStatus, TelegramDCStatus{
				int8(dcID),
				0,
				-1,
				0,
				0,
			})
		}
	}
	fmt.Println("\nStarted Telegram DC Checker!")
	return &TgCheckerClient{
		instance,
		listDCInfo,
		listStatus,
		0,
		true,
	}
}

func (tg *TgCheckerClient) readBackup() {
	r, err := os.ReadFile(backupFolder)
	if err == nil {
		var recovery []TelegramDCStatus
		_ = json.Unmarshal(r, &recovery)
		tg.statusDC = recovery
	}
}

func (tg *TgCheckerClient) doBackup() {
	r, _ := json.Marshal(tg.statusDC)
	_ = os.WriteFile(backupFolder, r, 0644)
}

func (tg *TgCheckerClient) runDownloadWithTimeout(fileId int32) int8 {
	waitChannel := make(chan int8, 1)
	start := time.Now()
	go func() {
		err := tg.client.DownloadFile(fileId)
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
		tg.client.CancelDownloadFile(fileId)
		return 0
	}
}

func (tg *TgCheckerClient) pingWithTimeout(address string) int64 {
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

func getIpById(id int8) string {
	switch id {
	case 1:
		return "149.154.175.50"
	case 2:
		return "149.154.167.50"
	case 3:
		return "149.154.175.100"
	case 4:
		return "149.154.167.91"
	case 5:
		return "91.108.56.100"
	}
	return "???"
}

func (tg *TgCheckerClient) Run() {
	tg.readBackup()
	updateFloodWait := int64(60 * 10)
	for {
		tg.isRefreshing = true
		var listStatus []TelegramDCStatus
		t := time.Now()
		for i := 0; i < len(tg.filesDC); i++ {
			canUpdate := false
			ipAddress := getIpById(tg.statusDC[i].Id)
			pingResult := tg.pingWithTimeout(ipAddress)
			if (t.Unix() - tg.statusDC[i].LastDown) >= updateFloodWait {
				canUpdate = true
			} else if (t.Unix() - tg.statusDC[i].LastLag) >= updateFloodWait {
				canUpdate = true
			}
			if canUpdate {
				res := tg.runDownloadWithTimeout(tg.filesDC[i].fileId)
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
					tg.statusDC[i].LastDown = t.Unix()
				} else if res == 2 {
					tg.statusDC[i].LastLag = t.Unix()
				}
				listStatus = append(listStatus, TelegramDCStatus{
					tg.filesDC[i].id,
					pingResult,
					res,
					tg.statusDC[i].LastDown,
					tg.statusDC[i].LastLag,
				})
				_ = os.Remove(tdSessionFiles + "/td_files/animations/st-" + strconv.Itoa(int(tg.filesDC[i].id)) + ".gif.mp4")
			} else {
				listStatus = append(listStatus, tg.statusDC[i])
			}
		}
		tg.statusDC = listStatus
		tg.isRefreshing = false
		t = time.Now()
		tg.lastRefresh = t.Unix()
		tg.doBackup()
		SendData(listStatus)
		time.Sleep(time.Second * time.Duration(60-t.Second()))
	}
}

func SendData(data []TelegramDCStatus) *string {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(apiEndpoint + "sendData")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	marshal, _ := json.Marshal(data)
	req.SetBody(marshal)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil
	}
	statusCode := resp.StatusCode()
	if statusCode != fasthttp.StatusOK {
		if statusCode == fasthttp.StatusTooManyRequests {
			log.Println("You have been banned from api.owlgram.org because you are flooding, wait and retry after 30 minutes")
		}
		return nil
	}
	body := string(resp.Body())
	return &body
}

type TgCheckerClient struct {
	client       *ClientContext
	filesDC      []TelegramDCInfo
	statusDC     []TelegramDCStatus
	lastRefresh  int64
	isRefreshing bool
}

type TelegramDCInfo struct {
	id     int8
	fileId int32
}

type TelegramDCStatus struct {
	Id       int8  `json:"dc_id"`
	Ping     int64 `json:"ping"`
	Status   int8  `json:"dc_status"`
	LastDown int64 `json:"last_down"`
	LastLag  int64 `json:"last_lag"`
}
