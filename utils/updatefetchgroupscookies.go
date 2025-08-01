package utils

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// CurlData represents the parsed curl request data
type CurlData struct {
	URL     string
	Headers map[string]string
	Cookies string
	Data    map[string]string
}

func UpdateFetchGroupsCookies() error {
	// Read the curl request from the cookies file
	curlRequest, err := readCurlRequest("cookies/fetchgroupscookies.txt")
	if err != nil {
		return fmt.Errorf("error reading curl request: %w", err)
	}

	// Parse the curl request
	parsedData, err := parseCurlRequest(curlRequest)
	if err != nil {
		return fmt.Errorf("error parsing curl request: %w", err)
	}

	// Update the fetchgroups.go file with new credentials
	err = updateFetchGroupsFile(parsedData)
	if err != nil {
		return fmt.Errorf("error updating fetchgroups.go: %w", err)
	}

	// Verify the update was successful
	err = verifyUpdate(parsedData)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Println("🎉 Facebook Groups cookies have been successfully updated!")

	return nil
}

func verifyUpdate(parsedData *CurlData) error {
	// Read the updated file to verify changes
	content, err := os.ReadFile("utils/fetchgroups.go")
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Verify critical parameters were updated
	criticalParams := []string{
		"av", "__user", "fb_dtsg", "lsd", "jazoest",
		"__req", "__rev", "__spin_r", "__spin_t",
	}

	for _, param := range criticalParams {
		if value, exists := parsedData.Data[param]; exists {
			// Check if the parameter was updated in the file
			expectedLine := fmt.Sprintf(`data.Set("%s", "%s")`, param, value)
			if !strings.Contains(contentStr, expectedLine) {
				return fmt.Errorf("critical parameter %s was not updated correctly", param)
			}
		}
	}

	// Verify cookies were updated
	if parsedData.Cookies != "" {
		expectedCookieLine := fmt.Sprintf(`req.Header.Set("cookie", "%s")`, parsedData.Cookies)
		if !strings.Contains(contentStr, expectedCookieLine) {
			return fmt.Errorf("cookies were not updated correctly")
		}
	}

	// Verify x-fb-lsd header was updated
	if lsdValue, exists := parsedData.Data["lsd"]; exists {
		expectedLsdLine := fmt.Sprintf(`req.Header.Set("x-fb-lsd", "%s")`, lsdValue)
		if !strings.Contains(contentStr, expectedLsdLine) {
			return fmt.Errorf("x-fb-lsd header was not updated correctly")
		}
	}

	// Verify critical headers were updated
	criticalHeaders := []string{"accept", "content-type", "origin", "user-agent", "x-fb-friendly-name"}
	for _, header := range criticalHeaders {
		if headerValue, exists := parsedData.Headers[header]; exists {
			expectedHeaderLine := fmt.Sprintf(`req.Header.Set("%s", "%s")`, header, headerValue)
			if !strings.Contains(contentStr, expectedHeaderLine) {
				return fmt.Errorf("critical header %s was not updated correctly", header)
			}
		}
	}

	fmt.Println("🔍 Verification passed: All critical parameters and headers updated successfully!")

	return nil
}

func readCurlRequest(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n"), scanner.Err()
}

