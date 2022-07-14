package main

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/protocol/webauthncose"
	"github.com/duo-labs/webauthn/webauthn"
)

type CreateOptions struct {
	PublicKey PublicKeyCredentialCreationOptions `json:"publicKey"`
}

type PublicKeyCredentialCreationOptions struct {
	Challenge                     protocol.URLEncodedBase64       `json:"challenge"`
	RelyingParty                  PublicKeyCredentialRpEntity     `json:"rp"`
	User                          PublicKeyCredentialUserEntity   `json:"user"`
	PublicKeyCredentialParameters []PublicKeyCredentialParameter  `json:"pubKeyCredParams,omitempty"`
	Timeout                       int                             `json:"timeout"`
	ExcludeCredentials            []PublicKeyCredentialDescriptor `json:"excludeCredentials,omitempty"`
	AuthenticatorSelection        protocol.AuthenticatorSelection `json:"authenticatorSelection"`
	Attestation                   protocol.ConveyancePreference   `json:"attestation"`
	Extensions                    map[string]interface{}          `json:"extensions,omitempty"`
}

type PublicKeyCredentialRpEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PublicKeyCredentialUserEntity struct {
	ID          protocol.URLEncodedBase64 `json:"id"`
	Name        string                    `json:"name"`
	DisplayName string                    `json:"displayName"`
}

type PublicKeyCredentialParameter struct {
	Type      protocol.CredentialType              `json:"type"`
	Algorithm webauthncose.COSEAlgorithmIdentifier `json:"alg"`
}

type PublicKeyCredentialDescriptor struct {
	Type       protocol.CredentialType   `json:"type"`
	ID         protocol.URLEncodedBase64 `json:"id"`
	Transports []string                  `json:"transports,omitempty"`
}

func MakeCreateOptions(handle *webauthn.WebAuthn, user *User) (*CreateOptions, error) {
	challenge, err := protocol.CreateChallenge()
	if err != nil {
		return nil, err
	}

	var exclude []PublicKeyCredentialDescriptor
	if credential, err := user.WebAuthnCredential(); err == nil {
		exclude = append(exclude, PublicKeyCredentialDescriptor{
			Type: protocol.PublicKeyCredentialType,
			ID:   protocol.URLEncodedBase64(credential.ID),
		})
	}

	return &CreateOptions{
		PublicKey: PublicKeyCredentialCreationOptions{
			Challenge: challenge,
			RelyingParty: PublicKeyCredentialRpEntity{
				ID:   handle.Config.RPID,
				Name: handle.Config.RPDisplayName,
			},
			User: PublicKeyCredentialUserEntity{
				ID:          []byte(user.ID),
				Name:        user.Name,
				DisplayName: user.Name,
			},
			// https://www.w3.org/TR/webauthn-2/#CreateCred-DetermineRpId
			// The default in the spec is ES256 and RS256.
			PublicKeyCredentialParameters: []PublicKeyCredentialParameter{
				{
					Type:      protocol.PublicKeyCredentialType,
					Algorithm: webauthncose.AlgES256,
				},
				{
					Type:      protocol.PublicKeyCredentialType,
					Algorithm: webauthncose.AlgRS256,
				},
			},
			Extensions: map[string]interface{}{
				// We want to know user verification method (uvm).
				// https://www.w3.org/TR/webauthn-2/#sctn-uvm-extension
				"uvm": true,
				// We want to know the credentials is client-side discoverable or not.
				// https://www.w3.org/TR/webauthn-2/#sctn-authenticator-credential-properties-extension
				"credProps": true,
			},
			AuthenticatorSelection: handle.Config.AuthenticatorSelection,
			Timeout:                handle.Config.Timeout,
			Attestation:            handle.Config.AttestationPreference,
			ExcludeCredentials:     exclude,
		},
	}, nil
}
