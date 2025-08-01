package main

import (
	"Facebook-Groups-GraphQL-Auto-Poster/utils"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	phonesPath := "/home/kwandapchumba/Pictures/SIMU"

	err := utils.UpdateFetchGroupsCookies()
	if err != nil {
		log.Fatal(err)
	}

	groups, err := utils.FetchGroups(r)
	if err != nil {
		log.Fatal(err)
	}

	r.Shuffle(len(groups), func(i, j int) {
		groups[i], groups[j] = groups[j], groups[i]
	})

	fmt.Println("Total groups: ", len(groups))

	for _, group := range groups {
		groupName := group.Name
		groupID := group.ID
		groupURL := group.URL

		fmt.Printf("Group ID: %s\nGroup Name: %s\nGroup URL: %s\n", groupID, groupName, groupURL)

		phones, err := utils.ExtractPhones(phonesPath)
		if err != nil {
			log.Fatal(err)
		}

		r.Shuffle(len(phones), func(i, j int) {
			phones[i], phones[j] = phones[j], phones[i]
		})

		phone := phones[0]

		photoIDs := []string{}

		err = utils.UpdateUploadPhotoCookies()
		if err != nil {
			log.Fatal(err)
		}

		for _, imagePath := range phone.ImagePaths {
			imageID, err := utils.UploadImage(imagePath)
			if err != nil {
				log.Fatal(err)
			}
			photoIDs = append(photoIDs, imageID)
		}

		fmt.Printf("Photo IDs: %v\n", photoIDs)

		post := utils.FacebookPost{
			MessageText: phone.Description,
			PhotoIDs:    photoIDs,
			GroupID:     group.ID,
		}

		err = utils.UpdateCreateGroupPostCookies()
		if err != nil {
			log.Fatal(err)
		}

		post, err = utils.CreateGroupPost(post)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v created\n Group Name: %s\n Group URL: %s\n", post, groupName, groupURL)

		os.Exit(0)

		duration := utils.ReturnRandomNumber(r, 2, 4)

		time.Sleep(time.Duration(duration) * time.Minute)

		fmt.Printf("Sleeping for %f minutes\n", duration)
	}
}
