package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var fetchGroupsCurl = `curl 'https://web.facebook.com/api/graphql/' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -b 'sb=E2sHZ8Dj0xZhKsq5e2frpYxk; ps_l=1; ps_n=1; datr=EZluZ78uE_dw-PpiG3FGLpyf; oo=v1; c_user=100016139237616; wd=1366x681; fr=1A8nFCnZNIAxIIUaz.AWdYb3ZlfY48CCQuck46A0USrH84R9Fq1L6oIKBkWMComAMQmsI.BoquF2..AAA.0.0.BoquF2.AWc_XA5JEwVwF7IioIG9NYQJDOc; xs=1%3AHCeFsY1N8T3k8w%3A2%3A1737623622%3A-1%3A-1%3A%3AAcXseJLaPuOV3SFAyLRrlUlfJb0CyeugC2sGNHKJjg; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1756029376368%2C%22v%22%3A1%7D' \
  -H 'origin: https://web.facebook.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://web.facebook.com/groups/joins/?nav_source=tab&ordering=viewer_added' \
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
  -H 'x-fb-friendly-name: GroupsCometAllJoinedGroupsSectionPaginationQuery' \
  -H 'x-fb-lsd: Hk3RdF0qVDvsVE3d67kML5' \
  --data-raw 'av=100016139237616&__aaid=0&__user=100016139237616&__a=1&__req=p&__hs=20324.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1026265415&__s=vwlvja%3Ahzfeis%3Awgvtk4&__hsi=7542088668936050370&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwnUW3q2ibwNw9G2Sawba1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0n24oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5o7G4-5pUfEe88o4Wm7-2K0-obUG2-azqwaW1jg2cwMwrUK2K2WEjxK2B08-269wqQ1FwgUjz89oeE-3WVU-4FqwIK6E4-mEbUaU2wwgo620XEaUcEK6Eqw&__csr=gugJ2sl6ctcB7llrRHsIAGdJn2k8FRidEDlltlbRBaQChsTkDnZJOipAIRfNheKAAyIhbSBuqVeTuqqQumhCile-cRLHly4XhuFER4x5a8yaXhJJ28SF8Z4y8niQEPyEyqmjKqGg999ufBzEyU8WwLxiq4EdbyKGUK5UiyKcxa2i-48jyVoo-bAxhpEmxm598S3XxadxifDxS7Erzoboyh1WE8F8swFwAzEdEmy84Si2GU-9K48Pg8U4u7E9o4J0mUKbwJG3-2uawPwwwmohBwuEWawjoeEy4Ey5E8E4O5UK1aGU88bonwj84C3a683yw8DiUC0FAh4zU3_xm1RyU0B62R09q1xgf-3eFQ3mym1ow46g8o1U8mw0kXe07CU14E0W7g0lohkl0ea00Vk84-089w49wf64Ubk093y40Qo0AO9804JS0lO0im09Sxy038kU0Gi9ykcDg0XJ019e4U0aGo0GC1fo10otxe19w0JmAw16y06jV80Pm1i80gOmE1uU0-GtQ&__hsdp=gUzcUs-q8hey3y9yV6lIk8EyEixx4iHGAN2WpAGeGp25IBJEmCWEyj2h7Fzxa33F4iaxKbdECh5eEHGAWJ7QA8JJQFV7Jd5gR2raFqUKOmhrnnbjhVt1ry65ESp5HUCx4gPnxiA8aT2tB92Yy49goaTeb6qA6h2fIshkCTfPbBfmh5d5TFiGNxBakiVpo-iClAgN3EKV8N1ybBhXGaxbG23uUCcS9jzdKVjuaJdkdjBAF9KFHxfSaBGG9FQEDy14kRp3QtGuQaQA8jhqmXwKpV42ma86VAfFyQFOy8zlh-4ECza6o-q4ukMB4DBACAy4mpdBAzqBihFl45ppoK5cE5S9wVgc8oemAqKX8hdGklalxa4m9m7Ura1eCzUuAAGczbyVpAi9u5Gy44pAifg6J6BwUwCgKkM-2eawSg2-KcgK4GwsE5G6Ueo1F88oG17AwiPUqoS5Ud98kwJD70goaEa8mxwYwfA2Wm5ojwv8Rwik5VyxafwTxAp7whURwxwbq1mwi87S5o3Hw8O09XV8do2uzOw8S09Jw5LyEGq1hwp8gKE2FBwmoO48do467U26AxC22i0ry0ME6u08LyU3Yw6Cw2fU1383bxC3q0Uo2Kw861uwbe2q1mw7Mw5jwnE1gU7u0gW&__hblp=09u0CUK0Au321eK260zU1TUcUco4Cu0Ro4C1uyo8EcU4i3S5EhyEfU6i1Bwxy8y0A8O0I-u11wo85m1lwiE2Pwo89EbA3l7wbG2O0Q85m08ky41gxK0FE5q2G4VE3Dwdq1Zgnw9q0lq0AU3Hw2g83ewRw9W0Eo7e4U5i0h23q19w7ezE7m1Swbu3C19w8q08hwbGhk9wJwKw9W15w5ayU3Ywd20QGw53wey4E5u1Bg4S0OUpwSwe60HEoxa1lwxwYDwaK2q1mw6Dwi411wbK1kwnE1gU16E4Gqfyk4o5e3a5VoaU5-0xo2JwaK14w&__sjsp=gUzcUs-q8hey3y9yV6lJP2a8GPEoh4GWFcgKQpazBpd2Drjhi9KG8AMAhWoULEceAmREIUZuGCCtecGXGQdgyh4F2F7Jd5gSaCiBHyX9gyWHOQcBG5K2-qbUsB9opgxyoJ7wz85oC2-5k9hUaQcgC789pEwLhEOO0wDgR0Ey8N1ybxmaxbG23uUCcUJecSddx14x1emiWpEux2mGJ5Aze4ouhQVV8N92oG-Xw-wByC2d1d2QEky94fwFxCfCwSByV8x1emidy4ESfzo9E5S9wWwo-V6hdy1lVm4EhoO7Ur81fCzUaUOcK4pAi9u5Gy464i1YhEfE7a0g644bxm1OwmErw7uwIwhU5y2G5U4W3ks0BomxwYwfA3i5o2uo4B0RwohAu17zm0S80kEwYwaUE0LK0m-ayFE2SxiE2FBwmoO16wgovw8qi6o898&__comet_req=15&fb_dtsg=NAfugHWuVdIrshH1K3diAPKmGyUKWyOnSE5wzWgPmy61pQcQP5pB6sw%3A1%3A1737623622&jazoest=25568&lsd=Hk3RdF0qVDvsVE3d67kML5&__spin_r=1026265415&__spin_b=trunk&__spin_t=1756029359&__crn=comet.fbweb.CometGroupsJoinsRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=GroupsCometAllJoinedGroupsSectionPaginationQuery&variables=%7B%22count%22%3A20%2C%22cursor%22%3A%22AQHSVHsPWzhQJnfw8f09SO4qNBWq4tzB255S1BQTTlp-zPZCO036Uav4uObrhOOGIXaXSvgaJRsfieAkkI74liLqBg%22%2C%22ordering%22%3A%5B%22viewer_added%22%5D%2C%22scale%22%3A1%7D&server_timestamps=true&doc_id=9974006939348139'`

