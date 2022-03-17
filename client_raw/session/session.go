package session

import (
	"encoding/json"
	"fmt"
	"os"

	"TelegramServerChecker/consts"
)

func InitSession() Data {
	if _, err := os.Stat(consts.TdSessionFiles); os.IsNotExist(err) {
		_ = os.Mkdir(consts.TdSessionFiles, 0644)
		_ = os.Mkdir(consts.TdSessionFiles+"/td_session", 0644)
		_ = os.Mkdir(consts.TdSessionFiles+"/td_files", 0644)
	}
	r, err := os.ReadFile(consts.SessionFolder)
	if err == nil {
		var session Data
		_ = json.Unmarshal(r, &session)
		return session
	} else {
		var apiId string
		var apiHash string
		fmt.Print("Insert your API ID: ")
		_, _ = fmt.Scanln(&apiId)
		fmt.Print("Insert your API HASH: ")
		_, _ = fmt.Scanln(&apiHash)
		session := Data{
			apiId,
			apiHash,
		}
		w, _ := json.Marshal(session)
		_ = os.WriteFile(consts.SessionFolder, w, 0644)
		return session
	}
}
