package style

import (
	"regexp"
)

func StripAllTag(text string) string {
	re := regexp.MustCompile(`\[(.*?)\]`)
	return string(re.ReplaceAll([]byte(text), []byte("")))
}


