package main

import (
	"github.com/duo-labs/webauthn/protocol"
)

type GetOptions struct {
	PublicKey PublicKeyCredentialGetOptions `json:"publicKey"`
	Mediation string                        `json:"mediation,omitempty"`
}

type PublicKeyCredentialGetOptions struct {
	Challenge        protocol.URLEncodedBase64            `json:"challenge"`
	Timeout          int                                  `json:"timeout"`
	RPID             string                               `json:"rpId"`
	UserVerification protocol.UserVerificationRequirement `json:"userVerification"`
	AllowCredentials []PublicKeyCredentialDescriptor      `json:"allowCredentials,omitempty"`
	Extensions       map[string]interface{}               `json:"extensions,omitempty"`
}

func MakeGetOptionsModal(config *WebAuthnConfig, credentialID string) (*GetOptions, error) {
	challenge, err := protocol.CreateChallenge()
	if err != nil {
		return nil, err
	}

	var allowCredentials []PublicKeyCredentialDescriptor
	if credentialID != "" {
		allowCredentials = append(allowCredentials, PublicKeyCredentialDescriptor{
			Type: protocol.PublicKeyCredentialType,
			ID:   Base64URLDecode(credentialID),
		})
	}

	return &GetOptions{
		PublicKey: PublicKeyCredentialGetOptions{
			Challenge:        challenge,
			Timeout:          config.MediationModalTimeout,
			RPID:             config.RPID,
			UserVerification: config.AuthenticatorSelection.UserVerification,
			AllowCredentials: allowCredentials,
			Extensions: map[string]interface{}{
				// We want to know user verification method (uvm).
				// https://www.w3.org/TR/webauthn-2/#sctn-uvm-extension
				"uvm": true,
			},
		},
	}, nil
}

func MakeGetOptionsConditional(config *WebAuthnConfig) (*GetOptions, error) {
	challenge, err := protocol.CreateChallenge()
	if err != nil {
		return nil, err
	}

	return &GetOptions{
		Mediation: "conditional",
		PublicKey: PublicKeyCredentialGetOptions{
			Challenge:        challenge,
			Timeout:          config.MediationConditionalTimeout,
			RPID:             config.RPID,
			UserVerification: config.AuthenticatorSelection.UserVerification,
			Extensions: map[string]interface{}{
				// We want to know user verification method (uvm).
				// https://www.w3.org/TR/webauthn-2/#sctn-uvm-extension
				"uvm": true,
			},
		},
	}, nil
}
