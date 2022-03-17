package client

import (
	"TelegramServerChecker/client/session"
	"TelegramServerChecker/consts"
	tdLib "github.com/Arman92/go-tdlib"
)

func New() Context {
	tdLib.SetLogVerbosityLevel(0)
	sessionInfo := session.InitSession()
	client := tdLib.NewClient(tdLib.Config{
		APIID:               sessionInfo.ApiID,
		APIHash:             sessionInfo.ApiHASH,
		SystemLanguageCode:  "en",
		DeviceModel:         "OwlGram Server Checker",
		SystemVersion:       consts.ClientVersion,
		ApplicationVersion:  "1.7.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   consts.TdSessionFiles + "/td_session",
		FileDirectory:       consts.TdSessionFiles + "/td_files",
		IgnoreFileNames:     false,
	})
	return Context{
		Client:   client,
		ChatID:   int64(-1001110310993),
		Username: "@connectivity_test",
	}
}
