package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// UploadConfig holds the parsed values from upload curl command
type UploadConfig struct {
	URL      string
	Headers  map[string]string
	FormData map[string]string
}

// FacebookUploadPayload represents the payload from Facebook's upload response
type FacebookUploadPayload struct {
	IsSpherical      bool   `json:"isSpherical"`
	Height           int    `json:"height"`
	ImageSrc         string `json:"imageSrc"`
	MediaLocation    string `json:"mediaLocation"`
	OriginalPhotoID  string `json:"originalPhotoID"`
	PhotoID          string `json:"photoID"`
	SphericalPhotoID string `json:"sphericalPhotoID"`
	ThumbSrc         string `json:"thumbSrc"`
	Width            int    `json:"width"`
	MediaTakenTime   string `json:"mediaTakenTime"`
}

// FacebookUploadResponse represents Facebook's full upload response
type FacebookUploadResponse struct {
	AR      int                   `json:"__ar"`
	RID     string                `json:"rid"`
	Payload FacebookUploadPayload `json:"payload"`
	LID     string                `json:"lid"`
}

// UploadResponse represents our processed upload response
type UploadResponse struct {
	Success          bool   `json:"success"`
	PhotoID          string `json:"photoId"`
	ImageSrc         string `json:"imageSrc"`
	ThumbSrc         string `json:"thumbSrc"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	IsSpherical      bool   `json:"isSpherical"`
	MediaLocation    string `json:"mediaLocation"`
	OriginalPhotoID  string `json:"originalPhotoId"`
	SphericalPhotoID string `json:"sphericalPhotoId"`
	MediaTakenTime   string `json:"mediaTakenTime"`
	RequestID        string `json:"rid"`
	LocalID          string `json:"lid"`
	Message          string `json:"message"`
	Error            string `json:"error,omitempty"`
	ImagePath        string `json:"imagePath"`
}

// Global upload config
var currentUploadConfig *UploadConfig

// parseUploadCurlCommand extracts URL and headers from the upload curl command
func parseUploadCurlCommand(uploadImageCurl string) (*UploadConfig, error) {
	config := &UploadConfig{
		Headers:  make(map[string]string),
		FormData: make(map[string]string),
	}

	// Extract URL using regex
	urlRegex := regexp.MustCompile(`curl '([^']+)'`)
	urlMatch := urlRegex.FindStringSubmatch(uploadImageCurl)
	if len(urlMatch) > 1 {
		config.URL = urlMatch[1]
	}

	// Extract cookies using regex (-b flag)
	cookieRegex := regexp.MustCompile(`-b '([^']+)'`)
	cookieMatch := cookieRegex.FindStringSubmatch(uploadImageCurl)
	if len(cookieMatch) > 1 {
		config.Headers["Cookie"] = cookieMatch[1]
	}

	// Extract headers using regex (-H flags)
	headerRegex := regexp.MustCompile(`-H '([^:]+):\s*([^']+)'`)
	headerMatches := headerRegex.FindAllStringSubmatch(uploadImageCurl, -1)
	for _, match := range headerMatches {
		if len(match) > 2 {
			headerName := strings.TrimSpace(match[1])
			headerValue := strings.TrimSpace(match[2])
			config.Headers[headerName] = headerValue
		}
	}

	// Parse form data from the multipart data
	config.FormData["source"] = "8"
	config.FormData["profile_id"] = "61553861467726"
	config.FormData["waterfallxapp"] = "comet"
	config.FormData["upload_id"] = "jsc_c_4"

	return config, nil
}

// UpdateUploadConfigFromCurl parses the upload curl command and updates global config
func UpdateUploadConfigFromCurl(uploadImageCurl string) error {
	fmt.Println("🔄 Updating upload configuration from curl command...")

	config, err := parseUploadCurlCommand(uploadImageCurl)
	if err != nil {
		return fmt.Errorf("error parsing upload curl command: %v", err)
	}

	currentUploadConfig = config

	fmt.Printf("✅ Upload configuration updated successfully!\n")
	fmt.Printf("   - Upload URL: %s\n", config.URL)
	fmt.Printf("   - Headers: %d items\n", len(config.Headers))
	fmt.Printf("   - Form Data: %d parameters\n", len(config.FormData))
	fmt.Println()

	return nil
}

// createMultipartForm creates a multipart form with the image file
func createMultipartForm(imagePath string) (*bytes.Buffer, string, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Add form fields
	if currentUploadConfig != nil {
		for key, value := range currentUploadConfig.FormData {
			err := writer.WriteField(key, value)
			if err != nil {
				return nil, "", fmt.Errorf("error writing field %s: %v", key, err)
			}
		}
	}

	// Add the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, "", fmt.Errorf("error opening image file: %v", err)
	}
	defer file.Close()

	// Create form file field
	fileName := filepath.Base(imagePath)
	part, err := writer.CreateFormFile("farr", fileName)
	if err != nil {
		return nil, "", fmt.Errorf("error creating form file: %v", err)
	}

	// Copy file content to form
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", fmt.Errorf("error copying file content: %v", err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("error closing multipart writer: %v", err)
	}

	return &buffer, writer.FormDataContentType(), nil
}

