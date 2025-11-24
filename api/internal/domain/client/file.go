package client

type FileClient interface {
	SaveImage(data []byte, filename string) (string, error)
}
