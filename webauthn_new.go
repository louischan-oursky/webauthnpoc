package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

func NewWebAuthn(rpDisplayName string, r *http.Request) (*webauthn.WebAuthn, error) {
	origin := url.URL{
		Scheme: GetProto(r),
		Host:   GetHost(r),
	}

	requireResidentKey := true
	config := &webauthn.Config{
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
		// Allow 5 minutes for the user to finish the registration process.
		Timeout: int((5 * time.Minute).Milliseconds()),
	}

	return webauthn.New(config)
}
