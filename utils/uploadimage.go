package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
)

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Define structs for the API responses (adapt to the actual responses)
type UploadResponse struct {
	Payload struct {
		IsSpherical      bool        `json:"isSpherical"`
		Height           int         `json:"height"`
		ImageSrc         string      `json:"imageSrc"`
		MediaLocation    interface{} `json:"mediaLocation"`
		OriginalPhotoID  string      `json:"originalPhotoID"`
		PhotoID          string      `json:"photoID"`
		SphericalPhotoID string      `json:"sphericalPhotoID"`
		ThumbSrc         string      `json:"thumbSrc"`
		Width            int         `json:"width"`
		MediaTakenTime   interface{} `json:"mediaTakenTime"`
	} `json:"payload"`
}

func UploadImage(imagePath string) (string, error) {
	fmt.Printf("Uploading image: %s\n", imagePath)

	// 1. Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("error opening image file: %w", err)
	}
	defer file.Close()

	// 2. Create a buffer to hold the request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 3. Add the form fields
	err = writer.WriteField("source", "8")
	if err != nil {
		return "", fmt.Errorf("error adding source field: %w", err)
	}
	err = writer.WriteField("profile_id", "61560452168137") // Replace with your profile ID
	if err != nil {
		return "", fmt.Errorf("error adding profile_id field: %w", err)
	}
	err = writer.WriteField("waterfallxapp", "comet")
	if err != nil {
		return "", fmt.Errorf("error adding waterfallxapp field: %w", err)
	}

	// 4. Create the form file field for the image
	fileField, err := writer.CreateFormFile("farr", filepath.Base(imagePath))
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	// 5. Copy the image data to the form file field
	_, err = io.Copy(fileField, file)
	if err != nil {
		return "", fmt.Errorf("error copying image data: %w", err)
	}

	// 6. Add the upload_id field with a unique timestamp and random component
	uploadID := fmt.Sprintf("jsc_c_%d_%s", time.Now().UnixNano(), fmt.Sprintf("%x", time.Now().UnixNano())[:6])
	err = writer.WriteField("upload_id", uploadID) // Generate unique upload_id for each upload
	if err != nil {
		return "", fmt.Errorf("error adding upload_id field: %w", err)
	}

	// 7. Close the multipart writer
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	// 8. Construct the URL with all the required query parameters from the curl request
	apiURL := "https://upload.facebook.com/ajax/react_composer/attachments/photo/upload"
	baseURL, err := url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("error parsing the baseURL: %w", err)
	}

	// Add query parameters from the working curl request
	params := url.Values{}
	params.Add("av", "61560452168137")
	params.Add("__aaid", "0")
	params.Add("__user", "61560452168137")
	params.Add("__a", "1")
	params.Add("__req", "1g")
	params.Add("__hs", "20299.HYP:comet_pkg.2.1...0")
	params.Add("dpr", "1")
	params.Add("__ccg", "GOOD")
	params.Add("__rev", "1025305868")
	params.Add("__s", "x4pc2s:972o5m:rd3gh9")
	params.Add("__hsi", "7532857753518398153")
	params.Add("__dyn", "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0-E4a3a4oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5oe8464-5pU9UmwUwxwjFovUaU3VwLyEbUGdG0HE88cA0z8c84q58jyUaUbGxe6Uak0zU8oC1hxB0qo4e4UO2m3G1eKnzUiBG2OUqwjVqwLwHwa211zU520XEaUcGy8qxG")
	params.Add("__csr", "g417gnigWyhcQrkj48G6Pihf4n2s8PNt8hkGqaP85tEDuhN5PirkzkjSDiQGdOFbiPGLILlCmsyGh8xqWemLrXOr-LnBQkHmijXrhkpfyQhXmq8hlQ-AoyqSyWRGiiF8pCylKSqXyVqiWjFBhbUyV7GiA_CGWGamGiyoABaLl2AVFumh28yFqzAi9BBy4qWyFUmxe5p8G8y9aV98CnyV9EB38KmdzuZ2prxGUrCDgkQ8K9prG8hEhzUzyUBoDCV8hAglAx28zFUG6QiU-eK5EyiivKqbKi9glGqFUpyEgy9qKHzWzA4oO78KaK4XgW69V8aFayo-9zHwECBG9xXx2dyEmwOCxq9yUG9z8x28rhE-i68cEOex6ayBBzEkwOybwCBwCAUa-qezUS2CU9U-2W48G4U88cUfUy3y6ryojxO-4UeUbawWK3658dUhwywi42NefxO3e2GE4mfBxi0xXwb60Wo2vxe2udIw46exm8wkogw8etkawxAg4m1xCix-fKGzd0RwdRlwyg2bGi6U4G0C83YwVwWw9O3aaw3dEOdgpw1me02ge1Dwky0zw1UKi2501BW0ahxu0dFw9q8g07iN00wKw0hUWxG0mq0gx07yo1F814qh8CFT4010G0Rm0l25EjAU3twgo5K0EjFGqfwe5027Gw2cFio1fFiw3SUow1kIU3fw1s25U5h07WxS0v-0cKw1ono0hq8m1jgym2e0dpy9S6Emw_w3EUYEoU0NW0eKw5rw20A8w1c2gE89uu9gO2nFwj80zO4U2QU88K1IWb4efrgobK09Nw3Q8")
	params.Add("__hsdp", "g49494FEOcxaF44EfaG56EjPhaaxaiEkxGraxO8wMG99k8EmwUgroEHcBh28KwB8hrWFJqf6AJIJEsKxJFkN2aAGkszEX9AB2Q5Q9EW8whaxkAAixRE4wJR1bMQQZ4GYisG5qezhp4kWP4FE94kpisW2REKSJlsliJEUygPXS8jhEF26imUxpoUgYxtuyNl5qam_x5xC-ny4F6uRAaJ9dGXGm4qcyCAaPHmYxZ7R6a82yikFtA4pbypAA5Uwz4SSie8RahBpyReGiYx9ZaC5SFP4t5gx8JA6Lt4Gy9FBh4AayoWgz77DADzK8oAx2d8zWGxp2XBzaRFu4O6ADhlBomiwGAghhKuui6E998RaCylz8OiFu64Sl3qhC1k8Sih2PXxuUWagFJiAQtyrmmHogAhpkdkBQloU9EGmq-F8hBgzy9AU4vg-68xpEaUh41u93gOm78ypouyCfgChiyAmjU4Z0EyU9EpoGdxK12oe8WiU8potwoK5o2CwNws88WxS1d20rwCU-4j0VhU692QxwiyO161IosGq8xScwFh8ra78lwtu2ooou9wzyqxmtwAwkQ3qax6um14a5A4A1-xKbwaC2O4E429w-xoUbE2mUC586O15w6hw69xe1Uwo81682tws40jC0kEE2VwsEqwfm1Hw4iK0DE6a0Cojg2nwdm0worw8C1Cw66wa-0sy1Uwe208xwKwDw960SU2nwto1_o1X84K0SE12U3ow")
	params.Add("__hblp", "08e8h8Wl0nU9UioqwuEnwwwh9Uco12o8qw8W0MEhga86Km2-E4eqtK68S0DFUkgOUrwTy8jxuawxCAyXUixTwKggwu89ouzU3cyEnDDxa8Dw8O5o6d1OAXxu7olxa4EcEjx-2i3Gp0nEc8eF8nyup4CzUizUcqDwxwxwfl4g4i1Exyq262-32q1QxO16wrU941YxK18BAw-woE4-2u5E4O9xK5Uaovkwsw8S1owr41kwSx-484-0Dorw9q1bwnofEmwh88o72582MUcE5u1lK9wjE3fw863C1bBxa5onwTgiwCwVwjE2Jxa8wo8f8uw47wRwzBwAw4AwyyQ1ZwNK1OxG1Pwn82gUnw9B0PwFK2a484aiUuwoFVem4EeEdEy9ocE5yagkBws84C6UiG1fxK0yotwywBwhFo7Kq3-1ewDw8S261hwaS2GER2oS13wLxS3im7V82rzox0DwwzE7i17wg42qu7U2gxa1mCwIyUgUK3a5oCibAwqoaohx6cwPwh8sx50gBU982UweibwBGQ2lxnghxi8Ua88UdF84u1cgrwCDUO9zEb9u0K46Wxi1pwDwCwKwtm1wDyuaByU9UhUvKawzypo4eU2awc6dDjHBgG1CxK22A78")
	params.Add("__sjsp", "g49494FEOcxaF44EfaG56EjPhaaxaiEkxGraxO8wMG99j7COEb-koKaal4KgylqsJiZ5XWHJ4tFN0xIJEAiKGsPFkn8yFaGTlaQWanhFqeEzzaTxfmKFC9KlxwNt91zgK21GduE9TiwAK54bwzgWy3j2axhdnh8QOzKashXJdi3jfxSbyF9o-5OB8nnEIkFmyyvx91C-mahURpmleSHCBx6z9aqgF3bgG4KWKm7HFDSi8AK9Cig88goGe8Qy4pksKyVkL8ivhm5Sac8hk8i8Q5C54Cl4ig8V25hQ8wKoAx9ABBBGeKmmfQhu7p9QmFm5AEaF44k8DDAw_AzkGq9mcxHx-1TroB2PBxKeyAah5ap24dmcolhpkdkBGl134h2aA8LGi4pEG484R0Ngy2K4l828d1a2mpouyQ48C5Aaw8C2q6mazo3jwzwoK0l-2q7o4Qs14U-0LoJ8o7y0vQ7U-1qwkE2u667wm60Po0gew3to7y0a8g0EgE")
	params.Add("__comet_req", "15")
	params.Add("fb_dtsg", "NAfv3z_uR5fDzXHQMYBsqBpHfDlhD6pXi92Av047coSCs_YuhnM3aRA:28:1753880016")
	params.Add("jazoest", "25483")
	params.Add("lsd", "lZAgVGi3u3_i95VHzIK5bj")
	params.Add("__spin_r", "1025305868")
	params.Add("__spin_b", "trunk")
	params.Add("__spin_t", "1753880119")
	params.Add("__crn", "comet.fbweb.CometGroupDiscussionRoute")

	baseURL.RawQuery = params.Encode()

	apiURL = baseURL.String()

	// 9. Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set headers from the working curl request
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Origin", "https://web.facebook.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://web.facebook.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("TE", "trailers")

	// Set the Cookies from the working curl request
	req.Header.Set("Cookie", "datr=nRWKaHtB5Qe7tY-SQljVzBUB; fr=0MyCxY2RnfzlTJwqw.AWcVTVz6UOctaxfcePE4msVBzaDfCLFKynKRr61TmZNxjLsWYfU.BoihWd..AAA.0.0.BoihXV.AWez3o_cmGNDOFvy8KG9HUdonEQ; ps_l=1; ps_n=1; sb=nRWKaCq7zly0YLI1NH88jXwV; wd=1366x683; c_user=61560452168137; xs=28%3A4gjwyfpBuuxjXg%3A2%3A1753880016%3A-1%3A-1; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1753880129801%2C%22v%22%3A1%7D")

	// 10. Make the API call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making API call: %w", err)
	}
	defer resp.Body.Close()

	// 11. Read the response body
	var respBody []byte
	var readErr error

	// Debug: Print response headers
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Content-Encoding: %s\n", resp.Header.Get("Content-Encoding"))
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("Response body: %s\n", resp.Body)

	// Try to handle different compression types
	contentEncoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		gzReader, gzErr := gzip.NewReader(resp.Body)
		if gzErr != nil {
			return "", fmt.Errorf("error creating gzip reader: %w", gzErr)
		}
		defer gzReader.Close()
		respBody, readErr = io.ReadAll(gzReader)
	} else if strings.Contains(contentEncoding, "br") {
		// Handle brotli compression
		respBody, readErr = io.ReadAll(resp.Body)
	} else if strings.Contains(contentEncoding, "deflate") {
		// Handle deflate compression
		respBody, readErr = io.ReadAll(resp.Body)
	} else if strings.Contains(contentEncoding, "zstd") {
		// Handle zstd compression
		reader, err := zstd.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error creating zstd reader: %w", err)
		}
		defer reader.Close()
		respBody, readErr = io.ReadAll(reader)
	} else {
		// No compression or unknown compression
		respBody, readErr = io.ReadAll(resp.Body)
	}

	if readErr != nil {
		return "", fmt.Errorf("error reading response body: %w", readErr)
	}

	// Debug: Print the raw response
	fmt.Printf("Raw response length: %d bytes\n", len(respBody))
	if len(respBody) > 0 {
		fmt.Printf("Raw response (first 500 chars): %s\n", string(respBody[:min(500, len(respBody))]))
	}

	// 12. Parse the JSON response
	stringRespBody := string(respBody)
	stringRespBody = strings.ReplaceAll(stringRespBody, "for (;;);", "")

	// Clean up any leading/trailing whitespace
	stringRespBody = strings.TrimSpace(stringRespBody)

	var uploadResponse UploadResponse
	err = json.Unmarshal([]byte(stringRespBody), &uploadResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w, response: %s", err, stringRespBody)
	}

	// 13. Check for errors in the response and validate real upload
	if uploadResponse.Payload.PhotoID == "" {
		return "", fmt.Errorf("upload failed, response: %s", string(respBody))
	}

	// Validate that we got real image URLs (not fake IDs)
	if uploadResponse.Payload.ImageSrc == "" || uploadResponse.Payload.ThumbSrc == "" {
		return "", fmt.Errorf("upload succeeded but got fake IDs - no image URLs in response: %s", string(respBody))
	}

	fmt.Printf("Uploaded image with ID: %s\n", uploadResponse.Payload.PhotoID)
	fmt.Printf("Image URL: %s\n", uploadResponse.Payload.ImageSrc)
	fmt.Printf("Thumbnail URL: %s\n", uploadResponse.Payload.ThumbSrc)
	fmt.Printf("Image dimensions: %dx%d\n", uploadResponse.Payload.Width, uploadResponse.Payload.Height)

	// 14. Return the image ID
	return uploadResponse.Payload.PhotoID, nil
}
