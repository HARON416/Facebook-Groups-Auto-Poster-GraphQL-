package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// FacebookPost represents the data needed for a Facebook post
type FacebookPost struct {
	MessageText string   `json:"message_text"`
	PhotoIDs    []string `json:"photo_ids"`
	GroupID     string   `json:"group_id"`
}

type InputVariables struct {
	Input struct {
		ComposerEntryPoint    string `json:"composer_entry_point"`
		ComposerSourceSurface string `json:"composer_source_surface"`
		ComposerType          string `json:"composer_type"`
		Logging               struct {
			ComposerSessionID string `json:"composer_session_id"`
		} `json:"logging"`
		Source  string `json:"source"`
		Message struct {
			Ranges []interface{} `json:"ranges"`
			Text   string        `json:"text"`
		} `json:"message"`
		WithTagsIDs        interface{}   `json:"with_tags_ids"`
		InlineActivities   []interface{} `json:"inline_activities"`
		TextFormatPresetID string        `json:"text_format_preset_id"`
		Attachments        []interface{} `json:"attachments"` // Use interface{} for now
		ComposedText       struct {
			BlockData    []string `json:"block_data"`
			BlockDepths  []int    `json:"block_depths"`
			BlockTypes   []int    `json:"block_types"`
			Blocks       []string `json:"blocks"`
			Entities     []string `json:"entities"`
			EntityMap    string   `json:"entity_map"`
			InlineStyles []string `json:"inline_styles"`
		} `json:"composed_text"`
		NavigationData struct {
			AttributionIDV2 string `json:"attribution_id_v2"`
		} `json:"navigation_data"`
		Tracking           []interface{} `json:"tracking"`
		EventShareMetadata struct {
			Surface string `json:"surface"`
		} `json:"event_share_metadata"`
		Audience struct {
			ToID string `json:"to_id"`
		} `json:"audience"`
		ActorID          string `json:"actor_id"`
		ClientMutationID string `json:"client_mutation_id"`
	} `json:"input"`
	FeedLocation                                                                 string      `json:"feedLocation"`
	FeedbackSource                                                               int         `json:"feedbackSource"`
	FocusCommentID                                                               interface{} `json:"focusCommentID"`
	GridMediaWidth                                                               interface{} `json:"gridMediaWidth"`
	GroupID                                                                      interface{} `json:"groupID"`
	Scale                                                                        int         `json:"scale"`
	PrivacySelectorRenderLocation                                                string      `json:"privacySelectorRenderLocation"`
	CheckPhotosToReelsUpsellEligibility                                          bool        `json:"checkPhotosToReelsUpsellEligibility"`
	RenderLocation                                                               string      `json:"renderLocation"`
	UseDefaultActor                                                              bool        `json:"useDefaultActor"`
	InviteShortLinkKey                                                           interface{} `json:"inviteShortLinkKey"`
	IsFeed                                                                       bool        `json:"isFeed"`
	IsFundraiser                                                                 bool        `json:"isFundraiser"`
	IsFunFactPost                                                                bool        `json:"isFunFactPost"`
	IsGroup                                                                      bool        `json:"isGroup"`
	IsEvent                                                                      bool        `json:"isEvent"`
	IsTimeline                                                                   bool        `json:"isTimeline"`
	IsSocialLearning                                                             bool        `json:"isSocialLearning"`
	IsPageNewsFeed                                                               bool        `json:"isPageNewsFeed"`
	IsProfileReviews                                                             bool        `json:"isProfileReviews"`
	IsWorkSharedDraft                                                            bool        `json:"isWorkSharedDraft"`
	Hashtag                                                                      interface{} `json:"hashtag"`
	CanUserManageOffers                                                          bool        `json:"canUserManageOffers"`
	RelayInternalPVCometUFIShareActionMigrationrelayprovider                     bool        `json:"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider"`
	RelayInternalPVGHLShouldChangeSponsoredDataFieldNamerelayprovider            bool        `json:"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider"`
	RelayInternalPVGHLShouldChangeAdIdFieldNamerelayprovider                     bool        `json:"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider"`
	RelayInternalPVCometUFI_dedicated_comment_routable_dialog_gkrelayprovider    bool        `json:"__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider"`
	RelayInternalPVIsWorkUserrelayprovider                                       bool        `json:"__relay_internal__pv__IsWorkUserrelayprovider"`
	RelayInternalPVCometUFIReactionsEnableShortNamerelayprovider                 bool        `json:"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider"`
	RelayInternalPVFBReels_deprecate_short_form_video_context_gkrelayprovider    bool        `json:"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider"`
	RelayInternalPVFeedDeepDiveTopicPillThreadViewEnabledrelayprovider           bool        `json:"__relay_internal__pv__FeedDeepDiveTopicPillThreadViewEnabledrelayprovider"`
	RelayInternalPVFBReels_enable_view_dubbed_audio_type_gkrelayprovider         bool        `json:"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider"`
	RelayInternalPVCometImmersivePhotoCanUserDisable3DMotionrelayprovider        bool        `json:"__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider"`
	RelayInternalPVWorkCometIsEmployeeGKProviderrelayprovider                    bool        `json:"__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider"`
	RelayInternalPVIsMergQAPollsrelayprovider                                    bool        `json:"__relay_internal__pv__IsMergQAPollsrelayprovider"`
	RelayInternalPVFBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider     bool        `json:"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider"`
	RelayInternalPVStoriesArmadilloReplyEnabledrelayprovider                     bool        `json:"__relay_internal__pv__StoriesArmadilloReplyEnabledrelayprovider"`
	RelayInternalPVFBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider        bool        `json:"__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider"`
	RelayInternalPVGHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider bool        `json:"__relay_internal__pv__GHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider"`
}

