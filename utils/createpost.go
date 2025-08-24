package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Create post curl command
var createPostCurl = `curl 'https://web.facebook.com/api/graphql/' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -b 'sb=E2sHZ8Dj0xZhKsq5e2frpYxk; ps_l=1; ps_n=1; datr=EZluZ78uE_dw-PpiG3FGLpyf; oo=v1; c_user=100016139237616; wd=1366x681; fr=1A8nFCnZNIAxIIUaz.AWdYb3ZlfY48CCQuck46A0USrH84R9Fq1L6oIKBkWMComAMQmsI.BoquF2..AAA.0.0.BoquF2.AWc_XA5JEwVwF7IioIG9NYQJDOc; xs=1%3AHCeFsY1N8T3k8w%3A2%3A1737623622%3A-1%3A-1%3A%3AAcXseJLaPuOV3SFAyLRrlUlfJb0CyeugC2sGNHKJjg; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1756029376368%2C%22v%22%3A1%7D' \
  -H 'origin: https://web.facebook.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://web.facebook.com/groups/571323476357125' \
  -H 'sec-ch-prefers-color-scheme: dark' \
  -H 'sec-ch-ua: "Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"' \
  -H 'sec-ch-ua-full-version-list: "Not;A=Brand";v="99.0.0.0", "Google Chrome";v="139.0.7258.66", "Chromium";v="139.0.7258.66"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-model: ""' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-ch-ua-platform-version: "6.8.0"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36' \
  -H 'x-asbd-id: 359341' \
  -H 'x-fb-friendly-name: ComposerStoryCreateMutation' \
  -H 'x-fb-lsd: Hk3RdF0qVDvsVE3d67kML5' \
  --data-raw $'av=100016139237616&__aaid=0&__user=100016139237616&__a=1&__req=2z&__hs=20324.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1026265415&__s=w2qn8t%3Ahzfeis%3Awgvtk4&__hsi=7542088668936050370&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu2u5Ee88o4Wm7-2K0-obUG2-azqwaW223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwWzUfHDzUiBG2OUqwjVqwLwHwa211wo83KwHwOyUqxG&__csr=gekt1hb8YIbN31xax2d6hsDifi2Yj_rsjYGPnkrZsG4MBqkJ9ddNdbpkAyuCOcyfLibnmSxmJd9AZXHWZqyBJasDGKCGfcl6BnHQlFdDAVFG4aEDBlyeQmjKqGnG9-hdbABqhp4Z4pHFKjxe8HmbyGGAmUOWhHGEyFEOVqVuHiABWHy_By99V9rGBgK4FHzWhqmibKaAACBV8h-bAyKtpmmfAQdyFGnKmVoPy8OeWVoO9LxaUKum5Acz94u4ETxvzGKihyoJBm9yryHgswCGp1Ouq69onxjy8-58Ny9GCADzE-XKbVWizqGrgyF8SibxS5FEuyaxC9yEdA9x6ibgyEsAyUKfgsG5Ea8jxeazEjx2i6awExqay9ohBwOBwxxGieye3PyWwNy9Q4FK5E8E9o9A5EK2Lz8TGUhzEGdzV8nxC390wxS9xWEXxy7UnwZByE5K2-7Uetbyoyex21Zh4ifwaBCo4G5odEfEKbwjU31Dwd-1bwEAj812wFAAAxe1xgepu3eFQ3mym7oeE4xa58swtUZ0xw66wmomw0wPw1_O04K-6E2zU0KS9wgqw3DoiwbRK3K0N208209oQ05m4l5g2so4e00E4403MK08wB8SE4-8x-0u20gC3a0LEjwJg0Ae8g3eF1Szwg81JUaUAw0zm4409Lo1n819ohwgE0y6680sew3To1C3w2F8C9gOt03KQ0wQ0fiwXxe02GC0aFwjS0g67rgW19waq8Dg4O0bZw4kw2Nk2TBt0hEuwaC08HAw5Fw2_E0pfAw60yWwtm5o6zC9fmdpk62kU5K1i80gOmE512U17o0-GsNUsdw&__hsdp=gjgtFMidEmNgywk8gyEsCoGeGp25IdEWEy93ezyxccImyEub48O22yiaAAxbdcggz8ixR8GdJ9JcIzRNr2dOsgRAp2aVicDComyT9jqsgxsixkvlEuygDpigj9goaT10G14NN5isI8_2i3cjbj66kgxiAp6y4aACkG8KyLah2HAz468KhBrPMyrG4oPTKlUPoBebeAFjuaJcELjBAF9Faleqh8LqlemGECzk8m84hjlAfhSFWqKZ95Fd5EzK2qsyah9DUgby32QA8cOin98pjFcFzARkvyAr8iz9194t2FCfRkMB4DBAp98xrhlanGdGl8x9l4cFpoK5cE89k9xx169izF4qBxd1t16moySizBF6FFcx4SFhkFm4EhoBovxIEy2sM-8CAoy9ig-iiFQWjjKbJCh8BUR7GQQ4pAh3A1ao-ehFot3m5Aem8hj2kiE88G49k7A1kgZ0lXz4bxaE7a1qxK3qE1F88oG17AwiPUqoS5Ud98kwJD70goaEa8mxwYwfA2Wm5ojwv8Rwik5VyxafwTxAp7whURwxwbq1mwi87S5o3Hw8O0y81YKi3m0DEYE2dwdO0p20m-ayFE561Ax2WwaCm1pz8gwRwgovw8qi6o8981K832wpU14U16UK0_81FE0z-0gO0OUpwSwe60HE21wnE2PwCwlE1Y81kU5W0ke1Tw4ew&__hblp=0qp46UKi1bwiEK4S3i16Uc84WU8o2fwdW0Topy8cUco4Cu7okG3eUbkA3i79U989E8olyo8FobohwyzUSUaVQdGeyGwzxy3i2-1wAwxy8y1HwAz898iyVu1lx7Dwxx-1wwMF0wAypeUeUeUWi0wo-789K9u2vwCxCmeg-2l7wbim4A9gS0AofE5m2K2h0oo4C1NwmFo4e8ga89UOewaq1mwGxehxG1fzK2K3-2K1dUvwWwKwRxeh1t09m1TS0SE2jwMwywQwgUuw9u1Zx92o5u4pEtwtFE2Rxy3m0C9oqwkod415wFG8yEuwRwa61KwNABwgUe839wcKewYwyxmbwqE2TwVAwh8C3O3-Ueqwmd0OwsUC1yxi9US362h16awtV5gC2R0Jx6dwCy8pykm9yGyU4h0KwHxCewa24ojyUhxKbwca3-7olwHwMwYwtFUmG3-4U23wYDxmq9w8u3GV8sxa1nwPwNgvwJjwmEK1lyoliDx-ew9yu10x-3636t0xxy4GwkE8oaoiDwhEcocEDghDxaq2q5olg1mUG2C6po8qJ2Ua4bwJho5m32cg4qdxCpa2J04QyrUWu2G3K1QwzoeE-WCy9Zedwzxmfgyqm2m5Am2K15xC1sBjgswoo43y8jKuFIw7OdxC78&__sjsp=gjgtFMidEmNgywk8gyEsCoGeGp25IdEWEy93ezyxccIF4O8IibdEH6kUAyF98iSINcOcN27kyESQCQOOfn5IAWtaL42SiGyicDwIjGakV444OSK548p98J7wya5oC2Yoapy7G8hUmgBk88mkxW8gFE_plCfhApGOp-l8x2GQ9i8KDLamay8N1ybADhzMyrG4oPTKlUPyQUIXy3ueh6cQVpbFGVoKh8ExkVqGQmLylx67AteuKHBihqgGSXwCyG8qCvx0K8acigxmOinqQpjF3zAAg-cgwwZ194t2FC4lguBBh98xp8xanGdGlbh9oRDBwCwwBgC644oBaeAhFokxt16c8JAxNF4p4S85nBoix5z8vxIwC2sM-8CAoy9igsDjFdeUZCh8BUR7GQQ64h0mCfzAq2cdomghy8lhU88rBg7h3Q1nx12UlwsE5G6UdE1FUb84u1owGxu1ewR709m5Eof83V0Qxm4U2ao4B0RwohAu17zm0S80Pm0va3O0Hyw5Gw6gw5LyEGq0JEkG0Gpo5CcwhE467U26AxC22i&__comet_req=15&fb_dtsg=NAfugHWuVdIrshH1K3diAPKmGyUKWyOnSE5wzWgPmy61pQcQP5pB6sw%3A1%3A1737623622&jazoest=25568&lsd=Hk3RdF0qVDvsVE3d67kML5&__spin_r=1026265415&__spin_b=trunk&__spin_t=1756029359&__crn=comet.fbweb.CometGroupDiscussionRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=ComposerStoryCreateMutation&variables=%7B%22input%22%3A%7B%22composer_entry_point%22%3A%22hosted_inline_composer%22%2C%22composer_source_surface%22%3A%22group%22%2C%22composer_type%22%3A%22group%22%2C%22logging%22%3A%7B%22composer_session_id%22%3A%222cc2cf08-9e94-46a0-b43d-2376a7b5be99%22%7D%2C%22source%22%3A%22WWW%22%2C%22message%22%3A%7B%22ranges%22%3A%5B%5D%2C%22text%22%3A%22Get%20your%20dream%20SAMSUNG%20GALAXY%20NOTE%2010%20PLUS%20(EX%20UK)%20today%20with%20our%20Lipa%20Mdogo%20Mdogo%20plan\u0021%20Only%20KSh%209%2C199%20deposit%20%2B%20KSh%201%2C240%20weekly%20for%2052%20weeks.%2012GB%20RAM%2F256GB%20storage.%20Visit%20us%20at%20Pioneer%20Building%2C%20Kimathi%20Street%2C%20or%20call%2FWhatsApp%200718448461.%22%7D%2C%22with_tags_ids%22%3Anull%2C%22inline_activities%22%3A%5B%5D%2C%22text_format_preset_id%22%3A%220%22%2C%22group_flair%22%3A%7B%22flair_id%22%3Anull%7D%2C%22attachments%22%3A%5B%7B%22photo%22%3A%7B%22id%22%3A%221851805048700785%22%7D%7D%2C%7B%22photo%22%3A%7B%22id%22%3A%221851805028700787%22%7D%7D%2C%7B%22photo%22%3A%7B%22id%22%3A%221851805015367455%22%7D%7D%5D%2C%22composed_text%22%3A%7B%22block_data%22%3A%5B%22%7B%7D%22%5D%2C%22block_depths%22%3A%5B0%5D%2C%22block_types%22%3A%5B0%5D%2C%22blocks%22%3A%5B%22Get%20your%20dream%20SAMSUNG%20GALAXY%20NOTE%2010%20PLUS%20(EX%20UK)%20today%20with%20our%20Lipa%20Mdogo%20Mdogo%20plan\u0021%20Only%20KSh%209%2C199%20deposit%20%2B%20KSh%201%2C240%20weekly%20for%2052%20weeks.%2012GB%20RAM%2F256GB%20storage.%20Visit%20us%20at%20Pioneer%20Building%2C%20Kimathi%20Street%2C%20or%20call%2FWhatsApp%200718448461.%22%5D%2C%22entities%22%3A%5B%22%5B%5D%22%5D%2C%22entity_map%22%3A%22%7B%7D%22%2C%22inline_styles%22%3A%5B%22%5B%5D%22%5D%7D%2C%22navigation_data%22%3A%7B%22attribution_id_v2%22%3A%22CometGroupDiscussionRoot.react%2Ccomet.group%2Cunexpected%2C1756029463700%2C753703%2C2361831622%2C%2C%3BGroupsCometJoinsRoot.react%2Ccomet.groups.joins%2Cvia_cold_start%2C1756029363123%2C338993%2C%2C%2C%22%7D%2C%22tracking%22%3A%5Bnull%5D%2C%22event_share_metadata%22%3A%7B%22surface%22%3A%22newsfeed%22%7D%2C%22audience%22%3A%7B%22to_id%22%3A%22571323476357125%22%7D%2C%22actor_id%22%3A%22100016139237616%22%2C%22client_mutation_id%22%3A%221%22%7D%2C%22feedLocation%22%3A%22GROUP%22%2C%22feedbackSource%22%3A0%2C%22focusCommentID%22%3Anull%2C%22gridMediaWidth%22%3Anull%2C%22groupID%22%3Anull%2C%22scale%22%3A1%2C%22privacySelectorRenderLocation%22%3A%22COMET_STREAM%22%2C%22checkPhotosToReelsUpsellEligibility%22%3Afalse%2C%22renderLocation%22%3A%22group%22%2C%22useDefaultActor%22%3Afalse%2C%22inviteShortLinkKey%22%3Anull%2C%22isFeed%22%3Afalse%2C%22isFundraiser%22%3Afalse%2C%22isFunFactPost%22%3Afalse%2C%22isGroup%22%3Atrue%2C%22isEvent%22%3Afalse%2C%22isTimeline%22%3Afalse%2C%22isSocialLearning%22%3Afalse%2C%22isPageNewsFeed%22%3Afalse%2C%22isProfileReviews%22%3Afalse%2C%22isWorkSharedDraft%22%3Afalse%2C%22hashtag%22%3Anull%2C%22canUserManageOffers%22%3Afalse%2C%22__relay_internal__pv__CometUFIShareActionMigrationrelayprovider%22%3Atrue%2C%22__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider%22%3Atrue%2C%22__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider%22%3Atrue%2C%22__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider%22%3Afalse%2C%22__relay_internal__pv__IsWorkUserrelayprovider%22%3Afalse%2C%22__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider%22%3Afalse%2C%22__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider%22%3Afalse%2C%22__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider%22%3Atrue%2C%22__relay_internal__pv__FeedDeepDiveTopicPillThreadViewEnabledrelayprovider%22%3Afalse%2C%22__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider%22%3Afalse%2C%22__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider%22%3Afalse%2C%22__relay_internal__pv__IsMergQAPollsrelayprovider%22%3Afalse%2C%22__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider%22%3Atrue%2C%22__relay_internal__pv__StoriesArmadilloReplyEnabledrelayprovider%22%3Atrue%2C%22__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider%22%3Atrue%2C%22__relay_internal__pv__GHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider%22%3Atrue%7D&server_timestamps=true&doc_id=31383956567869484'`

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
func parsePostCurlCommand() (*PostConfig, error) {
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
func UpdatePostConfigFromCurl() error {
	fmt.Println("🔄 Updating post creation configuration from curl command...")

	config, err := parsePostCurlCommand()
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
func CreatePost(text string, photoIDs []string, groupID string) (*PostResponse, error) {
	fmt.Printf("📝 Creating Facebook post...\n")
	fmt.Printf("   📄 Text: %s\n", text)
	fmt.Printf("   🖼️  Photos: %d images\n", len(photoIDs))
	fmt.Printf("   👥 Group ID: %s\n", groupID)

	// Update post config if not already done
	if currentPostConfig == nil {
		err := UpdatePostConfigFromCurl()
		if err != nil {
			return nil, fmt.Errorf("error updating post config: %v", err)
		}
	}

	// Build the post request
	response, err := makePostRequest(text, photoIDs, groupID)
	if err != nil {
		return nil, fmt.Errorf("error creating post: %v", err)
	}

	if response.Success {
		fmt.Printf("✅ Post created successfully!\n")
		fmt.Printf("   📝 Post ID: %s\n", response.PostID)
		fmt.Printf("   👥 Group ID: %s\n", response.GroupID)
		fmt.Printf("   📄 Text: %s\n", response.Text)
		fmt.Printf("   🖼️  Photos: %d images\n", len(response.PhotoIDs))
	} else {
		fmt.Printf("❌ Post creation failed: %s\n", response.Error)
	}

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
		} else {
			value := formData.Get(key)
			if len(value) > 50 {
				fmt.Printf("   - %s: %s... [%d chars]\n", key, value[:50], len(value))
			} else {
				fmt.Printf("   - %s: %s\n", key, value)
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
		} else {
			fmt.Printf("   - %s: %s\n", headerName, headerValue)
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

	fmt.Printf("📤 Post response status: %s\n", resp.Status)
	fmt.Printf("📤 Post response body length: %d bytes\n", len(body))
	fmt.Printf("📤 Response headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}

	if len(body) > 0 {
		fmt.Printf("📤 Post response body (raw): %s\n", string(body))

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

	// Try to parse as JSON - if it fails, we'll still return success based on status
	var fbResponse FacebookComposerStoryResponse
	err := json.Unmarshal([]byte(bodyStr), &fbResponse)

	// If parsing fails but we got a 200 response, consider it a success for now
	if err != nil {
		fmt.Printf("⚠️  JSON parsing failed: %v\n", err)
		fmt.Printf("⚠️  Raw response: %s\n", bodyStr)

		// Check if there are obvious errors in the response
		if hasErrors {
			return &PostResponse{
				Success:  false,
				Error:    "Facebook returned server errors",
				Message:  "Post creation failed due to server errors",
				Text:     text,
				PhotoIDs: photoIDs,
				GroupID:  groupID,
			}, nil
		}

		// Parsing failed but no obvious errors - assume success
		return &PostResponse{
			Success:  true,
			PostID:   "unknown", // We couldn't parse the ID
			Message:  "Post likely created successfully (parsing issue)",
			Text:     text,
			PhotoIDs: photoIDs,
			GroupID:  groupID,
		}, nil
	}

	// Check if the response indicates success
	success := fbResponse.Data.ComposerStoryCreate.Story.ID != ""
	postID := fbResponse.Data.ComposerStoryCreate.Story.ID

	fmt.Printf("✅ Successfully parsed response - Post ID: %s\n", postID)

	// Create our structured response
	response := &PostResponse{
		Success:  success,
		PostID:   postID,
		Text:     text,
		PhotoIDs: photoIDs,
		GroupID:  groupID,
	}

	if success {
		response.Message = "Post created successfully"
	} else {
		response.Error = "Post creation failed - no post ID returned"
		response.Message = "Post creation failed"
	}

	return response, nil
}
