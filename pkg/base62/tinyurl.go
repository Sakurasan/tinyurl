package base62

import (
	mh "github.com/spaolacci/murmur3"
)

func TinyUrl(in string) string {
	return EncodeBase62(int(mh.Sum32([]byte(in))))

}
