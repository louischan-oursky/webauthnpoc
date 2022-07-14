package main

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

type GetOptions struct {
	PublicKey PublicKeyCredentialGetOptions `json:"publicKey"`
}

type PublicKeyCredentialGetOptions struct {
	Challenge        protocol.URLEncodedBase64            `json:"challenge"`
	Timeout          int                                  `json:"timeout"`
	RPID             string                               `json:"rpId"`
	UserVerification protocol.UserVerificationRequirement `json:"userVerification"`
	AllowCredentials []PublicKeyCredentialDescriptor      `json:"allowCredentials,omitempty"`
	Extensions       map[string]interface{}               `json:"extensions,omitempty"`
}

func MakeGetOptions(handle *webauthn.WebAuthn) (*GetOptions, error) {
	challenge, err := protocol.CreateChallenge()
	if err != nil {
		return nil, err
	}

	return &GetOptions{
		PublicKey: PublicKeyCredentialGetOptions{
			Challenge:        challenge,
			Timeout:          handle.Config.Timeout,
			RPID:             handle.Config.RPID,
			UserVerification: handle.Config.AuthenticatorSelection.UserVerification,
			Extensions: map[string]interface{}{
				// We want to know user verification method (uvm).
				// https://www.w3.org/TR/webauthn-2/#sctn-uvm-extension
				"uvm": true,
			},
		},
	}, nil
}
