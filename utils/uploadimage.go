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
)

var uploadImageCurl = `curl 'https://upload.facebook.com/ajax/react_composer/attachments/photo/upload?av=100016139237616&__aaid=0&__user=100016139237616&__a=1&__req=2m&__hs=20324.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1026265415&__s=mn1ljd%3Ahzfeis%3Awgvtk4&__hsi=7542088668936050370&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu2u5Ee88o4Wm7-2K0-obUG2-azqwaW223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwWzUfHDzUiBG2OUqwjVqwLwHwa211wo83KwHwOyUqxG&__csr=gekt1hb8YIbN31xax2d6hsDifi2Yj_rsjYGPnkrZsG4MBqkJ9ddNdbpkAyuCOcyfLibnmSxmJd9AZXHWZqyBJasDGKCGfcl6BnHQlFdDAVFG4aEDBlyeQmjKqGnG9-hdbABqhp4Z4pHFKjxe8HmbyGGAmUOWhHGEyFEOVqVuHiABWHy_By99V9rGBgK4FHzWhqmibKaAACBV8h-bAyKtpmmfAQdyFGnKmVoPy8OeWVoO9LxaUKum5Acz94u4ETxvzGKihyoJBm9yryHgswCGp1Ouq69onxjy8-58Ny9GCADzE-XKbVWizqGrgyF8SibxS5FEuyaxC9yEdA9x6ibgyEsAyUKfgsG5Ea8jxeazEjx2i6awExqay9ohBwOBwxxGieye3PyWwNy9Q4FK5E8E9o9A5EK2Lz8TGUhzEGdzV8nxC390wxS9xWEXxy7UnwZByE5K2-7Uetbyoyex21Zh4ifwaBCo4G5odEfEKbwjU31Dwd-1bwEAj812wFAAAxe1xgepu3eFQ3mym7oeE4xa58swtUZ0xw66wmomw0wPw1_O04K-6E2zU0KS9wgqw3DoiwbRK3K0N208209oQ05m4l5g2so4e00E4403MK08wB8SE4-8x-0u20gC3a0LEjwJg0Ae8g3eF1Szwg81JUaUAw0zm4409Lo1n819ohwgE0y6680sew3To1C3w2F8C9gOt03KQ0wQ0fiwXxe02GC0aFwjS0g67rgW19waq8Dg4O0bZw4kw2Nk2TBt0hEuwaC08HAw5Fw2_E0pfAw60yWwtm5o6zC9fmdpk62kU5K1i80gOmE512U17o0-GsNUsdw&__hsdp=gjgtFMidEmNgywk8gyEsCoGeGp25IdEWEy93ezyxccImyEub48O22yiaAAxbdcggz8ixR8GdJ9JcIzRNr2dOsgRAp2aVicDComyT9jqsgxsixkvlEuygDpigj9goaT10G14NN5isI8_2i3cjbj66kgxiAp6y4aACkG8KyLah2HAz468KhBrPMyrG4oPTKlUPoBebeAFjuaJcELjBAF9Faleqh8LqlemGECzk8m84hjlAfhSFWqKZ95Fd5EzK2qsyah9DUgby32QA8cOin98pjFcFzARkvyAr8iz9194t2FCfRkMB4DBAp98xrhlanGdGl8x9l4cFpoK5cE89k9xx169izF4qBxd1t16moySizBF6FFcx4SFhkFm4EhoBovxIEy2sM-8CAoy9ig-iiFQWjjKbJCh8BUR7GQQ4pAh3A1ao-ehFot3m5Aem8hj2kiE88G49k7A1kgZ0lXz4bxaE7a1qxK3qE1F88oG17AwiPUqoS5Ud98kwJD70goaEa8mxwYwfA2Wm5ojwv8Rwik5VyxafwTxAp7whURwxwbq1mwi87S5o3Hw8O0y81YKi3m0DEYE2dwdO0p20m-ayFE561Ax2WwaCm1pz8gwRwgovw8qi6o8981K832wpU14U16UK0_81FE0z-0gO0OUpwSwe60HE21wnE2PwCwlE1Y81kU5W0ke1Tw4ew&__hblp=0qp46UKi1bwiEK4S3i16Uc84WU8o2fwdW0Topy8cUco4Cu7okG3eUbkA3i79U989E8olyo8FobohwyzUSUaVQdGeyGwzxy3i2-1wAwxy8y1HwAz898iyVu1lx7Dwxx-1wwMF0wAypeUeUeUWi0wo-789K9u2vwCxCmeg-2l7wbim4A9gS0AofE5m2K2h0oo4C1NwmFo4e8ga89UOewaq1mwGxehxG1fzK2K3-2K1dUvwWwKwRxeh1t09m1TS0SE2jwMwywQwgUuw9u1Zx92o5u4pEtwtFE2Rxy3m0C9oqwkod415wFG8yEuwRwa61KwNABwgUe839wcKewYwyxmbwqE2TwVAwh8C3O3-Ueqwmd0OwsUC1yxi9US362h16awtV5gC2R0Jx6dwCy8pykm9yGyU4h0KwHxCewa24ojyUhxKbwca3-7olwHwMwYwtFUmG3-4U23wYDxmq9w8u3GV8sxa1nwPwNgvwJjwmEK1lyoliDx-ew9yu10x-3636t0xxy4GwkE8oaoiDwhEcocEDghDxaq2q5olg1mUG2C6po8qJ2Ua4bwJho5m32cg4qdxCpa2J04QyrUWu2G3K1QwzoeE-WCy9Zedwzxmfgyqm2m5Am2K15xC1sBjgswoo43y8jKuFIw7OdxC78&__sjsp=gjgtFMidEmNgywk8gyEsCoGeGp25IdEWEy93ezyxccIF4O8IibdEH6kUAyF98iSINcOcN27kyESQCQOOfn5IAWtaL42SiGyicDwIjGakV444OSK548p98J7wya5oC2Yoapy7G8hUmgBk88mkxW8gFE_plCfhApGOp-l8x2GQ9i8KDLamay8N1ybADhzMyrG4oPTKlUPyQUIXy3ueh6cQVpbFGVoKh8ExkVqGQmLylx67AteuKHBihqgGSXwCyG8qCvx0K8acigxmOinqQpjF3zAAg-cgwwZ194t2FC4lguBBh98xp8xanGdGlbh9oRDBwCwwBgC644oBaeAhFokxt16c8JAxNF4p4S85nBoix5z8vxIwC2sM-8CAoy9igsDjFdeUZCh8BUR7GQQ64h0mCfzAq2cdomghy8lhU88rBg7h3Q1nx12UlwsE5G6UdE1FUb84u1owGxu1ewR709m5Eof83V0Qxm4U2ao4B0RwohAu17zm0S80Pm0va3O0Hyw5Gw6gw5LyEGq0JEkG0Gpo5CcwhE467U26AxC22i&__comet_req=15&fb_dtsg=NAfugHWuVdIrshH1K3diAPKmGyUKWyOnSE5wzWgPmy61pQcQP5pB6sw%3A1%3A1737623622&jazoest=25568&lsd=Hk3RdF0qVDvsVE3d67kML5&__spin_r=1026265415&__spin_b=trunk&__spin_t=1756029359&__crn=comet.fbweb.CometGroupDiscussionRoute' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundaryqAzuga36iQ65tS1b' \
  -b 'sb=E2sHZ8Dj0xZhKsq5e2frpYxk; ps_l=1; ps_n=1; datr=EZluZ78uE_dw-PpiG3FGLpyf; oo=v1; c_user=100016139237616; wd=1366x681; fr=1A8nFCnZNIAxIIUaz.AWdYb3ZlfY48CCQuck46A0USrH84R9Fq1L6oIKBkWMComAMQmsI.BoquF2..AAA.0.0.BoquF2.AWc_XA5JEwVwF7IioIG9NYQJDOc; xs=1%3AHCeFsY1N8T3k8w%3A2%3A1737623622%3A-1%3A-1%3A%3AAcXseJLaPuOV3SFAyLRrlUlfJb0CyeugC2sGNHKJjg; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1756029376368%2C%22v%22%3A1%7D' \
  -H 'origin: https://web.facebook.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://web.facebook.com/' \
  -H 'sec-ch-ua: "Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-site' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36' \
  --data-raw $'------WebKitFormBoundaryqAzuga36iQ65tS1b\r\nContent-Disposition: form-data; name="source"\r\n\r\n8\r\n------WebKitFormBoundaryqAzuga36iQ65tS1b\r\nContent-Disposition: form-data; name="profile_id"\r\n\r\n100016139237616\r\n------WebKitFormBoundaryqAzuga36iQ65tS1b\r\nContent-Disposition: form-data; name="waterfallxapp"\r\n\r\ncomet\r\n------WebKitFormBoundaryqAzuga36iQ65tS1b\r\nContent-Disposition: form-data; name="farr"; filename="449214697_122162211716199456_6629670485055728462_n.jpg"\r\nContent-Type: image/jpeg\r\n\r\n\r\n------WebKitFormBoundaryqAzuga36iQ65tS1b\r\nContent-Disposition: form-data; name="upload_id"\r\n\r\njsc_c_q\r\n------WebKitFormBoundaryqAzuga36iQ65tS1b--\r\n'`

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
func parseUploadCurlCommand() (*UploadConfig, error) {
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
func UpdateUploadConfigFromCurl() error {
	fmt.Println("🔄 Updating upload configuration from curl command...")

	config, err := parseUploadCurlCommand()
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
func UploadImage(imagePath string) (*UploadResponse, error) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file does not exist: %s", imagePath)
	}

	fmt.Printf("🖼️ Starting image upload for: %s\n", imagePath)

	// Update upload config if not already done
	if currentUploadConfig == nil {
		err := UpdateUploadConfigFromCurl()
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