func parseCurlRequest(curlRequest string) (*CurlData, error) {
	data := &CurlData{
		Headers: make(map[string]string),
		Data:    make(map[string]string),
	}

	lines := strings.Split(curlRequest, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract URL with query parameters
		if strings.Contains(line, "curl '") {
			urlMatch := regexp.MustCompile(`curl '([^']+)'`).FindStringSubmatch(line)
			if len(urlMatch) > 1 {
				fullURL := urlMatch[1]
				data.URL = fullURL
				// Parse URL to extract query parameters
				parsedURL, err := url.Parse(fullURL)
				if err == nil {
					// Extract query parameters
					for key, values := range parsedURL.Query() {
						if len(values) > 0 {
							data.Data[key] = values[0]
						}
					}
				}
			}
			continue
		}

		// Extract headers
		if strings.Contains(line, "-H '") {
			headerMatch := regexp.MustCompile(`-H '([^:]+): ([^']+)'`).FindStringSubmatch(line)
			if len(headerMatch) > 2 {
				headerName := headerMatch[1]
				headerValue := headerMatch[2]
				data.Headers[headerName] = headerValue
			}
			continue
		}

		// Extract cookies
		if strings.Contains(line, "-b '") {
			cookieMatch := regexp.MustCompile(`-b '([^']+)'`).FindStringSubmatch(line)
			if len(cookieMatch) > 1 {
				data.Cookies = cookieMatch[1]
			}
			continue
		}

		// Extract form data
		if strings.Contains(line, "--data-raw '") {
			dataRawMatch := regexp.MustCompile(`--data-raw '([^']+)'`).FindStringSubmatch(line)
			if len(dataRawMatch) > 1 {
				formData := dataRawMatch[1]
				// Parse the form data string
				pairs := strings.Split(formData, "&")
				for _, pair := range pairs {
					keyValue := strings.SplitN(pair, "=", 2)
					if len(keyValue) == 2 {
						key := keyValue[0]
						value := keyValue[1]
						// URL decode the value
						decodedValue, err := url.QueryUnescape(value)
						if err == nil {
							data.Data[key] = decodedValue
						} else {
							data.Data[key] = value
						}
					}
				}
			}
			continue
		}
	}

	return data, nil
}

