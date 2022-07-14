package main

import (
	"encoding/base64"

	"github.com/duo-labs/webauthn/protocol"
)

func Base64URLEncode(bytes protocol.URLEncodedBase64) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}
