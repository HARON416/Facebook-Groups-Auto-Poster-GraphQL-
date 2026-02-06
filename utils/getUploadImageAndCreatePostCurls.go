package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var post string = `✅ LIPA MDOGO MDOGO
✅ Samsung Galaxy S23 Ultra
✅ Deposit + First Week = Ksh 27,999
✅ Weekly Payment = Ksh 2,270
✅ Payment Period = 52 weeks
✅ Location: Pioneer Building, Kimathi Street
✅ Call/WhatsApp 0718448461`

var postImagePaths = []string{
	"./images/samsung1.jpg",
	"./images/samsung2.jpg",
	"./images/samsung3.jpg",
}

// escapeForDollarQuote escapes bytes for safe inclusion inside a $'...'
func escapeForDollarQuote(s string) string {
	b := []byte(s)
	var sb strings.Builder
	for _, c := range b {
		switch c {
		case '\r':
			sb.WriteString("\\r")
		case '\n':
			sb.WriteString("\\n")
		case '\\':
			sb.WriteString("\\\\")
		case '\'':
			sb.WriteString("\\'")
		default:
			if c < 0x20 || c > 0x7e {
				sb.WriteString(fmt.Sprintf("\\x%02x", c))
			} else {
				sb.WriteByte(c)
			}
		}
	}
	return sb.String()
}

