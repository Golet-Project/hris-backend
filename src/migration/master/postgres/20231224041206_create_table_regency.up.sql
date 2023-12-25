CREATE TABLE IF NOT EXISTS regency (
	id VARCHAR(5) NOT NULL,
	province_id VARCHAR(2) NOT NULL,
	name VARCHAR(255) NOT NULL,

	CONSTRAINT regency_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS regency_province_id_idx ON regency USING btree (province_id);