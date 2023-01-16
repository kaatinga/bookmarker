package bookmarker

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Bookmark struct {
	// Name of the bookmark
	Name string
	// URL of the bookmark
	URL string
}

// ParseExportedGoogleBookmarks parses the html document and returns the bookmarks.
func ParseExportedGoogleBookmarks(input string) (map[string][]Bookmark, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(input))

	var currentTag, currentCategory string
	var initNewBookmark = false
	var bookmarks = make(map[string][]Bookmark)
	var bookmark Bookmark
	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				return bookmarks, nil
			}
			// otherwise, there was an error tokenizing,
			// which likely means the HTML was malformed.
			return nil, err
		case tt == html.TextToken:
			t := tokenizer.Token()
			if newLine(t.Data) {
				continue
			}
			switch currentTag {
			case "a":
				if initNewBookmark {
					bookmark.Name = t.Data // get the bookmark name
					bookmarks[currentCategory] = append(bookmarks[currentCategory], bookmark)
					initNewBookmark = false
				}
			case "h3":
				currentCategory = t.Data
			}
		case tt == html.StartTagToken:
			t := tokenizer.Token()
			switch t.Data {
			case "h3":
				if isFolder(&t) {
					continue
				}
				currentTag = t.Data
			case "a":
				ok, url := getHref(&t)
				if ok {
					bookmark.URL = url
					currentTag = t.Data
					initNewBookmark = true
				}
			}
		}
	}
}

// getHref pulls the href attribute from a Token.
func getHref(t *html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// isFolder detects <h3> tag that contains a folder name.
func isFolder(t *html.Token) (ok bool) {
	for _, a := range t.Attr {
		if a.Key == "PERSONAL_TOOLBAR_FOLDER" {
			ok = true
		}
	}
	return
}

// newLine detects new line in the text.
func newLine(input string) bool {
	for _, char := range input {
		if char == 10 || char == 13 {
			return true
		}
	}
	return false
}
