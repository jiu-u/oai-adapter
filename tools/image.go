package tools

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ImageFileData struct {
	MIMEType string `json:"mimeType"`
	URL      string `json:"url"`
	Content  any    `json:"content"`
}

func NewImageFileData(url string, toB64 bool) (*ImageFileData, error) {
	if IsBase64Image(url) {
		return &ImageFileData{
			MIMEType: GetB64ImageMIMEType(url),
			URL:      url,
		}, nil
	}
	if IsNetImage(url) {
		b64Str, mimeType, err := GetNetImageB64(url)
		if err != nil {
			return nil, err
		}
		if toB64 {
			return &ImageFileData{
				MIMEType: mimeType,
				URL:      b64Str,
			}, nil
		} else {
			return &ImageFileData{
				MIMEType: mimeType,
				URL:      url,
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid image url: %s", url)
}

func IsBase64Image(url string) bool {
	if strings.HasPrefix(url, "data:image/") {
		return true
	}
	return false
}

func GetB64ImageMIMEType(url string) string {
	if strings.HasPrefix(url, "data:image/") {
		re := regexp.MustCompile(`^data:(.*?);base64,`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func IsNetImage(url string) bool {
	if strings.HasPrefix(url, "http") {
		return true
	}
	return false
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func getMimeType(data []byte) string {
	return http.DetectContentType(data)
}

func convertToBase64(data []byte, mimeType string) string {
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
}

func GetNetImageB64(url string) (b64 string, mimeType string, err error) {
	imgData, err := downloadImage(url)
	if err != nil {
		return "", "", err
	}
	mimeType = getMimeType(imgData)
	base64Image := convertToBase64(imgData, mimeType)
	return base64Image, mimeType, nil
}
