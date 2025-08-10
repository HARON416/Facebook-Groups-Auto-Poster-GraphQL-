// package utils

// import (
// 	"fmt"
// 	"net/url"
// 	"regexp"
// )

// // FacebookConfig holds all the dynamic values needed for Facebook API requests
// type FacebookConfig struct {
// 	// POST data fields
// 	AV      string
// 	User    string
// 	Req     string
// 	HS      string
// 	CCG     string
// 	Rev     string
// 	S       string
// 	HSI     string
// 	Dyn     string
// 	CSR     string
// 	HSDP    string
// 	HBLP    string
// 	SJSP    string
// 	FBDtsg  string
// 	Jazoest string
// 	LSD     string
// 	SpinT   string

// 	// Header fields
// 	Cookie                 string
// 	UserAgent              string
// 	SecChUA                string
// 	SecChUAFullVersionList string
// }

// // ParseCurlRequest extracts Facebook API parameters from a curl command
// func ParseCurlRequest(curlCommand string) (*FacebookConfig, error) {
// 	config := &FacebookConfig{}

// 	// Extract cookies from -b flag
// 	cookieRegex := regexp.MustCompile(`-b\s+'([^']+)'`)
// 	if matches := cookieRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
// 		config.Cookie = matches[1]
// 	}

// 	// Extract user-agent from -H flag
// 	userAgentRegex := regexp.MustCompile(`-H\s+'user-agent:\s*([^']+)'`)
// 	if matches := userAgentRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
// 		config.UserAgent = matches[1]
// 	}

// 	// Extract sec-ch-ua from -H flag
// 	secChUARegex := regexp.MustCompile(`-H\s+'sec-ch-ua:\s*([^']+)'`)
// 	if matches := secChUARegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
// 		config.SecChUA = matches[1]
// 	}

// 	// Extract sec-ch-ua-full-version-list from -H flag
// 	secChUAFullRegex := regexp.MustCompile(`-H\s+'sec-ch-ua-full-version-list:\s*([^']+)'`)
// 	if matches := secChUAFullRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
// 		config.SecChUAFullVersionList = matches[1]
// 	}

// 	// Extract POST data from --data-raw
// 	dataRegex := regexp.MustCompile(`--data-raw\s+'([^']+)'`)
// 	if matches := dataRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
// 		postData := matches[1]

// 		// Parse URL-encoded data
// 		values, err := url.ParseQuery(postData)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to parse POST data: %w", err)
// 		}

// 		// Extract individual values
// 		config.AV = values.Get("av")
// 		config.User = values.Get("__user")
// 		config.Req = values.Get("__req")
// 		config.HS = values.Get("__hs")
// 		config.CCG = values.Get("__ccg")
// 		config.Rev = values.Get("__rev")
// 		config.S = values.Get("__s")
// 		config.HSI = values.Get("__hsi")
// 		config.Dyn = values.Get("__dyn")
// 		config.CSR = values.Get("__csr")
// 		config.HSDP = values.Get("__hsdp")
// 		config.HBLP = values.Get("__hblp")
// 		config.SJSP = values.Get("__sjsp")
// 		config.FBDtsg = values.Get("fb_dtsg")
// 		config.Jazoest = values.Get("jazoest")
// 		config.LSD = values.Get("lsd")
// 		config.SpinT = values.Get("__spin_t")
// 	}

// 	return config, nil
// }

// // UpdateFetchGroupsCode generates updated Go code with new configuration
// func UpdateFetchGroupsCode(config *FacebookConfig, originalCode string) string {
// 	// Define replacements for POST data (only update if value is present)
// 	postDataReplacements := map[string]string{}

