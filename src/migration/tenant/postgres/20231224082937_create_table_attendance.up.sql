CREATE TABLE IF NOT EXISTS attendance (
	id BIGSERIAL NOT NULL,
	uid UUID NOT NULL DEFAULT uuid_generate_v4(),
	employee_uid UUID NOT NULL,
	timezone SMALLINT NOT NULL,
	latitude DOUBLE PRECISION NOT NULL,
	longitude DOUBLE PRECISION NOT NULL,
	radius INTEGER NOT NULL,
	approved_at TIMESTAMPTZ NULL,
	approved_by UUID NULL,
	checkin_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	checkout_time TIMESTAMPTZ NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NULL,

	CONSTRAINT attendance_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS attendance_employee_uid_idx ON attendance(employee_uid);