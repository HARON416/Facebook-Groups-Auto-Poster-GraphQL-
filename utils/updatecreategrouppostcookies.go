package utils

import (
	"encoding/json"
	"fmt"
	"net/url"

	"os"
	"regexp"
	"strings"
)

// CurlData holds the extracted data from a cURL request
type CurlData struct {
	UserID             string
	SessionID          string
	AttributionID      string
	LSD                string
	Cookies            string
	Headers            map[string]string
	FormData           map[string]string
	GroupID            string
	MessageText        string
	PhotoIDs           []string
	ComposerEntryPoint string
}

// UpdateCreateGroupPostFunctionFromCurl takes cURL command string and updates Go code file
func UpdateCreateGroupPostFunctionFromCurl(curlCommand string, goFilePath string) error {
	// Parse cURL command
	curlData, err := parseCurlCommand(curlCommand)
	if err != nil {
		return fmt.Errorf("error parsing cURL command: %v", err)
	}

	// Read existing Go file
	goContent, err := os.ReadFile(goFilePath)
	if err != nil {
		return fmt.Errorf("error reading Go file '%s': %v", goFilePath, err)
	}

	// Update Go code with new values
	updatedCode, err := updateCode(string(goContent), curlData)
	if err != nil {
		return fmt.Errorf("error updating Go code: %v", err)
	}

	// Write updated code back to file
	err = os.WriteFile(goFilePath, []byte(updatedCode), 0644)
	if err != nil {
		return fmt.Errorf("error writing updated Go file '%s': %v", goFilePath, err)
	}

	fmt.Printf("✅ Successfully updated %s with fresh cURL values\n", goFilePath)
	return nil
}

// UpdateCodeFromCurl reads a cURL command from a file and updates the Go code
// func UpdateCodeFromCurl(curlFilePath, goFilePath string) error {
// 	// Validate file paths
// 	if len(curlFilePath) > 100 {
// 		return fmt.Errorf("cURL file path too long (max 100 chars): %d chars", len(curlFilePath))
// 	}
// 	if len(goFilePath) > 100 {
// 		return fmt.Errorf("Go file path too long (max 100 chars): %d chars", len(goFilePath))
// 	}

// 	// Read cURL command from file
// 	curlContent, err := ioutil.ReadFile(curlFilePath)
// 	if err != nil {
// 		return fmt.Errorf("error reading cURL file '%s': %v", curlFilePath, err)
// 	}

// 	// Parse cURL command
// 	curlData, err := parseCurlCommand(string(curlContent))
// 	if err != nil {
// 		return fmt.Errorf("error parsing cURL command: %v", err)
// 	}

// 	// Read existing Go file
// 	goContent, err := ioutil.ReadFile(goFilePath)
// 	if err != nil {
// 		return fmt.Errorf("error reading Go file '%s': %v", goFilePath, err)
// 	}

// 	// Update Go code with new values
// 	updatedCode, err := updateGoCode(string(goContent), curlData)
// 	if err != nil {
// 		return fmt.Errorf("error updating Go code: %v", err)
// 	}

// 	// Write updated code back to file
// 	err = ioutil.WriteFile(goFilePath, []byte(updatedCode), 0644)
// 	if err != nil {
// 		return fmt.Errorf("error writing updated Go file '%s': %v", goFilePath, err)
// 	}

// 	fmt.Printf("✅ Successfully updated %s with fresh values from %s\n", goFilePath, curlFilePath)
// 	return nil
// }

