package types

import (
	tdLib "github.com/Arman92/go-tdlib"
)

type ClientContext struct {
	Client   *tdLib.Client
	ChatID   int64
	Username string
}
