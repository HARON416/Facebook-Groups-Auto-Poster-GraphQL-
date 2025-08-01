package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func UpdateCreateGroupPostCookies() error {
	// Read the curl command from file
	curlData, err := readCreateGroupPostCurlRequest("cookies/creategrouppostcookies.txt")
	if err != nil {
		return fmt.Errorf("error reading curl request: %w", err)
	}

	// Parse the curl command
	parsedData, err := parseCreateGroupPostCurlRequest(curlData)
	if err != nil {
		return fmt.Errorf("error parsing curl request: %w", err)
	}

	// Update the creategrouppost.go file
	err = updateCreateGroupPostFile(parsedData)
	if err != nil {
		return fmt.Errorf("error updating file: %w", err)
	}

	fmt.Println("✅ CreateGroupPost cookies updated successfully!")
	return nil
}

func readCreateGroupPostCurlRequest(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return strings.Join(lines, "\n"), nil
}

type CreateGroupPostParsedData struct {
	URLParams map[string]string
	Headers   map[string]string
	Cookies   map[string]string
	Variables map[string]interface{}
}

func parseCreateGroupPostCurlRequest(curlData string) (*CreateGroupPostParsedData, error) {
	data := &CreateGroupPostParsedData{
		URLParams: make(map[string]string),
		Headers:   make(map[string]string),
		Cookies:   make(map[string]string),
		Variables: make(map[string]interface{}),
	}

	lines := strings.Split(curlData, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse headers (-H) - handle both single and double quotes
		if strings.HasPrefix(line, "-H") {
			// Try single quotes first
			headerMatch := regexp.MustCompile(`-H\s+'([^:]+):\s*([^']*)'`).FindStringSubmatch(line)
			if len(headerMatch) == 3 {
				key := strings.TrimSpace(headerMatch[1])
				value := strings.TrimSpace(headerMatch[2])
				data.Headers[key] = value
				fmt.Printf("🔍 DEBUG: Parsed header %s = %s\n", key, value)
			} else {
				// Try double quotes
				headerMatch = regexp.MustCompile(`-H\s+"([^:]+):\s*([^"]*)"`).FindStringSubmatch(line)
				if len(headerMatch) == 3 {
					key := strings.TrimSpace(headerMatch[1])
					value := strings.TrimSpace(headerMatch[2])
					data.Headers[key] = value
					fmt.Printf("🔍 DEBUG: Parsed header %s = %s\n", key, value)
				}
			}
		}

		// Parse cookies (-b) - handle both single and double quotes
		if strings.HasPrefix(line, "-b") {
			// Try single quotes first
			cookieMatch := regexp.MustCompile(`-b\s+'([^']*)'`).FindStringSubmatch(line)
			if len(cookieMatch) == 2 {
				cookies := cookieMatch[1]
				parseCookies(cookies, data.Cookies)
			} else {
				// Try double quotes
				cookieMatch = regexp.MustCompile(`-b\s+"([^"]*)"`).FindStringSubmatch(line)
				if len(cookieMatch) == 2 {
					cookies := cookieMatch[1]
					parseCookies(cookies, data.Cookies)
				}
			}
		}

		// Parse URL parameters (from --data-raw) - handle both single and double quotes
		if strings.Contains(line, "--data-raw") {
			// Try single quotes first
			dataMatch := regexp.MustCompile(`--data-raw\s+'([^']*)'`).FindStringSubmatch(line)
			if len(dataMatch) == 2 {
				rawData := dataMatch[1]
				parseURLParams(rawData, data.URLParams)
			} else {
				// Try double quotes
				dataMatch = regexp.MustCompile(`--data-raw\s+"([^"]*)"`).FindStringSubmatch(line)
				if len(dataMatch) == 2 {
					rawData := dataMatch[1]
					parseURLParams(rawData, data.URLParams)
				}
			}
		}
	}

	return data, nil
}

func parseCookies(cookieString string, cookieMap map[string]string) {
	cookies := strings.Split(cookieString, ";")
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if cookie == "" {
			continue
		}

		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			cookieMap[key] = value
		}
	}
}

func parseURLParams(rawData string, paramMap map[string]string) {
	// Handle double encoding first - replace %% with % to get the original encoded value
	rawData = strings.ReplaceAll(rawData, "%%", "%")

	// Split by & to get individual parameters (don't decode the entire string)
	params := strings.Split(rawData, "&")
	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			// Handle double encoding in the value as well
			value = strings.ReplaceAll(value, "%%", "%")

			// Keep the value as-is (URL encoded) since that's what we want in the request body
			paramMap[key] = value
			fmt.Printf("🔍 DEBUG: Parsed param %s = %s\n", key, value)
		}
	}
}

