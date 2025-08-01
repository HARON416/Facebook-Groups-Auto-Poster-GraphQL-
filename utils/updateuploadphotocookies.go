package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func UpdateUploadPhotoCookies() error {
	// Read the curl request from the cookies file
	curlRequest, err := readCurlRequest("cookies/uploadimagecookies.txt")
	if err != nil {
		return fmt.Errorf("error reading curl request: %w", err)
	}

	// Parse the curl request
	parsedData, err := parseCurlRequest(curlRequest)
	if err != nil {
		return fmt.Errorf("error parsing curl request: %w", err)
	}

	// Update the uploadimage.go file with new credentials
	err = updateUploadImageFile(parsedData)
	if err != nil {
		return fmt.Errorf("error updating uploadimage.go: %w", err)
	}

	// Verify the update was successful
	err = verifyUploadUpdate(parsedData)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Println("🎉 Facebook Upload Photo cookies have been successfully updated!")

	return nil
}

func verifyUploadUpdate(parsedData *CurlData) error {
	// Read the updated file to verify changes
	content, err := os.ReadFile("utils/uploadimage.go")
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
			expectedLine := fmt.Sprintf(`params.Add("%s", "%s")`, param, value)
			if !strings.Contains(contentStr, expectedLine) {
				return fmt.Errorf("critical parameter %s was not updated correctly", param)
			}
		}
	}

	// Verify cookies were updated
	if parsedData.Cookies != "" {
		expectedCookieLine := fmt.Sprintf(`req.Header.Set("Cookie", "%s")`, parsedData.Cookies)
		if !strings.Contains(contentStr, expectedCookieLine) {
			return fmt.Errorf("cookies were not updated correctly")
		}
	}

	// Verify critical headers were updated
	criticalHeaders := []string{"Accept", "Origin", "User-Agent"}
	for _, header := range criticalHeaders {
		lowerHeader := strings.ToLower(header)
		if headerValue, exists := parsedData.Headers[lowerHeader]; exists {
			expectedHeaderLine := fmt.Sprintf(`req.Header.Set("%s", "%s")`, header, headerValue)
			if !strings.Contains(contentStr, expectedHeaderLine) {
				return fmt.Errorf("critical header %s was not updated correctly", header)
			}
		}
	}

	fmt.Println("🔍 Verification passed: All critical parameters and headers updated successfully!")
	return nil
}

