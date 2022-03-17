package main

import (
	"fmt"
)

func main() {
	fmt.Printf("TgServerChecker v%s, Copyright (C) 2021-2022 Laky-64 <https://github.com/Laky-64>\n"+
		"Licensed under the terms of the GNU Lesser General Public License v3 or later (LGPLv3+)",
		clientVersion,
	)

	tgClient := TelegramServerChecker()
	tgClient.Run()
}