func GetUploadImageAndCreatePostCurls(page *rod.Page) (string, string) {
	// Automatically accept all JavaScript dialogs (e.g., "Leave site?" confirmations)
	go page.EachEvent(func(e *proto.PageJavascriptDialogOpening) bool {
		_ = proto.PageHandleJavaScriptDialog{Accept: true}.Call(page)
		return false // continue listening for more dialogs
	})()

	url := "https://web.facebook.com/groups/1481181995617674"
	page = page.MustNavigate(url).MustWaitLoad().MustWaitDOMStable()
	page.Activate()

	var mu sync.Mutex
	var uploads []map[string]any
	var graphqlRequests []map[string]any
	var uploadImageCurl string
	var createPostCurl string

	spans := page.MustElements(`span.x1lliihq.x6ikm8r.x10wlt62.x1n2onr6`)

	var writeSomethingSpan *rod.Element

	for _, s := range spans {
		if s.MustText() == "Write something..." {
			writeSomethingSpan = s
			break
		}
	}

	writeSomethingSpan.MustClick()

	fmt.Println("Clicked on write something button")

	time.Sleep(10 * time.Second)

	dialog := page.MustElement(`div[role="dialog"]`)

	// start hijack router here so initial page load and element lookups are not affected
	router := page.HijackRequests()
	defer router.MustStop()

	router.MustAdd("*", func(ctx *rod.Hijack) {
		req := ctx.Request

		// Only XHR / Fetch (uploads are usually fetch/XHR)
		if req.Type() != proto.NetworkResourceTypeXHR && req.Type() != proto.NetworkResourceTypeFetch {
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}

		urlStr := req.URL().String()

		// normalize headers to map[string]string early so we can detect preflight
		origHeaders := req.Headers()
		norm := make(map[string]string)
		for k, v := range origHeaders {
			norm[strings.ToLower(k)] = v.String()
		}

		// Only capture the specific photo upload endpoint
		if strings.Contains(urlStr, "upload.facebook.com/ajax/react_composer/attachments/photo/upload") {
			// Skip preflight CORS requests: OPTIONS method or presence of Access-Control request headers
			if strings.EqualFold(req.Method(), "OPTIONS") {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}
			if _, ok := norm["access-control-request-method"]; ok {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}
			// Prefer only actual POST uploads
			if !strings.EqualFold(req.Method(), "POST") {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}

			body := req.Body()
			opName := norm["x-fb-friendly-name"]
			if opName == "" {
				fmt.Println("🔎 Upload friendly name: (none)")
			} else {
				fmt.Printf("🔎 Upload friendly name: %s\n", opName)
			}

			entry := map[string]any{
				"time":         time.Now().Format(time.RFC3339Nano),
				"url":          urlStr,
				"method":       req.Method(),
				"postData":     string(body),
				"headers":      norm,
				"friendlyName": opName,
			}
			mu.Lock()
			uploads = append(uploads, entry)
			mu.Unlock()
			fmt.Println("Captured upload request:", urlStr)
		}

		// Capture GraphQL requests (only ComposerStoryCreateMutation)
		if strings.Contains(urlStr, "/api/graphql") || strings.Contains(urlStr, "/graphql") {
			// Skip preflight OPTIONS
			if strings.EqualFold(req.Method(), "OPTIONS") {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
				return
			}

			// Only capture ComposerStoryCreateMutation
			if norm["x-fb-friendly-name"] == "ComposerStoryCreateMutation" {
				body := req.Body()
				gEntry := map[string]any{
					"time":     time.Now().Format(time.RFC3339Nano),
					"url":      urlStr,
					"method":   req.Method(),
					"postData": string(body),
					"headers":  norm,
				}
				mu.Lock()
				graphqlRequests = append(graphqlRequests, gEntry)
				mu.Unlock()
				fmt.Println("Captured GraphQL request: ComposerStoryCreateMutation")
			}
		}

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	fmt.Println(len(page.MustElements(`input[accept="image/*,image/heif,image/heic,video/*,video/mp4,video/x-m4v,video/x-matroska,.mkv"]`)))

	if dialog.MustHas(`input[type="file"]`) {
		fmt.Println("Dialog has file input")
		dialog.MustElement(`input[type="file"]`).MustSetFiles(postImagePaths...)
	} else if page.MustHas(`input[accept="image/*,image/heif,image/heic,video/*,video/mp4,video/x-m4v,video/x-matroska,.mkv"]`) {
		page.MustElement(`input[accept="image/*,image/heif,image/heic,video/*,video/mp4,video/x-m4v,video/x-matroska,.mkv"]`).MustSetFiles(postImagePaths...)
	} else {
		fmt.Println("Dialog does not have file input")
		if dialog.MustHas(`div[aria-label="Photo/video"]`) {
			dialog.MustElement(`div[aria-label="Photo/video"]`).MustClick()
			dialog.MustElement(`input[type="file"]`).MustSetFiles(postImagePaths...)

		} else {

			fmt.Println("Dialog does not have photo/video element")
			return "", ""

		}
	}

	fmt.Println("Images inserted")

	if dialog.MustHas(`div[role="textbox"][contenteditable="true"]`) {
		fmt.Println("Dialog has textbox for description")
		dialog.MustElement(`div[role="textbox"][contenteditable="true"]`).MustInput(post)
	} else {
		fmt.Println("Dialog does not have textbox, trying alternative selector")
		page.MustElement(`div[aria-placeholder="Create a public post…"]`).MustInput(post)
	}

	fmt.Println("Description inserted")

	// wait for network activity to settle
	time.Sleep(30 * time.Second)

	mu.Lock()
	if len(uploads) > 0 {
		f, err := os.Create("uploads.jsonl")
		if err == nil {
			enc := json.NewEncoder(f)
			for _, u := range uploads {
				_ = enc.Encode(u)
			}
			_ = f.Close()
			fmt.Println("Wrote", len(uploads), "upload request(s) to uploads.jsonl")
		} else {
			fmt.Println("Failed to create uploads.jsonl:", err)
		}
	} else {
		fmt.Println("No upload requests captured")
	}
	mu.Unlock()

	// If we have at least one upload, build a curl from the last captured entry
	if len(uploads) > 0 {
		last := uploads[len(uploads)-1]
		urlStr, _ := last["url"].(string)
		postData, _ := last["postData"].(string)
		headersIface, _ := last["headers"].(map[string]string)

		// Add sensible defaults for headers that may be missing so curl matches original template
		if headersIface == nil {
			headersIface = make(map[string]string)
		}
		if _, ok := headersIface["accept-language"]; !ok {
			headersIface["accept-language"] = "en-US,en;q=0.9"
		}
		if _, ok := headersIface["priority"]; !ok {
			headersIface["priority"] = "u=1, i"
		}
		if _, ok := headersIface["sec-fetch-dest"]; !ok {
			headersIface["sec-fetch-dest"] = "empty"
		}
		if _, ok := headersIface["sec-fetch-mode"]; !ok {
			headersIface["sec-fetch-mode"] = "cors"
		}
		if _, ok := headersIface["sec-fetch-site"]; !ok {
			headersIface["sec-fetch-site"] = "same-site"
		}

		// Prefer the captured request's Cookie header if present, else fall back to page cookies
		cookieStr := ""
		if v, ok := headersIface["cookie"]; ok && v != "" {
			cookieStr = v
		} else {
			cookies := page.MustCookies()
			for _, c := range cookies {
				if cookieStr != "" {
					cookieStr += "; "
				}
				cookieStr += c.Name + "=" + c.Value
			}
		}
		cookieLine := fmt.Sprintf("  -b '%s'", cookieStr)

		curlCmd := fmt.Sprintf("curl '%s' \\", urlStr)

		// header order to match original template
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

		for _, k := range headerOrder {
			if v, ok := headersIface[k]; ok && v != "" {
				if k == "content-type" {
					curlCmd += fmt.Sprintf("\n  -H '%s: %s' \\", k, v)
					if cookieStr != "" {
						curlCmd += fmt.Sprintf("\n%s \\\n+", cookieLine)
					}
				} else {
					curlCmd += fmt.Sprintf("\n  -H '%s: %s' \\", k, v)
				}
			}
		}

		// Append data-raw. Use $'...'' to preserve CRLF sequences where present.
		curlCmd += fmt.Sprintf("\n  --data-raw $'%s'", escapeForDollarQuote(postData))

		uploadImageCurl = curlCmd

		if err := os.WriteFile("uploads_curl.txt", []byte(curlCmd+"\n\n"), 0644); err != nil {
			fmt.Println("Failed to write uploads_curl.txt:", err)
		} else {
			fmt.Println("Wrote curl for last upload to uploads_curl.txt")
		}

		// sanitize the file we just wrote in case of stray characters
		// (remove any lines that are just a single '+')
		if data, err := os.ReadFile("uploads_curl.txt"); err == nil {
			s := strings.ReplaceAll(string(data), "\r", "")
			lines := strings.Split(s, "\n")
			var out []string
			for _, L := range lines {
				if strings.TrimSpace(L) == "+" {
					continue
				}
				out = append(out, L)
			}
			_ = os.WriteFile("uploads_curl.txt", []byte(strings.Join(out, "\n")+"\n"), 0644)
		}
	}

	if dialog.MustHas(`div[aria-label="Post"]`) {
		fmt.Println("Dialog has Post button")
		dialog.MustElement(`div[aria-label="Post"]`).MustClick()
	} else {
		fmt.Println("Dialog does not have Post button, trying page level selector")
		page.MustElement(`div[aria-label="Post"]`).MustClick()
	}

	fmt.Println("Post button clicked")

	for range 120 {
		if dialog.MustVisible() {
			if strings.Contains(dialog.MustText(), "We limit") {
				fmt.Println("Hit posting limit. Please wait and try again later.")
				return "", ""
			}
			fmt.Println("Posting...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("Ad posted successfully 🥳🥳🥳🥳🥳🥳🥳🥳🥳")
		break
	}

	// Write captured GraphQL requests to file
	mu.Lock()
	if len(graphqlRequests) > 0 {
		f, err := os.Create("graphql_requests.jsonl")
		if err == nil {
			enc := json.NewEncoder(f)
			for _, g := range graphqlRequests {
				_ = enc.Encode(g)
			}
			_ = f.Close()
			fmt.Println("Wrote", len(graphqlRequests), "GraphQL request(s) to graphql_requests.jsonl")
		} else {
			fmt.Println("Failed to create graphql_requests.jsonl:", err)
		}

		// Build curl command from the first GraphQL request
		if len(graphqlRequests) > 0 {
			req := graphqlRequests[0]
			urlStr, _ := req["url"].(string)
			postData, _ := req["postData"].(string)
			headersMap, _ := req["headers"].(map[string]string)

			// Add missing headers with defaults
			if _, ok := headersMap["accept-language"]; !ok {
				headersMap["accept-language"] = "en-US,en;q=0.9"
			}
			if _, ok := headersMap["priority"]; !ok {
				headersMap["priority"] = "u=1, i"
			}
			if _, ok := headersMap["sec-fetch-dest"]; !ok {
				headersMap["sec-fetch-dest"] = "empty"
			}
			if _, ok := headersMap["sec-fetch-mode"]; !ok {
				headersMap["sec-fetch-mode"] = "cors"
			}
			if _, ok := headersMap["sec-fetch-site"]; !ok {
				headersMap["sec-fetch-site"] = "same-origin"
			}

			// Build curl command with proper header order
			curlCmd := fmt.Sprintf("curl '%s' \\\n", urlStr)

			// Define header order to match the template
			headerOrder := []string{
				"accept",
				"accept-language",
				"content-type",
			}

			for _, key := range headerOrder {
				if val, ok := headersMap[key]; ok && val != "" {
					curlCmd += fmt.Sprintf("  -H '%s: %s' \\\n", key, val)
				}
			}

			// Add cookie header with -b flag
			if cookie, ok := headersMap["cookie"]; ok && cookie != "" {
				curlCmd += fmt.Sprintf("  -b '%s' \\\n", cookie)
			}

			// Continue with remaining headers
			remainingHeaders := []string{
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

			for _, key := range remainingHeaders {
				if val, ok := headersMap[key]; ok && val != "" {
					curlCmd += fmt.Sprintf("  -H '%s: %s' \\\n", key, val)
				}
			}

			// Add data-raw
			curlCmd += fmt.Sprintf("  --data-raw '%s'", postData)

			createPostCurl = curlCmd

			if err := os.WriteFile("graphql_curl.txt", []byte(curlCmd+"\n"), 0644); err != nil {
				fmt.Println("Failed to write graphql_curl.txt:", err)
			} else {
				fmt.Println("Wrote curl for GraphQL request to graphql_curl.txt")
			}
		}
	} else {
		fmt.Println("No GraphQL requests captured")
	}
	mu.Unlock()

	return uploadImageCurl, createPostCurl
}
