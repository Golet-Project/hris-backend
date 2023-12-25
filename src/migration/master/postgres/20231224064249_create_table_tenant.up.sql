CREATE TABLE IF NOT EXISTS tenant (
	id SERIAL NOT NULL,
	uid UUID NOT NULL DEFAULT uuid_generate_v4(),
	name VARCHAR(100) NOT NULL,
	domain VARCHAR(50) NOT NULL,
	latitude DOUBLE PRECISION NULL,
	longitude DOUBLE PRECISION NULL,
	timezone SMALLINT NULL,

	CONSTRAINT tenant_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS tenant_domain_key ON tenant USING btree (domain);
CREATE UNIQUE INDEX IF NOT EXISTS tenant_uid_key ON tenant USING btree (uid);