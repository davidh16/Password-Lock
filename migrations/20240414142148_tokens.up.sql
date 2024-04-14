CREATE TABLE IF NOT EXISTS tokens (
	uuid uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	/*----------------------------*/
    -- extras
    user_uuid     uuid,
    token          varchar,
    token_type     varchar,
    expire_at      timestamp with time zone not null,
    is_used        timestamp with time zone
);
