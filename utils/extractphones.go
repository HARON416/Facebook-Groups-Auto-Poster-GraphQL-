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

// GroupItem represents an item with images and description
type Phone struct {
	Description string   `json:"description"`
	ImagePaths  []string `json:"image_paths"`
}

// GetItemsForGroups reads items from a directory structure and returns both items and image paths
func ExtractPhones(path string) ([]Phone, error) {
	var phones []Phone
	var allImagePaths []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading root directory: %w", err)
	}

	for _, entry := range entries {
		subDir := filepath.Join(path, entry.Name())
		detailsFile := filepath.Join(subDir, "details.txt")

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

		file, err := os.Open(detailsFile)
		if err != nil {
			continue
		}
		defer file.Close()

		var description string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case strings.HasPrefix(line, "description:"):
				description = strings.ToUpper(strings.TrimSpace(line[len("description:"):]))
			}
		}

		parts := strings.Split(description, "...")

		for i := range parts {
			parts[i] = "✅ " + strings.TrimSpace(parts[i])
		}

		description = strings.Join(parts, "\n\n")

		phones = append(phones, Phone{
			Description: description,
			ImagePaths:  imageFiles,
		})
	}

	return phones, nil
}
