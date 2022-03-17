package client

func (context *Context) DownloadFile(fileId int32) error {
	_, err := context.Client.DownloadFile(fileId, 1, 0, 0, true)
	return err
}
