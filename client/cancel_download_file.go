package client

func (context *Context) CancelDownloadFile(fileId int32) {
	_, _ = context.Client.CancelDownloadFile(fileId, false)
}
