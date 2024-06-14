CREATE TABLE IF NOT EXISTS district (
	id VARCHAR(8) NOT NULL,
	regency_id VARCHAR(5) NOT NULL,
	name VARCHAR(255) NOT NULL,

	CONSTRAINT district_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS district_regency_id_idx ON district USING btree (regency_id);

