CREATE TABLE IF NOT EXISTS tenant_admin (
	id BIGSERIAL NOT NULL,
	uid UUID NOT NULL DEFAULT uuid_generate_v4(),
	email VARCHAR(63) NOT NULL,
	password TEXT NOT NULL,
	full_name VARCHAR(127) NOT NULL,
	domain VARCHAR(63) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ NULL,

	CONSTRAINT tenant_admin_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS admin_email_domain_deleted_at_key ON tenant_admin USING btree (email, domain, deleted_at);