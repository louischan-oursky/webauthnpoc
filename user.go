package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
)

var ErrInvalidSignCount = fmt.Errorf("invalid signCount detected")

type User struct {
	ID           string
	Name         string
	CredentialID string
	Credential   []byte
	SignCount    int64
	CreatedAt    time.Time
}

func (u *User) WebAuthnCredential() (*webauthn.Credential, error) {
	parsed, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(u.Credential))
	if err != nil {
		return nil, err
	}
	return webauthn.MakeNewCredential(parsed)
}

func (u *User) IncrementSignCount(signCount int64) error {
	if u.SignCount <= 0 && signCount <= 0 {
		return nil
	}

	if signCount <= u.SignCount {
		return ErrInvalidSignCount
	}

	u.SignCount = signCount
	return nil
}

func NewUserWithName(name string) *User {
	return &User{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
}

func NewUser(
	credential []byte,
	createOptions *CreateOptions,
	parsed *protocol.ParsedCredentialCreationData,
) *User {
	return &User{
		ID:           string(createOptions.PublicKey.User.ID),
		Name:         createOptions.PublicKey.User.Name,
		CredentialID: parsed.ID,
		Credential:   credential,
		SignCount:    int64(parsed.Response.AttestationObject.AuthData.Counter),
		CreatedAt:    time.Now().UTC(),
	}
}

func InsertUser(ctx context.Context, tx *sql.Tx, user *User) (err error) {
	_, err = tx.ExecContext(ctx, `
	INSERT INTO users (
		id,
		name,
		credential_id,
		credential,
		sign_count,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.Name, user.CredentialID, user.Credential, user.SignCount, user.CreatedAt)
	if err != nil {
		return
	}

	return
}

func UpdateUser(ctx context.Context, tx *sql.Tx, user *User) (err error) {
	_, err = tx.ExecContext(ctx, `
	UPDATE users
	SET sign_count = $1
	WHERE id = $2
	`, user.SignCount, user.ID)
	if err != nil {
		return
	}

	return
}

func FindUserWithName(ctx context.Context, tx *sql.Tx, name string) (*User, error) {
	row := tx.QueryRowContext(ctx, `
	SELECT id, name, credential_id, credential, sign_count, created_at
	FROM users
	WHERE name = $1
	`, name)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CredentialID,
		&user.Credential,
		&user.SignCount,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserWithCredentialID(ctx context.Context, tx *sql.Tx, userID string, credentialID string) (*User, error) {
	row := tx.QueryRowContext(ctx, `
	SELECT id, name, credential_id, credential, sign_count, created_at
	FROM users
	WHERE id = $1 AND credential_id = $2
	`, userID, credentialID)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CredentialID,
		&user.Credential,
		&user.SignCount,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
