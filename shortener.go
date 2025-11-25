package main

import (
	"fmt"
	"hash/adler32"
)

func GenerateCode(url string) string {
	checksum := adler32.Checksum([]byte(url))
	return fmt.Sprintf("%08x", checksum)
}

func GenerateShortUrl(code string) string {
	return fmt.Sprintf("%s/%s", Conf.BaseUrl, code)
}
