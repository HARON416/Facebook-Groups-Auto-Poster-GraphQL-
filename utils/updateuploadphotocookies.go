package utils

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
)

// Update Facebook Image Upload function from curl command
func UpdateImageUploadFunctionFromCurl(curlCommand, filePath string) error {
	// 1. Read the original file
	originalCode, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 2. Extract values from curl command
	config, err := extractImageUploadValuesFromCurl(curlCommand)
	if err != nil {
		return fmt.Errorf("failed to extract values from curl: %w", err)
	}

	// 3. Update the code
	updatedCode := updateImageUploadGoCode(string(originalCode), config)

	// 4. Create backup
	backupPath := filePath + ".backup"
	err = os.WriteFile(backupPath, originalCode, 0644)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// 5. Write updated code
	err = os.WriteFile(filePath, []byte(updatedCode), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	fmt.Printf("✅ Successfully updated %s\n", filePath)
	fmt.Printf("📁 Backup saved as %s\n", backupPath)
	fmt.Printf("🔄 Updated %d fields\n", len(config))

	return nil
}

// Extract all Facebook image upload values from curl command
func extractImageUploadValuesFromCurl(curlCommand string) (map[string]string, error) {
	config := make(map[string]string)

	// Extract cookies from -b flag
	cookieRegex := regexp.MustCompile(`-b\s+'([^']+)'`)
	if matches := cookieRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		config["cookie"] = matches[1]
	}

	// Extract user-agent
	userAgentRegex := regexp.MustCompile(`-H\s+'user-agent:\s*([^']+)'`)
	if matches := userAgentRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		config["user-agent"] = matches[1]
	}

	// Extract sec-ch-ua
	secChUARegex := regexp.MustCompile(`-H\s+'sec-ch-ua:\s*([^']+)'`)
	if matches := secChUARegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		config["sec-ch-ua"] = matches[1]
	}

	// Extract query parameters from the URL
	urlRegex := regexp.MustCompile(`curl\s+'([^']+)'`)
	if matches := urlRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		fullURL := matches[1]

		//Parse URL to get query parameters
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}

		queryParams := parsedURL.Query()

		// Extract all the Facebook query parameter fields
		fbQueryFields := []string{
			"av", "__user", "__req", "__hs", "__ccg", "__rev", "__s", "__hsi",
			"__dyn", "__csr", "__hsdp", "__hblp", "__sjsp", "fb_dtsg",
			"jazoest", "lsd", "__spin_t", "__spin_r",
		}

		for _, field := range fbQueryFields {
			if value := queryParams.Get(field); value != "" {
				config[field] = value
			}
		}
	}

	// Extract form data from --data-raw to get profile_id
	dataRegex := regexp.MustCompile(`--data-raw\s+\$'([^']+)'`)
	if matches := dataRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		formData := matches[1]

		// Extract profile_id from the multipart form data
		profileIDRegex := regexp.MustCompile(`name="profile_id".*?\\r\\n\\r\\n([^\\]+)`)
		if profileMatches := profileIDRegex.FindStringSubmatch(formData); len(profileMatches) > 1 {
			config["profile_id"] = profileMatches[1]
		}
	}

	// Validate we got the essential fields
	essential := []string{"av", "__user", "cookie", "lsd", "profile_id"}
	for _, field := range essential {
		if config[field] == "" {
			return nil, fmt.Errorf("missing essential field: %s", field)
		}
	}

	return config, nil
}

// Update the Go code with new values for image upload function
func updateImageUploadGoCode(originalCode string, config map[string]string) string {
	updatedCode := originalCode

	// Update query parameter fields in params.Add() calls
	queryFields := []string{
		"av", "__user", "__req", "__hs", "__ccg", "__rev", "__s", "__hsi",
		"__dyn", "__csr", "__hsdp", "__hblp", "__sjsp", "fb_dtsg",
		"jazoest", "lsd", "__spin_t", "__spin_r",
	}

	for _, field := range queryFields {
		if value, exists := config[field]; exists {
			// Handle URL-encoded values (like fb_dtsg)
			decodedValue, err := url.QueryUnescape(value)
			if err != nil {
				decodedValue = value // Use original if decoding fails
			}

			pattern := fmt.Sprintf(`params\.Add\("%s", "[^"]+"\)`, regexp.QuoteMeta(field))
			replacement := fmt.Sprintf(`params.Add("%s", "%s")`, field, decodedValue)
			re := regexp.MustCompile(pattern)
			updatedCode = re.ReplaceAllString(updatedCode, replacement)
		}
	}

	// Update profile_id in WriteField call
	if profileID, exists := config["profile_id"]; exists {
		re := regexp.MustCompile(`writer\.WriteField\("profile_id", "[^"]+"\)`)
		replacement := fmt.Sprintf(`writer.WriteField("profile_id", "%s")`, profileID)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	// Update headers
	if cookie, exists := config["cookie"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("Cookie", "[^"]+"\)`)
		replacement := fmt.Sprintf(`req.Header.Set("Cookie", "%s")`, cookie)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	if userAgent, exists := config["user-agent"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("User-Agent", "[^"]+"\)`)
		replacement := fmt.Sprintf(`req.Header.Set("User-Agent", "%s")`, userAgent)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	if secChUA, exists := config["sec-ch-ua"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("Sec-Ch-Ua", ` + "`" + `[^` + "`" + `]+` + "`" + `\)`)
		replacement := fmt.Sprintf(`req.Header.Set("Sec-Ch-Ua", `+"`"+`%s`+"`"+`)`, secChUA)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	return updatedCode
}

// Command line interface
// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Println("Usage: go run image_upload_updater.go '<curl-command>' <go-file-path>")
// 		fmt.Println()
// 		fmt.Println("Example:")
// 		fmt.Println(`go run image_upload_updater.go 'curl https://upload.facebook.com/ajax/react_composer/attachments/photo/upload...' image_upload.go`)
// 		os.Exit(1)
// 	}

// 	curlCommand := os.Args[1]
// 	filePath := os.Args[2]

// 	err := UpdateImageUploadFromCurl(curlCommand, filePath)
// 	if err != nil {
// 		log.Fatalf("❌ Update failed: %v", err)
// 	}

// 	fmt.Println("🎉 Image upload function updated successfully!")
// }

// Alternative: Update from string (no file needed)
func UpdateImageUploadCodeFromCurlAndOriginalCode(curlCommand, originalCode string) (string, error) {
	config, err := extractImageUploadValuesFromCurl(curlCommand)
	if err != nil {
		return "", err
	}

	updatedCode := updateImageUploadGoCode(originalCode, config)
	return updatedCode, nil
}

// Quick extraction function for debugging
func ExtractImageUploadValues(curlCommand string) error {
	config, err := extractImageUploadValuesFromCurl(curlCommand)
	if err != nil {
		return err
	}

	fmt.Println("📋 Extracted values from curl command:")
	fmt.Printf("   User ID: %s\n", config["__user"])
	fmt.Printf("   Profile ID: %s\n", config["profile_id"])
	fmt.Printf("   LSD: %s\n", config["lsd"])
	fmt.Printf("   FB DTSG: %s\n", config["fb_dtsg"])
	fmt.Printf("   Cookie length: %d characters\n", len(config["cookie"]))

	// Show all extracted fields
	fmt.Println("\n🔍 All extracted fields:")
	for key, value := range config {
		if len(value) > 50 {
			fmt.Printf("   %s: %s... (length: %d)\n", key, value[:50], len(value))
		} else {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}

	return nil
}
