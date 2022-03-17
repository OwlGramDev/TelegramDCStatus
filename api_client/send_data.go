package api_client

import (
	"encoding/json"
	"log"

	"TelegramServerChecker/consts"
	"TelegramServerChecker/types"
	"github.com/valyala/fasthttp"
)

func SendData(data []types.TelegramDCStatus) string {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(consts.ApiEndpoint + "sendData")
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
			log.Println("You have been banned from api_client.owlgram.org because you are flooding, wait and retry after 30 minutes")
		}
		return ""
	}

	return string(resp.Body())
}
