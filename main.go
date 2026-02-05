package main

import (
	"autoposter/utils"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func loginToFacebook() (*rod.Browser, *rod.Page) {
	dir := "./chrome"

	// u := launcher.New().UserDataDir(dir).Leakless(true).NoSandbox(true).Headless(true).MustLaunch()

	// browser := rod.New().ControlURL(u).MustConnect().NoDefaultDevice()

	// Create a more transparent launcher configuration
	l := launcher.NewUserMode().
		Leakless(false).  // Disable leakless mode which can trigger security
		NoSandbox(false). // Enable sandbox for better security
		Headless(false).  // Keep visible for transparency
		Devtools(false).  // Disable devtools for production
		UserDataDir(dir). // Use a local directory for user data
		MustLaunch()

	browser := rod.New().
		ControlURL(l).
		MustConnect().
		NoDefaultDevice()

	page := browser.MustPage("https://web.facebook.com/").MustWindowMaximize().MustWaitLoad()

	for {
		if page.MustInfo().Title == "Facebook – log in or sign up" || page.MustHas(`form[data-testid="royal_login_form"]`) || strings.Contains(page.MustHTML(), "Log in to Facebook") || strings.Contains(page.MustInfo().URL, "two_step_verification/authentication/") {
			fmt.Println("Please login to Facebook")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Login successful")
			break
		}
	}

	return browser, page
}

func main() {
	// Extract items once before the loop to ensure even distribution
	items, err := utils.ExtractItems()
	if err != nil {
		fmt.Printf("❌ Error extracting items: %v\n", err)
		return
	}

	if len(items) == 0 {
		fmt.Printf("❌ No items found\n")
		return
	}

	fmt.Printf("📦 Found %d items to post\n", len(items))

	// Shuffle items once at the start for variety
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

	r.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	browser, page := loginToFacebook()

	fetchGroupsCurl := utils.GetFetchGroupsCurl(page)

	uploadImageCurl, createPostCurl := utils.GetUploadImageAndCreatePostCurls(page)

	// Close browser and page immediately after getting curls
	page.MustClose()
	browser.MustClose()

	if fetchGroupsCurl == "" || uploadImageCurl == "" || createPostCurl == "" {
		fmt.Printf("❌ Failed to capture necessary curl commands\n")
		return
	}

	err = utils.UpdateFetchGroupsConfigFromCurl(fetchGroupsCurl)
	if err != nil {
		fmt.Printf("❌ Error updating config: %v\n", err)
		return
	}

	groups, err := utils.ExtractAllGroups()
	if err != nil {
		fmt.Printf("❌ Error extracting groups: %v\n", err)
		return
	}

	utils.PrintGroupResults(groups)

	jsonData, err := utils.SaveGroupsAsJSON(groups)
	if err != nil {
		fmt.Printf("❌ Error preparing JSON: %v\n", err)
		return
	}

	fmt.Printf("\n💾 JSON data ready (length: %d bytes)\n", len(jsonData))
	fmt.Println("🎉 All done!")
	fmt.Printf("📊 Total groups extracted: %d\n", len(groups))

	successfulPosts := 0
	totalGroups := len(groups)
	targetPosts := 250
	if totalGroups < targetPosts {
		targetPosts = totalGroups
	}
	itemIndex := 0 // Track current item for round-robin distribution

	r.Shuffle(len(groups), func(i, j int) {
		groups[i], groups[j] = groups[j], groups[i]
	})

	for i, group := range groups {
		fmt.Printf("\n📋 Processing group %d/%d: %s (%s)\n", i+1, totalGroups, group.Name, group.ID)

		// Get the current item using round-robin distribution
		item := items[itemIndex]
		itemIndex = (itemIndex + 1) % len(items) // Move to next item, wrap around if needed

		fmt.Printf("📄 Using item %d/%d from rotation\n", ((itemIndex-1)%len(items))+1, len(items))

		err = utils.UpdateUploadConfigFromCurl(uploadImageCurl)
		if err != nil {
			fmt.Printf("❌ Error updating upload config: %v\n", err)
			return
		}

		photoIDs := []string{}
		for _, imagePath := range item.ImagePaths {
			response, err := utils.UploadImage(imagePath, uploadImageCurl)
			if err != nil {
				fmt.Printf("❌ Error uploading image: %v\n", err)
				continue
			}

			if response.Success {
				fmt.Printf("🎉 Upload successful for: %s\n", response.ImagePath)
				fmt.Printf("   📷 Photo ID: %s\n", response.PhotoID)
				fmt.Printf("   🔗 Image URL: %s\n", response.ImageSrc)
				fmt.Printf("   🖼️  Thumbnail: %s\n", response.ThumbSrc)
				fmt.Printf("   📏 Size: %dx%d\n", response.Width, response.Height)
				fmt.Println()
				photoIDs = append(photoIDs, response.PhotoID)
			} else {
				fmt.Printf("❌ Upload failed for: %s - %s\n", response.ImagePath, response.Error)
			}
		}

		err = utils.UpdatePostConfigFromCurl(createPostCurl)
		if err != nil {
			fmt.Printf("❌ Error updating post config: %v\n", err)
			continue
		}

		postResponse, err := utils.CreatePost(item.Description, photoIDs, group.ID, createPostCurl)
		if err != nil {
			fmt.Printf("❌ Error creating post: %v\n", err)
			continue
		}

		if postResponse.Success {
			successfulPosts++
			fmt.Printf("🎉 Post created successfully!\n")
			fmt.Printf("   📝 Post ID: %s\n", postResponse.PostID)
			fmt.Printf("   👥 Group ID: %s\n", postResponse.GroupID)
			fmt.Printf("   📄 Text: %s\n", postResponse.Text)
			fmt.Printf("   🖼️  Photos: %d images\n", len(postResponse.PhotoIDs))
			fmt.Printf("   📊 Progress: %d/%d posts completed (%.1f%%)\n", successfulPosts, targetPosts, float64(successfulPosts)/float64(targetPosts)*100)
		} else {
			fmt.Printf("❌ Post creation failed: %v\n", postResponse)
		}

		if successfulPosts == targetPosts {
			fmt.Printf("\n🚀 Reached %d successful posts! Stopping further posts.\n", targetPosts)
			break
		}

		fmt.Println()
		time.Sleep(2 * time.Minute)
	}

	// Final summary
	fmt.Printf("\n🎯 Final Summary:\n")
	fmt.Printf("   ✅ Successful posts: %d/%d\n", successfulPosts, targetPosts)
	fmt.Printf("   📈 Success rate: %.1f%%\n", float64(successfulPosts)/float64(targetPosts)*100)
	fmt.Printf("   🏁 Auto-posting completed!\n")
}
