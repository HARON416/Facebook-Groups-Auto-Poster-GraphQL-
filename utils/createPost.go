package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
)

// PostConfig holds the parsed values from create post curl command
type PostConfig struct {
	URL      string
	Headers  map[string]string
	FormData map[string]string
}

// FacebookComposerStoryPayload represents the payload from ComposerStoryCreateMutation
type FacebookComposerStoryPayload struct {
	ComposerStoryCreate struct {
		Story struct {
			ID       string `json:"id"`
			LegacyID string `json:"legacy_id"`
			URL      string `json:"url"`
		} `json:"story"`
	} `json:"composer_story_create"`
}

// FacebookComposerStoryResponse represents Facebook's full ComposerStoryCreateMutation response
type FacebookComposerStoryResponse struct {
	Data       FacebookComposerStoryPayload `json:"data"`
	Extensions struct {
		IsFinal bool `json:"is_final"`
	} `json:"extensions"`
}

// PostResponse represents our processed post creation response
type PostResponse struct {
	Success   bool     `json:"success"`
	PostID    string   `json:"postId"`
	PostURL   string   `json:"postUrl"`
	Message   string   `json:"message"`
	Error     string   `json:"error,omitempty"`
	RequestID string   `json:"rid"`
	GroupID   string   `json:"groupId"`
	Text      string   `json:"text"`
	PhotoIDs  []string `json:"photoIds"`
}

// Global post config
var currentPostConfig *PostConfig

// parsePostCurlCommand extracts URL and headers from the create post curl command
func parsePostCurlCommand(createPostCurl string) (*PostConfig, error) {
	config := &PostConfig{
		Headers:  make(map[string]string),
		FormData: make(map[string]string),
	}

	// Extract URL using regex
	urlRegex := regexp.MustCompile(`curl '([^']+)'`)
	urlMatch := urlRegex.FindStringSubmatch(createPostCurl)
	if len(urlMatch) > 1 {
		config.URL = urlMatch[1]
	}

	// Extract cookies using regex (-b flag)
	cookieRegex := regexp.MustCompile(`-b '([^']+)'`)
	cookieMatch := cookieRegex.FindStringSubmatch(createPostCurl)
	if len(cookieMatch) > 1 {
		config.Headers["Cookie"] = cookieMatch[1]
	}

	// Extract headers using regex (-H flags)
	headerRegex := regexp.MustCompile(`-H '([^:]+):\s*([^']+)'`)
	headerMatches := headerRegex.FindAllStringSubmatch(createPostCurl, -1)
	for _, match := range headerMatches {
		if len(match) > 2 {
			headerName := strings.TrimSpace(match[1])
			headerValue := strings.TrimSpace(match[2])
			config.Headers[headerName] = headerValue
		}
	}

	// Extract form data from --data-raw (handles both ' and $' formats)
	dataRawRegex := regexp.MustCompile(`--data-raw \$?'([^']+)'`)
	dataRawMatch := dataRawRegex.FindStringSubmatch(createPostCurl)
	if len(dataRawMatch) > 1 {
		// Parse URL-encoded form data
		formData, err := url.ParseQuery(dataRawMatch[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing form data: %v", err)
		}

		// Convert to our format (taking first value for each key)
		for key, values := range formData {
			if len(values) > 0 {
				config.FormData[key] = values[0]
			}
		}
	}

	return config, nil
}

// UpdatePostConfigFromCurl parses the create post curl command and updates global config
func UpdatePostConfigFromCurl(createPostCurl string) error {
	fmt.Println("🔄 Updating post creation configuration from curl command...")

	config, err := parsePostCurlCommand(createPostCurl)
	if err != nil {
		return fmt.Errorf("error parsing post curl command: %v", err)
	}

	currentPostConfig = config

	fmt.Printf("✅ Post configuration updated successfully!\n")
	fmt.Printf("   - Post URL: %s\n", config.URL)
	fmt.Printf("   - Headers: %d items\n", len(config.Headers))
	fmt.Printf("   - Form Data: %d parameters\n", len(config.FormData))
	fmt.Println()

	return nil
}

// CreatePost creates a Facebook post with text and images
func CreatePost(text string, photoIDs []string, groupID, createPostCurl string) (*PostResponse, error) {
	fmt.Printf("📝 Creating Facebook post...\n")
	fmt.Printf("   📄 Text: %s\n", text)
	fmt.Printf("   🖼️ Photos: %d images\n", len(photoIDs))
	fmt.Printf("   👥 Group ID: %s\n", groupID)

	// Update post config if not already done
	if currentPostConfig == nil {
		err := UpdatePostConfigFromCurl(createPostCurl)
		if err != nil {
			return nil, fmt.Errorf("error updating post config: %v", err)
		}
	}

	// Build the post request
	response, err := makePostRequest(text, photoIDs, groupID)
	if err != nil {
		return nil, fmt.Errorf("error creating post: %v", err)
	}

	// if response.Success {
	// 	fmt.Printf("✅ Post created successfully!\n")
	// 	fmt.Printf("   📝 Post ID: %s\n", response.PostID)
	// 	fmt.Printf("   👥 Group ID: %s\n", response.GroupID)
	// 	fmt.Printf("   📄 Text: %s\n", response.Text)
	// 	fmt.Printf("   🖼️  Photos: %d images\n", len(response.PhotoIDs))
	// } else {
	// 	fmt.Printf("❌ Post creation failed: %s\n", response.Error)
	// }

	return response, nil
}

// extractUserIDFromConfig extracts the user ID from the current config
func extractUserIDFromConfig() string {
	if currentPostConfig == nil {
		return "61553861467726" // fallback
	}

	// Try to extract from form data first
	if userID, exists := currentPostConfig.FormData["__user"]; exists && userID != "" {
		return userID
	}

	// Try to extract from av parameter
	if av, exists := currentPostConfig.FormData["av"]; exists && av != "" {
		return av
	}

	fmt.Println("⚠️ Warning: Unable to extract user ID from config, using fallback value.")

	// Fallback to hardcoded value
	return "61553861467726"
}

// generateSessionID creates a UUID-like session ID
func generateSessionID() string {
	// Generate a simple UUID-like string for session ID
	// Using current timestamp for uniqueness
	now := time.Now().Unix()
	rand1 := rand.Int63n(0xFFFFFFFF)
	rand2 := rand.Int63n(0xFFFFFFFF)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		now&0xFFFFFFFF,
		rand1&0xFFFF,
		rand2&0xFFFF,
		(rand1>>16)&0xFFFF,
		rand2&0xFFFFFFFFFFFF)
}

