package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func GetFetchGroupsCurl(page *rod.Page) string {
	var capturedCurl string
	page = page.MustNavigate("https://web.facebook.com/groups/joins/?nav_source=tab&ordering=viewer_added").MustWaitLoad().MustWaitDOMStable()

	page.Activate()

	router := page.HijackRequests()
	defer router.MustStop()

	router.MustAdd("*", func(ctx *rod.Hijack) {
		req := ctx.Request

		// Only XHR / Fetch
		if req.Type() != proto.NetworkResourceTypeXHR &&
			req.Type() != proto.NetworkResourceTypeFetch {
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}

		urlStr := req.URL().String()

		// Only Facebook GraphQL
		if !strings.Contains(urlStr, "graphql") {
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}

		// Log the request
		fmt.Printf("➡️ GraphQL request: %s %s\n", req.Method(), urlStr)

		// Get request body
		body := req.Body()
		if len(body) > 0 {
			fmt.Printf("Request body (%d bytes)\n", len(body))
			// Get cookies from page
			cookies := page.MustCookies()
			var cookieStr string
			for _, c := range cookies {
				if cookieStr != "" {
					cookieStr += "; "
				}
				cookieStr += c.Name + "=" + c.Value
			}
			cookieLine := fmt.Sprintf("  -b '%s'", cookieStr)
			// Generate full curl command
			curlCmd := "curl 'https://web.facebook.com/api/graphql/' \\"
			// Normalize headers to map[string]string so we can safely add/modify values
			origHeaders := req.Headers()
			norm := make(map[string]string)
			for k, v := range origHeaders {
				norm[strings.ToLower(k)] = v.String()
			}
			// Add missing headers only if not present
			if _, ok := norm["accept-language"]; !ok {
				norm["accept-language"] = "en-US,en;q=0.9"
			}
			if _, ok := norm["priority"]; !ok {
				norm["priority"] = "u=1, i"
			}
			if _, ok := norm["sec-fetch-dest"]; !ok {
				norm["sec-fetch-dest"] = "empty"
			}
			if _, ok := norm["sec-fetch-mode"]; !ok {
				norm["sec-fetch-mode"] = "cors"
			}
			if _, ok := norm["sec-fetch-site"]; !ok {
				norm["sec-fetch-site"] = "same-origin"
			}

			// Log the friendly name if present
			opName := norm["x-fb-friendly-name"]
			if opName == "" {
				fmt.Println("🔎 GraphQL friendly name: (none)")
			} else {
				fmt.Printf("🔎 GraphQL friendly name: %s\n", opName)
			}

			// Only process the specific GraphQL operation we care about
			if opName != "GroupsCometAllJoinedGroupsSectionPaginationQuery" {
				// continue the request and skip curl generation
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}

			// Header order to match template
			headerOrder := []string{
				"accept",
				"accept-language",
				"content-type",
				"origin",
				"priority",
				"referer",
				"sec-ch-prefers-color-scheme",
				"sec-ch-ua",
				"sec-ch-ua-full-version-list",
				"sec-ch-ua-mobile",
				"sec-ch-ua-model",
				"sec-ch-ua-platform",
				"sec-ch-ua-platform-version",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"x-asbd-id",
				"x-fb-friendly-name",
				"x-fb-lsd",
			}
			var headerLines []string
			for _, displayKey := range headerOrder {
				if val, ok := norm[displayKey]; ok && val != "" {
					if displayKey == "content-type" {
						headerLines = append(headerLines, fmt.Sprintf("  -H '%s: %s' \\", displayKey, val))
						headerLines = append(headerLines, cookieLine+" \\")
					} else {
						headerLines = append(headerLines, fmt.Sprintf("  -H '%s: %s' \\", displayKey, val))
					}
				}
			}
			for _, h := range headerLines {
				curlCmd += "\n" + h
			}
			curlCmd += fmt.Sprintf("\n  --data-raw '%s'", string(body))

			// Capture the curl command
			capturedCurl = curlCmd

			// Save to file (overwrite existing contents)
			if err := os.WriteFile("graphql_curls.txt", []byte(curlCmd+"\n\n"), 0644); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}
		}

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	// page.MustElement(`div[aria-haspopup="menu"]`).MustClick()
	// page.MustWaitLoad().MustWaitDOMStable()
	// sortCard := page.MustElement(`div[aria-label="Sort joined groups"]`)
	// sortCard.MustElements(`div[role="menuitemradio"]`)[0].MustClick()
	// page.MustWaitLoad().MustWaitDOMStable()

	for range 3 {
		page.Mouse.MustScroll(0, 1000)
	}

	return capturedCurl
}
