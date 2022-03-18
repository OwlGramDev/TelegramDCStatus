package client_raw

import (
	"fmt"

	tdLib "github.com/Arman92/go-tdlib"
)

func (context *Context) Login() {
	for {
		currentState, _ := context.Client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("\nEnter Phone Number: ")
			var number string
			_, _ = fmt.Scanln(&number)
			_, err := context.Client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("\nError sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateWaitCodeType {
			fmt.Print("\nEnter Code: ")
			var code string
			_, _ = fmt.Scanln(&code)
			_, err := context.Client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("\nError sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateWaitPasswordType {
			fmt.Print("\nEnter Password: ")
			var password string
			_, _ = fmt.Scanln(&password)
			_, err := context.Client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("\nError sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdLib.AuthorizationStateReadyType {
			break
		}
	}
}
