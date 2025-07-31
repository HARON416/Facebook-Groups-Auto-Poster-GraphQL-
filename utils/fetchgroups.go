package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Group represents a Facebook group
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	JoinState   string `json:"viewer_join_state"`
	LastVisited int64  `json:"viewer_last_visited_time"`
}

// PageInfo represents pagination information
type PageInfo struct {
	EndCursor   string `json:"end_cursor"`
	HasNextPage bool   `json:"has_next_page"`
}

// GroupsResponse represents the Facebook GraphQL response
type GroupsResponse struct {
	Data struct {
		Viewer struct {
			AllJoinedGroups struct {
				TabGroupsList struct {
					Edges []struct {
						Node Group `json:"node"`
					} `json:"edges"`
					PageInfo PageInfo `json:"page_info"`
				} `json:"tab_groups_list"`
			} `json:"all_joined_groups"`
		} `json:"viewer"`
	} `json:"data"`
}

// FetchGroups makes a request to Facebook's GraphQL API to fetch all joined groups with pagination
func FetchGroups() ([]Group, error) {
	var allGroups []Group
	var cursor *string
	pageCount := 0

	for {
		pageCount++
		fmt.Printf("Fetching page %d...\n", pageCount)

		// Build the variables JSON
		variables := map[string]interface{}{
			"count":    20,
			"ordering": []string{"name"},
			"scale":    1,
		}
		if cursor != nil {
			variables["cursor"] = *cursor
		} else {
			variables["cursor"] = nil
		}

		variablesJSON, err := json.Marshal(variables)
		if err != nil {
			return nil, fmt.Errorf("error marshaling variables: %w", err)
		}

		data := url.Values{}
		data.Set("av", "61560452168137")
		data.Set("__aaid", "0")
		data.Set("__user", "61560452168137")
		data.Set("__a", "1")
		data.Set("__req", "14")
		data.Set("__hs", "20299.HYP:comet_pkg.2.1...0")
		data.Set("dpr", "1")
		data.Set("__ccg", "GOOD")
		data.Set("__rev", "1025302125")
		data.Set("__s", "6wmevx:wrk89z:mdljwg")
		data.Set("__hsi", "7532770811035525370")
		data.Set("__dyn", "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwnUW3q2ibwNw9G2Sawba1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0n24oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5o7G4-5pUfEe88o4Wm7-2K0-obUG2-azqwaW1jg2cwMwrUK2K2WEjxK2B08-269wqQ1FwgUjz89oeE4WVU-4FqwIK6E4-mEbUaU2wwgo620XEaUcGy8qxG")
		data.Set("__csr", "gngVb5MOyT2IcEIW94iERlsJgBfnOmBZWukGImGetsYiWFZZ88N298N5sIkIH9aKOaJHZifQFHmuhvH9ZHtrCmiGnJQiQiEGGyJ7VKK8ScUyoCqVemnjAABBxaXrgGmiWhoyu4oy5HyFpoG16AgmzuaCJvhUCE8XorwJz8yfyp8Ou2ym3OVpUyq4AfxCQ9xi19BxeexmcCByE9FEGfCGUb85GE_x269UCcDyWwi8-XzE8FoWfwpooxa2a48W14UqGEbod8K3CFEdE4S2qm8DwCzrwpEy2-jwZxG1kwPDwnobE3fxCi2ui13wUho2Fw5ww4Swlo460zGQEK1RxF1Sl3810EtAwyxm0X8ckp1i2XU3ixzDxC1hBxu0D8SfwNzd0Gw9e12w921mwHxG5UG5p8O2S01a8w8u03xC062E03vuxy3m0fog2oo0Gi0e8o1gEzga80CK04gA15G096Ow0FuC9g0pMg0hYwHw6Yw1LOC0fnxO04k80aXE3Ww0zVykcy80JuU0Cq0fgw")
		data.Set("__hsdp", "gwwAE98wV75xnaF6iGaHqiiD5a8i7HGhayD4L7F8GhqAW8CPaCzmDqIbKGAhxy8gey546Kp69jp4AAC9GWymsSq8Uyy-uH8xLFhbl2iaiBSgWJpllBz-noKEaAgGeIOGIK4mFhDE44oAFO9d84P5iJNiAyTGB5QNhzkAMymV3jy19pqGgPd55j69TeGmF3yiKO4AjiGUaUR28giaXDwxx4OGbggyki4t96qhFqzk9kF3A4226xa4EGKjzSviF4G47Y9sNrAhuBAGzDzDHxmCq5q-6kmi58uohGG7AGODyF9ayQ5CCz2HVrGuQaihosDqCCp3Jsxp8ixcxAiEC2ta8yWgSt1cw4e2Kq2gCao844Gojjy9Amp7KexZ163GApe6YWgwiTcEgxTBxp4BGKjBhA5AEOhDwkEVpC262mVKuiax2u2CJ14E8o762bpU4a2K1Pyk2Oq0-E2my8eB60uo830QgqwIxqhf5wHyXFJ1ebwh85i2qm3R1m2GE55xN29xj70sorwj8vxJ0Mxe1KwSxWQ1NxG3IxE5eu0vq2C0wA6U7C0jy7EC0pS0v-0Qqw4FwiU3QwgA0nei6E0xC1ZyU2Twn81UU1WU0wC08iw7bwi81D81zEhyGxe0gO0WU1eo720ckw")
		data.Set("__hblp", "08W4VU2cwww8y488V8fEy2m1Bw2bErwJG15DwcK0FEapA1EwNxu5EhxW7o4W1-yEy4oK1qxG0iu698uwWwdu3W7A0wu2G3a5UowIUb87S0h6086w8a3m1jg2gwJzo6i2q586N05mwr8swMw8C1eK0i-0lu1IwNwaS0Bo3kw7twibw4hxem5UfE5S08uwiEbo3iwSg7W0Z98qwcG0hy3m1ZyU2Twn88U1K9E5i0QU4e1-wjU1Kp816k0YE1IEoy8CbwFwJw5Kw6egnw8W0C8sw9O360jC1UwvE5y8xK4EfESm1xwDwg81Zo")
		data.Set("__sjsp", "gwwAE98wV75xnaF6iGaHqiiD52QxWWAiEFPXNtKyd5GjEygCGqdqrGMKWGcYgh495hExb8uahAV4leAAaGq9gyQ8yEyy-uHmH4V598QEAxWnpOaRBgOfyEOEaUm9BHu2aleEC7pAEB2Q7Ey_gpAz9ohxK68W5m2mi8DxtlBy9dyOD89h9ofk8xlrKu264jaEK4oB1ybqhFohykCkV10wwb8-jxDwwpovKiqhpaEVUZqUlwxAUpgvwRwHyF98V1p34exF2AA2a8yEqABAwNAG9wDiy8KAdy8j81LCwFwh64V8yp5DyUW1MhF8rPF21bqa4UtVonABBVedxy561FGpwhXCyp8rwrU2kSu12wHwsUB0PwfG0DUeB60uo2zxqhf5wXFK58521kwCBwZg4e1hosgyokNM760F40IU0haw6OxW9w6tw7_wd6E1ao1g804J2aG4U0qUw")
		data.Set("__comet_req", "15")
		data.Set("fb_dtsg", "NAfs8fVe8u0XAj8gIJR9U2FTHlLyVB8fIlGbHJKSvIv_Yl0ZzPgOU8Q:8:1735880128")
		data.Set("jazoest", "25300")
		data.Set("lsd", "IN0jSIIXfq_xlEmqLUKjhe")
		data.Set("__spin_r", "1025302125")
		data.Set("__spin_b", "trunk")
		data.Set("__spin_t", "1753859876")
		data.Set("__crn", "comet.fbweb.CometGroupsJoinsRoute")
		data.Set("fb_api_caller_class", "RelayModern")
		data.Set("fb_api_req_friendly_name", "GroupsCometAllJoinedGroupsSectionPaginationQuery")
		data.Set("variables", string(variablesJSON))
		data.Set("server_timestamps", "true")
		data.Set("doc_id", "9974006939348139")

		encodedData := data.Encode()

		req, err := http.NewRequest("POST", "https://web.facebook.com/api/graphql/", strings.NewReader(encodedData))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("cookie", "sb=XWoHZ-GdXVwlvXZrSaFe7gwz; ps_l=1; ps_n=1; datr=78xqZ8Fc-vdpp5Ii5nGp2A0P; c_user=61560452168137; wd=1366x681; fr=18xPTThxW1UIlQusI.AWcRHtl3cn_NKIw6QFzOPwSsIlb_bA9uoucmRDHBbkppadBXHUA.BoicQ4..AAA.0.0.BoicQ4.AWfW-QydIVuLuVA3N44d62vPRKY; xs=8%3AeMpQ9UiVlPPogg%3A2%3A1735880128%3A-1%3A-1%3AqSC9VrelgmIhZA%3AAcXfb2eUWBgbwzFM_XwWkqqSCL0VeZVoUUJ3VBOvlDrA; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1753859884424%2C%22v%22%3A1%7D")
		req.Header.Set("origin", "https://web.facebook.com")
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("referer", "https://web.facebook.com/groups/joins/?nav_source=tab&ordering=viewer_added")
		req.Header.Set("sec-ch-prefers-color-scheme", "light")
		req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
		req.Header.Set("sec-ch-ua-full-version-list", `"Not)A;Brand";v="8.0.0.0", "Chromium";v="138.0.7204.168", "Google Chrome";v="138.0.7204.168"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-model", `""`)
		req.Header.Set("sec-ch-ua-platform", `"Linux"`)
		req.Header.Set("sec-ch-ua-platform-version", `"6.8.0"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
		req.Header.Set("x-asbd-id", "359341")
		req.Header.Set("x-fb-friendly-name", "GroupsCometAllJoinedGroupsSectionPaginationQuery")
		req.Header.Set("x-fb-lsd", "IN0jSIIXfq_xlEmqLUKjhe")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}

		// Parse the response
		var groupsResp GroupsResponse
		if err := json.Unmarshal(body, &groupsResp); err != nil {
			return nil, fmt.Errorf("error parsing response: %w", err)
		}

		// Extract groups from this page
		for _, edge := range groupsResp.Data.Viewer.AllJoinedGroups.TabGroupsList.Edges {
			allGroups = append(allGroups, edge.Node)
		}

		fmt.Printf("Page %d: Found %d groups (Total: %d)\n", pageCount, len(groupsResp.Data.Viewer.AllJoinedGroups.TabGroupsList.Edges), len(allGroups))

		// Check if there are more pages
		pageInfo := groupsResp.Data.Viewer.AllJoinedGroups.TabGroupsList.PageInfo
		if !pageInfo.HasNextPage {
			fmt.Printf("No more pages. Total groups fetched: %d\n", len(allGroups))
			break
		}

		// Set cursor for next page
		cursor = &pageInfo.EndCursor
	}

	return allGroups, nil
}
