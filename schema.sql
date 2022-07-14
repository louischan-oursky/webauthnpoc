CREATE TABLE users (
	id text PRIMARY KEY,
	name text NOT NULL,
	credential_id text NOT NULL,
	credential jsonb NOT NULL,
	sign_count bigint NOT NULL,
	created_at timestamp with time zone NOT NULL,
	UNIQUE(credential_id)
);

CREATE TABLE sessions (
	challenge text PRIMARY KEY,
	create_options jsonb,
	get_options jsonb,
	created_at timestamp with time zone NOT NULL
)
