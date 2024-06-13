CREATE TABLE IF NOT EXISTS village (
	id VARCHAR(13) NOT NULL,
	district_id VARCHAR(8) NOT NULL,
	name VARCHAR(255) NOT NULL,

	CONSTRAINT village_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS village_district_id_idx ON village USING btree (district_id);

