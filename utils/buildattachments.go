package utils

import "fmt"

// Helper function to build attachments JSON
func buildAttachmentsJSON(photoIDs []string) string {
	if len(photoIDs) == 0 {
		return "[]"
	}

	attachments := "["
	for i, photoID := range photoIDs {
		if i > 0 {
			attachments += ","
		}
		attachments += fmt.Sprintf(`{"photo":{"id":"%s"}}`, photoID)
	}
	attachments += "]"

	return attachments
}
