package scraping

import (
	"path"
	"strings"
)

func normalizePath(u string, l string) string {
	dir, _ := path.Split(u)

	return path.Clean(path.Join(dir, l))
}

func isHTMLPage(u string) bool {
	return strings.HasSuffix(u, "html") || strings.HasSuffix(u, "htm")
}

func isExternalLink(u string) bool {
	return strings.HasPrefix(u, "http")
}
