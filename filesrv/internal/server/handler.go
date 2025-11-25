package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Handler struct {
	UploadDir string
	BaseUrl   string
}

func NewHandler(uploadDir string, baseUrl string) *Handler {
	return &Handler{
		UploadDir: uploadDir,
		BaseUrl:   baseUrl,
	}
}

type DogApiResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (h *Handler) RandomDogHandler(w http.ResponseWriter, r *http.Request) {
	const dogApiUrl = "https://dog.ceo/api/breeds/image/random"

	var resp DogApiResponse
	res, err := http.Get(dogApiUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching dog image: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.Message))
}

func (h *Handler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form data: %v", err), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file content into memory", http.StatusInternalServerError)
		return
	}

	pngBytes, err := ConvertToPNG(fileBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting image format: %v", err), http.StatusBadRequest)
		return
	}

	baseName := filepath.Base(header.Filename)
	ext := filepath.Ext(baseName)
	uniqueName := fmt.Sprintf("%d_%s.png", time.Now().UnixNano(), baseName[:len(baseName)-len(ext)])

	filePath := filepath.Join(h.UploadDir, uniqueName)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating file on server: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = dst.Write(pngBytes); err != nil {
		http.Error(w, fmt.Sprintf("Error writing PNG data to file: %v", err), http.StatusInternalServerError)
		return
	}

	fileURL := h.BaseUrl + fmt.Sprintf("/images/%s", uniqueName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully", "url": "%s"}`, fileURL)))
}

func ConvertToPNG(inputData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(inputData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode input image data: %w", err)
	}

	output := new(bytes.Buffer)

	err = png.Encode(output, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to PNG: %w", err)
	}

	return output.Bytes(), nil
}
