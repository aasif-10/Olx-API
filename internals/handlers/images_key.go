package handlers

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	maxImageBytes      = 5 * 1024 * 1024
	maxImagePerListing = 10
	uploadPrefix       = "uploads"
)

var allowdContentTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

func mintUploadKey(userId uuid.UUID, ext string) string {
	return fmt.Sprintf("%s%s%s.%s", uploadPrefix, userId, uuid.New(), ext)
}