// parseCurlCommand extracts relevant data from a cURL command
func parseCurlCommand(curlCmd string) (*CurlData, error) {
	data := &CurlData{
		Headers:  make(map[string]string),
		FormData: make(map[string]string),
	}

	// Remove line breaks and extra spaces
	curlCmd = strings.ReplaceAll(curlCmd, "\\\n", " ")
	curlCmd = regexp.MustCompile(`\s+`).ReplaceAllString(curlCmd, " ")

	// Extract cookies
	cookieRegex := regexp.MustCompile(`-b\s+'([^']+)'`)
	if matches := cookieRegex.FindStringSubmatch(curlCmd); len(matches) > 1 {
		data.Cookies = matches[1]

		// Extract user ID from cookies
		userRegex := regexp.MustCompile(`c_user=(\d+)`)
		if userMatches := userRegex.FindStringSubmatch(data.Cookies); len(userMatches) > 1 {
			data.UserID = userMatches[1]
		}
	}

	// Extract headers
	headerRegex := regexp.MustCompile(`-H\s+'([^:]+):\s*([^']+)'`)
	headerMatches := headerRegex.FindAllStringSubmatch(curlCmd, -1)
	for _, match := range headerMatches {
		if len(match) > 2 {
			data.Headers[match[1]] = match[2]
		}
	}

	// Extract LSD from headers
	if lsd, exists := data.Headers["x-fb-lsd"]; exists {
		data.LSD = lsd
	}

	// Extract form data
	dataRegex := regexp.MustCompile(`--data-raw\s+'([^']+)'`)
	if matches := dataRegex.FindStringSubmatch(curlCmd); len(matches) > 1 {
		formDataStr := matches[1]

		// Parse URL-encoded form data
		formValues, err := url.ParseQuery(formDataStr)
		if err == nil {
			for key, values := range formValues {
				if len(values) > 0 {
					data.FormData[key] = values[0]
				}
			}
		}

		// Extract variables JSON if present
		if variablesStr, exists := data.FormData["variables"]; exists {
			err := extractVariablesData(variablesStr, data)
			if err != nil {
				fmt.Printf("Warning: Could not parse variables JSON: %v\n", err)
			}
		}
	}

	return data, nil
}

// extractVariablesData parses the variables JSON to extract additional data
func extractVariablesData(variablesStr string, data *CurlData) error {
	// URL decode the variables string
	decodedVars, err := url.QueryUnescape(variablesStr)
	if err != nil {
		return err
	}

	// Parse JSON
	var variables map[string]interface{}
	err = json.Unmarshal([]byte(decodedVars), &variables)
	if err != nil {
		return err
	}

	// Extract input data
	if input, ok := variables["input"].(map[string]interface{}); ok {
		// Extract composer entry point
		if entryPoint, ok := input["composer_entry_point"].(string); ok {
			data.ComposerEntryPoint = entryPoint
		}

		// Extract session ID
		if logging, ok := input["logging"].(map[string]interface{}); ok {
			if sessionID, ok := logging["composer_session_id"].(string); ok {
				data.SessionID = sessionID
			}
		}

		// Extract attribution ID
		if navData, ok := input["navigation_data"].(map[string]interface{}); ok {
			if attrID, ok := navData["attribution_id_v2"].(string); ok {
				data.AttributionID = attrID
			}
		}

		// Extract group ID
		if audience, ok := input["audience"].(map[string]interface{}); ok {
			if groupID, ok := audience["to_id"].(string); ok {
				data.GroupID = groupID
			}
		}

		// Extract message text
		if message, ok := input["message"].(map[string]interface{}); ok {
			if text, ok := message["text"].(string); ok {
				data.MessageText = text
			}
		}

		// Extract photo IDs
		if attachments, ok := input["attachments"].([]interface{}); ok {
			for _, attachment := range attachments {
				if attachMap, ok := attachment.(map[string]interface{}); ok {
					if photo, ok := attachMap["photo"].(map[string]interface{}); ok {
						if photoID, ok := photo["id"].(string); ok {
							data.PhotoIDs = append(data.PhotoIDs, photoID)
						}
					}
				}
			}
		}
	}

	return nil
}

