package client_telegram

import (
	"time"
)

func (tg *Client) runDownloadWithTimeout(fileId int32) int8 {
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
