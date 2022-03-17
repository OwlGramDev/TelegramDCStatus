package main

import (
	"fmt"
)

const clientVersion = "1.0.0"
const tdSessionFiles = "./DcCheckerSession"
const backupFolder = tdSessionFiles + "/tmpDcStatus.json"
const sessionFolder = tdSessionFiles + "/session.json"
const apiEndpoint = "https://api.owlgram.org/"

func main() {
	fmt.Printf("TgServerChecker v%s, Copyright (C) 2021-2022 Laky-64 <https://github.com/Laky-64>\n"+
		"Licensed under the terms of the GNU Lesser General Public License v3 or later (LGPLv3+)",
		clientVersion,
	)

	tgClient := TelegramServerChecker()
	tgClient.Run()
}
