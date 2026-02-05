package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Item struct {
	Description string
	ImagePaths  []string
}

func ExtractItems() ([]Item, error) {
	path := "/home/kibet/Downloads/Phones"
	var items []Item
	var allImagePaths []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading root directory: %w", err)
	}

	for _, entry := range entries {
		subDir := filepath.Join(path, entry.Name())
		descriptionFile := filepath.Join(subDir, "description.txt")

		subEntries, err := os.ReadDir(subDir)
		if err != nil {
			continue
		}

		var imageFiles []string
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() && filepath.Ext(subEntry.Name()) != ".txt" {
				filePath := filepath.Join(subDir, subEntry.Name())
				imageFiles = append(imageFiles, filePath)
				allImagePaths = append(allImagePaths, filePath)
			}
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		r.Shuffle(len(imageFiles), func(i, j int) {
			imageFiles[i], imageFiles[j] = imageFiles[j], imageFiles[i]
		})

		if len(imageFiles) > 6 {
			imageFiles = imageFiles[:6]
		}

		file, err := os.Open(descriptionFile)
		if err != nil {
			continue
		}
		defer file.Close()

		var description string
		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				lines = append(lines, line)
			}
		}
		description = strings.ToUpper(strings.Join(lines, "\n"))

		fmt.Println("=====Description=====")
		fmt.Println(description)

		items = append(items, Item{
			Description: description,
			ImagePaths:  imageFiles,
		})
	}

	return items, nil
}
