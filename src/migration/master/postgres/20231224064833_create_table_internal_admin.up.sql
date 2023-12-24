CREATE TABLE IF NOT EXISTS internal_admin (
	id BIGSERIAL NOT NULL,
	uid UUID NOT NULL DEFAULT uuid_generate_v4(),
	email VARCHAR(100) NOT NULL,
	password TEXT NOT NULL,
	full_name VARCHAR(100) NOT NULL,
	birth_date DATE NOT NULL,
	profile_picture TEXT NOT NULL,

	CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS internal_admin_email_key ON internal_admin (email)