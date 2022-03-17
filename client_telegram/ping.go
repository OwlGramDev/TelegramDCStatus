package client_telegram

import (
	"log"
	"time"

	"github.com/go-ping/ping"
)

func (tg *Client) pingWithTimeout(address string) int64 {
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