// 	if config.AV != "" {
// 		postDataReplacements[`data.Set\("av", "[^"]+"\)`] = fmt.Sprintf(`data.Set("av", "%s")`, config.AV)
// 	}
// 	if config.User != "" {
// 		postDataReplacements[`data.Set\("__user", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__user", "%s")`, config.User)
// 	}
// 	if config.Req != "" {
// 		postDataReplacements[`data.Set\("__req", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__req", "%s")`, config.Req)
// 	}
// 	if config.HS != "" {
// 		postDataReplacements[`data.Set\("__hs", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__hs", "%s")`, config.HS)
// 	}
// 	if config.CCG != "" {
// 		postDataReplacements[`data.Set\("__ccg", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__ccg", "%s")`, config.CCG)
// 	}
// 	if config.Rev != "" {
// 		postDataReplacements[`data.Set\("__rev", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__rev", "%s")`, config.Rev)
// 	}
// 	if config.S != "" {
// 		postDataReplacements[`data.Set\("__s", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__s", "%s")`, config.S)
// 	}
// 	if config.HSI != "" {
// 		postDataReplacements[`data.Set\("__hsi", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__hsi", "%s")`, config.HSI)
// 	}
// 	if config.Dyn != "" {
// 		postDataReplacements[`data.Set\("__dyn", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__dyn", "%s")`, config.Dyn)
// 	}
// 	if config.CSR != "" {
// 		postDataReplacements[`data.Set\("__csr", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__csr", "%s")`, config.CSR)
// 	}
// 	if config.HSDP != "" {
// 		postDataReplacements[`data.Set\("__hsdp", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__hsdp", "%s")`, config.HSDP)
// 	}
// 	if config.HBLP != "" {
// 		postDataReplacements[`data.Set\("__hblp", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__hblp", "%s")`, config.HBLP)
// 	}
// 	if config.SJSP != "" {
// 		postDataReplacements[`data.Set\("__sjsp", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__sjsp", "%s")`, config.SJSP)
// 	}
// 	if config.FBDtsg != "" {
// 		postDataReplacements[`data.Set\("fb_dtsg", "[^"]+"\)`] = fmt.Sprintf(`data.Set("fb_dtsg", "%s")`, config.FBDtsg)
// 	}
// 	if config.Jazoest != "" {
// 		postDataReplacements[`data.Set\("jazoest", "[^"]+"\)`] = fmt.Sprintf(`data.Set("jazoest", "%s")`, config.Jazoest)
// 	}
// 	if config.LSD != "" {
// 		postDataReplacements[`data.Set\("lsd", "[^"]+"\)`] = fmt.Sprintf(`data.Set("lsd", "%s")`, config.LSD)
// 	}
// 	if config.SpinT != "" {
// 		postDataReplacements[`data.Set\("__spin_t", "[^"]+"\)`] = fmt.Sprintf(`data.Set("__spin_t", "%s")`, config.SpinT)
// 	}

// 	// Define replacements for headers (only update if value is present)
// 	headerReplacements := map[string]string{}

// 	if config.Cookie != "" {
// 		headerReplacements[`req.Header.Set\("cookie", "[^"]+"\)`] = fmt.Sprintf(`req.Header.Set("cookie", "%s")`, config.Cookie)
// 	}
// 	if config.UserAgent != "" {
// 		headerReplacements[`req.Header.Set\("user-agent", "[^"]+"\)`] = fmt.Sprintf(`req.Header.Set("user-agent", "%s")`, config.UserAgent)
// 	}
// 	if config.SecChUA != "" {
// 		headerReplacements[`req.Header.Set\("sec-ch-ua", `+"`"+`[^`+"`"+`]+`+"`"+`\)`] = fmt.Sprintf(`req.Header.Set("sec-ch-ua", `+"`"+`%s`+"`"+`)`, config.SecChUA)
// 	}
// 	if config.SecChUAFullVersionList != "" {
// 		headerReplacements[`req.Header.Set\("sec-ch-ua-full-version-list", `+"`"+`[^`+"`"+`]+`+"`"+`\)`] = fmt.Sprintf(`req.Header.Set("sec-ch-ua-full-version-list", `+"`"+`%s`+"`"+`)`, config.SecChUAFullVersionList)
// 	}
// 	if config.LSD != "" {
// 		headerReplacements[`req.Header.Set\("x-fb-lsd", "[^"]+"\)`] = fmt.Sprintf(`req.Header.Set("x-fb-lsd", "%s")`, config.LSD)
// 	}

// 	updatedCode := originalCode

// 	// Apply all replacements
// 	allReplacements := make(map[string]string)
// 	for k, v := range postDataReplacements {
// 		allReplacements[k] = v
// 	}
// 	for k, v := range headerReplacements {
// 		allReplacements[k] = v
// 	}

// 	for pattern, replacement := range allReplacements {
// 		re := regexp.MustCompile(pattern)
// 		updatedCode = re.ReplaceAllString(updatedCode, replacement)
// 	}

// 	return updatedCode
// }

// // ValidateConfig checks if we have the minimum required fields
// func ValidateConfig(config *FacebookConfig) error {
// 	missing := []string{}

