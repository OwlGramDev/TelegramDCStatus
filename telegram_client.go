package main

import (
	"fmt"

	tdLib "github.com/Arman92/go-tdlib"
)

func Client() *ClientContext {
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
	return &ClientContext{
		client,
		int64(-1001110310993),
		"@connectivity_test",
	}
}

func (context *ClientContext) Login() {
	for {
		currentState, _ := context.Client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("\nEnter Phone Number: ")
			var number string
			_, _ = fmt.Scanln(&number)
			_, err := context.Client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateWaitCodeType {
			fmt.Print("\nEnter Code: ")
			var code string
			_, _ = fmt.Scanln(&code)
			_, err := context.Client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateWaitPasswordType {
			fmt.Print("\nEnter Password: ")
			var password string
			_, _ = fmt.Scanln(&password)
			_, err := context.Client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateReadyType {
			break
		}
	}
}

func (context *ClientContext) DownloadFile(fileId int32) error {
	_, err := context.Client.DownloadFile(fileId, 1, 0, 0, true)
	return err
}

func (context *ClientContext) CancelDownloadFile(fileId int32) {
	_, _ = context.Client.CancelDownloadFile(fileId, false)
}

func (context *ClientContext) GetMessageList() []tdLib.Message {
	_, err := context.Client.GetChat(context.ChatID)
	if err != nil {
		_, _ = context.Client.SearchPublicChat(context.Username)
	}
	lastMessage, _ := context.Client.GetChatHistory(context.ChatID, 0, 0, 100, false)
	messagesList, _ := context.Client.GetChatHistory(context.ChatID, lastMessage.Messages[0].ID, 0, 100, false)
	messagesList.Messages = append(messagesList.Messages, lastMessage.Messages[0])
	return messagesList.Messages
}

type ClientContext struct {
	Client   *tdLib.Client
	ChatID   int64
	Username string
}
