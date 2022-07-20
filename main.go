package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
)

const RPDisplayName = "Authgear"

func main() {
	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/create-options", func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		err = r.ParseForm()
		if err != nil {
			return
		}

		name := r.FormValue("name")

		config, err := NewWebAuthnConfig(RPDisplayName, r)
		if err != nil {
			return
		}

		var createOptions *CreateOptions
		err = WithTx(r.Context(), db, func(tx *sql.Tx) error {
			user, err := FindUserWithName(r.Context(), tx, name)
			if err != nil {
				err = nil
				user = NewUserWithName(name)
			}

			createOptions, err = MakeCreateOptions(config, user)
			if err != nil {
				return err
			}

			session, err := NewSessionWithCreateOptions(createOptions)
			if err != nil {
				return err
			}

			err = SaveSession(r.Context(), tx, session)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return
		}

		w.Header().Set("content-type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(createOptions)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var err error
		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		parsed, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(body))
		if err != nil {
			return
		}

		challenge := parsed.Response.CollectedClientData.Challenge

		err = WithTx(r.Context(), db, func(tx *sql.Tx) error {
			session, err := GetSession(r.Context(), tx, challenge)
			if err != nil {
				return err
			}

			user := NewUser(body, session.CreateOptions, parsed)

			err = InsertUser(r.Context(), tx, user)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return
		}
	})

	http.HandleFunc("/get-options-modal", func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		err = r.ParseForm()
		if err != nil {
			return
		}

		name := r.FormValue("name")

		config, err := NewWebAuthnConfig(RPDisplayName, r)
		if err != nil {
			return
		}

		var getOptions *GetOptions
		err = WithTx(r.Context(), db, func(tx *sql.Tx) error {
			user, err := FindUserWithName(r.Context(), tx, name)
			if err != nil {
				err = nil
			}

			var credentialID string
			if user != nil {
				credentialID = user.CredentialID
			}

			getOptions, err = MakeGetOptionsModal(config, credentialID)
			if err != nil {
				return err
			}

			session, err := NewSessionWithGetOptions(getOptions)
			if err != nil {
				return err
			}

			return SaveSession(r.Context(), tx, session)
		})
		if err != nil {
			return
		}

		w.Header().Set("content-type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(getOptions)
	})

	http.HandleFunc("/get-options-conditional", func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		config, err := NewWebAuthnConfig(RPDisplayName, r)
		if err != nil {
			return
		}

		getOptions, err := MakeGetOptionsConditional(config)
		if err != nil {
			return
		}

		session, err := NewSessionWithGetOptions(getOptions)
		if err != nil {
			return
		}

		err = WithTx(r.Context(), db, func(tx *sql.Tx) error {
			return SaveSession(r.Context(), tx, session)
		})
		if err != nil {
			return
		}

		w.Header().Set("content-type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(getOptions)
	})

	http.HandleFunc("/sign-in", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var err error
		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		config, err := NewWebAuthnConfig(RPDisplayName, r)
		if err != nil {
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		parsed, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(body))
		if err != nil {
			return
		}

		challenge := parsed.Response.CollectedClientData.Challenge
		userHandle := parsed.Response.UserHandle
		credentialID := parsed.ID

		err = WithTx(r.Context(), db, func(tx *sql.Tx) error {
			_, err := GetSession(r.Context(), tx, challenge)
			if err != nil {
				return err
			}

			user, err := FindUserWithCredentialID(r.Context(), tx, string(userHandle), credentialID)
			if err != nil {
				return err
			}

			credential, err := user.WebAuthnCredential()
			if err != nil {
				return err
			}

			err = parsed.Verify(
				challenge,
				config.RPID,
				config.RPOrigin,
				"",    // We do not support FIDO AppID extension
				false, // user verification is preferred so we do not require user verification here.
				credential.PublicKey,
			)
			if err != nil {
				return err
			}

			err = user.IncrementSignCount(int64(parsed.Response.AuthenticatorData.Counter))
			if err != nil {
				return err
			}

			err = UpdateUser(r.Context(), tx, user)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return
		}
	})

	http.ListenAndServeTLS(":443", "tls-cert.pem", "tls-key.pem", nil)
}
