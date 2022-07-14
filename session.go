package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Session struct {
	Challenge     string
	CreateOptions *CreateOptions
	GetOptions    *GetOptions
	CreatedAt     time.Time
}

func NewSessionWithCreateOptions(createOptions *CreateOptions) (*Session, error) {
	challenge := Base64URLEncode(createOptions.PublicKey.Challenge)
	session := &Session{
		Challenge:     challenge,
		CreateOptions: createOptions,
		CreatedAt:     time.Now().UTC(),
	}
	return session, nil
}

func NewSessionWithGetOptions(getOptions *GetOptions) (*Session, error) {
	challenge := Base64URLEncode(getOptions.PublicKey.Challenge)
	session := &Session{
		Challenge:  challenge,
		GetOptions: getOptions,
		CreatedAt:  time.Now().UTC(),
	}
	return session, nil
}

func SaveSession(ctx context.Context, tx *sql.Tx, session *Session) (err error) {
	createBytes := []byte("null")
	getBytes := []byte("null")

	if session.CreateOptions != nil {
		createBytes, err = json.Marshal(session.CreateOptions)
		if err != nil {
			return
		}
	}
	if session.GetOptions != nil {
		getBytes, err = json.Marshal(session.GetOptions)
		if err != nil {
			return
		}
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO sessions (
		challenge,
		create_options,
		get_options,
		created_at
	) VALUES ($1, $2, $3, $4)
	`, session.Challenge, createBytes, getBytes, session.CreatedAt)
	if err != nil {
		return
	}

	return
}

func GetSession(ctx context.Context, tx *sql.Tx, challenge string) (*Session, error) {
	row := tx.QueryRowContext(ctx, `
	SELECT challenge, create_options, get_options, created_at
	FROM sessions
	WHERE challenge = $1
	`, challenge)

	var session Session
	var createBytes []byte
	var getBytes []byte

	err := row.Scan(
		&session.Challenge,
		&createBytes,
		&getBytes,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if createBytes != nil {
		err = json.Unmarshal(createBytes, &session.CreateOptions)
		if err != nil {
			return nil, err
		}
	}

	if getBytes != nil {
		err = json.Unmarshal(getBytes, &session.GetOptions)
		if err != nil {
			return nil, err
		}
	}

	return &session, nil
}