func updateCreateGroupPostFile(parsedData *CreateGroupPostParsedData) error {
	// Read the current file
	content, err := os.ReadFile("utils/creategrouppost.go")
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fmt.Printf("🔍 DEBUG: Starting update with %d form data items and %d headers\n", len(parsedData.URLParams), len(parsedData.Headers))

	contentStr := string(content)

	// Ensure proper file structure before updates
	contentStr = ensureProperFileStructure(contentStr)

	// 1. Update the request body with exact values from curl
	contentStr = updateRequestBody(contentStr, parsedData.URLParams)

	// 2. Update all headers with exact values from curl
	contentStr = updateAllHeaders(contentStr, parsedData.Headers)

	// 3. Update cookies
	contentStr = updateCookies(contentStr, parsedData.Cookies)

	// 4. Update variables JSON values
	contentStr = updateAllVariables(contentStr, parsedData)

	// Write the updated content back
	err = os.WriteFile("utils/creategrouppost.go", []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func updateRequestBody(content string, urlParams map[string]string) string {
	// Build the exact request body string from the curl command with proper URL encoding
	// and escape % characters for fmt.Sprintf format string

	// Define required parameters and their default values
	requiredParams := map[string]string{
		"av":                       "",
		"__aaid":                   "",
		"__user":                   "",
		"__a":                      "",
		"__req":                    "",
		"__hs":                     "",
		"dpr":                      "",
		"__ccg":                    "",
		"__rev":                    "",
		"__s":                      "",
		"__hsi":                    "",
		"__dyn":                    "",
		"__csr":                    "",
		"__hsdp":                   "",
		"__hblp":                   "",
		"__sjsp":                   "",
		"__comet_req":              "",
		"fb_dtsg":                  "",
		"jazoest":                  "",
		"lsd":                      "",
		"__spin_r":                 "",
		"__spin_b":                 "",
		"__spin_t":                 "",
		"__crn":                    "",
		"fb_api_caller_class":      "",
		"fb_api_req_friendly_name": "",
		"server_timestamps":        "",
		"doc_id":                   "",
	}

	// Use values from curl request, fallback to defaults if missing
	for key := range requiredParams {
		if value, exists := urlParams[key]; exists && value != "" {
			requiredParams[key] = value
			fmt.Printf("✅ Using curl value for %s: %s\n", key, value)
		} else {
			fmt.Printf("⚠️  Missing or empty value for %s, using default\n", key)
		}
	}

	// Debug: Check for encoding issues in specific parameters
	encodingIssues := []string{"__hs", "__s", "fb_dtsg"}
	for _, param := range encodingIssues {
		if value, exists := requiredParams[param]; exists && value != "" {
			if strings.Contains(value, "%%") {
				fmt.Printf("🔧 Fixed double encoding in %s: %s\n", param, value)
			} else {
				fmt.Printf("✅ No encoding issues in %s: %s\n", param, value)
			}
		}
	}

	params := []string{
		fmt.Sprintf("av=%s", fixURLEncoding(requiredParams["av"])),
		fmt.Sprintf("__aaid=%s", fixURLEncoding(requiredParams["__aaid"])),
		fmt.Sprintf("__user=%s", fixURLEncoding(requiredParams["__user"])),
		fmt.Sprintf("__a=%s", fixURLEncoding(requiredParams["__a"])),
		fmt.Sprintf("__req=%s", fixURLEncoding(requiredParams["__req"])),
		fmt.Sprintf("__hs=%s", fixURLEncoding(requiredParams["__hs"])),
		fmt.Sprintf("dpr=%s", fixURLEncoding(requiredParams["dpr"])),
		fmt.Sprintf("__ccg=%s", fixURLEncoding(requiredParams["__ccg"])),
		fmt.Sprintf("__rev=%s", fixURLEncoding(requiredParams["__rev"])),
		fmt.Sprintf("__s=%s", fixURLEncoding(requiredParams["__s"])),
		fmt.Sprintf("__hsi=%s", fixURLEncoding(requiredParams["__hsi"])),
		fmt.Sprintf("__dyn=%s", fixURLEncoding(requiredParams["__dyn"])),
		fmt.Sprintf("__csr=%s", fixURLEncoding(requiredParams["__csr"])),
		fmt.Sprintf("__hsdp=%s", fixURLEncoding(requiredParams["__hsdp"])),
		fmt.Sprintf("__hblp=%s", fixURLEncoding(requiredParams["__hblp"])),
		fmt.Sprintf("__sjsp=%s", fixURLEncoding(requiredParams["__sjsp"])),
		fmt.Sprintf("__comet_req=%s", fixURLEncoding(requiredParams["__comet_req"])),
		fmt.Sprintf("fb_dtsg=%s", fixURLEncoding(requiredParams["fb_dtsg"])),
		fmt.Sprintf("jazoest=%s", fixURLEncoding(requiredParams["jazoest"])),
		fmt.Sprintf("lsd=%s", fixURLEncoding(requiredParams["lsd"])),
		fmt.Sprintf("__spin_r=%s", fixURLEncoding(requiredParams["__spin_r"])),
		fmt.Sprintf("__spin_b=%s", fixURLEncoding(requiredParams["__spin_b"])),
		fmt.Sprintf("__spin_t=%s", fixURLEncoding(requiredParams["__spin_t"])),
		fmt.Sprintf("__crn=%s", fixURLEncoding(requiredParams["__crn"])),
		fmt.Sprintf("fb_api_caller_class=%s", fixURLEncoding(requiredParams["fb_api_caller_class"])),
		fmt.Sprintf("fb_api_req_friendly_name=%s", fixURLEncoding(requiredParams["fb_api_req_friendly_name"])),
		"variables=%s",
		fmt.Sprintf("server_timestamps=%s", fixURLEncoding(requiredParams["server_timestamps"])),
		fmt.Sprintf("doc_id=%s", fixURLEncoding(requiredParams["doc_id"])),
	}

	requestBody := strings.Join(params, "&")

	// Replace the entire requestBody line with proper formatting
	pattern := regexp.MustCompile(`requestBody := fmt\.Sprintf\(".*", encodedVariables\)`)
	replacement := fmt.Sprintf(`	requestBody := fmt.Sprintf("%s", encodedVariables)`, requestBody)

	updated := pattern.ReplaceAllString(content, replacement)
	fmt.Printf("✅ Updated request body with exact values from curl (with proper URL encoding and %% escaping)\n")

	// Validate the encoding fix
	validateEncodingFix(updated)

	return updated
}

// escapePercentSigns escapes % characters in URL-encoded values for fmt.Sprintf format strings
// escapePercentSigns escapes % characters for fmt.Sprintf format strings
// but preserves valid URL encoding sequences like %3A, %2F, etc.
func escapePercentSigns(value string) string {
	// For fmt.Sprintf, we need to escape % characters, but we want to preserve the original encoding
	// So we'll use the original value if it's already properly encoded, otherwise encode it
	if strings.Contains(value, "%") && !strings.Contains(value, "%%") {
		// Value is already URL encoded, just escape % for fmt.Sprintf
		return strings.ReplaceAll(value, "%", "%%")
	} else {
		// Value needs to be URL encoded first, then escaped for fmt.Sprintf
		return strings.ReplaceAll(value, "%", "%%")
	}
}

// fixURLEncoding fixes the double encoding issue by properly handling URL encoding
func fixURLEncoding(value string) string {
	// The value from curl request is already URL encoded, we just need to escape % for fmt.Sprintf
	// Handle double encoding first - replace %% with % to get the original encoded value
	value = strings.ReplaceAll(value, "%%", "%")

	// Simply escape % for fmt.Sprintf without re-encoding
	return escapePercentSigns(value)
}

// validateEncodingFix checks if the encoding issues have been resolved
func validateEncodingFix(content string) {
	// Check for proper fmt.Sprintf escaping - %%3A is correct for fmt.Sprintf
	patterns := []string{
		`__hs=.*%%3A`,
		`__s=.*%%3A`,
		`fb_dtsg=.*%%3A`,
	}

	foundIssues := false
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(content) {
			fmt.Printf("✅ Found correct fmt.Sprintf escaping in pattern: %s\n", pattern)
		} else {
			fmt.Printf("⚠️  Missing proper fmt.Sprintf escaping in pattern: %s\n", pattern)
			foundIssues = true
		}
	}

	if !foundIssues {
		fmt.Printf("🎉 All encoding issues have been resolved!\n")
	} else {
		fmt.Printf("🔧 Some encoding issues still need attention\n")
	}
}

func updateAllHeaders(content string, headers map[string]string) string {
	// Find the header section and replace all headers at once
	headerStart := strings.Index(content, "// Set headers from the curl request")
	if headerStart == -1 {
		return content
	}

	// Find the end of the header section (before cookies)
	headerEnd := strings.Index(content[headerStart:], "// Set cookies from the curl request")
	if headerEnd == -1 {
		return content
	}
	headerEnd += headerStart

	// Build all header lines with proper formatting
	var headerLines []string
	headerLines = append(headerLines, "	// Set headers from the curl request - exact values")

	// Always add Content-Type header first
	headerLines = append(headerLines, "	req.Header.Set(\"Content-Type\", \"application/x-www-form-urlencoded\")")

	for key, value := range headers {
		// Skip Content-Type as we already added it above
		if strings.ToLower(key) == "content-type" {
			continue
		}

		// Use backticks for headers with quotes to avoid escaping issues
		if strings.Contains(value, `"`) {
			headerLines = append(headerLines, fmt.Sprintf("	req.Header.Set(\"%s\", `%s`)", key, value))
		} else {
			headerLines = append(headerLines, fmt.Sprintf("	req.Header.Set(\"%s\", \"%s\")", key, value))
		}
		fmt.Printf("✅ Updated header: %s = %s\n", key, value)
	}

	// Add a blank line before the cookie section
	headerLines = append(headerLines, "")

	// Replace the entire header section
	oldHeaderSection := content[headerStart:headerEnd]
	newHeaderSection := strings.Join(headerLines, "\n")

	updated := strings.Replace(content, oldHeaderSection, newHeaderSection, 1)
	fmt.Printf("✅ Updated all %d headers\n", len(headers))

	return updated
}

func updateCookies(content string, cookies map[string]string) string {
	cookieString := buildCookieString(cookies)

	// Replace the cookie line with proper formatting
	pattern := regexp.MustCompile(`req\.Header\.Set\("Cookie",\s*"[^"]*"\)`)
	replacement := fmt.Sprintf("	req.Header.Set(\"Cookie\", \"%s\")", cookieString)

	updated := pattern.ReplaceAllString(content, replacement)
	fmt.Printf("✅ Updated cookies\n")

	return updated
}

func updateAllVariables(content string, parsedData *CreateGroupPostParsedData) string {
	// Extract variables from the parsed data
	variablesStr, exists := parsedData.URLParams["variables"]
	if !exists {
		return content
	}

	// URL decode the variables
	decoded, err := url.QueryUnescape(variablesStr)
	if err != nil {
		return content
	}

	var variables map[string]interface{}
	if json.Unmarshal([]byte(decoded), &variables) != nil {
		return content
	}

	// Extract values from the variables JSON
	input, ok := variables["input"].(map[string]interface{})
	if !ok {
		return content
	}

	// Update ActorID from variables JSON (not from av parameter)
	if input, ok := variables["input"].(map[string]interface{}); ok {
		if actorID, ok := input["actor_id"].(string); ok {
			pattern := regexp.MustCompile(`ActorID:\s*"[^"]*"`)
			replacement := fmt.Sprintf(`ActorID: "%s"`, actorID)
			content = pattern.ReplaceAllString(content, replacement)
			fmt.Printf("✅ Updated ActorID: %s\n", actorID)
		}
	}

	// Update ClientMutationID from variables JSON (not from __req parameter)
	if input, ok := variables["input"].(map[string]interface{}); ok {
		if mutationID, ok := input["client_mutation_id"].(string); ok {
			pattern := regexp.MustCompile(`ClientMutationID:\s*"[^"]*"`)
			replacement := fmt.Sprintf(`ClientMutationID: "%s"`, mutationID)
			content = pattern.ReplaceAllString(content, replacement)
			fmt.Printf("✅ Updated ClientMutationID: %s\n", mutationID)
		}
	}

	// Update ComposerSessionID
	if logging, ok := input["logging"].(map[string]interface{}); ok {
		if sessionID, ok := logging["composer_session_id"].(string); ok {
			pattern := regexp.MustCompile(`ComposerSessionID:\s*"[^"]*"`)
			replacement := fmt.Sprintf(`ComposerSessionID: "%s"`, sessionID)
			content = pattern.ReplaceAllString(content, replacement)
			fmt.Printf("✅ Updated ComposerSessionID: %s\n", sessionID)
		}
	}

	// Update ComposerEntryPoint
	if entryPoint, ok := input["composer_entry_point"].(string); ok {
		pattern := regexp.MustCompile(`ComposerEntryPoint:\s*"[^"]*"`)
		replacement := fmt.Sprintf(`ComposerEntryPoint: "%s"`, entryPoint)
		content = pattern.ReplaceAllString(content, replacement)
		fmt.Printf("✅ Updated ComposerEntryPoint: %s\n", entryPoint)
	}

	// Update AttributionIDV2
	if navigationData, ok := input["navigation_data"].(map[string]interface{}); ok {
		if attributionID, ok := navigationData["attribution_id_v2"].(string); ok {
			pattern := regexp.MustCompile(`AttributionIDV2:\s*"[^"]*"`)
			replacement := fmt.Sprintf(`AttributionIDV2: "%s"`, attributionID)
			content = pattern.ReplaceAllString(content, replacement)
			fmt.Printf("✅ Updated AttributionIDV2: %s\n", attributionID)
		}
	}

	// Update EventShareMetadata.Surface from variables JSON
	if input, ok := variables["input"].(map[string]interface{}); ok {
		if eventShareMetadata, ok := input["event_share_metadata"].(map[string]interface{}); ok {
			if surface, ok := eventShareMetadata["surface"].(string); ok {
				pattern := regexp.MustCompile(`Surface:\s*"[^"]*"`)
				replacement := fmt.Sprintf(`Surface: "%s"`, surface)
				content = pattern.ReplaceAllString(content, replacement)
				fmt.Printf("✅ Updated Surface: %s\n", surface)
			}
		}
	}

	// Update static values
	updates := map[string]string{
		`ComposerSourceSurface:\s*"[^"]*"`: `ComposerSourceSurface: "group"`,
		`ComposerType:\s*"[^"]*"`:          `ComposerType: "group"`,
		`Source:\s*"[^"]*"`:                `Source: "WWW"`,
	}

	for pattern, replacement := range updates {
		re := regexp.MustCompile(pattern)
		content = re.ReplaceAllString(content, replacement)
	}

	// 🔧 FIX STRUCTURAL ISSUES
	// Fix 1: Remove ComposedText structure from struct definition
	composedTextStructPattern := regexp.MustCompile(`(?s)ComposedText\s+struct\s*\{[^}]*\}\s*` + "`json:\"composed_text\"`")
	content = composedTextStructPattern.ReplaceAllString(content, "")
	fmt.Printf("🔧 Removed ComposedText structure from struct definition\n")

	// Fix 2: Remove ComposedText field assignment (more comprehensive pattern)
	composedTextAssignmentPattern := regexp.MustCompile(`(?s)ComposedText:\s*struct\s*\{[^}]*\}\s*\{[^}]*\},?`)
	content = composedTextAssignmentPattern.ReplaceAllString(content, "")
	fmt.Printf("🔧 Removed ComposedText field assignment\n")

	// Fix 3: Fix TextFormatPresetID to be "0" instead of empty string
	textFormatPattern := regexp.MustCompile(`TextFormatPresetID:\s*""`)
	content = textFormatPattern.ReplaceAllString(content, `TextFormatPresetID: "0"`)
	fmt.Printf("🔧 Fixed TextFormatPresetID to be \"0\"\n")

	// Fix 4: Ensure WithTagsIDs is nil (not empty array)
	withTagsPattern := regexp.MustCompile(`WithTagsIDs:\s*\[\]interface\{\}\{\}`)
	content = withTagsPattern.ReplaceAllString(content, `WithTagsIDs: nil`)
	fmt.Printf("🔧 Fixed WithTagsIDs to be nil\n")

	// Fix 5: Ensure InlineActivities is empty array
	inlineActivitiesPattern := regexp.MustCompile(`InlineActivities:\s*nil`)
	content = inlineActivitiesPattern.ReplaceAllString(content, `InlineActivities: []interface{}{}`)
	fmt.Printf("🔧 Fixed InlineActivities to be empty array\n")

	// Fix 6: Remove any remaining ComposedText references (catch-all)
	remainingComposedTextPattern := regexp.MustCompile(`(?s)ComposedText:\s*struct\s*\{[^}]*\}\s*\{[^}]*\}`)
	content = remainingComposedTextPattern.ReplaceAllString(content, "")
	fmt.Printf("🔧 Removed any remaining ComposedText references\n")

	fmt.Printf("✅ Updated all variables JSON values and fixed structural issues\n")
	return content
}

func buildCookieString(cookies map[string]string) string {
	var cookiePairs []string
	for key, value := range cookies {
		cookiePairs = append(cookiePairs, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(cookiePairs, "; ")
}

func ensureProperFileStructure(content string) string {
	// Ensure proper header section structure
	headerSectionPattern := regexp.MustCompile(`// Set headers from the curl request[^\n]*\n([^\n]*\n)*?// Set cookies from the curl request`)
	if !headerSectionPattern.MatchString(content) {
		// If header section is malformed, fix it
		content = strings.ReplaceAll(content, "// Set cookies from the curl request - empty placeholder", "\n	// Set cookies from the curl request - exact values")
	}

	return content
}
