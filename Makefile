.PHONY: mkcert
mkcert:
	mkcert -cert-file tls-cert.pem -key-file tls-key.pem "::1" "127.0.0.1" localhost webauthn.com
