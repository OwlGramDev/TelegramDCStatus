package main

import (
	"fmt"
	"strconv"
	"strings"

	"TelegramServerChecker/client"
	"TelegramServerChecker/types"
	tdLib "github.com/Arman92/go-tdlib"
)

func TelegramServerChecker() *types.TelegramCheckerClient {
	instance := client.New()
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
		Client:       &instance,
		FilesDC:      listDCInfo,
		StatusDC:     listStatus,
		IsRefreshing: true,
	}
}
