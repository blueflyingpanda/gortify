package main

import "regexp"

var urlRegex = regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(/\S*)?$`)

func isValidURL(urlString string) bool {
	return urlRegex.MatchString(urlString)
}