// 	if config.AV == "" {
// 		missing = append(missing, "av")
// 	}
// 	if config.User == "" {
// 		missing = append(missing, "__user")
// 	}
// 	if config.Cookie == "" {
// 		missing = append(missing, "cookie")
// 	}
// 	if config.LSD == "" {
// 		missing = append(missing, "lsd")
// 	}

// 	if len(missing) > 0 {
// 		return fmt.Errorf("missing essential fields from curl command: %v", missing)
// 	}

// 	// Warn about optional but important missing fields
// 	warnings := []string{}
// 	if config.CSR == "" {
// 		warnings = append(warnings, "__csr")
// 	}
// 	if config.HSDP == "" {
// 		warnings = append(warnings, "__hsdp")
// 	}
// 	if config.HBLP == "" {
// 		warnings = append(warnings, "__hblp")
// 	}
// 	if config.SJSP == "" {
// 		warnings = append(warnings, "__sjsp")
// 	}

// 	if len(warnings) > 0 {
// 		fmt.Printf("Warning: Some fields missing from curl command (will keep existing values): %v\n", warnings)
// 	}

// 	return nil
// }

// // UpdateFetchGroupsFromCurl updates the code but warns about missing fields
// func UpdateFetchGroupsFromCurl(curlCommand string) error {
// 	// Parse the curl command
// 	config, err := ParseCurlRequest(curlCommand)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse curl command: %w", err)
// 	}

// 	// Validate configuration and warn about missing fields
// 	if err := ValidateConfig(config); err != nil {
// 		return err
// 	}

// 	fmt.Printf("✓ Extracted configuration:\n")
// 	fmt.Printf("  User ID: %s\n", config.User)
// 	fmt.Printf("  AV: %s\n", config.AV)
// 	fmt.Printf("  LSD: %s\n", config.LSD)
// 	fmt.Printf("  Cookie length: %d characters\n", len(config.Cookie))

// 	fmt.Println("✅ Configuration extracted successfully!")
// 	return nil
// }

// // ExtractEssentialValues helper function to extract just the essential values for quick updates
// func ExtractEssentialValues(curlCommand string) (userID, av, lsd, cookie string, err error) {
// 	config, err := ParseCurlRequest(curlCommand)
// 	if err != nil {
// 		return "", "", "", "", err
// 	}

// 	return config.User, config.AV, config.LSD, config.Cookie, nil
// }

// // QuickUpdateFetchGroups provides a simple way to update just the most critical values
// func QuickUpdateFetchGroups(curlCommand, originalCode string) (string, error) {
// 	userID, av, lsd, cookie, err := ExtractEssentialValues(curlCommand)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Quick replacements for the most essential values
// 	updatedCode := originalCode

// 	// Replace user ID in both av and __user fields
// 	updatedCode = regexp.MustCompile(`data\.Set\("av", "[^"]+"\)`).ReplaceAllString(updatedCode, fmt.Sprintf(`data.Set("av", "%s")`, av))
// 	updatedCode = regexp.MustCompile(`data\.Set\("__user", "[^"]+"\)`).ReplaceAllString(updatedCode, fmt.Sprintf(`data.Set("__user", "%s")`, userID))

// 	// Replace LSD in both data and header
// 	updatedCode = regexp.MustCompile(`data\.Set\("lsd", "[^"]+"\)`).ReplaceAllString(updatedCode, fmt.Sprintf(`data.Set("lsd", "%s")`, lsd))
// 	updatedCode = regexp.MustCompile(`req\.Header\.Set\("x-fb-lsd", "[^"]+"\)`).ReplaceAllString(updatedCode, fmt.Sprintf(`req.Header.Set("x-fb-lsd", "%s")`, lsd))

// 	// Replace cookie
// 	updatedCode = regexp.MustCompile(`req\.Header\.Set\("cookie", "[^"]+"\)`).ReplaceAllString(updatedCode, fmt.Sprintf(`req.Header.Set("cookie", "%s")`, cookie))

// 	return updatedCode, nil
// }

package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
)

