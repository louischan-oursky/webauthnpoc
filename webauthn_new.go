package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/duo-labs/webauthn/protocol"
)

type WebAuthnConfig struct {
	RPID                        string
	RPOrigin                    string
	RPDisplayName               string
	AttestationPreference       protocol.ConveyancePreference
	AuthenticatorSelection      protocol.AuthenticatorSelection
	MediationModalTimeout       int
	MediationConditionalTimeout int
}

func NewWebAuthnConfig(rpDisplayName string, r *http.Request) (*WebAuthnConfig, error) {
	origin := url.URL{
		Scheme: GetProto(r),
		Host:   GetHost(r),
	}

	requireResidentKey := true
	return &WebAuthnConfig{
		RPDisplayName: rpDisplayName,

		// The RPID must be a domain only.
		RPID: origin.Hostname(),
		// Origin must be the actual origin as observed by the browser.
		RPOrigin: origin.String(),

		AttestationPreference: protocol.PreferDirectAttestation,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			// AuthenticatorAttachment is intentionally left blank so that the user
			// can choose "platform" or "cross-platform" attachment.
			// This means the authenticator can either be on-device or off-device.
			// AuthenticatorAttachment:,

			// ResidentKey is "required" means we want the credential to be client-side discoverable.
			// This implies we do not need to identity the user first and find out allowCredentials.
			// The outcome is that we can present a flow that the user signs in by selecting
			// credentials, without the need of first entering their email address.
			ResidentKey: protocol.ResidentKeyRequirementRequired,
			// RequireResidentKey is a deprecated field.
			// https://www.w3.org/TR/webauthn-2/#dom-authenticatorselectioncriteria-requireresidentkey
			// It MUST BE true if ResidentKey is "required".
			RequireResidentKey: &requireResidentKey,

			// https://www.w3.org/TR/webauthn-2/#user-verification
			// Per the WWDC video https://developer.apple.com/videos/play/wwdc2022/10092/ at 19:12
			// UserVerification MUST be kept as preferred for the best user experience
			// regardless of whether biometric is available.
			UserVerification: protocol.VerificationPreferred,
		},

		// For modal, the timeout is 5 minutes which is relatively short.
		MediationModalTimeout: int((5 * time.Minute).Milliseconds()),

		// For conditional, the timeout is 1 hour which is long.
		MediationConditionalTimeout: int((1 * time.Hour).Milliseconds()),
	}, nil
}
