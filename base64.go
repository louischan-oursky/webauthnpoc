package main

import (
	"encoding/base64"

	"github.com/duo-labs/webauthn/protocol"
)

func Base64URLEncode(bytes protocol.URLEncodedBase64) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func Base64URLDecode(s string) protocol.URLEncodedBase64 {
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bytes
}
