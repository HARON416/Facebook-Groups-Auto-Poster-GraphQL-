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
	err = writer.WriteField("profile_id", "61553016608952") // Will be updated by UpdateUploadPhotoCookies()
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

	// 6. Add the upload_id field - use fixed value like curl request
	err = writer.WriteField("upload_id", "jsc_c_4") // Use fixed upload_id like curl request
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

	// Add query parameters - will be updated by UpdateUploadPhotoCookies()
	params := url.Values{}
	params.Add("av", "61553016608952")
	params.Add("__aaid", "0")
	params.Add("__user", "61553016608952")
	params.Add("__a", "1")
	params.Add("__req", "4f")
	params.Add("__hs", "20301.HYP:comet_pkg.2.1...0")
	params.Add("dpr", "1")
	params.Add("__ccg", "EXCELLENT")
	params.Add("__rev", "1025398783")
	params.Add("__s", "f4fr92:ybawqk:hg1d1b")
	params.Add("__hsi", "7533483886033570164")
	params.Add("__dyn", "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUKbgS3q2ibwNwnof8boG0x8bo6u3y4o2Gwfi0LVEtwMw6ywIK1Rwwwg8a8462mcwfG12wOx62G5UswoEcE7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzE8FU5e3ym2SU4i5oe8464-5pU9UmwUwxwjFovUaU3qxW2-awLyESE7i3C223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwUU-4rwEKufxamEbbxG1fBG2-2K0E846fwk83KwHwOG8xG6E")
	params.Add("__csr", "gfYl1n5gjhYAYp34AAOkYA54p2InkAjOjR2c8P8mBk98CAL6ZPMAJ8hi8yvkP4lsGOT4QWqFbYCBXi9kChR9WOnQhaJnh_FSgyiQmJbHFGJv4ehGxihvlqqSGBZkBqmqGhaGV7Z4p9AAAQZlyR8ySn4XQllmIx4y9lABF6Hm-8CO5DmLq98ASVi4A9hmAvugPGBnF4qm8AApqEBmmLyVehWBV4bgCWjmGiDKKCVaiRGZCCGWQirWAGmimQuQp6ykJaCGWCZaCBjCAALCy4uQFppdemFJ5GGJ9LgFqx2fGkziAG9mmudzbjAp98HDJrBh-rXgCV4FoSqGAZ6CBAgJBBBQEF4CGKdzHVpQWAFBKqFGhFoyq9yaLhbGrwIyUgypu5UC9BCjCV9FVp4icjUz-Feim5o-9xKcA-qmmi5FUgDmaHGV9EybDCzAUGrzEGUgyECF8gUSdGbyokz9oiUq-ex6ey8Shqx-cDxhuEiBBxCaUSGDx3yUObhUC2-2a6UJoGaFeit2Xx2l39V-2am6UoxWex6bwVGi48qBjxaq7F-2S2GEtxe485vwEVqG16wa60Gorz89Eng2qwjVUnwh8WU9Fo943q1sAy8d8WmU6eqi8xm2a7e5E88V2r-1gw9uax2221VwXwjFE-21WBwCwXwFwsU20wwU4Cim7swnzeucyofo22wZxCbCwGg2Pwt20GQ2au8wzyE9poaA7Upwwz80s9w3So6S06xEKUtgc43fye06Rp87t0gU0llwo8564UoQU0wVk1jQ686WbQ18w5Jwbu01QpwMwnUG054o04DS3G0lC0l67842t0gesUoxFUiG11wuShwjUy3a0Y82owfIOh44p207Gw8qaW5U5h0si0Fwai0a5wgo2uwqS365E8810GwbG9w8C2a050E0JO680Jq640dcwkEC8m8w15B07qo0IK1HwJwaQUb4U2xwf10Fw2sEy9g0vkw7Lg5m14yEog4e9x10Qg2xwie1tyd2Qbh81YO019q0gV7xa1rgnxm320Qrw8y0gS0N8-0s60A80rFpKzxq67blvg4Ew3pyECFC2oK2bJo4e16o6t1acCwRpE1f8rwLwkAE9E")
	params.Add("__hsdp", "g495A2i8ggy94p712ACzAr0wGxp0iHhhB9eg5ky4OHC5Hp3EH2cpcgmx5gwxqkIIIx6B3O7lNbq6Gwo8W5E4ZFEh585Q8hNc4J2NceY9B68SAIaAtI9cH7Eoxs9mx54Eleh9h4Nmz4m98ndhaq7FAkn9iFB68yrZ4MHcpqmwCyGXF1ayBhkzFCh6BhRKqRdSz8zJLGXcIZRileCGsp8QqtjAB4n5agFkRyijORvlGlb8y2xHlyQF4Je989hYgOsAPaxmkRjY8kNpx1TDax9kO8Sl98NQIYFEVA6Ezct2zey2xuHiiqiOJHpRoS4aB8h4LoB4VK58DQiiEJeABQlde4Pe48xkcxd7KFaCUF4LgFx6FFUO2t3GyEJcF8F0OBGaAF2-AAoxaF7BhfBgDULkPha4b4PCK8Dch7uUy8iia4kyF6iFa9mpEF5QXhR0xEQsiyyoS9cwKox6X3FVkm9Fd1S4XmGm9Oz4Au594i9zQ494dgSFkaogzu9xgUhzOyURDcmA8xe7ElwBzhqGifxqul3UighF11a4Q0C43218wNwo89osGCdzjai48Z0BgyUixyubwExFxRwPz44Uqy8vz8762ebe940SAjogwDz8UBcgNF46EXz8mgC22q48gwXrc4ywiofkpd0Ga25wLx23fwVwYB8AtxO361OyA78hxsU4ui6oO2q0xU99o7B1a0h388wam1Ewww4owxwCwJBDw25o1nU1JEbqwiE2zAwQwiU6e19wUxmufwywywGxGfwjqwat0zwIwRwbq0UXwfi1DGE3fwno9U3Uwam0wU423C223e0OE6C1HwcC0_EcEy3G0lqcwiUK0TU3fwfW0P86C0_E4a0iHypoG7pE3kw")
	params.Add("__hblp", "0oAbEkAp2ogV8S1nDxutxK3y2W3268O2CdAxW1gBwgofo4K2q1AxG0Ko88sVo5iu5A2CucxV0FHUK3m6VXCxW1mxa6EaEy3G8xx3oiVFlUgGbyUozFEhy8G3abg-3u26dQ1cwMG4o551SUjCG3q2zw8y1uG2mEmzGgK6peE-7UoBwyJ0IGEqBBz8C4Q33wEwYG5GDl0AwVyKcwQxi3m2e4EoxO2i6EK2cC7EeU4Cq4o4i0wU5CuewGw-Bz84O291auqEuxm-dK2O3ubwKwQBxu58HyU9E9U7SVEcocUigry468S3i6UW2qU4-2u221nwDwAwqUsABwJw8icDojxmE8p87G4UcE4Hxei5EgCgC1lyEdUpy8iwIwFxG1Lx6dK2S19yE523e2y4obU-5EnDgrwk986u2Gm8AwxCzGAy8Kbxm6EdU7u9geUbEpxm1awpUfE8Ejwywg8K1VXwQy82LwBwGxi3iVEszovG3y4olwCwFwYAwCK9wLxu-4ojCKuuUsGjDwGz8yaxifCwy-58y2uawWzUdonG686J1mbixLG48rwPVUmghxeEfUbk1jm4ku12ghDKfG2am2e13x2cgeoG1gxqGCy8lxG17wVz89UcZ1eQdwXxi8CK8xidJe36m2qcxG3u15wgHwxDxa4Ud8sxi58KeCz8gz88U7zCAxa2C3-m3l2rzUgx29zF8hxe3q5o9U-bwRw_wYAgS5_x654ay8ngW4E1dbgjCzEbUqh9K4UmwzGmi6oc9ob-2im9K26HABBwkoyawhUhwvk6A5u16wUz8pQ9y8S8xqyAxy5m3yErCwFyrxe448GdABCGexOm48ym5RK6Vo24AigbU5-q9UG8z82tyoiG")
	params.Add("__sjsp", "g495A2i8ggy94p712ACzAr0wGxp0iHhhB9ehkgQbR8IYGVN6JAeyIG5clO4fq4q225FiEDsxkBq2cB99LlqcblboUm4t4ArEmQChoOD9BjxGcZ3pqxG9wNOyK4KJz8x128ziCh9AbBp7hAhp4p78dyryx97AiIxi6Wcic8sOAFEuChhsB5VxeCSt62INBKhgwBiKp1ayBhkzFz6BhRKqRdSJ48XrAKbLtWil136i4CDzoWaqB8qRyijRiSRoCJAoO8jlyQF4Je98aDOQOcyaOElBdkXsmANpx1Ty28gqlypria8iUF6zAukGhmazteh2OiAJ97gjpQAdwxhbQaherxidh9UJeWjhkQUjlK6J38jhXGurxjgFxiqucwDgWE984iEGiAbWgExWhVEWa-43d4F8yNf8bx0N4tUpIw9kHz8CgFrjx215hA2Oaixt6G4gQul5yqgeGGm7bhUkAh8pggAgsGl2Q7K9xgU8mdwyz8jwQwBzhqG58942916A440iK0AE4aGoS6E6y1Moto5m13z8762YMEg3Je1m9j4cqm1fgC0xmP18E4C3R6jg10U1UE1kQ0if809hwxw0hvE09JE")
	params.Add("__comet_req", "15")
	params.Add("fb_dtsg", "NAfu4ikzYjWt_U5CvSMrLJavDntdTRnEYkzfEaunOgxMEeIY3Ua5AVA:30:1737396043")
	params.Add("jazoest", "25706")
	params.Add("lsd", "Ojhd_SK0sSq0wM2RpBbKAc")
	params.Add("__spin_r", "1025398783")
	params.Add("__spin_b", "trunk")
	params.Add("__spin_t", "1754025902")
	params.Add("__crn", "comet.fbweb.CometGroupDiscussionRoute")

	baseURL.RawQuery = params.Encode()

	apiURL = baseURL.String()

	// 9. Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set headers - will be updated by UpdateUploadPhotoCookies()
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://web.facebook.com")
	req.Header.Set("Referer", "https://web.facebook.com/")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

	// Set the Cookies - will be updated by UpdateUploadPhotoCookies()
	req.Header.Set("Cookie", "ps_l=1; ps_n=1; sb=EyQJZ_p60kcMqjhfiPbKiQun; datr=OY-OZwqAfui6fV8_zZnYl_zQ; c_user=61553016608952; wd=1366x681; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1754025907720%2C%22v%22%3A1%7D; fr=1HwRYFMccYQ60Ti6R.AWd6h7zAp-sNZqN_F07CLmWjgvAaVuXWRlNHmjKeQjIZOgLS0As.BojE-z..AAA.0.0.BojE-z.AWcBkjQ1eLuFkKOsxLH8zBzuVKo; xs=30%3AdG3sK9-gG0z1jg%3A2%3A1737396043%3A-1%3A-1%3A%3AAcUaHk4sWBPmVZFPMl8DUMhRFWcpUnwnlwumSKhzas8")

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