// buildPostVariables creates the variables JSON for Facebook's ComposerStoryCreateMutation
func buildPostVariables(text string, photoIDs []string, groupID string) (string, error) {
	// Build attachments array from photo IDs
	attachments := make([]map[string]interface{}, len(photoIDs))
	for i, photoID := range photoIDs {
		attachments[i] = map[string]interface{}{
			"photo": map[string]interface{}{
				"id": photoID,
			},
		}
	}

	// Create the variables structure matching Facebook's expected format
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"composer_entry_point":    "inline_composer",
			"composer_source_surface": "group",
			"composer_type":           "group",
			"logging": map[string]interface{}{
				"composer_session_id": generateSessionID(),
			},
			"source": "WWW",
			"message": map[string]interface{}{
				"ranges": []interface{}{},
				"text":   text,
			},
			"with_tags_ids":         nil,
			"inline_activities":     []interface{}{},
			"text_format_preset_id": "0",
			"group_flair": map[string]interface{}{
				"flair_id": nil,
			},
			"attachments": attachments,
			"composed_text": map[string]interface{}{
				"block_data":    []string{"{}"},
				"block_depths":  []int{0},
				"block_types":   []int{0},
				"blocks":        []string{text},
				"entities":      []string{"[]"},
				"entity_map":    "{}",
				"inline_styles": []string{"[]"},
			},
			"navigation_data": map[string]interface{}{
				"attribution_id_v2": "CometGroupDiscussionRoot.react,comet.group,via_cold_start,1755966072871,86686,2361831622,,",
			},
			"tracking": []interface{}{nil},
			"event_share_metadata": map[string]interface{}{
				"surface": "newsfeed",
			},
			"audience": map[string]interface{}{
				"to_id": groupID,
			},
			"actor_id":           extractUserIDFromConfig(),
			"client_mutation_id": "1",
		},
		"feedLocation":                        "GROUP",
		"feedbackSource":                      0,
		"focusCommentID":                      nil,
		"gridMediaWidth":                      nil,
		"groupID":                             nil,
		"scale":                               1,
		"privacySelectorRenderLocation":       "COMET_STREAM",
		"checkPhotosToReelsUpsellEligibility": false,
		"renderLocation":                      "group",
		"useDefaultActor":                     false,
		"inviteShortLinkKey":                  nil,
		"isFeed":                              false,
		"isFundraiser":                        false,
		"isFunFactPost":                       false,
		"isGroup":                             true,
		"isEvent":                             false,
		"isTimeline":                          false,
		"isSocialLearning":                    false,
		"isPageNewsFeed":                      false,
		"isProfileReviews":                    false,
		"isWorkSharedDraft":                   false,
		"hashtag":                             nil,
		"canUserManageOffers":                 false,
		"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider":                     true,
		"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider":            true,
		"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider":                     true,
		"__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider":    false,
		"__relay_internal__pv__IsWorkUserrelayprovider":                                       false,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":                 false,
		"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":         false,
		"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider":    true,
		"__relay_internal__pv__FeedDeepDiveTopicPillThreadViewEnabledrelayprovider":           false,
		"__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider":        false,
		"__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider":                    false,
		"__relay_internal__pv__IsMergQAPollsrelayprovider":                                    false,
		"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider":     true,
		"__relay_internal__pv__StoriesArmadilloReplyEnabledrelayprovider":                     true,
		"__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider":        true,
		"__relay_internal__pv__GHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider": true,
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(variables)
	if err != nil {
		return "", fmt.Errorf("error marshaling variables to JSON: %v", err)
	}

	return string(jsonBytes), nil
}

