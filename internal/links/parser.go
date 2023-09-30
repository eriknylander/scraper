package links

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// Parser describes an interface for parsing links.
type Parser interface {
	ParseLinks(b []byte) ([]string, error)
}

// htmlParser implements the Parser interface.
type htmlParser struct{}

func NewHTMLParser() Parser {
	return htmlParser{}
}

// ParseLinks takes an html.Node and extracts any links from it or it's children.
func (p htmlParser) ParseLinks(b []byte) ([]string, error) {
	root, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("could not parse HTML, %w", err)
	}

	links := make([]string, 0)

	var crawl func(node *html.Node)
	crawl = func(node *html.Node) {
		switch node.Data {
		case "a", "link":
			if s, ok := getAttribute(node, "href"); ok {
				links = append(links, s)
			}
		case "img", "script":
			if s, ok := getAttribute(node, "src"); ok {
				links = append(links, s)
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawl(child)
		}
	}
	crawl(root)

	return links, nil
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val, true
		}
	}

	return "", false
}
