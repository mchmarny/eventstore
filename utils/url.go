package utils

import (
	"fmt"
	"log"
	"strings"
)

const (
	imageNameSufix = "/photo.jpg"
)

// ServerSizeResizePlusPic resizes Plus image on server size
func ServerSizeResizePlusPic(picURL string, size int) string {

	log.Printf("URL: %s", picURL)
	if !strings.HasSuffix(picURL, imageNameSufix) {
		log.Printf("Not valid profile picture format")
		return picURL
	}

	sizedImageName := fmt.Sprintf("/s%d%s", size, imageNameSufix)
	log.Printf("Sized image name: %s", sizedImageName)

	return strings.Replace(picURL, imageNameSufix, sizedImageName, -1)

}