func CreateGroupPost(post FacebookPost) (FacebookPost, error) {
	// Create a post with random test data (text only)
	// post := FacebookPost{
	// 	MessageText: "✅ LIPA MDOGO MDOGO\n\n✅ SAMSUNG GALAXY S22 [EX UK]\n\n✅ KSH. 11,599 DEPOSIT + KSH.1,560/WEEK [52 WEEKS]\n\n✅ TUKO PIONEER BUILDING KIMATHI STREET\n\n✅ CALL/WHATSAPP 0718448461",
	// 	PhotoIDs:    []string{},
	// 	GroupID:     "",
	// }

	//Populate the InputVariables struct
	variables := InputVariables{
		Input: struct {
			ComposerEntryPoint    string `json:"composer_entry_point"`
			ComposerSourceSurface string `json:"composer_source_surface"`
			ComposerType          string `json:"composer_type"`
			Logging               struct {
				ComposerSessionID string `json:"composer_session_id"`
			} `json:"logging"`
			Source  string `json:"source"`
			Message struct {
				Ranges []interface{} `json:"ranges"`
				Text   string        `json:"text"`
			} `json:"message"`
			WithTagsIDs        interface{}   `json:"with_tags_ids"`
			InlineActivities   []interface{} `json:"inline_activities"`
			TextFormatPresetID string        `json:"text_format_preset_id"`
			Attachments        []interface{} `json:"attachments"` // Use interface{} for now
			ComposedText       struct {
				BlockData    []string `json:"block_data"`
				BlockDepths  []int    `json:"block_depths"`
				BlockTypes   []int    `json:"block_types"`
				Blocks       []string `json:"blocks"`
				Entities     []string `json:"entities"`
				EntityMap    string   `json:"entity_map"`
				InlineStyles []string `json:"inline_styles"`
			} `json:"composed_text"`
			NavigationData struct {
				AttributionIDV2 string `json:"attribution_id_v2"`
			} `json:"navigation_data"`
			Tracking           []interface{} `json:"tracking"`
			EventShareMetadata struct {
				Surface string `json:"surface"`
			} `json:"event_share_metadata"`
			Audience struct {
				ToID string `json:"to_id"`
			} `json:"audience"`
			ActorID          string `json:"actor_id"`
			ClientMutationID string `json:"client_mutation_id"`
		}{
			ComposerEntryPoint:    "hosted_inline_composer",
			ComposerSourceSurface: "group",
			ComposerType:          "group",
			Logging: struct {
				ComposerSessionID string `json:"composer_session_id"`
			}{
				ComposerSessionID: "4665e32c-1aa1-43bb-8309-b6f2b1be223a",
			},
			Source: "WWW",
			Message: struct {
				Ranges []interface{} `json:"ranges"`
				Text   string        `json:"text"`
			}{
				Ranges: []interface{}{},
				Text:   post.MessageText,
			},
			WithTagsIDs:        nil,
			InlineActivities:   []interface{}{},
			TextFormatPresetID: "0",
			Attachments:        []interface{}{}, // Handle attachments later
			ComposedText: struct {
				BlockData    []string `json:"block_data"`
				BlockDepths  []int    `json:"block_depths"`
				BlockTypes   []int    `json:"block_types"`
				Blocks       []string `json:"blocks"`
				Entities     []string `json:"entities"`
				EntityMap    string   `json:"entity_map"`
				InlineStyles []string `json:"inline_styles"`
			}{
				BlockData:    []string{"{}"},
				BlockDepths:  []int{0},
				BlockTypes:   []int{0},
				Blocks:       []string{post.MessageText},
				Entities:     []string{"[]"},
				EntityMap:    "{}",
				InlineStyles: []string{"[]"},
			},
			NavigationData: struct {
				AttributionIDV2 string `json:"attribution_id_v2"`
			}{
				AttributionIDV2: "CometGroupDiscussionRoot.react,comet.group,unexpected,1753865409319,736839,2361831622,,;GroupsCometJoinsRoot.react,comet.groups.joins,unexpected,1753865393835,296351,,,,;GroupsCometCrossGroupFeedRoot.react,comet.groups.feed,tap_bookmark,1753865391472,591435,2361831622,,",
			},
			Tracking: []interface{}{nil},
			EventShareMetadata: struct {
				Surface string `json:"surface"`
			}{
				Surface: "newsfeed",
			},
			Audience: struct {
				ToID string `json:"to_id"`
			}{
				ToID: post.GroupID,
			},
			ActorID:          "61560452168137",
			ClientMutationID: "3",
		},
		FeedLocation:                        "GROUP",
		FeedbackSource:                      0,
		FocusCommentID:                      nil,
		GridMediaWidth:                      nil,
		GroupID:                             nil,
		Scale:                               1,
		PrivacySelectorRenderLocation:       "COMET_STREAM",
		CheckPhotosToReelsUpsellEligibility: false,
		RenderLocation:                      "group",
		UseDefaultActor:                     false,
		InviteShortLinkKey:                  nil,
		IsFeed:                              false,
		IsFundraiser:                        false,
		IsFunFactPost:                       false,
		IsGroup:                             true,
		IsEvent:                             false,
		IsTimeline:                          false,
		IsSocialLearning:                    false,
		IsPageNewsFeed:                      false,
		IsProfileReviews:                    false,
		IsWorkSharedDraft:                   false,
		Hashtag:                             nil,
		CanUserManageOffers:                 false,
		RelayInternalPVCometUFIShareActionMigrationrelayprovider:                     true,
		RelayInternalPVGHLShouldChangeSponsoredDataFieldNamerelayprovider:            true,
		RelayInternalPVGHLShouldChangeAdIdFieldNamerelayprovider:                     true,
		RelayInternalPVCometUFI_dedicated_comment_routable_dialog_gkrelayprovider:    false,
		RelayInternalPVIsWorkUserrelayprovider:                                       false,
		RelayInternalPVCometUFIReactionsEnableShortNamerelayprovider:                 false,
		RelayInternalPVFBReels_deprecate_short_form_video_context_gkrelayprovider:    true,
		RelayInternalPVFeedDeepDiveTopicPillThreadViewEnabledrelayprovider:           false,
		RelayInternalPVFBReels_enable_view_dubbed_audio_type_gkrelayprovider:         false,
		RelayInternalPVCometImmersivePhotoCanUserDisable3DMotionrelayprovider:        false,
		RelayInternalPVWorkCometIsEmployeeGKProviderrelayprovider:                    false,
		RelayInternalPVIsMergQAPollsrelayprovider:                                    false,
		RelayInternalPVFBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider:     true,
		RelayInternalPVStoriesArmadilloReplyEnabledrelayprovider:                     true,
		RelayInternalPVFBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider:        true,
		RelayInternalPVGHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider: false,
	}
	//variables.Input.Attachments = []interface{}{}

	attachmentsJSON := buildAttachmentsJSON(post.PhotoIDs)
	var attachments []interface{}
	if err := json.Unmarshal([]byte(attachmentsJSON), &attachments); err != nil {
		return post, fmt.Errorf("Error parsing attachments JSON: %v", err)
	}
	variables.Input.Attachments = attachments

	// Debug: Print the attachments being sent
	fmt.Printf("Attachments being sent: %+v\n", attachments)
	fmt.Printf("Photo IDs: %v\n", post.PhotoIDs)

	// Convert the variables struct to JSON
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return post, fmt.Errorf("Error marshaling JSON: %v", err)
	}
	// Use proper URL encoding
	encodedVariables := url.QueryEscape(string(variablesJSON))
	// Create the request body with dynamic data
	requestBody := fmt.Sprintf("av=61560452168137&__aaid=0&__user=61560452168137&__a=1&__req=6g&__hs=20299.HYP%%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1025302125&__s=9n5i35%%3Awrk89z%%3As2qa8l&__hsi=7532786109695679640&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0-E4a3a4oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5oe8464-5pU9UmwUwxwjFovUaU3VwLyEbUGdG0HE88cA0z8c84q58jyUaUbGxe6Uak0zU8oC1hxB0qo4e4UO2m3G1eKnzUiBG2OUqwjVqwLwHwa211zU520XEaUcGy8qxG&__csr=g4N3c6cfND9q69NQri4NYp8DfkdNi1148rESyROidTh4ycDFsAhbkTpdnayl9mDvnKgynltAbnOr54Ltp5W8FsDQGWRQBGhrmDSOnbm8JAaCiF5h9vvpOQqqHyFkQCJk8auKlRGAuhCRmRzqgB-qJPibHGQEHy4F4QC8SGhaF6FkEzoxuVvXKQjvKiaVt6ADy25WQiF9HAQl4ZDKFrSfmiaoyhp9bhtd2oSqiiq8qXGbhXA-mhoKFXWWyuaXmnDgW49FHiF3FJel1KqAu9KiV8gjV9VExbrjRx24Ah2FoHB-Vl_xaAHVFuUZemhXzlACyrCiXJ5zUOiUSrxiEO4AETGjzGxiaUWcy8lz98j-UCcGFULCGjCBBzUoyUlK7WBjAwEAKufDDyrzF9qzAFHBy9UKnx3zrxW3q4bAwLAUjyUSEO9AyQ6k3y6oOcDgiAxeu58F7wHGHxa5UK64m2G5oc8pAwDAwgUGA2BrglDyHyEG1tFe2y3m8wxwvK1lwCzEowIzXw9jC-K15zUuwyw8WJabLjxe3zxK6A7pkcwUgdFo8UG4U4y7p88Elwc62G356gkwK-2LwBwvk6eu6o4Cum5U7W7oSfwNzd0GG0A84a0A85q2K6EnGu5p8O2S0bww9Aw1o80pgwgE0lNwgE10E5S1zyoF3o0uixNe06p81U82kxK5U0S20E9Q1hw0rEA0jG02aC03RS0g-0EUe12qxy3mfwfR0Rw2AQ0C62O1xwGwfWzOd0lrxS1wxy0dLo1gEzga82xwaG0hW480gNg4mE0Arc04eE4q0m60qsMvg0GKoB0rU1580Yu5EO0_Q04jEboaU6alw5hw1LOC0dlwxwn21O3S4EO4E0QK5VS5-1nw3Hy5yo7-1vwlQ08Zw8210w2YU3Ww2u8C9w13O48bHa1Agyzwp-S0ZEy9may85Z0Xg3vwhUlwtjilJ0srw2pE0Z2&__hsdp=g494448gykAzA1CAx9xIuxa1J1ix2P8oy41MFCkBc-AxOP8W9B44cgiNVaIy2a6diGOaxH1kOT8yQwai8Ib8dE44riD8AQwpOONiAyTQB7skNyj2aPEvRN9q8GAcMygy68YWEJcqOKPaTqiJDAjyqeiy8X4EG6KWc8gw85aKzBDF4yYrCowyjFFcCWzi5Ragl422BAzoN4N4CjzSn8yJml27Y9sNrAhiitazDKasD6miDAxhbQy5dbjp2A99re9O58TIAclHeAS8AFGbhrz6Cz5AvrFauO4iihoN3FSFF2W6sxtinDtGc8p4G9ze64ExkoMTh4dr88Kq8ghKsHy8nyFAiu-hyioFUiDJoXFy8zqy9yep7KiQi3a7z3lhamhe6YWhxsIR125kyVqy8EyLECgJ24S44EOhDwkEVklBgR2S2eVKuiah8Gu9mezp49cwOdG2-bFeQpxO4Aeanzk3d0GwkagtgdVEggW3ZVUiwp8hKaxp0sUy3Fhw7C20MaoF1G2O5F4YmUS5rJ6KCQ8yEK2i7UtwMBzQ5kim3t1m6AeG68dbpkhzA8C5cs1NyFUy1cx-eybgc8G8zo63Fxu8yWx2dJ1R0kE5oxE57Dw5wU2jwFg72446U8o5u5E4wg1FwpUuyo1rUbo1gU2Xwa-26E1ao4K0C85K12gC0mW6E1Bo5e2-3u2e26bwTwvU5O0BEowd-5ovxS2-78hw5swIG0xE2Zwci2y0gS0qS6o4W0A84y1cwaK12wo8C0oi4oGEjwcS3W0WU1eo724o3ow2ao17U&__sjsp=g494448gykAzA1CAx9xIuxa1J1ix2P8oy47NV2l6ymkChiqi7b8hSQFlagNcX6kiH8wwR98Qy68lEBDfNBfaoJHLyialrgiox2pQyrWTxC8BjGii7pda9gJ1W8LQ6p8Oh164Ey6aYGhoznoXxijCCm5lmmiQZEPasNDjKSuieEVO8X4EHqFrKz24821iHEVoGmbEo8oBdeCAOrGeptiA4sg8aaxF4yQFe6mh168povKh5n9QGeuUJqZ2SiDxBbQy5diQSkx2imgMDpydX8yQ44p7joyiCehrz4it5BmAqf8h99193EyaAh8-inkBVQCEnAG9ze64EyAsEThoVcwyVEx16VOwAyFAibV6ezEizEjoy8-EyozDyUx4wkj117mii6YWhxsGa4UlibBG56Frujt2944Eko6CFCqdgN0zKq9AyV8G4lzEtcsa3ocEKAXgy2Z3JDzk3e2K1gF1W3W1Iw8S0DUeB60uo2zxqhf5wLgGCUkwk8twSwCBwZg4e1hosgyokNM760F40IVk04gU1IEuyo1Do1gU2Xwd6E1ao1g804J2aG4U0qUw&__comet_req=15&fb_dtsg=NAftSKSWhPgbN1VZDx3QB3kYbkfockXiJv5syWOVfkruSHXYikFe6cw%%3A8%%3A1735880128&jazoest=25682&lsd=7yjkTI33tHoY5_gr2ktfw_&__spin_r=1025302125&__spin_b=trunk&__spin_t=1753877982&__crn=comet.fbweb.CometGroupDiscussionRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=ComposerStoryCreateMutation&variables=%s&server_timestamps=true&doc_id=24343131148645455", encodedVariables)

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://web.facebook.com/api/graphql/", strings.NewReader(requestBody))
	if err != nil {
		return post, fmt.Errorf("Error creating request: %v", err)
	}

	// Set headers from the curl request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://web.facebook.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://web.facebook.com/groups/"+post.GroupID)
	req.Header.Set("Sec-Ch-Prefers-Color-Scheme", "light")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua-Full-Version-List", `"Not)A;Brand";v="8.0.0.0", "Chromium";v="138.0.7204.168", "Google Chrome";v="138.0.7204.168"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Model", `""`)
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Platform-Version", `"6.8.0"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("X-Asbd-Id", "359341")
	req.Header.Set("X-Fb-Friendly-Name", "ComposerStoryCreateMutation")
	req.Header.Set("X-Fb-Lsd", "7yjkTI33tHoY5_gr2ktfw_")

	// Set cookies from the curl request
	req.Header.Set("Cookie", "sb=XWoHZ-GdXVwlvXZrSaFe7gwz; ps_l=1; ps_n=1; datr=78xqZ8Fc-vdpp5Ii5nGp2A0P; c_user=61560452168137; fr=1a15IgfoL9YD8UJdb.AWeGji5rDy2rm4SEmo8-nOt83EMVrk9hwCYHYfy806YGW-UXyaA.BoidJJ..AAA.0.0.BoidJJ.AWfz2onld5nimQfpPC3ByT3JxNk; xs=8%3AeMpQ9UiVlPPogg%3A2%3A1735880128%3A-1%3A-1%3AqSC9VrelgmIhZA%3AAcX82BDvQ8ITJ-PMf5_tPa7D8oEZNhi45oevhYbng0dV; wd=1366x681; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1753877990211%2C%22v%22%3A1%7D")

	// Create HTTP client
	client := &http.Client{}

	// Make the request
	fmt.Println("Making Facebook GraphQL request...")
	resp, err := client.Do(req)
	if err != nil {
		return post, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	return post, nil
}