func updateUploadImageFile(parsedData *CurlData) error {
	// Read the current uploadimage.go file
	content, err := os.ReadFile("utils/uploadimage.go")
	if err != nil {
		return err
	}

	contentStr := string(content)

	fmt.Printf("🔍 DEBUG: Starting update with %d query parameters and %d headers\n", len(parsedData.Data), len(parsedData.Headers))

	// Update cookies
	if parsedData.Cookies != "" {
		cookiePattern := regexp.MustCompile(`req\.Header\.Set\("Cookie", "[^"]*"\)`)
		newCookieLine := fmt.Sprintf(`req.Header.Set("Cookie", "%s")`, parsedData.Cookies)
		oldContent := contentStr
		contentStr = cookiePattern.ReplaceAllString(contentStr, newCookieLine)
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for cookies\n")
		} else {
			fmt.Printf("✅ Updated cookies: %s\n", parsedData.Cookies[:min(50, len(parsedData.Cookies))])
		}
	}

	// Update query parameters - ALL parameters from the curl request
	queryParamUpdates := map[string]string{
		"av":          parsedData.Data["av"],
		"__aaid":      parsedData.Data["__aaid"],
		"__user":      parsedData.Data["__user"],
		"__a":         parsedData.Data["__a"],
		"__req":       parsedData.Data["__req"],
		"__hs":        parsedData.Data["__hs"],
		"dpr":         parsedData.Data["dpr"],
		"__ccg":       parsedData.Data["__ccg"],
		"__rev":       parsedData.Data["__rev"],
		"__s":         parsedData.Data["__s"],
		"__hsi":       parsedData.Data["__hsi"],
		"__dyn":       parsedData.Data["__dyn"],
		"__csr":       parsedData.Data["__csr"],
		"__hsdp":      parsedData.Data["__hsdp"],
		"__hblp":      parsedData.Data["__hblp"],
		"__sjsp":      parsedData.Data["__sjsp"],
		"__comet_req": parsedData.Data["__comet_req"],
		"fb_dtsg":     parsedData.Data["fb_dtsg"],
		"jazoest":     parsedData.Data["jazoest"],
		"lsd":         parsedData.Data["lsd"],
		"__spin_r":    parsedData.Data["__spin_r"],
		"__spin_b":    parsedData.Data["__spin_b"],
		"__spin_t":    parsedData.Data["__spin_t"],
		"__crn":       parsedData.Data["__crn"],
	}

	for key, value := range queryParamUpdates {
		// Update even if value is empty (to replace empty placeholders)
		pattern := regexp.MustCompile(fmt.Sprintf(`params\.Add\("%s", "[^"]*"\)`, key))
		newLine := fmt.Sprintf(`params.Add("%s", "%s")`, key, value)
		oldContent := contentStr
		contentStr = pattern.ReplaceAllString(contentStr, newLine)
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for query parameter: %s\n", key)
		} else {
			fmt.Printf("✅ Updated query parameter: %s = %s\n", key, value[:min(20, len(value))])
		}
	}

	// Update form data profile_id
	if userID, exists := parsedData.Data["av"]; exists {
		profileIDPattern := regexp.MustCompile(`writer\.WriteField\("profile_id", "[^"]*"\)`)
		newProfileIDLine := fmt.Sprintf(`writer.WriteField("profile_id", "%s")`, userID)
		oldContent := contentStr
		contentStr = profileIDPattern.ReplaceAllString(contentStr, newProfileIDLine)
		if oldContent == contentStr {
			fmt.Printf("⚠️  WARNING: No match found for profile_id\n")
		} else {
			fmt.Printf("✅ Updated profile_id: %s\n", userID)
		}
	}

	// Update headers from curl request
	headerUpdates := map[string]string{
		"Accept":          parsedData.Headers["accept"],
		"Accept-Language": parsedData.Headers["accept-language"],
		"Origin":          parsedData.Headers["origin"],
		"Referer":         parsedData.Headers["referer"],
		"User-Agent":      parsedData.Headers["user-agent"],
	}

	for headerName, headerValue := range headerUpdates {
		// Try to update existing headers (both empty placeholders and hardcoded values)
		emptyPattern := regexp.MustCompile(fmt.Sprintf(`req\.Header\.Set\("%s", ""\)`, headerName))
		existingPattern := regexp.MustCompile(fmt.Sprintf(`req\.Header\.Set\("%s", "[^"]*"\)`, headerName))

		oldContent := contentStr
		updated := false

		// Use backticks for headers that contain quotes to avoid escaping issues
		if strings.Contains(headerValue, `"`) {
			newHeaderLine := fmt.Sprintf(`req.Header.Set("%s", `+"`%s`"+`)`, headerName, headerValue)
			// Try empty placeholder first
			contentStr = emptyPattern.ReplaceAllString(contentStr, newHeaderLine)
			if oldContent != contentStr {
				updated = true
			} else {
				// Try existing hardcoded value
				contentStr = existingPattern.ReplaceAllString(contentStr, newHeaderLine)
				if oldContent != contentStr {
					updated = true
				}
			}
		} else {
			newHeaderLine := fmt.Sprintf(`req.Header.Set("%s", "%s")`, headerName, headerValue)
			// Try empty placeholder first
			contentStr = emptyPattern.ReplaceAllString(contentStr, newHeaderLine)
			if oldContent != contentStr {
				updated = true
			} else {
				// Try existing hardcoded value
				contentStr = existingPattern.ReplaceAllString(contentStr, newHeaderLine)
				if oldContent != contentStr {
					updated = true
				}
			}
		}

		if !updated {
			fmt.Printf("⚠️  WARNING: No match found for header: %s\n", headerName)
		} else {
			fmt.Printf("✅ Updated header: %s = %s\n", headerName, headerValue[:min(20, len(headerValue))])
		}
	}

	// Add missing headers that are required (get from parsed curl request)
	missingHeaders := map[string]string{
		"Priority":           parsedData.Headers["priority"],
		"Sec-Ch-Ua":          parsedData.Headers["sec-ch-ua"],
		"Sec-Ch-Ua-Mobile":   parsedData.Headers["sec-ch-ua-mobile"],
		"Sec-Ch-Ua-Platform": parsedData.Headers["sec-ch-ua-platform"],
		"Sec-Fetch-Dest":     parsedData.Headers["sec-fetch-dest"],
		"Sec-Fetch-Mode":     parsedData.Headers["sec-fetch-mode"],
		"Sec-Fetch-Site":     parsedData.Headers["sec-fetch-site"],
	}

	for headerName, headerValue := range missingHeaders {
		// Check if header already exists (both with quotes and backticks)
		existingPattern := regexp.MustCompile(fmt.Sprintf(`req\.Header\.Set\("%s", "[^"]*"\)`, headerName))
		backtickPattern := regexp.MustCompile(fmt.Sprintf(`req\.Header\.Set\("%s", `+"`[^`]*`"+`\)`, headerName))

		oldContent := contentStr
		updated := false

		if existingPattern.MatchString(contentStr) || backtickPattern.MatchString(contentStr) {
			// Update existing header
			if strings.Contains(headerValue, `"`) {
				newHeaderLine := fmt.Sprintf(`req.Header.Set("%s", `+"`%s`"+`)`, headerName, headerValue)
				contentStr = existingPattern.ReplaceAllString(contentStr, newHeaderLine)
				contentStr = backtickPattern.ReplaceAllString(contentStr, newHeaderLine)
			} else {
				newHeaderLine := fmt.Sprintf(`req.Header.Set("%s", "%s")`, headerName, headerValue)
				contentStr = existingPattern.ReplaceAllString(contentStr, newHeaderLine)
				contentStr = backtickPattern.ReplaceAllString(contentStr, newHeaderLine)
			}
			if oldContent != contentStr {
				updated = true
			}
		}

		if !updated {
			// Add new header after the existing headers
			insertPoint := `req.Header.Set("Sec-Fetch-Site", "same-site")`
			if strings.Contains(headerValue, `"`) {
				newHeaderLine := fmt.Sprintf(`	req.Header.Set("%s", `+"`%s`"+`)`, headerName, headerValue)
				contentStr = strings.Replace(contentStr, insertPoint, insertPoint+"\n"+newHeaderLine, 1)
			} else {
				newHeaderLine := fmt.Sprintf(`	req.Header.Set("%s", "%s")`, headerName, headerValue)
				contentStr = strings.Replace(contentStr, insertPoint, insertPoint+"\n"+newHeaderLine, 1)
			}
			if oldContent != contentStr {
				updated = true
			}
		}

		if updated {
			fmt.Printf("✅ Updated/Added header: %s = %s\n", headerName, headerValue)
		} else {
			fmt.Printf("⚠️  WARNING: Could not update/add header: %s\n", headerName)
		}
	}

	// Remove problematic headers that should not be present
	headersToRemove := []string{
		`req.Header.Set("Accept-Encoding", "")`,
		`req.Header.Set("Connection", "")`,
		`req.Header.Set("TE", "")`,
	}

	for _, headerToRemove := range headersToRemove {
		oldContent := contentStr
		contentStr = strings.Replace(contentStr, headerToRemove+"\n", "", 1)
		contentStr = strings.Replace(contentStr, headerToRemove, "", 1)
		if oldContent != contentStr {
			fmt.Printf("✅ Removed problematic header: %s\n", headerToRemove)
		}
	}

	// Write the updated content back to the file
	err = os.WriteFile("utils/uploadimage.go", []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
