package main

import (
	"fmt"
)

var ClientVersion = "1.0.0"
var TdSessionFiles = "./DcCheckerSession"
var BackupFolder = TdSessionFiles + "/tmpDcStatus.json"
var SessionFolder = TdSessionFiles + "/session.json"
var ApiEndpoint = "https://api.owlgram.org/"

func main() {
	fmt.Println("TgServerChecker v" + ClientVersion +
		", Copyright (C) 2021-2022 Laky-64 <https://github.com/Laky-64>\n" +
		"Licensed under the terms of the GNU Lesser General Public License v3 or later (LGPLv3+)",
	)
	tgClient := TelegramServerChecker()
	tgClient.Run()
}
