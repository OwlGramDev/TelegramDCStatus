package client_raw

import (
	tdLib "github.com/Arman92/go-tdlib"
)

type Context struct {
	Client   *tdLib.Client
	ChatID   int64
	Username string
}
