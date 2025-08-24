package main

import (
	"facebook-autoposter/utils"
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {
	err := utils.UpdateConfigFromCurl()
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

	itemsPath := "/home/kwandapchumba/Pictures/HOT_PHONES"
	successfulPosts := 0
	totalGroups := len(groups)

	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

	for i, group := range groups {
		fmt.Printf("\n📋 Processing group %d/%d: %s (%s)\n", i+1, totalGroups, group.Name, group.ID)

		items, err := utils.ExtractItems(itemsPath)
		if err != nil {
			fmt.Printf("❌ Error extracting items: %v\n", err)
			return
		}

		r.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})

		item := items[0]

		err = utils.UpdateUploadConfigFromCurl()
		if err != nil {
			fmt.Printf("❌ Error updating upload config: %v\n", err)
			return
		}

		photoIDs := []string{}
		for _, imagePath := range item.ImagePaths {
			response, err := utils.UploadImage(imagePath)
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

		err = utils.UpdatePostConfigFromCurl()
		if err != nil {
			fmt.Printf("❌ Error updating post config: %v\n", err)
			continue
		}

		postResponse, err := utils.CreatePost(item.Text, photoIDs, group.ID)
		if err != nil {
			fmt.Printf("❌ Error creating post: %v\n", err)
			continue
		}

		successfulPosts++
		fmt.Printf("🎉 Post created successfully!\n")
		fmt.Printf("   📝 Post ID: %s\n", postResponse.PostID)
		fmt.Printf("   👥 Group ID: %s\n", postResponse.GroupID)
		fmt.Printf("   📄 Text: %s\n", postResponse.Text)
		fmt.Printf("   🖼️  Photos: %d images\n", len(postResponse.PhotoIDs))
		fmt.Printf("   📊 Progress: %d/%d posts completed (%.1f%%)\n", successfulPosts, totalGroups, float64(successfulPosts)/float64(totalGroups)*100)
		fmt.Println()
		time.Sleep(2 * time.Minute)
	}

	// Final summary
	fmt.Printf("\n🎯 Final Summary:\n")
	fmt.Printf("   ✅ Successful posts: %d/%d\n", successfulPosts, totalGroups)
	fmt.Printf("   📈 Success rate: %.1f%%\n", float64(successfulPosts)/float64(totalGroups)*100)
	fmt.Printf("   🏁 Auto-posting completed!\n")
}