// Facebook GraphQL response structures
type FacebookResponse struct {
	Data struct {
		Viewer struct {
			AllJoinedGroups struct {
				TabGroupsList struct {
					Edges []struct {
						Node struct {
							ID                  string `json:"id"`
							Name                string `json:"name"`
							URL                 string `json:"url"`
							ViewerJoinState     string `json:"viewer_join_state"`
							ViewerLastVisitTime int64  `json:"viewer_last_visited_time"`
						} `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						HasNextPage bool   `json:"has_next_page"`
						EndCursor   string `json:"end_cursor"`
					} `json:"page_info"`
				} `json:"tab_groups_list"`
			} `json:"all_joined_groups"`
		} `json:"viewer"`
	} `json:"data"`
	Extensions struct {
		ServerMetadata struct {
			RequestID string `json:"request_id"`
		} `json:"server_metadata"`
		IsFinal bool `json:"is_final"`
	} `json:"extensions"`
}

// Group represents a Facebook group
type Group struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	URL                 string `json:"url"`
	ViewerJoinState     string `json:"viewer_join_state"`
	ViewerLastVisitTime int64  `json:"viewer_last_visited_time"`
}

// CurlConfig holds the parsed values from curl command
type CurlConfig struct {
	Cookies   string
	Headers   map[string]string
	FormData  map[string]string
	UserAgent string
	XFBLsd    string
	XASBDID   string
}

// Global config to store parsed curl values
var currentConfig *CurlConfig

// parseCurlCommand extracts all values from the fetchGroupsCurl variable
func parseCurlCommand() (*CurlConfig, error) {
	config := &CurlConfig{
		Headers:  make(map[string]string),
		FormData: make(map[string]string),
	}

	// Extract cookies using regex
	cookieRegex := regexp.MustCompile(`-b '([^']+)'`)
	cookieMatch := cookieRegex.FindStringSubmatch(fetchGroupsCurl)
	if len(cookieMatch) > 1 {
		config.Cookies = cookieMatch[1]
	}

	// Extract headers using regex
	headerRegex := regexp.MustCompile(`-H '([^:]+):\s*([^']+)'`)
	headerMatches := headerRegex.FindAllStringSubmatch(fetchGroupsCurl, -1)
	for _, match := range headerMatches {
		if len(match) > 2 {
			headerName := strings.TrimSpace(match[1])
			headerValue := strings.TrimSpace(match[2])
			config.Headers[headerName] = headerValue

			// Store specific headers we need
			switch headerName {
			case "user-agent":
				config.UserAgent = headerValue
			case "x-fb-lsd":
				config.XFBLsd = headerValue
			case "x-asbd-id":
				config.XASBDID = headerValue
			}
		}
	}

	// Extract form data using regex
	dataRegex := regexp.MustCompile(`--data-raw '([^']+)'`)
	dataMatch := dataRegex.FindStringSubmatch(fetchGroupsCurl)
	if len(dataMatch) > 1 {
		formDataString := dataMatch[1]

		// Parse URL-encoded form data
		formValues, err := url.ParseQuery(formDataString)
		if err != nil {
			return nil, fmt.Errorf("error parsing form data: %v", err)
		}

		// Convert to map[string]string
		for key, values := range formValues {
			if len(values) > 0 {
				config.FormData[key] = values[0]
			}
		}
	}

	return config, nil
}

// UpdateConfigFromCurl parses the curl command and updates global config
func UpdateConfigFromCurl() error {
	fmt.Println("🔄 Updating configuration from curl command...")

	config, err := parseCurlCommand()
	if err != nil {
		return fmt.Errorf("error parsing curl command: %v", err)
	}

	currentConfig = config

	fmt.Printf("✅ Configuration updated successfully!\n")
	fmt.Printf("   - User ID: %s\n", config.FormData["__user"])
	fmt.Printf("   - Cookies: %d characters\n", len(config.Cookies))
	fmt.Printf("   - Headers: %d items\n", len(config.Headers))
	fmt.Printf("   - Form Data: %d parameters\n", len(config.FormData))
	fmt.Println()

	return nil
}

func makeRequest(apiURL, encodedData string) (*FacebookResponse, error) {
	// Create the request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(encodedData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Use current config if available, otherwise fall back to defaults
	if currentConfig != nil {
		// Set all headers from parsed curl command
		for headerName, headerValue := range currentConfig.Headers {
			req.Header.Set(headerName, headerValue)
		}

		// Set cookies from parsed curl command
		if currentConfig.Cookies != "" {
			req.Header.Set("Cookie", currentConfig.Cookies)
		}
	} else {
		// Fallback to hardcoded headers if no config is available
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Origin", "https://web.facebook.com")
		req.Header.Set("Priority", "u=1, i")
		req.Header.Set("Referer", "https://web.facebook.com/groups/joins/?nav_source=tab&ordering=viewer_added")
		req.Header.Set("sec-ch-prefers-color-scheme", "dark")
		req.Header.Set("sec-ch-ua", `"Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"`)
		req.Header.Set("sec-ch-ua-full-version-list", `"Not;A=Brand";v="99.0.0.0", "Google Chrome";v="139.0.7258.66", "Chromium";v="139.0.7258.66"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-model", `""`)
		req.Header.Set("sec-ch-ua-platform", `"Linux"`)
		req.Header.Set("sec-ch-ua-platform-version", `"6.8.0"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
		req.Header.Set("x-asbd-id", "359341")
		req.Header.Set("x-fb-friendly-name", "GroupsCometAllJoinedGroupsSectionPaginationQuery")
		req.Header.Set("x-fb-lsd", "QDQ3UTMzvV6T0JNh7cF5sn")
		req.Header.Set("Cookie", "sb=E2sHZ8Dj0xZhKsq5e2frpYxk; ps_l=1; ps_n=1; datr=EZluZ78uE_dw-PpiG3FGLpyf; oo=v1; c_user=100016139237616; wd=1365x680; xs=1%3AHCeFsY1N8T3k8w%3A2%3A1737623622%3A-1%3A-1%3A%3AAcXitPQA2LhbiqQyJt3lEX8W2aoV7nQLCvLr68wOkA; fr=1SZotKB8bY6KmprvI.AWeWlBVnA-jeNZ8-qYlwnj1k7l8BazQJ9_txVdw0G_w86pheqpU.BoqbAI..AAA.0.0.BoqbNw.AWfhAahwH3Cs3x6OJ_0FG_fa_0g; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1755952478211%2C%22v%22%3A1%7D")
	}

	// Create HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Parse JSON response
	var fbResponse FacebookResponse
	if err := json.Unmarshal(body, &fbResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return &fbResponse, nil
}

func buildRequestData(cursor *string) string {
	// Build the variables JSON with cursor for pagination
	variables := map[string]interface{}{
		"count":    20,
		"cursor":   cursor, // nil for first request, then use end_cursor for subsequent requests
		"ordering": []string{"name"},
		"scale":    1,
	}

	variablesJSON, _ := json.Marshal(variables)
	variablesEncoded := url.QueryEscape(string(variablesJSON))

	// Use current config if available, otherwise fall back to hardcoded data
	if currentConfig != nil && len(currentConfig.FormData) > 0 {
		// Build form data from parsed config
		formValues := url.Values{}

		// Copy all form data from config
		for key, value := range currentConfig.FormData {
			formValues.Set(key, value)
		}

		// Override the variables parameter with our pagination cursor
		formValues.Set("variables", string(variablesJSON))

		return formValues.Encode()
	}

	// Fallback to hardcoded form data if no config is available
	baseData := "av=100016139237616&__aaid=0&__user=100016139237616&__a=1&__req=x&__hs=20323.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=GOOD&__rev=1026255062&__s=kxs8p3%3Abnqw2h%3Ay4v0lc&__hsi=7541758441371950575&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwnUW3q2ibwNw9G2Sawba1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0n24oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5o7G4-5pUfEe88o4Wm7-2K0-obUG2-azqwaW1jg2cwMwrUK2K2WEjxK2B08-269wqQ1FwgUjz89oeE-3WVU-4FqwIK6E4-mEbUaU2wwgo620XEaUcEK6Eqw&__csr=gaY8NssB6QxdNkuG8x26TdkISziiSxclqbTHERV2iPIyjXvv8D8Bj9jayiOaG-Cyi9nqV5jHEzF4yF9K9RQhDjpGRGOKAue-oRG-GGmtyfgOUKVBiVfBLACGVbxarGVqUKmmvDV4VV8a9WQEG7UhDBDGVpElzEgx24FU8-ayKq48Gii9xaAUWexfDACxC2Gl6xOdwExOEZ4xWifBzVE8EkGEsKHgy325equEfe4EaokxG4EaojAy8C4pUCi0GEbo-59V8PxC7U5ei3q3mcxG6UC9nzUuw9a323Sm2S8yErK22u2-3a2C686Ocwn41gy8DG2q68pAxGaw7twa-26EmAzA1bw6cgC0q27FVo2ja1xU1M8cWwOw5XwPqK4A6uFqCiwiojwaS0w9e6qw0kB-05FXwfi1CG0erwMwPw1eW58PBg2Tw0e9tw4ZcAU7qm09Vw2xQ4k0sgXF02eo6J28C08MGaw6ro1rUdE0X902-80m1wdi0Pk0PcM0Jkw4a0L405Ek04cU7i02tC0aswlO0ea6E4m029S09PK04PU0sOxO088S0-o1pWg1CE0Exxe2R0_w&__hsdp=gyww2sg88kBUke9Arjc8qeCGFFEmp5EPfpaFRh6AA88Bl9A695ayIt8J9abFR3b8gyA4NgFRExreRDKGiBjaiihaiiF9sGhq4iAl2pmQIzm42BcBGz9bgGRCym2yEGmVO4EunO2qhExljsrZA1nlAs5zN0IoBkmW227b8sABvvTiqaA8W2y8JAsJIwx94iqgAxqGuiGmEOl2Auh1V38Ix8Be5opCOUTVkfAzagJelazKVP4imC9gR5Iwx6hpttqHNkmqEwWAA8mua8gCcHCAmfxh5rG8xB1fooQAbBEzueJ1agPEIkbgyp5VxeUybx2EiVrcmi4kfGaqDBGEV6BhhHGUOfW9BwkofEjAxC350Fogj4xua-y1Og8zG8qUzh99888rBwyxqmeymp5jypAm4p8G6GyHykm79821GVUCiA3qfQ8c0XEkwIo-dwoU2LAg17V85q7PwzzUeS9wTwtovK2xooIq5oG486e1wQhBsweEtg5-bxNwKwSB87E-3IE4y78B1a5Ac2wlwgES0Fohw8i0zE2-wjU7G0jG5U1vU4aawhy0Sw5lwdp0b-1cwoU5S11Dw5owio7y19wjomw2BE2rwGwbO3a04UE1p80SC0M83Iwdy0tW0rS1vwfi0hm&__hblp=09C0Colw8qi2K1dDx-0Zo1voiwjE4Sq0MU6m322i2G3m1AUiwyxW2m1Mwo876u3S1kxq0BXg2fyo882pzo4GU2TyUfoK4ElwfK3W2q1aw7qwlk6FU34wpUbUe8kxK1nxO0Qo5y0LU-0C8420sm0kS1Fwo81eEnw57wnUO687Gm0vi0-84O1zxO7k263im1wwXw7nwkE39wmE28wfG0CUox60L8cE2_w861IDyE15U7a0Yo2zw8K1VwnEowo8foak1nwQwcK2614wtU4W3S3S14waKER0Ky8621KwAwdC0TQ2m12woU3LwnA0Z40VocUrg20zUa8C3iu0xUvw4iwfS2y36&__sjsp=gyww2sg88kBUke9Arj9y6zFGGqq5ChqcPSiGtkhKA8pdkCacAkGaNQyQAEKLkcIxdaoOPyXmm9lAmuWKFklyp4ECEHaAmx4F5gClKXjm42y6AiLih4aAVEdUsUWq8QbF6yqQh39aBwxiUeA2zyQU4d2roC2it2y7yQAiUYwyayA2Oh1WdyUhjxm6pIKfVE-icF1unaXDalfFykdhpEgCyFaAF2k2ygxoiGmUObBU-26E8UlooQA4q8S4bgsK7k8DBzQm8x-2W4k4EOFUyE-ql1mHxPDBw924V84K19gnzey1OgYi4i2Ezh998f9oepoopAl3pA78G6GyER0aCq4pag1hE45xS1zwa-h075zU4u0PrwEm71ElyE3lQhBsweEtg28obEf20bu78B1a5Ac2wlwgES0Fohw8i0zE2-wjU7G0jG5U1vU4aawhy0Sw5lwdp0b-1cwoU5S11Dw5owio7y19wjomw2BE2rwGwbO3a04UE1p80SC0M83Iwdy0tW0rS1vwfi0hm&__hblp=09C0Colw8qi2K1dDx-0Zo1voiwjE4Sq0MU6m322i2G3m1AUiwyxW2m1Mwo876u3S1kxq0BXg2fyo882pzo4GU2TyUfoK4ElwfK3W2q1aw7qwlk6FU34wpUbUe8kxK1nxO0Qo5y0LU-0C8420sm0kS1Fwo81eEnw57wnUO687Gm0vi0-84O1zxO7k263im1wwXw7nwkE39wmE28wfG0CUox60L8cE2_w861IDyE15U7a0Yo2zw8K1VwnEowo8foak1nwQwcK2614wtU4W3S3S14waKER0Ky8621KwAwdC0TQ2m12woU3LwnA0Z40VocUrg20zUa8C3iu0xUvw4iwfS2y36&__sjsp=gyww2sg88kBUke9Arj9y6zFGGqq5ChqcPSiGtkhKA8pdkCacAkGaNQyQAEKLkcIxdaoOPyXmm9lAmuWKFklyp4ECEHaAmx4F5gClKXjm42y6AiLih4aAVEdUsUWq8QbF6yqQh39aBwxiUeA2zyQU4d2roC2it2y7yQAiUYwyayA2Oh1WdyUhjxm6pIKfVE-icF1unaXDalfFykdhpEgCyFaAF2k2ygxoiGmUObBU-26E8UlooQA4q8S4bgsK7k8DBzQm8x-2W4k4EOFUyE-ql1mHxPDBw924V84K19gnzey1OgYi4i2Ezh998f9oepoopAl3pA78G6GyER0aCq4pag1hE45xS1zwa-h075zU4u0PrwEm71ElyE3lQhBsweEtg28obEf20bu78B1a5Ac2wlwgES0Fohw8i0zE2-wjU7G0jG5U1vU4aawhy0Sw5lwdp0b-1cwoU5S11Dw5owio7y19wjomw2BE2rwGwbO3a04UE1p80SC0M83Iwdy0tW0rS1vwfi0hm&__hblp=09C0Colw8qi2K1dDx-0Zo1voiwjE4Sq0MU6m322i2G3m1AUiwyxW2m1Mwo876u3S1kxq0BXg2fyo882pzo4GU2TyUfoK4ElwfK3W2q1aw7qwlk6FU34wpUbUe8kxK1nxO0Qo5y0LU-0C8420sm0kS1Fwo81eEnw57wnUO687Gm0vi0-84O1zxO7k263im1wwXw7nwkE39wmE28wfG0CUox60L8cE2_w861IDyE15U7a0Yo2zw8K1VwnEowo8foak1nwQwcK2614wtU4W3S3S14waKER0Ky8621KwAwdC0TQ2m12woU3LwnA0Z40VocUrg20zUa8C3iu0xUvw4iwfS2y36&__sjsp=gyww2sg88kBUke9Arj9y6zFGGqq5ChqcPSiGtkhKA8pdkCacAkGaNQyQAEKLkcIxdaoOPyXmm9lAmuWKFklyp4ECEHaAmx4F5gClKXjm42y6AiLih4aAVEdUsUWq8QbF6yqQh39aBwxiUeA2zyQU4d2roC2it2y7yQAiUYwyayA2Oh1WdyUhjxm6pIKfVE-icF1unaXDalfFykdhpEgCyFaAF2k2ygxoiGmUObBU-26E8UlooQA4q8S4bgsK7k8DBzQm8x-2W4k4EOFUyE-ql1mHxPDBw924V84K19gnzey1OgYi4i2Ezh998f9oepoopAl3pA78G6GyER0aCq4pag1hE45xS1zwa-h075zU4u0PrwEm71ElyE3lQhBsweEtg28obEf20bu78B1a5Ac2wlwgES0Fohw8i0zE2-wjU7G0jG5U1vU4aawhy0Sw5lwdp0b-1cwoU5S11Dw5owio7y19wjomw2BE2rwGwbO3a04UE1p80SC0M83Iwdy0tW0rS1vwfi0hm&__comet_req=15&fb_dtsg=NAfuwC7u4hCmbEGoAgFJlvkXW358ehEshAJ7auKp2lpjg242SNTWJpw%3A1%3A1737623622&jazoest=25421&lsd=QDQ3UTMzvV6T0JNh7cF5sn&__spin_r=1026255062&__spin_b=trunk&__spin_t=1755952472&__crn=comet.fbweb.CometGroupsJoinsRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=GroupsCometAllJoinedGroupsSectionPaginationQuery&variables=" + variablesEncoded + "&server_timestamps=true&doc_id=9974006939348139"

	return baseData
}

// ExtractAllGroups fetches all Facebook groups with pagination
func ExtractAllGroups() ([]Group, error) {
	apiURL := "https://web.facebook.com/api/graphql/"

	fmt.Println("🚀 Starting Facebook Groups extraction with pagination...")
	fmt.Println(strings.Repeat("=", 60))

	var allGroups []Group
	var cursor *string
	page := 1

	for {
		fmt.Printf("📄 Fetching page %d...\n", page)

		// Build request data with current cursor
		requestData := buildRequestData(cursor)

		// Make the API request
		response, err := makeRequest(apiURL, requestData)
		if err != nil {
			return nil, fmt.Errorf("error making request for page %d: %v", page, err)
		}

		// Extract groups from response
		edges := response.Data.Viewer.AllJoinedGroups.TabGroupsList.Edges
		pageGroups := make([]Group, 0, len(edges))

		for _, edge := range edges {
			group := Group{
				ID:                  edge.Node.ID,
				Name:                edge.Node.Name,
				URL:                 edge.Node.URL,
				ViewerJoinState:     edge.Node.ViewerJoinState,
				ViewerLastVisitTime: edge.Node.ViewerLastVisitTime,
			}
			pageGroups = append(pageGroups, group)
		}

		allGroups = append(allGroups, pageGroups...)
		fmt.Printf("✅ Found %d groups on page %d\n", len(pageGroups), page)

		// Check if there are more pages
		pageInfo := response.Data.Viewer.AllJoinedGroups.TabGroupsList.PageInfo
		if !pageInfo.HasNextPage {
			fmt.Printf("🏁 Reached end of pagination\n")
			break
		}

		// Update cursor for next page
		cursor = &pageInfo.EndCursor
		page++

		// Add delay to avoid rate limiting
		fmt.Printf("⏳ Waiting 2 seconds before next request...\n")
		time.Sleep(2 * time.Second)
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📊 EXTRACTION COMPLETE!\n")
	fmt.Printf("Total pages fetched: %d\n", page)
	fmt.Printf("Total groups found: %d\n", len(allGroups))
	fmt.Println(strings.Repeat("=", 60))

	return allGroups, nil
}

// PrintGroupResults displays the groups in a formatted way
func PrintGroupResults(groups []Group) {
	fmt.Printf("\n📋 ALL GROUP IDs AND NAMES:\n")
	fmt.Println(strings.Repeat("-", 80))

	for i, group := range groups {
		fmt.Printf("%3d. ID: %-18s | %s\n", i+1, group.ID, group.Name)
	}

	fmt.Printf("\n🔢 GROUP IDs ONLY (for easy copying):\n")
	fmt.Println(strings.Repeat("-", 50))
	for _, group := range groups {
		fmt.Println(group.ID)
	}
}

// SaveGroupsAsJSON converts groups to JSON format
func SaveGroupsAsJSON(groups []Group) ([]byte, error) {
	jsonData, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	return jsonData, nil
}
