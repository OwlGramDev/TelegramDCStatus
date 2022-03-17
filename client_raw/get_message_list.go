package client_raw

import (
	tdLib "github.com/Arman92/go-tdlib"
)

func (context *Context) GetMessageList() []tdLib.Message {
	_, err := context.Client.GetChat(context.ChatID)
	if err != nil {
		_, _ = context.Client.SearchPublicChat(context.Username)
	}
	lastMessage, _ := context.Client.GetChatHistory(context.ChatID, 0, 0, 100, false)
	messagesList, _ := context.Client.GetChatHistory(context.ChatID, lastMessage.Messages[0].ID, 0, 100, false)
	messagesList.Messages = append(messagesList.Messages, lastMessage.Messages[0])
	return messagesList.Messages
}
