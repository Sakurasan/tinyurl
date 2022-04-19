package tinyurl

import (
	"fmt"
	"tinyurl/pkg/base62"

	"github.com/flamego/flamego"
)

func TinyurlHandler(c flamego.Context) string {
	var originurl string
	if len(c.Request().URL.RawQuery) > 1 {
		originurl = c.Param("url") + "?" + c.Request().URL.RawQuery
	} else {
		originurl = c.Param("url")
	}
	return fmt.Sprintf(
		"TinyUrl , %s to %s",
		originurl,
		base62.TinyUrl(originurl),
	)
}