func updateFetchGroupsFile(parsedData *CurlData) error {
	// Read the current fetchgroups.go file
	content, err := os.ReadFile("utils/fetchgroups.go")
	if err != nil {
		return err
	}

	contentStr := string(content)

	fmt.Printf("🔍 DEBUG: Starting update with %d form data items and %d headers\n", len(parsedData.Data), len(parsedData.Headers))

	// Update cookies
	if parsedData.Cookies != "" {
		cookiePattern := regexp.MustCompile(`req\.Header\.Set\("cookie", ".*"\)`)
		newCookieLine := fmt.Sprintf(`req.Header.Set("cookie", "%s")`, parsedData.Cookies)
		oldContent := contentStr
		contentStr = cookiePattern.ReplaceAllString(contentStr, newCookieLine)
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for cookies\n")
		} else {
			fmt.Printf("✅ Updated cookies: %s\n", parsedData.Cookies[:min(50, len(parsedData.Cookies))])
		}
	}

	// Update form data parameters - ALL parameters from the curl request
	formDataUpdates := map[string]string{
		"av":                       parsedData.Data["av"],
		"__aaid":                   parsedData.Data["__aaid"],
		"__user":                   parsedData.Data["__user"],
		"__a":                      parsedData.Data["__a"],
		"__req":                    parsedData.Data["__req"],
		"__hs":                     parsedData.Data["__hs"],
		"dpr":                      parsedData.Data["dpr"],
		"__ccg":                    parsedData.Data["__ccg"],
		"__rev":                    parsedData.Data["__rev"],
		"__s":                      parsedData.Data["__s"],
		"__hsi":                    parsedData.Data["__hsi"],
		"__dyn":                    parsedData.Data["__dyn"],
		"__csr":                    parsedData.Data["__csr"],
		"__hsdp":                   parsedData.Data["__hsdp"],
		"__hblp":                   parsedData.Data["__hblp"],
		"__sjsp":                   parsedData.Data["__sjsp"],
		"__comet_req":              parsedData.Data["__comet_req"],
		"fb_dtsg":                  parsedData.Data["fb_dtsg"],
		"jazoest":                  parsedData.Data["jazoest"],
		"lsd":                      parsedData.Data["lsd"],
		"__spin_r":                 parsedData.Data["__spin_r"],
		"__spin_b":                 parsedData.Data["__spin_b"],
		"__spin_t":                 parsedData.Data["__spin_t"],
		"__crn":                    parsedData.Data["__crn"],
		"fb_api_caller_class":      parsedData.Data["fb_api_caller_class"],
		"fb_api_req_friendly_name": parsedData.Data["fb_api_req_friendly_name"],
		"server_timestamps":        parsedData.Data["server_timestamps"],
		"doc_id":                   parsedData.Data["doc_id"],
	}

	for key, value := range formDataUpdates {
		// Update even if value is empty (to replace empty placeholders)
		pattern := regexp.MustCompile(fmt.Sprintf(`data\.Set\("%s", ".*"\)`, key))
		newLine := fmt.Sprintf(`data.Set("%s", "%s")`, key, value)
		oldContent := contentStr
		contentStr = pattern.ReplaceAllString(contentStr, newLine)
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for form data key: %s\n", key)
		} else {
			fmt.Printf("✅ Updated form data: %s = %s\n", key, value[:min(20, len(value))])
		}
	}

	// Update x-fb-lsd header
	if lsdValue, exists := parsedData.Data["lsd"]; exists {
		lsdPattern := regexp.MustCompile(`req\.Header\.Set\("x-fb-lsd", ".*"\)`)
		newLsdLine := fmt.Sprintf(`req.Header.Set("x-fb-lsd", "%s")`, lsdValue)
		oldContent := contentStr
		contentStr = lsdPattern.ReplaceAllString(contentStr, newLsdLine)
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for x-fb-lsd\n")
		} else {
			fmt.Printf("✅ Updated x-fb-lsd: %s\n", lsdValue)
		}
	}

	// Update other headers that might change
	headerUpdates := map[string]string{
		"accept":                      parsedData.Headers["accept"],
		"accept-language":             parsedData.Headers["accept-language"],
		"content-type":                parsedData.Headers["content-type"],
		"origin":                      parsedData.Headers["origin"],
		"priority":                    parsedData.Headers["priority"],
		"referer":                     parsedData.Headers["referer"],
		"sec-ch-prefers-color-scheme": parsedData.Headers["sec-ch-prefers-color-scheme"],
		"sec-ch-ua":                   parsedData.Headers["sec-ch-ua"],
		"sec-ch-ua-full-version-list": parsedData.Headers["sec-ch-ua-full-version-list"],
		"sec-ch-ua-mobile":            parsedData.Headers["sec-ch-ua-mobile"],
		"sec-ch-ua-model":             parsedData.Headers["sec-ch-ua-model"],
		"sec-ch-ua-platform":          parsedData.Headers["sec-ch-ua-platform"],
		"sec-ch-ua-platform-version":  parsedData.Headers["sec-ch-ua-platform-version"],
		"sec-fetch-dest":              parsedData.Headers["sec-fetch-dest"],
		"sec-fetch-mode":              parsedData.Headers["sec-fetch-mode"],
		"sec-fetch-site":              parsedData.Headers["sec-fetch-site"],
		"user-agent":                  parsedData.Headers["user-agent"],
		"x-asbd-id":                   parsedData.Headers["x-asbd-id"],
		"x-fb-friendly-name":          parsedData.Headers["x-fb-friendly-name"],
	}

	for headerName, headerValue := range headerUpdates {
		// Update even if value is empty (to replace empty placeholders)
		pattern := regexp.MustCompile(fmt.Sprintf(`req\.Header\.Set\("%s", ".*"\)`, headerName))
		oldContent := contentStr
		// Use backticks for headers that contain quotes to avoid escaping issues
		if strings.Contains(headerValue, `"`) {
			newHeaderLine := fmt.Sprintf(`req.Header.Set("%s", `+"`%s`"+`)`, headerName, headerValue)
			contentStr = pattern.ReplaceAllString(contentStr, newHeaderLine)
		} else {
			newHeaderLine := fmt.Sprintf(`req.Header.Set("%s", "%s")`, headerName, headerValue)
			contentStr = pattern.ReplaceAllString(contentStr, newHeaderLine)
		}
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for header: %s\n", headerName)
		} else {
			fmt.Printf("✅ Updated header: %s = %s\n", headerName, headerValue[:min(20, len(headerValue))])
		}
	}

	// Write the updated content back to the file
	err = os.WriteFile("utils/fetchgroups.go", []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
