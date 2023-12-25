CREATE TABLE IF NOT EXISTS employee (
	id BIGSERIAL NOT NULL,
	uid UUID NOT NULL DEFAULT uuid_generate_v4(),
	email VARCHAR(63) NOT NULL,
	full_name VARCHAR(255) NOT NULL,
	gender VARCHAR(1) NOT NULL,
	birth_date DATE NOT NULL,
	profile_picture TEXT NOT NULL DEFAULT '',
	address TEXT NOT NULL DEFAULT '',
	province_id VARCHAR(2) NOT NULL DEFAULT '',
	regency_id VARCHAR(5) NOT NULL DEFAULT '',
	district_id VARCHAR(8) NOT NULL DEFAULT '',
	village_id VARCHAR(13) NOT NULL DEFAULT '',
	join_date DATE NOT NULL,
	end_date DATE NULL,
	employee_status VARCHAR(20) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT employee_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS employee_email_key ON employee(email);
CREATE UNIQUE INDEX IF NOT EXISTS employee_uid_key ON employee(uid);