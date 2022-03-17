package client

import (
	"fmt"

	tdLib "github.com/Arman92/go-tdlib"
)

func Client() Context {
	tdLib.SetLogVerbosityLevel(0)
	sessionInfo := initSession()
	client := tdLib.NewClient(tdLib.Config{
		APIID:               sessionInfo.ApiID,
		APIHash:             sessionInfo.ApiHASH,
		SystemLanguageCode:  "en",
		DeviceModel:         "OwlGram Server Checker",
		SystemVersion:       clientVersion,
		ApplicationVersion:  "1.7.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   tdSessionFiles + "/td_session",
		FileDirectory:       tdSessionFiles + "/td_files",
		IgnoreFileNames:     false,
	})
	return Context{
		Client:   client,
		ChatID:   int64(-1001110310993),
		Username: "@connectivity_test",
	}
}