// updateGoCode replaces values in the Go code with fresh data from cURL
func updateCode(goCode string, data *CurlData) (string, error) {
	updatedCode := goCode

	// Build complete request body from form data with double-escaped format for fmt.Sprintf
	if len(data.FormData) > 0 {
		// First, let's replace the entire fmt.Sprintf requestBody pattern
		requestBodyPattern := regexp.MustCompile(`requestBody\s*:=\s*fmt\.Sprintf\("([^"]*)",\s*encodedVariables\)`)

		// Build new request body from extracted form data with double-escaping
		var formPairs []string
		paramOrder := []string{
			"av", "__aaid", "__user", "__a", "__req", "__hs", "dpr", "__ccg", "__rev", "__s", "__hsi",
			"__dyn", "__csr", "__hsdp", "__hblp", "__sjsp", "__comet_req", "fb_dtsg", "jazoest",
			"lsd", "__spin_r", "__spin_b", "__spin_t", "__crn", "fb_api_caller_class",
			"fb_api_req_friendly_name", "server_timestamps", "doc_id",
		}

		// Add parameters in order with proper double-escaping for fmt.Sprintf
		for _, param := range paramOrder {
			if value, exists := data.FormData[param]; exists && param != "variables" {
				// Double-escape URL-encoded characters for fmt.Sprintf
				escapedValue := doubleEscapeForFmtSprintf(value)
				formPairs = append(formPairs, fmt.Sprintf("%s=%s", param, escapedValue))
			}
		}

		// Add any missing parameters
		for param, value := range data.FormData {
			if param != "variables" {
				found := false
				for _, orderedParam := range paramOrder {
					if param == orderedParam {
						found = true
						break
					}
				}
				if !found {
					escapedValue := doubleEscapeForFmtSprintf(value)
					formPairs = append(formPairs, fmt.Sprintf("%s=%s", param, escapedValue))
				}
			}
		}

		// Build the new request body (without variables, we'll add that with %s)
		newRequestBodyTemplate := strings.Join(formPairs, "&") + "&variables=%s"

		// Replace the entire requestBody assignment
		newRequestBodyLine := fmt.Sprintf(`requestBody := fmt.Sprintf("%s", encodedVariables)`, newRequestBodyTemplate)
		updatedCode = requestBodyPattern.ReplaceAllString(updatedCode, newRequestBodyLine)
	}

	// Update user/actor ID in struct
	if data.UserID != "" {
		actorIDRegex := regexp.MustCompile(`ActorID:\s*"[^"]*"`)
		updatedCode = actorIDRegex.ReplaceAllString(updatedCode, fmt.Sprintf(`ActorID: "%s"`, data.UserID))
	}

	// Update composer entry point
	if data.ComposerEntryPoint != "" {
		entryPointRegex := regexp.MustCompile(`ComposerEntryPoint:\s*"[^"]*"`)
		updatedCode = entryPointRegex.ReplaceAllString(updatedCode, fmt.Sprintf(`ComposerEntryPoint: "%s"`, data.ComposerEntryPoint))
	}

	// Update session ID
	if data.SessionID != "" {
		sessionIDRegex := regexp.MustCompile(`ComposerSessionID:\s*"[^"]*"`)
		updatedCode = sessionIDRegex.ReplaceAllString(updatedCode, fmt.Sprintf(`ComposerSessionID: "%s"`, data.SessionID))
	}

	// Update attribution ID
	if data.AttributionID != "" {
		attrIDRegex := regexp.MustCompile(`AttributionIDV2:\s*"[^"]*"`)
		updatedCode = attrIDRegex.ReplaceAllString(updatedCode, fmt.Sprintf(`AttributionIDV2: "%s"`, data.AttributionID))
	}

	// Update LSD token in header
	if data.LSD != "" {
		lsdHeaderRegex := regexp.MustCompile(`req\.Header\.Set\("x-fb-lsd",\s*"[^"]*"\)`)
		updatedCode = lsdHeaderRegex.ReplaceAllString(updatedCode, fmt.Sprintf(`req.Header.Set("x-fb-lsd", "%s")`, data.LSD))
	}

	// Update cookies
	if data.Cookies != "" {
		cookieRegex := regexp.MustCompile(`req\.Header\.Set\("Cookie",\s*"[^"]*"\)`)
		updatedCode = cookieRegex.ReplaceAllString(updatedCode, fmt.Sprintf(`req.Header.Set("Cookie", "%s")`, data.Cookies))
	}

	return updatedCode, nil
}

