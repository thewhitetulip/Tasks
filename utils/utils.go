package utils

import "strings"

func GetRedirectUrl(referer string) string {
	var redirectUrl string
	url := strings.Split(referer, "/")

	if len(url) > 4 {
		redirectUrl = "/" + strings.Join(url[3:], "/")
	} else {
		redirectUrl = "/"
	}
	return redirectUrl
}