// Simple function to update Facebook Groups code from curl
func UpdateFacebookGroupsFromCurl(curlCommand, filePath string) error {
	// 1. Read the original file
	originalCode, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 2. Extract values from curl command
	config, err := extractValuesFromCurl(curlCommand)
	if err != nil {
		return fmt.Errorf("failed to extract values from curl: %w", err)
	}

	// 3. Update the code
	updatedCode := updateGoCode(string(originalCode), config)

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

// Extract all Facebook values from curl command
func extractValuesFromCurl(curlCommand string) (map[string]string, error) {
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

	// Extract sec-ch-ua-full-version-list
	secChUAFullRegex := regexp.MustCompile(`-H\s+'sec-ch-ua-full-version-list:\s*([^']+)'`)
	if matches := secChUAFullRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		config["sec-ch-ua-full-version-list"] = matches[1]
	}

	// Extract POST data
	dataRegex := regexp.MustCompile(`--data-raw\s+'([^']+)'`)
	if matches := dataRegex.FindStringSubmatch(curlCommand); len(matches) > 1 {
		postData := matches[1]

		values, err := url.ParseQuery(postData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse POST data: %w", err)
		}

		// Extract all the Facebook fields
		fbFields := []string{
			"av", "__user", "__req", "__hs", "__ccg", "__rev", "__s", "__hsi",
			"__dyn", "__csr", "__hsdp", "__hblp", "__sjsp", "fb_dtsg",
			"jazoest", "lsd", "__spin_t",
		}

		for _, field := range fbFields {
			if value := values.Get(field); value != "" {
				config[field] = value
			}
		}
	}

	// Validate we got the essential fields
	essential := []string{"av", "__user", "cookie", "lsd"}
	for _, field := range essential {
		if config[field] == "" {
			return nil, fmt.Errorf("missing essential field: %s", field)
		}
	}

	return config, nil
}

// Update the Go code with new values
func updateGoCode(originalCode string, config map[string]string) string {
	updatedCode := originalCode

	// Update POST data fields
	postFields := []string{
		"av", "__user", "__req", "__hs", "__ccg", "__rev", "__s", "__hsi",
		"__dyn", "__csr", "__hsdp", "__hblp", "__sjsp", "fb_dtsg",
		"jazoest", "lsd", "__spin_t",
	}

	for _, field := range postFields {
		if value, exists := config[field]; exists {
			pattern := fmt.Sprintf(`data\.Set\("%s", "[^"]+"\)`, regexp.QuoteMeta(field))
			replacement := fmt.Sprintf(`data.Set("%s", "%s")`, field, value)
			re := regexp.MustCompile(pattern)
			updatedCode = re.ReplaceAllString(updatedCode, replacement)
		}
	}

	// Update headers
	if cookie, exists := config["cookie"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("cookie", "[^"]+"\)`)
		replacement := fmt.Sprintf(`req.Header.Set("cookie", "%s")`, cookie)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	if userAgent, exists := config["user-agent"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("user-agent", "[^"]+"\)`)
		replacement := fmt.Sprintf(`req.Header.Set("user-agent", "%s")`, userAgent)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	if secChUA, exists := config["sec-ch-ua"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("sec-ch-ua", ` + "`" + `[^` + "`" + `]+` + "`" + `\)`)
		replacement := fmt.Sprintf(`req.Header.Set("sec-ch-ua", `+"`"+`%s`+"`"+`)`, secChUA)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	if secChUAFull, exists := config["sec-ch-ua-full-version-list"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("sec-ch-ua-full-version-list", ` + "`" + `[^` + "`" + `]+` + "`" + `\)`)
		replacement := fmt.Sprintf(`req.Header.Set("sec-ch-ua-full-version-list", `+"`"+`%s`+"`"+`)`, secChUAFull)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	if lsd, exists := config["lsd"]; exists {
		re := regexp.MustCompile(`req\.Header\.Set\("x-fb-lsd", "[^"]+"\)`)
		replacement := fmt.Sprintf(`req.Header.Set("x-fb-lsd", "%s")`, lsd)
		updatedCode = re.ReplaceAllString(updatedCode, replacement)
	}

	return updatedCode
}

// Command line interface
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run updater.go '<curl-command>' <go-file-path>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println(`go run updater.go 'curl https://web.facebook.com/api/graphql/ ...' facebook_groups.go`)
		os.Exit(1)
	}

	curlCommand := os.Args[1]
	filePath := os.Args[2]

	err := UpdateFacebookGroupsFromCurl(curlCommand, filePath)
	if err != nil {
		log.Fatalf("❌ Update failed: %v", err)
	}

	fmt.Println("🎉 Update completed successfully!")
}

// Alternative: Update from string (no file needed)
func UpdateCodeFromCurl(curlCommand, originalCode string) (string, error) {
	config, err := extractValuesFromCurl(curlCommand)
	if err != nil {
		return "", err
	}

	updatedCode := updateGoCode(originalCode, config)
	return updatedCode, nil
}