// doubleEscapeForFmtSprintf converts URL-encoded strings to double-escaped format for fmt.Sprintf
func doubleEscapeForFmtSprintf(value string) string {
	// If the value is already URL-encoded, we need to double-escape the % characters
	// so that fmt.Sprintf doesn't interpret them as format verbs

	// First, if it's not URL-encoded, encode it
	if !strings.Contains(value, "%") {
		value = url.QueryEscape(value)
	}

	// Double-escape % characters for fmt.Sprintf
	// %3A becomes %%3A, %2C becomes %%2C, etc.
	doubleEscaped := strings.ReplaceAll(value, "%", "%%")

	return doubleEscaped
}

// PrintCurlData prints the extracted data for debugging
// func PrintCurlData(data *CurlData) {
// 	fmt.Println("=== Extracted cURL Data ===")
// 	fmt.Printf("User ID: %s\n", data.UserID)
// 	fmt.Printf("Session ID: %s\n", data.SessionID)
// 	fmt.Printf("Attribution ID: %s\n", data.AttributionID)
// 	fmt.Printf("LSD Token: %s\n", data.LSD)
// 	fmt.Printf("Composer Entry Point: %s\n", data.ComposerEntryPoint)
// 	fmt.Printf("Group ID: %s\n", data.GroupID)
// 	fmt.Printf("Message Text: %s\n", data.MessageText)
// 	fmt.Printf("Photo IDs: %v\n", data.PhotoIDs)
// 	fmt.Printf("Cookies: %s\n", data.Cookies)
// 	fmt.Println("Headers:")
// 	for key, value := range data.Headers {
// 		fmt.Printf("  %s: %s\n", key, value)
// 	}
// 	fmt.Println("Form Data:")
// 	for key, value := range data.FormData {
// 		if len(value) > 100 {
// 			fmt.Printf("  %s: %s...\n", key, value[:100])
// 		} else {
// 			fmt.Printf("  %s: %s\n", key, value)
// 		}
// 	}
// 	fmt.Println("========================")
// }

// Example usage function
// func ExampleUsage() {
// 	// Save your cURL command to a file called "curl.txt" (keep it short!)
// 	// Then call this function to update your Go code
// 	err := UpdateCodeFromCurl("curl.txt", "post.go")
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Code updated successfully!")
// }

// UpdateCodeFromCurlString takes cURL as string and updates Go code file
// func UpdateCodeFromCurlString(curlCommand string, goFilePath string) error {
// 	// Parse cURL command
// 	curlData, err := parseCurlCommand(curlCommand)
// 	if err != nil {
// 		return fmt.Errorf("error parsing cURL command: %v", err)
// 	}

// 	// Read existing Go file
// 	goContent, err := os.ReadFile(goFilePath)
// 	if err != nil {
// 		return fmt.Errorf("error reading Go file '%s': %v", goFilePath, err)
// 	}

// 	// Update Go code with new values
// 	updatedCode, err := updateCode(string(goContent), curlData)
// 	if err != nil {
// 		return fmt.Errorf("error updating Go code: %v", err)
// 	}

// 	// Write updated code back to file
// 	err = os.WriteFile(goFilePath, []byte(updatedCode), 0644)
// 	if err != nil {
// 		return fmt.Errorf("error writing updated Go file '%s': %v", goFilePath, err)
// 	}

// 	fmt.Printf("✅ Successfully updated %s with fresh cURL values\n", goFilePath)
// 	return nil
// }

// QuickUpdate is a helper function with very short file names
// func QuickUpdate() {
// 	err := UpdateCodeFromCurl("c.txt", "p.go")
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}
// 	fmt.Println("✅ Updated successfully!")
// }

// ParseAndPrintCurlData is a helper function to just parse and display data without updating files
// func ParseAndPrintCurlData(curlFilePath string) error {
// 	curlContent, err := ioutil.ReadFile(curlFilePath)
// 	if err != nil {
// 		return fmt.Errorf("error reading cURL file: %v", err)
// 	}

// 	curlData, err := parseCurlCommand(string(curlContent))
// 	if err != nil {
// 		return fmt.Errorf("error parsing cURL command: %v", err)
// 	}

// 	PrintCurlData(curlData)
// 	return nil
// }
