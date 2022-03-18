package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"

	"TelegramServerChecker/consts"
	"TelegramServerChecker/types"
)

func checkScore() {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(consts.ApiEndpoint + "getScore")
	err := fasthttp.Do(req, resp)
	if err != nil {
		panic("No Internet Connection")
	}
	statusCode := resp.StatusCode()
	if statusCode != fasthttp.StatusOK {
		if statusCode == fasthttp.StatusTooManyRequests {
			panic("You have been banned from api_client.owlgram.org because you are flooding, wait and retry after 30 minutes")
		}
	}
	var result types.ScoreResult
	_ = json.Unmarshal(resp.Body(), &result)
	fmt.Printf("\n\nYour Score is of %d%%\n", result.Score)
}