// makeUploadRequest performs the actual upload request
func makeUploadRequest(imagePath string) (*UploadResponse, error) {
	if currentUploadConfig == nil {
		return nil, fmt.Errorf("upload config not initialized - call UpdateUploadConfigFromCurl first")
	}

	// Create multipart form data
	formBuffer, contentType, err := createMultipartForm(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error creating multipart form: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", currentUploadConfig.URL, formBuffer)
	if err != nil {
		return nil, fmt.Errorf("error creating upload request: %v", err)
	}

	// Set headers from config
	for headerName, headerValue := range currentUploadConfig.Headers {
		// Skip Content-Type as it's set by multipart writer
		if headerName != "Content-Type" {
			req.Header.Set(headerName, headerValue)
		}
	}

	// Set the correct Content-Type for multipart form
	req.Header.Set("Content-Type", contentType)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making upload request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading upload response: %v", err)
	}

	// decompress if needed
	if resp.Header.Get("Content-Encoding") == "zstd" {
		fmt.Println("Decompressing zstd response...")

		dec, err := zstd.NewReader(nil)
		if err != nil {
			panic(err)
		}
		defer dec.Close()

		body, err = dec.DecodeAll(body, nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(body)) // readable HTML
	} else {
		fmt.Println("Uknown encoding type:", string(body))
	}
	// end decompression code ---

	fmt.Printf("📤 Upload response status: %s\n", resp.Status)
	fmt.Printf("📤 Upload response body: %s\n", string(body))

	// Parse the response
	return parseUploadResponse(body, imagePath)
}

// parseUploadResponse parses Facebook's response and extracts the upload data
func parseUploadResponse(body []byte, imagePath string) (*UploadResponse, error) {
	bodyStr := string(body)

	// Facebook responses start with "for (;;);" - remove it
	bodyStr = strings.TrimPrefix(bodyStr, "for (;;);")

	// Parse the JSON response
	var fbResponse FacebookUploadResponse
	err := json.Unmarshal([]byte(bodyStr), &fbResponse)
	if err != nil {
		return &UploadResponse{
			Success:   false,
			Error:     fmt.Sprintf("Failed to parse JSON response: %v", err),
			Message:   "JSON parsing error",
			ImagePath: imagePath,
		}, nil
	}

	// Check if the response indicates success
	success := fbResponse.AR == 1 && fbResponse.Payload.PhotoID != ""

	// Create our structured response
	response := &UploadResponse{
		Success:          success,
		PhotoID:          fbResponse.Payload.PhotoID,
		ImageSrc:         fbResponse.Payload.ImageSrc,
		ThumbSrc:         fbResponse.Payload.ThumbSrc,
		Width:            fbResponse.Payload.Width,
		Height:           fbResponse.Payload.Height,
		IsSpherical:      fbResponse.Payload.IsSpherical,
		MediaLocation:    fbResponse.Payload.MediaLocation,
		OriginalPhotoID:  fbResponse.Payload.OriginalPhotoID,
		SphericalPhotoID: fbResponse.Payload.SphericalPhotoID,
		MediaTakenTime:   fbResponse.Payload.MediaTakenTime,
		RequestID:        fbResponse.RID,
		LocalID:          fbResponse.LID,
		ImagePath:        imagePath,
	}

	if success {
		response.Message = "Image uploaded successfully"
	} else {
		response.Error = "Upload failed - no photo ID returned"
		response.Message = "Upload failed"
	}

	return response, nil
}

// UploadImage uploads an image to Facebook and returns the upload response
func UploadImage(imagePath, uploadImageCurl string) (*UploadResponse, error) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file does not exist: %s", imagePath)
	}

	fmt.Printf("🖼️ Starting image upload for: %s\n", imagePath)

	// Update upload config if not already done
	if currentUploadConfig == nil {
		err := UpdateUploadConfigFromCurl(uploadImageCurl)
		if err != nil {
			return nil, fmt.Errorf("error updating upload config: %v", err)
		}
	}

	// Perform the upload
	response, err := makeUploadRequest(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error uploading image: %v", err)
	}

	if response.Success {
		fmt.Printf("✅ Image upload completed successfully!\n")
		fmt.Printf("   📷 Photo ID: %s\n", response.PhotoID)
		fmt.Printf("   🔗 Image URL: %s\n", response.ImageSrc)
		fmt.Printf("   📏 Dimensions: %dx%d\n", response.Width, response.Height)
	} else {
		fmt.Printf("❌ Image upload failed: %s\n", response.Error)
	}

	return response, nil
}
