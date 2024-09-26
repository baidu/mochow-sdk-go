// string.go - define the string util function

// Package util defines the common utilities including string and time.
package util

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func URIEncode(uri string, encodeSlash bool) string {
	var byteBuf bytes.Buffer
	for _, b := range []byte(uri) {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' || b == '~' || (b == '/' && !encodeSlash) {
			byteBuf.WriteByte(b)
		} else {
			byteBuf.WriteString(fmt.Sprintf("%%%02X", b))
		}
	}
	return byteBuf.String()
}

func NewUUID() string {
	var buf [16]byte
	for {
		if _, err := rand.Read(buf[:]); err == nil {
			break
		}
	}
	buf[6] = (buf[6] & 0x0f) | (4 << 4)
	buf[8] = (buf[8] & 0xbf) | 0x80

	res := make([]byte, 36)
	hex.Encode(res[0:8], buf[0:4])
	res[8] = '-'
	hex.Encode(res[9:13], buf[4:6])
	res[13] = '-'
	hex.Encode(res[14:18], buf[6:8])
	res[18] = '-'
	hex.Encode(res[19:23], buf[8:10])
	res[23] = '-'
	hex.Encode(res[24:], buf[10:])
	return string(res)
}

func NewRequestID() string {
	return NewUUID()
}
