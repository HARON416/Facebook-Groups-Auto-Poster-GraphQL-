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
	Text       string   `json:"text"`
	ImagePaths []string `json:"image_paths"`
}

func ExtractItems(path string) ([]Item, error) {
	var items []Item
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

		var text string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case strings.HasPrefix(line, "description:"):
				text = strings.TrimSpace(line[len("description:"):])
			}
		}

		items = append(items, Item{
			Text:       text,
			ImagePaths: imageFiles,
		})
	}

	return items, nil
}
