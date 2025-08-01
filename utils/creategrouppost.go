package utils

import (
	"encoding/json"
	"fmt"
	"io"
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
			ComposerEntryPoint:    "inline_composer",
			ComposerSourceSurface: "group",
			ComposerType:          "group",
			Logging: struct {
				ComposerSessionID string `json:"composer_session_id"`
			}{
				ComposerSessionID: "14f915c8-de15-4d07-bf0c-a1971c52b41a",
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
			NavigationData: struct {
				AttributionIDV2 string `json:"attribution_id_v2"`
			}{
				AttributionIDV2: "CometGroupDiscussionRoot.react,comet.group,via_cold_start,1753963901694,240833,2361831622,,",
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
			ActorID:          "61555590462485",
			ClientMutationID: "1",
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

	requestBody := fmt.Sprintf("av=61555590462485&__aaid=0&__user=61555590462485&__a=1&__req=2h&__hs=20300.HYP%%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1025357550&__s=9diw46%%3Afbxq0i%%3A3qq4wm&__hsi=7533217585352648528&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu2u5Ee88o4Wm7-2K0-obUG2-azqwaW223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwWwjHBU-4FqwIK6E4-mEbUaU2wwgo-1gweW2K3aEy6Eqw&__csr=gaI7QagXkj2Jh5hQ4Yh3_NWf5JW6q6PHbbYriPN2tPT5Fiq9iZH8CDQCT4FZp7NivQCN7cWIGl4mDvAV9eAO9kQKECDh94KkCmJXBZ5A-CnZ2Bi4aLGV5h-rjheF4Ah68HyVWy9uLBgx4ASmQKWnOayXAGi4ebheWgRdq-jzejGQFXXCGqbDAgnF2FEK8iyXKmheKqfVEW8UyeURu9GuUmACyoOS8y8mKBDCx6QQlpGFG2au9xiEiKVV8uy-iRxq6oCezEFrKVoKF8kxO9G68y6ooiAzooGUK9J2-74Eom6pppQ498G4V42-qicQ8Gm9Uy4Fd12uj-qqdwzhU9EW58gmfxqq8xq6oybwRz8mxmm48hzqyk8AyUiBxe9wPzVU-8xC2Gh4xJ2UpXxi2acxOi3Oq14BwxxS8xu5odoG6UK4Uym8AwExa4EvCCyopyUhx-EdUnwq82fyOa0HA0g21jg4S1ywbu3C0zo28oBw8SeByo4N3WyAQU6Z5wn8zwa61qy85B05Sl9dedBU5O16RgB0g8BK2J0XwWAo5-3101h4yw5yw0CJw1QV6m7o0kkw2fk3Oez87-0ky09Jw11q04w8dU569U02dHg0cxA0j-E0p2Pxy76UqyVU2-xO0re06U8fpA0t-toao0xK0A43q0rS08uo0T-2-0Vo5K1Lw8CkK4k1iw1du04ok0qi0i1eXh4bOU0bF84m5onU3EwcG328w_81DwgQ1Pg0LG9gKt08a3C2O13w2lE8-7U3jw2EA0e1wmU3_w7xwso30waK1aweW0PEgweC6U2vy80zKiA0tS1Rya8upWw7PghZi4ngco3fwae3e9a46dwzU&__hsdp=g4fq4FEF28iGA4E42ax1123sygwa8ilEO1sWfE4yEO7bmshji5Gy9N5EslmJ8xaf5FAjPH69_T9p4AOh29ux4z8B88aAtTe9FGC52sGMi8j9n5ihcuzcdhvewogx4A22NcmHEO2z922ax4qCNAjh7Mwl8j4hjREy4leKBgMIxi94yUOHh6msxKyJNYFA8i4qz-9By9ohKGyykAy9hOSqUyy6kFzyi9HVaH9XJQSKwzFEg8ybqqiimeazIKSqQSBBoJhz8Qb8vhk9qOAbl25yYWCgxy-IBNN2CujykgyoCmhAiajh8AxEI9yoC9yV4gEVzQhGhEpcpy4Wqya-OUOl4x2fhQ4EsJoC58kgoxeWXcn8bKUgyQagN1LiooBwwpA957B8tzkA4ocbyy2pozQ4O6Ji8yB4h9OHhyBGuu8VFFp4ihaiulF5xcxFFU4PmodzE9ongcteh4c9gZd0Qxi3GHAz8qgoyQ3ydoOfweOargkBgixS1oxy1ww5Xx6Er2wCUjAgsgQwx3A4UR3hMbie749wQzE88fU6124i259wJU8C9wiq3gy1syEdU8pEqxe2219x2dwu8owA8U98f87a0BUlxa0H8vwTwrUiAxyi0NGwJ86oV38vxy1OGyxiazE4m0iW0mC08rwemew2JUC1Hwjo4C23w8Cq0pu0BU6C0W834wwwl81H82Pwe20z832w2Gotwm833wh85K0m2bw2yE3Vwr83Vw5Lw9i18w&__hblp=08CqQbwoE9UkoowgUaUqwEwXyEbE138gwvo3Axi2CdwXCwXUcKdmUeEK0xEyicyoa8aUtzUbEkzoPyouyoSry8hwUG1xwXwxwYwm8vUnx2czp84e7E5O3SbxnGegpG4EcbwiocFVUdob4Q2eewmbWQm0x-9xiUhwiEfoycw44ig5K5E98cUnwjo5qu2J0Ew8O1Gwnoc8owoWK5omwgo9Ud86-362uq9wEFx91C321KwlpE62axG0RF85m0w87a0EEC7E7K2e2O0FoogeomxW1NQbx-26ezopxW1RwlWwww8210wwwiU9E521owzww81gxm4otwVwmUhw821kwlU5C1exG26ewso3Ox-mfw8e1OwLwMwWwiGxKdAxGaxnxW8wopGwSyo2SmUtzU6vwkU9ooG1sz8oy827z8-4oKUiwoE4V1-1ky824wdG7UK5U7u48aElwZwTwl8kzp8bEPBwwCxmbG1KxGcw920Foa8521NwsUiwko98e9EgyU4O5By47oG4oybyFomxW2S4ouwJCgggaE2eG8wHDyEco45e0JTU4a267U8Hw_wEiK48bpUaU7WbK17xd0bC36djwl8V7xe4U-axS1jwrUuxueg9E4yu12AzEiQ789U4ela486qE2iwp61xxJBw&__sjsp=g4fq4FEF28iGA4E42ax1123sygwa8ilEO1sWcGV0WGZ8sJ2h5daAmG8Ax5EIwBmGJDiq595jEONyvWitAiAH48ymzBIyJy2G9lOxe6kYFja293jqp5yoXLJ1ebCjP3Up62ojmEKEC26hyalcGcQOxegR7rgqFVqWaWW8zdsxql12a8A48pg-8W8hNYFECg-E_ypoy5XGEExkh9b5fybya8vaoUCGqKiECahKi5aUC8y94CAvzyybbJCJdxWpIzgNDQmaqEx2RgxoQUB1rH8Ujl6uibh2byoJ17h8Dgj46oaF4h151rGhJJ1JeCEKV-cBh8gzQt0LJouxh1y4XGn8n8bwBwvkt2hhVi7gV90gEUw-8J1cxHEy8FjEitV4sxqzUy5p4ihaieKq4O0puodwQxt0UAxmfh8d85qi2t0nERz8-0hem361Nw7sx6E8hMS3p3i242edgUq2Qw841dwd61Wo5ewQ0OE1ko0Iq0E83ewJ80g6awl8&__comet_req=15&fb_dtsg=NAftO4MpRGstU5SjER8WWyyBlN8MerGcI11T233Ilt27F74_XYNgv7A%%3A23%%3A1737397894&jazoest=25217&lsd=8JXVbwiNDdewcoTUr19-1Z&__spin_r=1025357550&__spin_b=trunk&__spin_t=1753963898&__crn=comet.fbweb.CometGroupDiscussionRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=ComposerStoryCreateMutation&variables=%s&server_timestamps=true&doc_id=31469424502648938", encodedVariables)

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://web.facebook.com/api/graphql/", strings.NewReader(requestBody))
	if err != nil {
		return post, fmt.Errorf("Error creating request: %v", err)
	}

	// Set headers from the curl request - exact values
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-fb-lsd", "8JXVbwiNDdewcoTUr19-1Z")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("origin", "https://web.facebook.com")
	req.Header.Set("sec-ch-prefers-color-scheme", "light")
	req.Header.Set("sec-ch-ua-full-version-list", `"Not)A;Brand";v="8.0.0.0", "Chromium";v="138.0.7204.168", "Google Chrome";v="138.0.7204.168"`)
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-ch-ua-platform-version", `"6.8.0"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("sec-ch-ua-model", `""`)
	req.Header.Set("x-asbd-id", "359341")
	req.Header.Set("x-fb-friendly-name", "ComposerStoryCreateMutation")
	req.Header.Set("accept", "*/*")
	req.Header.Set("referer", "https://web.facebook.com/groups/894063942126928")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	// Set cookies from the curl request
	req.Header.Set("Cookie", "ar_debug=1; ps_n=1; datr=bpaOZ6y9AuJ3LuuspFRSlBZF; c_user=61555590462485; wd=1366x681; sb=TE0JZ7NYzOmh0Gpb3d5BFjVf; ps_l=1; fr=15wX5RIbJZj770eqv.AWe6cw0dSholVKUkqtSgGhdbH4iC4evl2HvnXiwhz7jt370aZ24.Boi1M-..AAA.0.0.Boi1M-.AWfyUDb8Pl2vd1OBsfpj1cC27mo; xs=23%3AvABQDiJTzL7XLA%3A2%3A1737397894%3A-1%3A-1%3A%3AAcXTgo2kLkb6NfDrwvJegyyOvpL_Op5sgQXHXg4KfoXy; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1753963908991%2C%22v%22%3A1%7D")

	// Create HTTP client
	client := &http.Client{}

	// Make the request
	fmt.Println("Making Facebook GraphQL request...")
	resp, err := client.Do(req)
	if err != nil {
		return post, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response headers:", resp.Header)

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
	} else {
		fmt.Printf("Response body (raw): %s\n", string(bodyBytes))
		fmt.Printf("Response body length: %d bytes\n", len(bodyBytes))
	}

	return post, nil
}
