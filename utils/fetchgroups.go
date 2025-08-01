package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
func FetchGroups(r *rand.Rand) ([]Group, error) {
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
		// Dynamic values - will be updated by UpdateFetchGroupsCookies()
		data.Set("av", "61553016608952")
		data.Set("__aaid", "0")
		data.Set("__user", "61553016608952")
		data.Set("__a", "1")
		data.Set("__req", "3d")
		data.Set("__hs", "20301.HYP:comet_pkg.2.1...0")
		data.Set("dpr", "1")
		data.Set("__ccg", "EXCELLENT")
		data.Set("__rev", "1025398783")
		data.Set("__s", "psiadt:ybawqk:hg1d1b")
		data.Set("__hsi", "7533483886033570164")
		data.Set("__dyn", "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUKbgS3q2ibwNwnof8boG0x8bo6u3y4o2Gwfi0LVEtwMw6ywIK1Rwwwg8a8462mcwfG12wOx62G5UswoEcE7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzE8FU5e3ym2SU4i5oe8464-5pUfEe88o4Wm7-2K0SEuwLyEbUGdG1QwVwwwOg2cwMwhEkxebwHwKG4UrwFg2fwxyo566k1FwgUjz89oeefx6UabDzUiBG2OUqwjVqwLwHwa211wo83KwHwOG8xG6E")
		data.Set("__csr", "gfYl1n5gjhYAYp34AAOkYA54p2InkAjOjR2c8P8mBk98CAL6ZPMAGEhir9ZjchnOeJv4QWqFaYCBXi9kChR9WOnQhaJnh_F-gxJ5GGWWqHnN3AqERF5ZlFHqGnRilFpGKiGKrZ4p9AAAQllyQCyF5TanhllqSl9KBABF6Hhu8CSq8WZEzijh9t4A9hAAv-j8WLnF5im8AApplRpq-bAV7GnAgJ2rFdqF2HHFKiAJqLpFGCGiaGiFp8B7J6hEBbiFGKFLiFFkaAALCy4uQFppdemFJ5GFACl2GG48-FidaiEB7DzoOQV6iiaVV9pkvC-Q9KhamdCGFbgSmh2SmmniyAiqGUSeLBDjGiCmVGCF6By9EC8GZ2aCUb8K48Cnxu9yppXCCCDyAicjUz-Feim5o-9xKcA-qmmi5FUogOKqby8KuqejyFK6rx2aypUgUS78C26m4K6LzEhzEydAmEvz8snG69opyKdGEqyUObghwLwyxKbGayGjADzbx2qcDDU8E8ooxWex615Gi48qBjxaq7F-2S2GEtxe1EUaemGwhE2xwaC6UO2q5Q0CE4-u5U4ieK2qm2h0Swn98y3ieBK1zCAy8lwyxPxq22egC_wk82nyEgwwwuoeU4W5E87Gm2q3K0Do20wwU4Cim7swnze58C3S0wE5ubCwHwbe1Q82Hg8FUy2eawBBwGgvxC07dU0h7w1FaU523fye06Rp87t01pC1wwkocU0wVk1lgowrE5u0mS0JU07hC321vyE0khw0ivoeE1mo1koswg9Q12De68qu4GwbRAo4-8wOwf20C83XcAh16gw1WE26yKxu1kg74wao2Aw2xo460DE6JwNxq220gaE2Wyo29wyw1ga0bsxy0bmxx03j85a9y5y80hpg1SC0bbwqUbo2Je2Ne0Eo3Mgao0Da8yk07R80Aq1fxx0nA0RU4G1tyAbgJ4w7P804BE13A6E5J1u5oc83hK0y81QE-0s60A80rFpF3xq67blu1d80V8CFC2oK2bKE3Qgiz9Edmq0jO6U22iw")
		data.Set("__hsdp", "g495A2i8ggy94p712ACzAr1Syaxp0iHhhB9eg5ky4OHC5Hp3EH2cpcgmx5gwxqkIIIx6B3O7lNbq6Gwo8W5E4ZFEh585Q8hNc4J2NceY9B68SAIaAtI9cH7Eoxs9mx54Eleh9mDclEN5yi5PkiCG69Akn9iFB68yrZ4MHcpqmwCyGXF1ayBhkzFCh6BhRKqRcSDi8tLGXcIZRileCGsp8QqtjAB4n5agFkRyij8yFZmFkIy8a6JmbiAiQUAwB7B8mGKiOEy9BcoL20x4y9FTDazalcydCCnyaPTCx2gkG8gFQa8mCbY-EJCOJHylzUgy24hbS9herxi9Z4AGbjF9o9ze48x3UjhXGiEyahbGaohGez89Q78JcEZ0OBGaAF16AoxaF7BhfBgDULgF4EgIjeqUyswhXy8x96x58GgGAEBpCbHjJ7k26ABq4AEECdyj8bC8hKS26l5yqjgtxeQtoDacz8ao8p4dgS8gFx2dUC53x6422URDcmA8xe7ElwBzhjG58mDBg-4A4qggixe0OU4y361wwBxOGoSdcF8gzQ2l2bxa69UK2y6C7m3ecgjxG2ycwsobPyh0dF4S489UOe9j4cqh1GeU8Q9wioeSP18E4C3R6jgaywxgc8gwrUf9i97oswNwsEK787u2K2q0xU99o7B1a0h388wam1Ewww4owxwCwJBDw25o1nU1JEbqwiE2zAwQwiU6e19wjFU-2a2a2G6E-1dG0FQ2e2O3m0JE3zK0Z86uGwc-1twDwfy0Fo23wg8eo88cU3awqo6K0Oo3-wOy8eE1lEO1byU3vwc-0_E3cwqo3-wgE1aK9ByEtCwdi")
		data.Set("__hblp", "0oAbEkAp2ogV8S1nDxutxK3y2W3268O2CdAxW1gBwgofo4K2q1AxG0Ko88sVo5iu5A2C58ugf8do8XCxW1mxa6EaEy3G8xx3oiVFkUgzEK68W68yawOz84u26dQ1cwMG4o551SUjCG3q2zw8y1uG2mEmzGgK6opx-11J0LG6EKcyojgce2y3OEpDwEwVyKcwQxi3m2e4EoxO2i6EK2cC7EeU4C1tw8e1pDzEaEfFoO1cwygizE8olKfK2O3ubwKwQxS58HyU9E9U7SUe8cUigrwMwQxKewCK1fwDwu89U986K799obo24z9S4UlG26i1Wxe3a1aUjAxq48hwloG3u6oy4Eb84i1LwzwJwioG1gwPwEx62-fxq5Uy6U521IwGBy988pEWF8ybyUlxG3u1Tyk1GxC5o4G1Dw-wyxe2a10yU7DK3i8wa-2m3-3iU8US7WwUx63O2C3Oi2qdwLxuUjxeqVVXxOFeu2Gcxe58evxi3215zUdonG686J1mbixLG48rwPVUnwBwgEbo6B7wgEiDK4Ea88U4e48N0VyE525GGq8xm6E4u3CcwDwPQ68S3K7pHy88HjwNBwCz8qwTwho4aU8pUixe3i3658KeCz8gz88U7zCxu2C3-3J2rzUgx29zF8hxe3q3S6U7m11gaohxiby8ngW4E1dbgjCzEbUqh9K4Ufpo4-m2-2qmdwxCyFo568yE4u4o7R1F1nwhE4m6t2omx-yAxy5m3yErCwFyonxS4FpGzE8Egy9onmUrBw4uCxq5o2tyoiG")
		data.Set("__sjsp", "g495A2i8ggy94p712ACzAr1Syaxp0iHhhB9ehkgQbR8IYGVN6JAeyIG5clO4fq4q225FiEDsxkBq2cB99LlqcblboUm4t4ArEmQChoOD9BjxGcZ3pqxG9wNOyK4KJz8x128ziCh9AbBp7hAhp4p78dyryx97AiJ78rEN8MxPaiCG69Akn9huojFJDhwHcprAk89kHCgiEFkl8WoNFktrCJjdAB8tKiUKZTF9k4cp8iqudzEFGkxHm99cBBrlyqShz8xdmbiAiQUAwGunhpoyOEy9BcoHsm9h8yqtUwyc8BoCdmaHyVEnBiF6gGdCAzvai2m9ghwxBLgF4VK58R4DyQXF0ElK6E-4QuWDy8kGaokzEO2t0uWyFaghyy7F7CzEHUphai8IjO2Ug84u6ocbz8C4AUgwuUFa5QqEcFkm9F0WyUsz84Kh1O8gJ1Xyoke293o8EO4Ud89oQkWxi2h0yghF1104Hw9a12GCdxG1Ews67m1lwgUO1NwLca40XjwlykN36BwjQ9w8lIMia19wZhAQ0ge0ua0ld04zO02ko8o04nW02rq")
		data.Set("__comet_req", "15")
		data.Set("fb_dtsg", "NAfu4ikzYjWt_U5CvSMrLJavDntdTRnEYkzfEaunOgxMEeIY3Ua5AVA:30:1737396043")
		data.Set("jazoest", "25706")
		data.Set("lsd", "Ojhd_SK0sSq0wM2RpBbKAc")
		data.Set("__spin_r", "1025398783")
		data.Set("__spin_b", "trunk")
		data.Set("__spin_t", "1754025902")
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

		// Dynamic headers - will be updated by UpdateFetchGroupsCookies()
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("cookie", "ps_l=1; ps_n=1; sb=EyQJZ_p60kcMqjhfiPbKiQun; datr=OY-OZwqAfui6fV8_zZnYl_zQ; c_user=61553016608952; wd=1366x681; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1754025907720%2C%22v%22%3A1%7D; fr=1HwRYFMccYQ60Ti6R.AWd6h7zAp-sNZqN_F07CLmWjgvAaVuXWRlNHmjKeQjIZOgLS0As.BojE-z..AAA.0.0.BojE-z.AWcBkjQ1eLuFkKOsxLH8zBzuVKo; xs=30%3AdG3sK9-gG0z1jg%3A2%3A1737396043%3A-1%3A-1%3A%3AAcUaHk4sWBPmVZFPMl8DUMhRFWcpUnwnlwumSKhzas8")
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
		req.Header.Set("x-fb-lsd", "Ojhd_SK0sSq0wM2RpBbKAc")

		// Make the request
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