// makePostRequest performs the actual post creation request
func makePostRequest(text string, photoIDs []string, groupID string) (*PostResponse, error) {
	if currentPostConfig == nil {
		return nil, fmt.Errorf("post config not initialized - call UpdatePostConfigFromCurl first")
	}

	// Build the variables JSON
	variablesJSON, err := buildPostVariables(text, photoIDs, groupID)
	if err != nil {
		return nil, fmt.Errorf("error building post variables: %v", err)
	}

	// Build form data for post creation
	formData := url.Values{}

	// Copy base form data from config
	for key, value := range currentPostConfig.FormData {
		formData.Set(key, value)
	}

	// Override the variables with our custom post data
	formData.Set("variables", variablesJSON)

	// Log the request details for debugging
	fmt.Printf("🔍 Request debugging:\n")
	fmt.Printf("   URL: %s\n", currentPostConfig.URL)
	fmt.Printf("   Method: POST\n")
	fmt.Printf("   Form data keys: %d\n", len(formData))
	for key := range formData {
		if key == "variables" {
			fmt.Printf("   - %s: [custom JSON - %d chars]\n", key, len(variablesJSON))
			//fmt.Println("Key 'variables' contains custom JSON data - omitted for brevity")
		} else {
			value := formData.Get(key)
			if len(value) > 50 {
				fmt.Printf("   - %s: %s... [%d chars]\n", key, value[:50], len(value))
				//fmt.Printf("Key '%s' contains long data - omitted for brevity\n", key)
			} else {
				fmt.Printf("   - %s: %s\n", key, value)
				//fmt.Println("Key '" + key + "' with short data - value omitted for brevity")
			}
		}
	}

	previewLen := 200
	if len(variablesJSON) < previewLen {
		previewLen = len(variablesJSON)
	}
	fmt.Printf("   Variables JSON preview: %s...\n", variablesJSON[:previewLen])

	// Create the HTTP request
	req, err := http.NewRequest("POST", currentPostConfig.URL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating post request: %v", err)
	}

	// Set headers from config
	fmt.Printf("   Request headers:\n")
	for headerName, headerValue := range currentPostConfig.Headers {
		req.Header.Set(headerName, headerValue)
		if len(headerValue) > 100 {
			fmt.Printf("   - %s: %s... [%d chars]\n", headerName, headerValue[:100], len(headerValue))
			//fmt.Printf("   - %s: [long value omitted for brevity]\n", headerName)
		} else {
			fmt.Printf("   - %s: %s\n", headerName, headerValue)
			//fmt.Printf("   - %s: [short value omitted for brevity]\n", headerName)
		}
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making post request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading post response: %v", err)
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

		//fmt.Println(string(body)) // readable HTML
	} else {
		fmt.Println("Uknown encoding type:", string(body))
	}
	// end decompression code ---

	fmt.Printf("📤 Post response status: %s\n", resp.Status)
	fmt.Printf("📤 Post response body length: %d bytes\n", len(body))
	fmt.Printf("📤 Response headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}

	log.Println("RATE LIMITING DETECTED", strings.Contains(string(body), "We limit how often you"))

	if strings.Contains(string(body), "We limit how often you") {
		fmt.Printf("⚠️  Facebook rate limiting detected in response body\n")
		os.Exit(0)
	}

	if len(body) > 0 {
		//fmt.Printf("📤 Post response body (raw): %s\n", string(body))

		// Also show as hex if it contains non-printable characters
		hasNonPrintable := false
		for _, b := range body {
			if b < 32 && b != 9 && b != 10 && b != 13 {
				hasNonPrintable = true
				break
			}
		}
		if hasNonPrintable {
			fmt.Printf("📤 Post response body (hex): %x\n", body)
		}
	} else {
		fmt.Printf("📤 Post response body is completely empty\n")
	}

	// Parse the response
	return parsePostResponse(body, text, photoIDs, groupID)
}

// parsePostResponse parses Facebook's response and extracts the post data
func parsePostResponse(body []byte, text string, photoIDs []string, groupID string) (*PostResponse, error) {
	bodyStr := string(body)

	// Handle empty response - this might indicate success in some cases
	if len(bodyStr) == 0 {
		fmt.Printf("⚠️  Empty response received - this might indicate success\n")
		return &PostResponse{
			Success:   true,
			PostID:    "unknown-empty-response",
			Message:   "Post likely created successfully (empty response)",
			Error:     "",
			RequestID: "",
			GroupID:   groupID,
			Text:      text,
			PhotoIDs:  photoIDs,
		}, nil
	}

	// Facebook responses may start with "for (;;);" - remove it
	bodyStr = strings.TrimPrefix(bodyStr, "for (;;);")

	// Facebook may return multiple JSON objects or complex responses
	// For now, let's check if the response contains errors
	hasErrors := strings.Contains(bodyStr, `"errors":[`) || strings.Contains(bodyStr, `"severity":"ERROR"`)

	fmt.Println("Response contains errors:", hasErrors)

	dec := json.NewDecoder(strings.NewReader(bodyStr))
	dec.UseNumber() // optional, preserves numbers accurately

	var objects []map[string]any

	for {
		var obj map[string]any
		if err := dec.Decode(&obj); err != nil {
			// EOF means we're done
			if err.Error() == "EOF" {
				break
			}
			panic(err)
		}
		objects = append(objects, obj)
	}

	// fmt.Println(objects[0])

	postIDs, urls := findPosts(objects[0])

	if len(postIDs) == 0 || len(urls) == 0 {
		//fmt.Println("Post not successful")
		response := &PostResponse{
			Success:  false,
			PostID:   "",
			PostURL:  "",
			Text:     text,
			PhotoIDs: photoIDs,
			GroupID:  groupID,
		}

		return response, nil
	} else {
		//fmt.Println("Post successful")
		var postID, postURL string
		for i := range postIDs {
			fmt.Printf("post_id:%s url:%s\n", postIDs[i], urls[i])
			postID = postIDs[i]
			postURL = urls[i]
		}

		response := &PostResponse{
			Success:  true,
			PostID:   postID,
			PostURL:  postURL,
			Text:     text,
			PhotoIDs: photoIDs,
			GroupID:  groupID,
		}

		return response, nil
	}

	// Usage
	// postID, permalinkURL := getPostData(objects[0])
	// if postID != "" && permalinkURL != "" {
	// 	fmt.Printf("Post ID: %s\n", postID)
	// 	fmt.Printf("Permalink URL: %s\n", permalinkURL)
	// 	response := &PostResponse{
	// 		Success:  true,
	// 		PostID:   postID,
	// 		PostURL:  permalinkURL,
	// 		Text:     text,
	// 		PhotoIDs: photoIDs,
	// 		GroupID:  groupID,
	// 	}

	// 	return response, nil
	// } else {
	// 	fmt.Printf("⚠️  Post ID or Permalink URL not found in response object\n")

	// 	if hasErrors {
	// 		return &PostResponse{
	// 			Success:  false,
	// 			Error:    "Facebook returned server errors",
	// 			Message:  "Post creation failed due to server errors",
	// 			Text:     text,
	// 			PhotoIDs: photoIDs,
	// 			GroupID:  groupID,
	// 		}, nil
	// 	}

	// 	response := &PostResponse{
	// 		Success:  false,
	// 		PostID:   postID,
	// 		Text:     text,
	// 		PhotoIDs: photoIDs,
	// 		GroupID:  groupID,
	// 	}

	// 	return response, nil
	// }

	// // Try to parse as JSON - if it fails, we'll still return success based on status
	// var fbResponse FacebookComposerStoryResponse
	// err := json.Unmarshal([]byte(bodyStr), &fbResponse)

	// // If parsing fails but we got a 200 response, consider it a success for now
	// if err != nil {
	// 	fmt.Printf("⚠️  JSON parsing failed: %v\n", err)
	// 	fmt.Printf("⚠️  Raw response: %s\n", bodyStr)

	// 	// Check if there are obvious errors in the response
	// 	if hasErrors {
	// 		return &PostResponse{
	// 			Success:  false,
	// 			Error:    "Facebook returned server errors",
	// 			Message:  "Post creation failed due to server errors",
	// 			Text:     text,
	// 			PhotoIDs: photoIDs,
	// 			GroupID:  groupID,
	// 		}, nil
	// 	}

	// 	// Parsing failed but no obvious errors - assume success
	// 	return &PostResponse{
	// 		Success:  true,
	// 		PostID:   "unknown", // We couldn't parse the ID
	// 		Message:  "Post likely created successfully (parsing issue)",
	// 		Text:     text,
	// 		PhotoIDs: photoIDs,
	// 		GroupID:  groupID,
	// 	}, nil
	// }

	// // Check if the response indicates success
	// success := fbResponse.Data.ComposerStoryCreate.Story.ID != ""
	// postID = fbResponse.Data.ComposerStoryCreate.Story.ID

	// fmt.Printf("✅ Successfully parsed response - Post ID: %s\n", postID)

	// // Create our structured response
	// response := &PostResponse{
	// 	Success:  success,
	// 	PostID:   postID,
	// 	Text:     text,
	// 	PhotoIDs: photoIDs,
	// 	GroupID:  groupID,
	// }

	// if success {
	// 	response.Message = "Post created successfully"
	// } else {
	// 	response.Error = "Post creation failed - no post ID returned"
	// 	response.Message = "Post creation failed"
	// }

	// return response, nil
}

// Recursive function to collect post_id and url strings
func findPosts(data interface{}) (postIDs []string, urls []string) {
	switch v := data.(type) {
	case map[string]interface{}:
		// Check if current object has post_id and url
		if postID, okID := v["post_id"].(string); okID {
			if url, okURL := v["url"].(string); okURL && strings.Contains(url, "permalink") {
				postIDs = append(postIDs, postID)
				urls = append(urls, url)
			}
		}
		// Recurse into all values
		for _, value := range v {
			subIDs, subURLs := findPosts(value)
			postIDs = append(postIDs, subIDs...)
			urls = append(urls, subURLs...)
		}
	case []interface{}:
		for _, item := range v {
			subIDs, subURLs := findPosts(item)
			postIDs = append(postIDs, subIDs...)
			urls = append(urls, subURLs...)
		}
	}
	return
}
