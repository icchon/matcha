package fileclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type filesrvClient struct {
	imageUploadEndpoint string
}
var _ repo.FileClient = (*filesrvClient)(nil)

func NewFilesrvClient(imageUploadEndpoint string) *filesrvClient {
	return &filesrvClient{imageUploadEndpoint: imageUploadEndpoint}
}

type uploadResponse struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

func (c *filesrvClient) SaveImage(data []byte, filename string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file part: %w", err)
	}

	_, err = part.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to write data to form file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.imageUploadEndpoint, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to file service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("file service returned non-OK status code %d: %s", resp.StatusCode, respBody)
	}

	var result uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}
	return result.URL, nil
}
