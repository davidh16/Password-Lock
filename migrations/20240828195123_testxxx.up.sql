CREATE TABLE IF NOT EXISTS testxxx (
	uuid uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	/*----------------------------*/
    -- extras

);
